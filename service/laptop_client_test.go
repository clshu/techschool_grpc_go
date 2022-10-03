package service_test

import (
	"context"
	"io"
	"learngrpc/pcbook/pb"
	sample "learngrpc/pcbook/samples"
	"learngrpc/pcbook/serializer"
	"learngrpc/pcbook/service"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestClientCreateLaptop(t *testing.T) {
	t.Parallel()

	laptopServer, serverAddress := startTestLaptopServer(t, service.NewInMemoryLaptopStore())
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
	other, err := laptopServer.Store.Find(res.Id)
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
		MinMemory: &pb.Memory{
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
			laptop.Memory = &pb.Memory{
				Value: 4096,
				Unit: pb.Memory_MEGABYTE,
			}
		case 4:
			  laptop.PriceUsd = 1999
				laptop.Cpu.NumCores = 4
				laptop.Cpu.MinGhz = 2.5
				laptop.Cpu.MaxGhz = 4.5
				laptop.Memory = &pb.Memory{
					Value: 16,
					Unit: pb.Memory_GIGABYTE,
				}
				expectedIDs[laptop.Id] = true
		case 5:
			laptop.PriceUsd = 2000
			laptop.Cpu.NumCores = 6
				laptop.Cpu.MinGhz = 2.8
				laptop.Cpu.MaxGhz = 5.0
				laptop.Memory = &pb.Memory{
					Value: 64,
					Unit: pb.Memory_GIGABYTE,
				}
				expectedIDs[laptop.Id] = true
		}

		err := store.Save(laptop)
		require.NoError(t, err)
	}

	_, serverAddress := startTestLaptopServer(t, store)
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

func startTestLaptopServer(t *testing.T, store service.LaptopStore) (*service.LaptopServer, string) {
	laptopServer := service.NewLaptopServer(store)
	grpcServer := grpc.NewServer()
	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)

	listener, err := net.Listen("tcp", ":0") // random available port
	require.Nil(t, err)

	go grpcServer.Serve(listener) // blocking call

	return laptopServer, listener.Addr().String()
	
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