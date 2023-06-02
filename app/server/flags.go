package server

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var server_uri string
var database_uri string

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("geoip2")

	fs.StringVar(&server_uri, "server-uri", "http://localhost:8080", "A valid aaronland/go-http-server URI string.")
	fs.StringVar(&database_uri, "database-uri", "embed://", "Valid options are: 'embed://' to use the embedded GeoIP database or the path to an alternate .mmdb database file.")

	return fs
}
