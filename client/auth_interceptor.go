package client

import (
	"context"
	"learngrpc/pcbook/service"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// AuthInterceptor is the client interceptor for authentication and authorization.
type AuthInterceptor struct {
	authClient *AuthClient
	authMethods map[string]bool
	accessToken string
}

// NewAuthInterceptor creates a new AuthInterceptor.
func NewAuthInterceptor(authClient *AuthClient, authMethods map[string]bool, refreshDuration time.Duration) (*AuthInterceptor, error) {
	interceptor := &AuthInterceptor{
		authClient,
		authMethods,
		"",
	}	

	err := interceptor.scheduleRefreshToken(refreshDuration)
	if err != nil {
		return nil, err
	}

	return interceptor, nil
}

func (interceptor *AuthInterceptor) scheduleRefreshToken(refreshDuration time.Duration) error {
	err := interceptor.refreshToken()
	if err != nil {
		return err
	}

	go func() {
		wait := refreshDuration
		for {
			time.Sleep(wait)
			err := interceptor.refreshToken()
			if err != nil {
				wait = time.Second
			} else {
				wait = refreshDuration
			}
		}
	}()

	return nil
}

func (interceptor *AuthInterceptor) refreshToken() error {
	token, err := interceptor.authClient.Login()
	if err != nil {
		return err
	}

	interceptor.accessToken = token

	log.Printf("refreshed token: %v", token)
	return nil
}

// Unary returns a new unary client interceptor for authentication and authorization.
func (interceptor *AuthInterceptor) Unary() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context, 
		method string, 
		req, reply interface{}, 
		cc *grpc.ClientConn, 
		invoker grpc.UnaryInvoker, 
		opts ...grpc.CallOption,
		) error {
		log.Printf("--> unary interceptor: %v", method)

		if interceptor.authMethods[method] {
			return invoker(interceptor.attachToken(ctx), method, req, reply, cc, opts...)
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}	
}

func (interceptor *AuthInterceptor) attachToken(ctx context.Context) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", interceptor.accessToken)
}

// Stream returns a new stream client interceptor for authentication and authorization.
func (interceptor *AuthInterceptor) Stream() grpc.StreamClientInterceptor {
	return func(
		ctx context.Context, 
		desc *grpc.StreamDesc, 
		cc *grpc.ClientConn, 
		method string, 
		streamer grpc.Streamer, 
		opts ...grpc.CallOption,
		) (grpc.ClientStream, error) {
			log.Printf("--> stream interceptor: %v", method)
			
			if interceptor.authMethods[method] {
				return streamer(interceptor.attachToken(ctx), desc, cc, method, opts...)
			}
			return streamer(ctx, desc, cc, method, opts...)
	}
}

// NewAuthMethods returns a map of methods that require authentication.
// The map is created from the service definition of accessible roles.
func NewAuthMethods() map[string]bool {
	accessibleRoles := service.NewAccessibleRoles()
	authMethods := make(map[string]bool)
	for method := range accessibleRoles {
		authMethods[method] = true
	}
	return authMethods
}