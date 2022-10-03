package main

import (
	"context"
	"flag"
	"io"
	"learngrpc/pcbook/pb"
	sample "learngrpc/pcbook/samples"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func createLaptop(laptopClient pb.LaptopServiceClient) {
	laptop := sample.NewLaptop()
		laptop.Id = ""

		req := &pb.CreateLaptopRequest{
			Laptop: laptop,
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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

func searchLaptop(laptopClient pb.LaptopServiceClient, filter *pb.LaptopFilter) {
	log.Printf("search filter: %v", filter)
	req := &pb.SearchLaptopRequest{Filter: filter}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := laptopClient.SearchLaptop(ctx, req)
	if err != nil {
		log.Fatalf("cannot search laptop: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatalf("cannot receive response: %v", err)
		}	

		logLaptop(res.GetLaptop())
	}
}

func logLaptop(laptop *pb.Laptop) {
	log.Print("_ found: ", laptop.GetId())
	log.Print(" + brand: ", laptop.GetBrand())
	log.Print(" + name: ", laptop.GetName())
	log.Print(" + cpu cores: ", laptop.GetCpu().GetNumCores())
	log.Print(" + cpu minGhz: ", laptop.GetCpu().GetMinGhz())
	log.Print(" + cpu maxGhz: ", laptop.GetCpu().GetMaxGhz())
	log.Print(" + ram: ", laptop.GetMemory().GetValue(), laptop.GetMemory().GetUnit())
	log.Print(" + price: ", laptop.GetPriceUsd(), "usd")
}

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
	
		// create multiple new laptops on the server side
		for i := 0; i < 10; i++ {
			createLaptop(laptopClient)
		}

		// search for laptops
		filter := &pb.LaptopFilter{
			MaxPriceUsd: 3000,
			MinCpuCores: 4,
			MinCpuGhz: 2.5,
			MinMemory: &pb.Memory{
				Value: 8,
				Unit: pb.Memory_GIGABYTE,
			},
		}

		searchLaptop(laptopClient, filter)
}