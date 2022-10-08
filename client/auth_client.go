package client

import (
	"context"
	"learngrpc/pcbook/pb"

	"google.golang.org/grpc"
)

// AuthClient is the client for authentication service.
type AuthClient struct {
	service pb.AuthServiceClient
	username string
	password string
}

// NewAuthClient creates a new AuthClient.
func NewAuthClient(cc *grpc.ClientConn, username, password string) *AuthClient {
	service := pb.NewAuthServiceClient(cc)
	return &AuthClient{
		service, username, password,
	}
}

// Login logs in and returns the token.
func (client *AuthClient) Login() (string, error) {
	req := &pb.LoginRequest{
		Username: client.username,
		Password: client.password,
	}

	res, err := client.service.Login(context.Background(), req)
	if err != nil {
		return "", err
	}

	return res.GetToken(), nil
}