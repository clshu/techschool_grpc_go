package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"learngrpc/pcbook/client"
	"learngrpc/pcbook/pb"
	sample "learngrpc/pcbook/samples"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func testCreateLaptop(laptopClient *client.LaptopClient) {
	laptopClient.CreateLaptop(sample.NewLaptop())
}

func testSearchLaptop(laptopClient *client.LaptopClient) {
	for i := 0; i < 10; i++ {
		laptopClient.CreateLaptop(sample.NewLaptop())
	}
	// search for laptops
	filter := &pb.LaptopFilter{
		MaxPriceUsd: 3000,
		MinCpuCores: 4,
		MinCpuGhz:   2.5,
		MinRam: &pb.Memory{
			Value: 8,
			Unit:  pb.Memory_GIGABYTE,
		},
	}

	laptopClient.SearchLaptop(filter)
}

func testUploadImage(laptopClient *client.LaptopClient) {
	imagePath := "from/macbook-air-gold-2015-16.jpg"
	laptop := sample.NewLaptop()
	laptopClient.CreateLaptop(laptop)
	laptopClient.UploadImage(laptop.GetId(), imagePath)
}

func testRateLaptop(laptopClient *client.LaptopClient) {
	n := 3
	laptopIDs := make([]string, n)
	for i := 0; i < n; i++ {
		laptop := sample.NewLaptop()
		laptopIDs[i] = laptop.GetId()
		laptopClient.CreateLaptop(laptop)
	}
	scores := make([]float64, n)

	for i := 0; i < n; i++ {
		fmt.Print("rate laptop (y/n)? ")
		var answer string
		fmt.Scanln(&answer)
		if answer != "y" {
			break
		}
		for i := 0; i < n; i++ {
			scores[i] = sample.RandomLaptopScore()
		}
		err := laptopClient.RateLaptop(laptopIDs, scores)
		if err != nil {
			log.Fatalf("cannot rate laptop: %v", err)
		}
	}
}

const (
	// change username to user1 to test persion denied
	username        = "admin1"
	password        = "secret"
	refreshDuration = 30 * time.Second
	caCertFile      = "cert/ca-cert.pem"
	certFile        = "cert/client-cert.pem"
	keyFile         = "cert/client-key.pem"
)

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// load CA certificate
	pemServerCA, err := ioutil.ReadFile(caCertFile)
	if err != nil {
		return nil, err
	}

	// create a certificate pool from CA certificate
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("cannot add server CA's certificate")
	}

	// load cleint certificate and private key
	// cleintCert, err := tls.LoadX509KeyPair(certFile, keyFile)
	// if err != nil {
	// 	return nil, err
	// }

	config := &tls.Config{
		// Certificates: []tls.Certificate{cleintCert},
		RootCAs: certPool,
	}

	return credentials.NewTLS(config), nil
}

func main() {

	serverAddress := flag.String("address", "", "the server address")
	enableTLS := flag.Bool("tls", false, "enable TLS for RPC")
	flag.Parse()
	log.Printf("dial server %s TLS = %t", *serverAddress, *enableTLS)

	transportOption := grpc.WithInsecure()
	if *enableTLS {
		tlsCredentials, err := loadTLSCredentials()
		if err != nil {
			log.Fatalf("cannot load TLS credentials: %v", err)
		}
		transportOption = grpc.WithTransportCredentials(tlsCredentials)
	}

	// create a new gRPC client
	conn, err := grpc.Dial(*serverAddress, transportOption)
	if err != nil {
		log.Fatalf("cannot dial server: %v", err)
	}

	authMethods := client.NewAuthMethods()
	authClient := client.NewAuthClient(conn, username, password)

	interceptor, err := client.NewAuthInterceptor(authClient, authMethods, refreshDuration)
	if err != nil {
		log.Fatalf("cannot create auth interceptor: %v", err)
	}

	opts := []grpc.DialOption{
		grpc.WithUnaryInterceptor(interceptor.Unary()),
		grpc.WithStreamInterceptor(interceptor.Stream()),
		transportOption,
	}

	conn2, err := grpc.Dial(*serverAddress, opts...)

	if err != nil {
		log.Fatalf("cannot dial server: %v", err)
	}

	// create a new LaptopClient
	laptopClient := client.NewLaptopClient(conn2)
	// testCreateLaptop(laptopClient)
	// testSearchLaptop(laptopClient)
	// testUploadImage(laptopClient)
	testRateLaptop(laptopClient)
}
