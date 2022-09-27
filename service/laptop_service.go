package service

import (
	"context"
	"learngrpc/pcbook/pb"
	"log"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LaptopServer is a service that provides laptop services.
type LaptopServer struct {
	Store LaptopStore
	pb.UnimplementedLaptopServiceServer
}

// NewLaptopServer creates a new LaptopServer.
func NewLaptopServer(store LaptopStore) *LaptopServer {
	return &LaptopServer{store, pb.UnimplementedLaptopServiceServer{}}
}

// CreateLaptop creates a new laptop.
func (s *LaptopServer) CreateLaptop(
	ctx context.Context,
	req *pb.CreateLaptopRequest,
	) (*pb.CreateLaptopResponse, error) {
		laptop := req.GetLaptop()
		id := laptop.GetId()
		log.Printf("received a create-laptop request with id: %s", id)
		
		if len(id) > 0 {
			// check if the laptop id is valid
			_, err := uuid.Parse(id)
			if err != nil {
				return nil, status.Error(codes.InvalidArgument, "laptop id is not valid UUID: %v")
			}
		} else {
			id, err := uuid.NewRandom()
			if err != nil {
				return nil, status.Error(codes.Internal, "cannot generate laptop ID: %v")
			}
			laptop.Id = id.String()

		}
		// save the laptop to the store
		// ... in memory store for now
		err := s.Store.Save(laptop)
		if err != nil {
			code := codes.Internal
			if err == ErrAlreadyExists {
				code = codes.AlreadyExists
			}
			return nil, status.Error(code, "cannot save laptop to store: %v")
		}
		log.Printf("saved laptop with id: %s", laptop.Id)
		// create a response
		res := &pb.CreateLaptopResponse{
			Id: laptop.Id,
		}
		return res, nil
}
