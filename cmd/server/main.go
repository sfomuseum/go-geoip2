package main

import (
	"context"
	"log"

	"github.com/sfomuseum/go-geoip2/app/server"
)

func main() {

	ctx := context.Background()
	logger := log.Default()

	err := server.Run(ctx, logger)

	if err != nil {
		logger.Fatalf("Failed to run server, %w", err)
	}
}
