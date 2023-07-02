//go:build grpc

package main

import (
	bootstrap "art-sso/internal/bootstrap"
	"log"
)

func main() {
	err := bootstrap.InitGRPCServer()

	if err != nil {
		log.Fatalf("Failed to initialize GRPC server: %v", err)
	}
}
