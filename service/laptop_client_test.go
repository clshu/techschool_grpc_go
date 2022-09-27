package service_test

import (
	"context"
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

	laptopServer, serverAddress := startTestLaptopServer(t)
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

func startTestLaptopServer(t *testing.T) (*service.LaptopServer, string) {
	laptopServer := service.NewLaptopServer(service.NewInMemoryLaptopStore())
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