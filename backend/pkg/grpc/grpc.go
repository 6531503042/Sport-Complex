package grpc

import (
	"context"
	"errors"
	"log"
	"main/config"
	"main/pkg/jwt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type (
	GrpcClientFactoryHandler interface {
		
	}

	grpcClientFactory struct {
		client *grpc.ClientConn
	}

	grpcAuth struct {
		secretKey string
	}
)

func NewGrpcClient(host string) (GrpcClientFactoryHandler, error) {
    opts := make([]grpc.DialOption, 0)

    opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

    ClientConn, err := grpc.NewClient(host, opts...)
    if err != nil {
        log.Printf("Error: Grpc Client connection failed %s", err.Error())
        return nil, errors.New("error: grpc client connection failed")
    }

    return &grpcClientFactory{
        client: ClientConn,
    }, nil
}

func NewGrpcServer(cfg *config.Jwt, host string) (*grpc.Server, net.Listener) {
	opts := make([]grpc.ServerOption, 0)

	grpcAuth := &grpcAuth{
		secretKey: cfg.AccessSecretKey,
	}

	opts = append(opts, grpc.UnaryInterceptor(grpcAuth.unaryAuthorization))

	grpcServer := grpc.NewServer(opts...)

	lis, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	return grpcServer, lis
}

func (g *grpcAuth) unaryAuthorization (ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Printf("Error: Metadata not found")
		return nil, errors.New("metadata not found")
	}

	authHeader, ok := md["auth"]
	if !ok {
		log.Printf("Error: Auth header not found")
		return nil, errors.New("auth header not found")
	}

	if len(authHeader) == 0 {
		log.Printf("Error: Auth header not found")
		return nil, errors.New("auth header not found")
	}

	claims, err := jwt.ParseToken(g.secretKey, string(authHeader[0]))
	if err != nil {
		log.Printf("Error: ParseToken: %s", err.Error())
		return nil, errors.New("parse token failed")
	}
	log.Printf("claims: %v", claims)
	return handler(ctx, req)
}


