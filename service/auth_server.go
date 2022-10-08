package service

import (
	"context"
	"learngrpc/pcbook/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AuthServer is the server for authentication service.
type AuthServer struct {
	userStore UserStore
	jwtManager *JWTManager
	pb.UnimplementedAuthServiceServer
}

// NewAuthServer creates a new AuthServer.
func NewAuthServer(userStore UserStore, jwtManager *JWTManager) *AuthServer {
	return &AuthServer{
		userStore: userStore,
		jwtManager: jwtManager,
	}
}

// Login is a unary RPC to login.
func (server *AuthServer)	Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := server.userStore.Find(req.GetUsername())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot find user: %v", err)
	}

	if user == nil || !user.IsCorrectPassword(req.GetPassword()) {
		return nil, status.Errorf(codes.NotFound, "invalid username or password")
	}

	token, err := server.jwtManager.Generate(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot generate token: %v", err)
	}

	res := &pb.LoginResponse{
		Token: token,
	}	

	return res, nil
}
