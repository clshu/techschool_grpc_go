package client

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"learngrpc/pcbook/pb"
	"log"
	"os"
	"path/filepath"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LaptopClient is the client for laptop service.
type LaptopClient struct {
	service pb.LaptopServiceClient
}

// NewLaptopClient creates a new LaptopClient.
func NewLaptopClient(cc *grpc.ClientConn) *LaptopClient {
	service := pb.NewLaptopServiceClient(cc)
	return &LaptopClient{service}
}

// CreateLaptop creates a new laptop and send it to the server.
func (laptopClient *LaptopClient) CreateLaptop(laptop *pb.Laptop) {
		
		req := &pb.CreateLaptopRequest{
			Laptop: laptop,
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		res, err := laptopClient.service.CreateLaptop(ctx, req)
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

// SearchLaptop searches for laptops that match the filter.
func (laptopClient *LaptopClient) SearchLaptop(filter *pb.LaptopFilter) {
	log.Printf("search filter: %v", filter)
	req := &pb.SearchLaptopRequest{Filter: filter}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := laptopClient.service.SearchLaptop(ctx, req)
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
	log.Print(" + ram: ", laptop.GetRam().GetValue(), laptop.GetRam().GetUnit())
	log.Print(" + price: ", laptop.GetPriceUsd(), "usd")
}

// UploadImage uploads an image for a laptop.
func (laptopClient *LaptopClient) UploadImage(laptopID string, imagePath string) {
	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatalf("cannot open image file: %v", err)
	}
	defer file.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := laptopClient.service.UploadImage(ctx)
	if err != nil {
		log.Fatalf("cannot upload image: %v", err)
	}

	req := &pb.UploadImageRequest{
		Data: &pb.UploadImageRequest_Info{
			Info: &pb.ImageInfo{
				LaptopId: laptopID,
				ImageType: filepath.Ext(imagePath),
			},
		},
	}

	err = stream.Send(req)
	if err != nil {
		log.Fatalf("cannot send image info: %v", err)
	}

	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)
	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("cannot read image file: %v", err)
		}

		req := &pb.UploadImageRequest{
			Data: &pb.UploadImageRequest_ChunkData{
				ChunkData: buffer[:n],
			},
		}

		err = stream.Send(req)
		if err != nil {
			err2 := stream.RecvMsg(nil)
			log.Fatalf("cannot send image chunk: %v| %v", err, err2)
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("cannot receive response: %v", err)
	}

	log.Printf("image uploaded with id: %s", res.GetId())	
}

// RateLaptop rates a laptop.
func (laptopClient *LaptopClient) RateLaptop(laptopIDs []string, scores []float64) error  {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := laptopClient.service.RateLaptop(ctx)
	if err != nil {
		log.Fatalf("cannot rate laptop: %v", err)
	}

	waitResponse := make(chan error)
	// go routine to receive async responses
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				log.Print("no more response")
				waitResponse <- nil
				return
			}
			if err != nil {
				waitResponse <- fmt.Errorf("cannot receive response: %v", err)
				return
			} 
			log.Printf("received response: %v", res)
		}
	}()

	// send requests
	for i, laptopID := range laptopIDs {
		req := &pb.RateLaptopRequest{
			LaptopId: laptopID,
			Score: scores[i],
		}

		err := stream.Send(req)
		if err != nil {
			log.Fatalf("cannot send rate request: %v, %v", err, stream.RecvMsg(nil))
		}
		log.Printf("sent request: %v", req)
	}

	err = stream.CloseSend()
	if err != nil {
		return fmt.Errorf("cannot close send: %v", err)
	}

	// Runtime return from go routines
	err = <- waitResponse
	return err
}
