package service_test

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"learngrpc/pcbook/pb"
	sample "learngrpc/pcbook/samples"
	"learngrpc/pcbook/serializer"
	"learngrpc/pcbook/service"
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestClientCreateLaptop(t *testing.T) {
	t.Parallel()

	laptopStore := service.NewInMemoryLaptopStore()
	serverAddress := startTestLaptopServer(t, laptopStore, nil)
	laptopClient := newTestLaptopClient(t, serverAddress)

	laptop := sample.NewLaptop()
	expectedID := laptop.Id
	req := &pb.CreateLaptopRequest{
		Laptop: laptop,
	}

	res, err := laptopClient.CreateLaptop(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, expectedID, res.Id)

	// check if the laptop is saved to the store
	other, err := laptopStore.Find(res.Id)
	require.NoError(t, err)
	require.NotNil(t, other)

	// check if the other laptop is the same as the one we sent
	requireSameLaptop(t, laptop, other)
}

func TestClientSearchLaptop(t *testing.T) {
	t.Parallel()

	filter := &pb.LaptopFilter{
		MaxPriceUsd: 2000,
		MinCpuCores: 4,
		MinCpuGhz: 2.2,
		MinRam: &pb.Memory{
			Value: 8,
			Unit: pb.Memory_GIGABYTE,
		},
	}

	store := service.NewInMemoryLaptopStore()
	expectedIDs := make(map[string]bool)

	for i:=0; i<6; i++ {
		laptop := sample.NewLaptop()

		switch i {
		case 0:
			laptop.PriceUsd = 2500
		case 1:
			laptop.Cpu.NumCores = 2
		case 2:
			laptop.Cpu.MinGhz = 2.0
		case 3:
			// 4 Gigabytes
			laptop.Ram = &pb.Memory{
				Value: 4096,
				Unit: pb.Memory_MEGABYTE,
			}
		case 4:
			  laptop.PriceUsd = 1999
				laptop.Cpu.NumCores = 4
				laptop.Cpu.MinGhz = 2.5
				laptop.Cpu.MaxGhz = 4.5
				laptop.Ram = &pb.Memory{
					Value: 16,
					Unit: pb.Memory_GIGABYTE,
				}
				expectedIDs[laptop.Id] = true
		case 5:
			laptop.PriceUsd = 2000
			laptop.Cpu.NumCores = 6
				laptop.Cpu.MinGhz = 2.8
				laptop.Cpu.MaxGhz = 5.0
				laptop.Ram = &pb.Memory{
					Value: 64,
					Unit: pb.Memory_GIGABYTE,
				}
				expectedIDs[laptop.Id] = true
		}

		err := store.Save(laptop)
		require.NoError(t, err)
	}

	serverAddress := startTestLaptopServer(t, store, nil)
	laptopClient := newTestLaptopClient(t, serverAddress)

	req := &pb.SearchLaptopRequest{
		Filter: filter,
	}
	stream, err := laptopClient.SearchLaptop(context.Background(), req)
	require.NoError(t, err)

	found := 0
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		require.Contains(t, expectedIDs, res.GetLaptop().GetId())

		found++
	}
	require.Equal(t, len(expectedIDs), found)
}

func TestClientUploadImage(t *testing.T) {
	t.Parallel()

	fileName := "macbook-air-gold-2015-16.jpg"
	from := "../from/" + fileName
	to := "../tmp/" + fileName

	// Copy the file
	err := copyFile(from, to)
	require.NoError(t, err)

	testImageFolder := "../tmp"
	imageStore := service.NewDiskImageStore(testImageFolder)
	laptopStore := service.NewInMemoryLaptopStore()

	laptop := sample.NewLaptop()
	err = laptopStore.Save(laptop)
	require.NoError(t, err)

	serverAddress := startTestLaptopServer(t, laptopStore, imageStore)
	laptopClient := newTestLaptopClient(t, serverAddress)

	

	uploadImageTest(t, laptopClient, laptop.Id, to, testImageFolder)

	// clean up
	require.NoError(t, os.Remove(to))
}

func startTestLaptopServer(t *testing.T, laptopStore service.LaptopStore, imageStore service.ImageStore) string {
	laptopServer := service.NewLaptopServer(laptopStore, imageStore)
	grpcServer := grpc.NewServer()
	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)

	listener, err := net.Listen("tcp", ":0") // random available port
	require.Nil(t, err)

	go grpcServer.Serve(listener) // blocking call

	return listener.Addr().String()
	
}

func newTestLaptopClient(t *testing.T, serverAddress string) pb.LaptopServiceClient {
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	require.Nil(t, err)

	return pb.NewLaptopServiceClient(conn)
}

func requireSameLaptop(t *testing.T, expected, actual *pb.Laptop) {
	json1, err := serializer.ProtobufToJSON(expected)
	require.NoError(t, err)
	require.NotEmpty(t, json1)
	json2, err := serializer.ProtobufToJSON(actual)
	require.NoError(t, err)
	require.NotEmpty(t, json2)
	
	require.Equal(t, json1, json2)
}

func uploadImageTest(t *testing.T, laptopClient pb.LaptopServiceClient, laptopID string, imagePath string, testImageFolder string) {
	file, err := os.Open(imagePath)
	defer file.Close()
	require.NoError(t, err)

	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	stream, err := laptopClient.UploadImage(context.Background())
	require.NoError(t, err)

	imageType := filepath.Ext(imagePath)

	req := &pb.UploadImageRequest{
		Data: &pb.UploadImageRequest_Info{
			Info: &pb.ImageInfo{
				LaptopId: laptopID,
				ImageType: imageType,
			},
		},
	}

	err = stream.Send(req)
	require.NoError(t, err)

	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)
	size := 0

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		size += n

		req := &pb.UploadImageRequest{
			Data: &pb.UploadImageRequest_ChunkData{
				ChunkData: buffer[:n],
			},
		}

		err = stream.Send(req)
		require.NoError(t, err)
	}

	res, err := stream.CloseAndRecv()
	require.NoError(t, err)
	require.NotEmpty(t, res.GetId())
	require.Equal(t, size, int(res.GetSize()))

	savedImaePath := fmt.Sprintf("%s/%s%s", testImageFolder, res.GetId(), imageType)
	require.FileExists(t, savedImaePath)
	require.NoError(t, os.Remove(savedImaePath))
}

func copyFile(from, to string) error {
	  sfi, err := os.Stat(from)
		if err != nil {
			return err
		}
		if !sfi.Mode().IsRegular() {
        // cannot copy non-regular files (e.g., directories,
        // symlinks, devices, etc.)
        return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
    }
		dfi, err := os.Stat(to)
		if err != nil {
			if os.IsNotExist(err) {
				return os.Link(from, to)
			} else {
				os.Remove(to)
				return os.Link(from, to)
			}
		} else {
			if !(dfi.Mode().IsRegular()) {
				return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
			}
			if os.SameFile(sfi, dfi) {
				return nil
			}
		}
		return nil
}