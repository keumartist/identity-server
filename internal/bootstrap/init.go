package bootstrap

import (
	"log"
	"os"
)

func InitApp() {
	serverType := os.Getenv("SERVER_TYPE")

	var err error

	switch serverType {
	case "grpc":
		err = InitGRPCServer()
	case "http":
		err = InitHTTPServer()
	default:
		log.Fatalf("Unknown server type: %s", serverType)
	}

	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}
}
