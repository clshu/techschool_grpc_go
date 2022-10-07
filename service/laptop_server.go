package service

import (
	"bytes"
	"context"
	"io"
	"learngrpc/pcbook/pb"
	"log"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const maxImageSize = 1 << 20 // 1 MB

// LaptopServer is a service that provides laptop services.
type LaptopServer struct {
	laptopStore LaptopStore
	imageStore ImageStore
	ratingStore RatingStore
	pb.UnimplementedLaptopServiceServer
}

// NewLaptopServer creates a new LaptopServer.
func NewLaptopServer(laptopStore LaptopStore, imageStore ImageStore, ratingStore RatingStore) *LaptopServer {
	return &LaptopServer{laptopStore, imageStore, ratingStore, pb.UnimplementedLaptopServiceServer{}}
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
		
		err := contexError(ctx)
		if err != nil {
			return nil, err
		}

		// save the laptop to the store
		// ... in memory store for now
		err = s.laptopStore.Save(laptop)
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
		err := s.laptopStore.Search(
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

// UploadImage is client streaming RPC to upload a laptop image.
func (s *LaptopServer) UploadImage(stream pb.LaptopService_UploadImageServer) error {
	req, err := stream.Recv()
	if err != nil {
		return logError(status.Errorf(codes.Unknown, "cannot receive image info: %v", err))
	}

	laptopID := req.GetInfo().GetLaptopId()
	imageType := req.GetInfo().GetImageType()
	log.Printf("received an ipload-image request for laptop %s with imageType %s", laptopID, imageType)

	laptop, err := s.laptopStore.Find(laptopID)
	if err != nil {
		return logError(status.Errorf(codes.Internal, "laptop store internal error: %v", err))
	}

	if (laptop == nil) {
		return logError(status.Errorf(codes.NotFound, "laptop not found: %v", laptopID))
	}

	imageData := bytes.Buffer{}
	imageSize := 0
	bulkSize := 10 * 1024

	for {
		// check if the context is canceled or deadline is exceeded
		err := contexError(stream.Context())
		if err != nil {
			return err
		}

		if imageSize % bulkSize == 0 {
			log.Printf("image size: %d", imageSize)
		}
		
		req, err := stream.Recv()
	
			if err == io.EOF {
				log.Printf("finished receiving image data: %d", imageSize)
				break
			}
			if err != nil {
				return logError(status.Errorf(codes.Unknown, "cannot receive image data: %v", err))
			}

			chunk := req.GetChunkData()
			size := len(chunk)
			imageSize += size
			if imageSize > maxImageSize {
				return logError(status.Errorf(codes.InvalidArgument, "image size is too large: %d > %d", imageSize, maxImageSize))
			}

			// write slowly
			// time.Sleep(1 * time.Second)

			_, err = imageData.Write(chunk)
			if err != nil {
				return logError(status.Errorf(codes.Internal, "cannot write image data: %v", err))
			}
	}

	imageID, err := s.imageStore.Save(laptopID, imageType, imageData)
	if err != nil {
		return logError(status.Errorf(codes.Internal, "cannot save image to store: %v", err))
	}

	res := &pb.UploadImageResponse{
		Id: imageID,
		Size: uint32(imageSize),
	}

	err = stream.SendAndClose(res)
	if err != nil {
		return logError(status.Errorf(codes.Unknown, "cannot send response: %v", err))
	}
	log.Printf("saved image with id: %s size: %d", imageID, imageSize)
	return nil;
}

// RateLaptop is a server streaming RPC to rate a laptop.
func (s *LaptopServer) RateLaptop(stream pb.LaptopService_RateLaptopServer) error {
	for {
		err := contexError(stream.Context())
		if err != nil {
			return err
		}

		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return logError(status.Errorf(codes.Unknown, "cannot receive client stream request: %v", err))
		}

		laptopID := req.GetLaptopId()
		score := req.GetScore()
		log.Printf("received rate-laptop request with score: %f for laptop %s", score, laptopID)

		found, err := s.laptopStore.Find(laptopID)
		if err != nil {
			return logError(status.Errorf(codes.Internal, "laptop store internal error: %v", err))
		}
		if (found == nil) {
			return logError(status.Errorf(codes.NotFound, "laptop not found: %v", laptopID))
		}

		rating, err := s.ratingStore.Add(laptopID, score)
		if err != nil {
			return logError(status.Errorf(codes.Internal, "rating store internal error: %v", err))
		}

		res := &pb.RateLaptopResponse{
			LaptopId: laptopID,
			RatedCount: rating.Count,
			AverageScore: float64(rating.Sum) / float64(rating.Count),
		}

		err = stream.Send(res)
		if err != nil {
			return logError(status.Errorf(codes.Unknown, "cannot send rate-laptop response: %v", err))
		}
		log.Printf("sent rate-laptop response with count %d, average score: %f", res.RatedCount, res.AverageScore)
	}


	return nil
}

func logError(err error) error {
	if err != nil {
		log.Print(err)
	}

	return err
}

func contexError(ctx context.Context) error {
	switch ctx.Err() {
	case context.Canceled:
		return logError(status.Error(codes.Canceled, "request is canceled"))
	case context.DeadlineExceeded:
		return logError(status.Error(codes.DeadlineExceeded, "deadline is exceeded"))
	default:
		return nil
	}
}	

