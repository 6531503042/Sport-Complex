// package grpc

// import (
// 	"google.golang.org/grpc"
// )

// type (
// 	GrpcClientFactoryHandler interface {
	
// 	}

// 	GrpcClientFactory struct {
// 		client *grpc.ClientConn
// 	}

// 	grpcAuth struct {
// 		secretKey string
// 	}
// )

// func (g * GrpcClientFactory) Auth() authPb.AuthGrpcServiceClient {
// 	return authPb.NewAuthGrpcServiceClient(g.client)
// }