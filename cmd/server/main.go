package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"learngrpc/pcbook/pb"
	"learngrpc/pcbook/service"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

func sendUsers(userStore service.UserStore) error {
	// create 2 users
	err := createUser(userStore, "admin1", "secret", "admin")
	if err != nil {
		return err
	}
	return createUser(userStore, "user1", "secret", "user")
}

func createUser(userStore service.UserStore, username, password, role string) error {
	user, err := service.NewUser(username, password, role)
	if err != nil {
		return err
	}

	return userStore.Save(user)
}

const (
	secretKey     = "longjohnsilver"
	tokenDuration = 15 * time.Minute
	certFile      = "cert/server-cert.pem"
	keyFile       = "cert/server-key.pem"
	caCertFile    = "cert/ca-cert.pem"
)

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// load CA certificate
	// pemClientCA, err := ioutil.ReadFile(caCertFile)
	// if err != nil {
	// 	return nil, err
	// }

	// // create a certificate pool from CA certificate
	// certPool := x509.NewCertPool()
	// if !certPool.AppendCertsFromPEM(pemClientCA) {
	// 	return nil, fmt.Errorf("cannot add client CA's certificate")
	// }

	// load server certificate and private key
	serverCert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	// create credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
		// ClientCAs:    certPool,
	}
	return credentials.NewTLS(config), nil
}

func main() {
	port := flag.Int("port", 0, "the server port")
	enableTLS := flag.Bool("tls", false, "enable TLS for RPC")
	flag.Parse()
	log.Printf("start server on port %d TLS = %t ", *port, *enableTLS)

	userStore := service.NewInMemoryUserStore()
	jwtManager := service.NewJWTManager(secretKey, tokenDuration)
	authServer := service.NewAuthServer(userStore, jwtManager)
	accessibleRoles := service.NewAccessibleRoles()
	authInterceptor := service.NewAuthInterceptor(jwtManager, accessibleRoles)

	laptopStore := service.NewInMemoryLaptopStore()
	imageStore := service.NewDiskImageStore("img")
	ratingStore := service.NewInMemoryRatingStore()
	err := sendUsers(userStore)
	if err != nil {
		log.Fatal(err)
	}

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(authInterceptor.Unary()),
		grpc.StreamInterceptor(authInterceptor.Stream()),
	}
	if *enableTLS {
		tlsCredentials, err := loadTLSCredentials()
		if err != nil {
			log.Fatal("cannot load TLS credentials: ", err)
		}
		opts = append(opts, grpc.Creds(tlsCredentials))
	}
	// create a new server
	laptopServer := service.NewLaptopServer(laptopStore, imageStore, ratingStore)
	grpcServer := grpc.NewServer(opts...)

	pb.RegisterAuthServiceServer(grpcServer, authServer)
	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)
	reflection.Register(grpcServer)

	address := fmt.Sprintf("localhost:%d", *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("cannot start server: %v", err)
	}

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("cannot start grpc server: %v", err)
	}

}
