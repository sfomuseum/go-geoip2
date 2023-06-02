package functionurl

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/aaronland/go-aws-lambda/functionurl"
	"github.com/jtacoma/uritemplates"
)

// LocationForRequest will attempt to derive a `time.Location` instance for the "RemoteAddr" property of 'req' by querying an instance of the
// `sfomuseum/go-geoip2/app/server` application that has been deployed as a AWS Lambda Function URL identified by 'functionurl_uri_template'
// URI template. The request is made using a `go-aws-lambda/functionurl.Client` instance which is instantiated using 'functionurl_client_uri'.
func LocationForRequest(ctx context.Context, functionurl_client_uri string, functionurl_uri_template string, req *http.Request) (*time.Location, error) {
	return LocationForAddress(ctx, functionurl_client_uri, functionurl_uri_template, req.RemoteAddr)
}

// LocationForAddress will attempt to derive a `time.Location` instance for 'address' by querying an instance of the `sfomuseum/go-geoip2/app/server`
// application that has been deployed as a AWS Lambda Function URL identified by 'functionurl_uri_template' URI template. The request is made using
// a `go-aws-lambda/functionurl.Client` instance which is instantiated using 'functionurl_client_uri'.
func LocationForAddress(ctx context.Context, functionurl_client_uri string, functionurl_uri_template string, address string) (*time.Location, error) {

	tz_uritemplate, err := uritemplates.Parse(functionurl_uri_template)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse timezone URI template, %w", err)
	}

	tz_vars := map[string]interface{}{
		"address": address,
	}

	tz_uri, err := tz_uritemplate.Expand(tz_vars)

	if err != nil {
		return nil, fmt.Errorf("Failed to expand timezone URI template variables, %w", err)
	}

	cl, err := functionurl.NewClient(ctx, functionurl_client_uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new timezone function URL client, %w", err)
	}

	tz_rsp, err := cl.Get(ctx, tz_uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to GET timezone function URL, %w", err)
	}

	defer tz_rsp.Body.Close()

	tz_buf, err := io.ReadAll(tz_rsp.Body)

	if err != nil {
		return nil, fmt.Errorf("Failed to read timezone function URL response, %w", err)
	}

	str_tz := string(tz_buf)

	loc, err := time.LoadLocation(str_tz)

	if err != nil {
		return nil, fmt.Errorf("Failed to load location for timezone '%s', %w", str_tz, err)
	}

	return loc, nil
}
