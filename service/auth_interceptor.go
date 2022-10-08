package service

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// AuthInterceptor is a server interceptor for authentication and authrization.
type AuthInterceptor struct {
	jwtManager *JWTManager
	accessibleRols map[string][]string
}

// NewAuthInterceptor creates a new AuthInterceptor.
func NewAuthInterceptor(jwtManager *JWTManager, accessibleRols map[string][]string) *AuthInterceptor {
	return &AuthInterceptor{
		jwtManager,
		accessibleRols,
	}
}

// Unary returns a new unary server interceptor for authentication and authorization.
func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func (
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		log.Println("--> unary interceptor: ", info.FullMethod)

		err = interceptor.authorize(ctx, info.FullMethod) 
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

// Stream returns a new stream server interceptor for authentication and authorization.
func (interceptor *AuthInterceptor) Stream() grpc.StreamServerInterceptor {
	return func (
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		log.Println("--> stream interceptor: ", info.FullMethod)

		err := interceptor.authorize(stream.Context(), info.FullMethod) 
		if err != nil {
			return  err
		}

		return handler(srv, stream)
	}
}

func (interceptor *AuthInterceptor) authorize(ctx context.Context, fullMethod string) error {

	// check if the user has the role to access the RPC
	accessibleRoles, ok := interceptor.accessibleRols[fullMethod]
	if !ok {
		// no role is specified, public RPC
		return nil
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}
	
	accessibletoken := values[0]
	claims, err := interceptor.jwtManager.Verify(accessibletoken)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "accessibletoken is invalid: %v", err)
	}

	for _, role := range accessibleRoles {
		if role == claims.Role {
			return nil
		}
	}

	return status.Errorf(codes.PermissionDenied, "role %s cannot access %s", claims.Role, fullMethod)
}	

// NewAccessibleRoles creates a new map of accessible roles.
func NewAccessibleRoles() map[string][]string {
	packagePath := "techschool.pcbook"
	laptopServicePath := fmt.Sprintf("/%s.LaptopService/", packagePath)
	// if RPC calls are not in the same package, we need to use the full path
	// if calls are not in the map, they are public to everyone
	return map[string][]string{
		laptopServicePath + "CreateLaptop": {"admin"},
		laptopServicePath + "UploadImage": {"admin"},
		laptopServicePath + "RateLaptop": {"admin", "user"},
	}
}