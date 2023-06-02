package main

/*

$> go run cmd/location/main.go \
	-client-uri 'aws://?credentials={CREDENTIALS}&region={REGION}' \
	-uri-template 'https://(FUNCTION_URL_ID).lambda-url.us-(AWS_REGION).on.aws/timezone?address={address}' \
	219.3.47.6

Asia/Tokyo

*/

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/sfomuseum/go-geoip2/functionurl"
)

func main() {

	functionurl_client_uri := flag.String("client-uri", "", "A valid aaronland/go-aws-lambda/functionurl.Client URI in the form of 'functionurl://?credentials={CREDENTIALS}&region={AWS_REGION}' where {CREDENTIALS} is a valid aaronland/go-aws-session credentials strings as described in https://github.com/aaronland/go-aws-session#credentials.")
	functionurl_uri_template := flag.String("uri-template", "", "A valid URI template where an instance of the `sfomuseum/go-geoip2/app/server` application that has been deployed as a AWS Lambda Function URL.")

	flag.Parse()

	ctx := context.Background()

	for _, addr := range flag.Args() {

		loc, err := functionurl.LocationForAddress(ctx, *functionurl_client_uri, *functionurl_uri_template, addr)

		if err != nil {
			log.Fatalf("Failed to derive location for '%s', %v", addr, err)
		}

		fmt.Println(loc)
	}
}
