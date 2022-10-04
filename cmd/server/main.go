package main

import (
	"flag"
	"fmt"
	"learngrpc/pcbook/pb"
	"learngrpc/pcbook/service"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
		port := flag.Int("port", 0, "the server port")
		flag.Parse()
		log.Printf("start server on port %d", *port)

		// create a new server
		laptopServer := service.NewLaptopServer(service.NewInMemoryLaptopStore(), nil)
		grpcServer := grpc.NewServer()
		pb.RegisterLaptopServiceServer(grpcServer, laptopServer)

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