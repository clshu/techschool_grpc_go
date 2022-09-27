package main

import (
	"context"
	"flag"
	"learngrpc/pcbook/pb"
	sample "learngrpc/pcbook/samples"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

		// create a new gRPC client
		laptopClient := pb.NewLaptopServiceClient(conn)
	
		// create a new laptop
		laptop := sample.NewLaptop()
		laptop.Id = ""

		req := &pb.CreateLaptopRequest{
			Laptop: laptop,
		}

		// set timeout to 1 seconds
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		res, err := laptopClient.CreateLaptop(ctx, req)
		if err != nil {
			st, ok := status.FromError(err)
			if ok && st.Code() == codes.AlreadyExists {
				log.Printf("laptop already exists")
			}	 else {
				log.Fatalf("cannot create laptop: %v", err)
			}
			return
		}
		
		log.Printf("created laptop with id: %s", res.Id)
}