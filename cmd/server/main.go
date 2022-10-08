package main

import (
	"context"
	"flag"
	"fmt"
	"learngrpc/pcbook/pb"
	"learngrpc/pcbook/service"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func unaryInterceotor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	log.Println("--> unary interceptor: ", info.FullMethod)
	return handler(ctx, req)
}

func streamInterceptor(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	log.Println("--> stream interceptor: ", info.FullMethod)
	return handler(srv, ss)
}

func main() {
		port := flag.Int("port", 0, "the server port")
		flag.Parse()
		log.Printf("start server on port %d", *port)

		laptopStore := service.NewInMemoryLaptopStore()
		imageStore := service.NewDiskImageStore("img")
		ratingStore := service.NewInMemoryRatingStore()

		// create a new server
		laptopServer := service.NewLaptopServer(laptopStore, imageStore, ratingStore)
		grpcServer := grpc.NewServer(
			grpc.UnaryInterceptor(unaryInterceotor),
			grpc.StreamInterceptor(streamInterceptor),
		)
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