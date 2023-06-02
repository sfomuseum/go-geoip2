package server

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/aaronland/go-http-server"
	"github.com/oschwald/geoip2-golang"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-geoip2/data"
	"github.com/sfomuseum/go-geoip2/http/api"
)

func Run(ctx context.Context, logger *log.Logger) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs, logger)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet, logger *log.Logger) error {

	flagset.Parse(fs)

	err := flagset.SetFlagsFromEnvVars(fs, "SFOMUSEUM")

	if err != nil {
		return fmt.Errorf("Failed to assign flags from environment variables, %w", err)
	}

	// If/when we ever support more than just the timezone endpoint move all the handler
	// code around to use aaronland/go-http-server/handler.RouteHandler

	var db *geoip2.Reader

	switch database_uri {
	case "embed://":
		db, err = geoip2.FromBytes(data.GeoLite2City)
	default:
		db, err = geoip2.Open(database_uri)
	}

	if err != nil {
		return fmt.Errorf("Failed to initialize database, %w", err)
	}

	tz_handler := api.TimeZoneHandler(db)

	mux := http.NewServeMux()
	mux.Handle("/timezone", tz_handler)

	s, err := server.NewServer(ctx, server_uri)

	if err != nil {
		return fmt.Errorf("Failed to create server, %w", err)
	}

	log.Printf("Listening on %s\n", s.Address())

	err = s.ListenAndServe(ctx, mux)

	if err != nil {
		return fmt.Errorf("Failed to serve requests, %w", err)
	}

	return nil
}
