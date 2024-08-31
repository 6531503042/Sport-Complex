package grpc

import (
	"context"
	"errors"
	"log"
	"main/pkg/jwt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/internal/metadata"
)

type (
	GrpcClientFactoryHndler interface {

	}

	grpcClientFactory struct {
		client *grpc.ClientConn
	}

	grpcAuth struct {
		secretKey string
	}
)

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