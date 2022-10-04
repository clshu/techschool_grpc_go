package main

import (
	"bufio"
	"context"
	"flag"
	"io"
	"learngrpc/pcbook/pb"
	sample "learngrpc/pcbook/samples"
	"log"
	"os"
	"path/filepath"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func createLaptop(laptopClient pb.LaptopServiceClient, laptop *pb.Laptop) {
		
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
	log.Print(" + ram: ", laptop.GetRam().GetValue(), laptop.GetRam().GetUnit())
	log.Print(" + price: ", laptop.GetPriceUsd(), "usd")
}

func uploadImage(laptopClient pb.LaptopServiceClient, laptopID string, imagePath string) {
	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatalf("cannot open image file: %v", err)
	}
	defer file.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := laptopClient.UploadImage(ctx)
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

func testCreateLaptop(laptopClient pb.LaptopServiceClient) {
	createLaptop(laptopClient, sample.NewLaptop())
}

func testSearchLaptop(laptopClient pb.LaptopServiceClient) {
	for i := 0; i < 10; i++ {
		createLaptop(laptopClient, sample.NewLaptop())
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

	searchLaptop(laptopClient, filter)
}

func testUploadImage(laptopClient pb.LaptopServiceClient) {
	imagePath := "from/macbook-air-gold-2015-16.jpg"
	laptop := sample.NewLaptop()
	createLaptop(laptopClient, laptop)
	uploadImage(laptopClient, laptop.GetId(), imagePath)
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
	
		// testCreateLaptop(laptopClient)
		// testSearchLaptop(laptopClient)
		testUploadImage(laptopClient)		
}