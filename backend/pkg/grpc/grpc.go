package grpc

import (
	"context"
	"errors"
	"log"
	"main/config"
	authPb "main/modules/auth/proto"
	bookingPb "main/modules/booking/proto"
	userPb "main/modules/user/proto"
	paymentPb "main/modules/payment/proto"
	"main/pkg/jwt"
	"net"

	facilityPb "main/modules/facility/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type (
	GrpcClientFactoryHandler interface {
		User() userPb.UserGrpcServiceClient
		Auth() authPb.AuthGrpcServiceClient
		Booking() bookingPb.BookingServiceClient
		Facility() facilityPb.FacilityServiceClient
		Payment() paymentPb.PaymentServiceClient
	}

	grpcClientFactory struct {
		client *grpc.ClientConn
	}

	grpcAuth struct {
		secretKey string
	}
)

func (g *grpcClientFactory) User() userPb.UserGrpcServiceClient {
	return userPb.NewUserGrpcServiceClient(g.client)
}

func (g *grpcClientFactory) Auth() authPb.AuthGrpcServiceClient {
	return authPb.NewAuthGrpcServiceClient(g.client)
}

func (g *grpcClientFactory) Booking() bookingPb.BookingServiceClient {
	return bookingPb.NewBookingServiceClient(g.client)
}

func (g *grpcClientFactory) Facility() facilityPb.FacilityServiceClient {
	return facilityPb.NewFacilityServiceClient(g.client)
}

func (g *grpcClientFactory) Payment() paymentPb.PaymentServiceClient { // Implement Payment service client
	return paymentPb.NewPaymentServiceClient(g.client)
}

func NewGrpcClient(host string) (GrpcClientFactoryHandler, error) {
    opts := make([]grpc.DialOption, 0)
    opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

    ctx := context.Background()
    clientConn, err := grpc.DialContext(ctx, host, opts...)
    if err != nil {
        log.Printf("Error: gRPC client connection to %s failed: %s", host, err.Error())
        return nil, errors.New("error: grpc client connection failed")
    }

    log.Printf("gRPC client connected to %s", host)
    return &grpcClientFactory{
        client: clientConn,
    }, nil
}


func NewGrpcServer(cfg *config.Jwt, host string) (*grpc.Server, net.Listener) {
	opts := make([]grpc.ServerOption, 0)

	grpcAuth := &grpcAuth{
		secretKey: cfg.ApiSecretKey,
	}

	opts = append(opts, grpc.UnaryInterceptor(grpcAuth.unaryAuthorization))

	grpcServer := grpc.NewServer(opts...)

	lis, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatalf("Error: Failed to listen: %v", err)
	}

	return grpcServer, lis
}

func (g *grpcAuth) unaryAuthorization(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
    md, ok := metadata.FromIncomingContext(ctx)
    if !ok {
        log.Printf("Error: Metadata not found")
        return nil, errors.New("error: metadata not found")
    }

    authHeader, ok := md["auth"]
    if !ok || len(authHeader) == 0 {
        log.Printf("Error: Auth metadata not found")
        return nil, errors.New("error: metadata not found")
    }

	if len (authHeader) == 0 {
		log.Printf("Error: Auth metadata not found")
		return nil, errors.New("error: metadata not found")
	}

    claims, err := jwt.ParseToken(g.secretKey, string(authHeader[0]))
    if err != nil {
        log.Printf("Error: Parse token failed for auth header: %s, error: %s", authHeader[0], err.Error())
        return nil, errors.New("error: token is invalid")
    }
    log.Printf("Token claims: %v", claims)

    // Call the handler
    return handler(ctx, req)
}

