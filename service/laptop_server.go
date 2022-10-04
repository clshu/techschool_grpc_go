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
		log.Printf("received a create-laptop request with id: %v", id)
		
		if len(id) > 0 {
			// check if the laptop id is valid
			_, err := uuid.Parse(id)
			if err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "laptop id is not valid UUID: %v", id)
			}
		} else {
			id, err := uuid.NewRandom()
			if err != nil {
				return nil, status.Error(codes.Internal, "cannot generate laptop ID")
			}
			laptop.Id = id.String()

		}

		// some heavy processing
		// time.Sleep(6 * time.Second)

		if ctx.Err() == context.Canceled {
			log.Print("request is canceled")
			return nil, status.Error(codes.Canceled, "request is canceled")
		}

		if ctx.Err() == context.DeadlineExceeded {
			log.Print("deadline is exceeded")
			return nil, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
		}


		// save the laptop to the store
		// ... in memory store for now
		err := s.Store.Save(laptop)
		if err != nil {
			code := codes.Internal
			if err == ErrAlreadyExists {
				code = codes.AlreadyExists
			}
			return nil, status.Errorf(code, "Duplicate ID: %v", laptop.Id)
		}
		log.Printf("saved laptop with id: %s", laptop.Id)
		// create a response
		res := &pb.CreateLaptopResponse{
			Id: laptop.Id,
		}
		return res, nil
}

// SearchLaptop searches for a laptop with filter.
func (s *LaptopServer) SearchLaptop(
	req *pb.SearchLaptopRequest,
	stream pb.LaptopService_SearchLaptopServer,
	) error {
		log.Print("received a search-laptop request")
		filter := req.GetFilter()
		err := s.Store.Search(
			stream.Context(),
			filter,
			func(laptop *pb.Laptop) error {
			res := &pb.SearchLaptopResponse{
				Laptop: laptop,
			}
			err := stream.Send(res)
			if err != nil {
				return err
			}
			log.Printf("sent laptop with id: %s", laptop.Id)
			return nil
		})
		
		if err != nil {
			return err
		}
		return nil
}

