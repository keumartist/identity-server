package bootstrap

import (
	api "art-sso/internal/grpc/api/proto"
	authhandler "art-sso/internal/handler/auth/grpc"
	tokenservice "art-sso/internal/service/token"
	"fmt"
	"net"

	"google.golang.org/grpc"
)

func InitGRPCServer() error {

	privateKey, secretKey, keyErr := getKeys()
	if keyErr != nil {
		return fmt.Errorf("Could not init app: %v", keyErr)
	}

	issuer, issuerErr := getIssuer()
	if issuerErr != nil {
		return fmt.Errorf("Could not init app: %v", issuerErr)
	}
	tokenService := tokenservice.NewTokenService(privateKey, secretKey, issuer)
	grpcHandler := authhandler.NewAuthHandler(tokenService)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		return fmt.Errorf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	api.RegisterAuthServiceServer(grpcServer, grpcHandler)

	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("Failed to serve gRPC server: %v", err)
	}

	return nil
}
