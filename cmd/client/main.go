package main

import (
	"flag"
	"fmt"
	"learngrpc/pcbook/client"
	"learngrpc/pcbook/pb"
	sample "learngrpc/pcbook/samples"
	"log"
	"time"

	"google.golang.org/grpc"
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
		MinCpuGhz: 2.5,
		MinRam: &pb.Memory{
			Value: 8,
			Unit: pb.Memory_GIGABYTE,
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
	username = "admin1"
	password = "secret"
	refreshDuration = 30 * time.Second
)

func main() {
		
		serverAddress := flag.String("address", "", "the server address")
		flag.Parse()
		log.Printf("dial server %s", *serverAddress)

		// create a new gRPC client
		conn, err := grpc.Dial(*serverAddress, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("cannot dial server: %v", err)
		}

		authMethods := client.NewAuthMethods()
		authClient := client.NewAuthClient(conn, username, password)
		interceptor, err := client.NewAuthInterceptor(authClient, authMethods, refreshDuration)
		if err != nil {
			log.Fatalf("cannot create auth interceptor: %v", err)
		}

		conn2, err := grpc.Dial(
			*serverAddress, 
			grpc.WithUnaryInterceptor(interceptor.Unary()),
			grpc.WithStreamInterceptor(interceptor.Stream()),
			grpc.WithInsecure())
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