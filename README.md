# go-geoip2

Opinionated Go package providing tools for working with MaxMind GeoLite2 databases. This package is not specific to SFO Museum but it is specific to our needs.

## Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/sfomuseum/go-geoip2.svg)](https://pkg.go.dev/github.com/sfomuseum/go-geoip2)

## Motivation

We want to be able to localize dates and times for people visiting SFO Museum websites and in order to do that we need to know "where" they are. Not precisely where they are but just what timezone they are in.

To accomplish this we are using the openly-licensed MaxMind `GeoLite2City` IP lookup database. That database is _bundled_ with this package and embedded in the binary tools it produces. This results in very large binary tools but is considered an acceptible trade-off because it allows us to deploy a HTTP endpoint to map IP addresses to their timezones as a simple AWS Lambda Function URL without the need for filesystems, databases or containers (to manage filesystems or databases).

This works for us but it may not work for you.

## Tools

```
$> make cli
go build -mod readonly -ldflags="-s -w" -o bin/ip2city cmd/ip2city/main.go
go build -mod readonly -ldflags="-s -w" -o bin/location cmd/location/main.go
go build -mod readonly -ldflags="-s -w" -o bin/server cmd/server/main.go
```

### ip2city

`ip2city` returns the JSON-encoded result for one or more IP addresses in the embedded GeoLite2City database.

```
$> ./bin/ip2city -h
Usage of ./bin/ip2city:
  -json
    	Output data as a JSON array.
```

For example:

```
$> ./bin/ip2city -json 219.3.47.6 | jq

[
  {
    "City": {
      "GeoNameID": 1850147,
      "Names": {
        "de": "Tokio",
        "en": "Tokyo",
        "es": "Tokio",
        "fr": "Tokyo",
        "ja": "東京",
        "pt-BR": "Tóquio",
        "ru": "Токио",
        "zh-CN": "东京"
      }
    },
    "Continent": {
      "Code": "AS",
      "GeoNameID": 6255147,
      "Names": {
        "de": "Asien",
        "en": "Asia",
        "es": "Asia",
        "fr": "Asie",
        "ja": "アジア",
        "pt-BR": "Ásia",
        "ru": "Азия",
        "zh-CN": "亚洲"
      }
    },
    "Country": {
      "GeoNameID": 1861060,
      "IsInEuropeanUnion": false,
      "IsoCode": "JP",
      "Names": {
        "de": "Japan",
        "en": "Japan",
        "es": "Japón",
        "fr": "Japon",
        "ja": "日本",
        "pt-BR": "Japão",
        "ru": "Япония",
        "zh-CN": "日本"
      }
    },
    "Location": {
      "AccuracyRadius": 10,
      "Latitude": 35.6837,
      "Longitude": 139.6805,
      "MetroCode": 0,
      "TimeZone": "Asia/Tokyo"
    },
    "Postal": {
      "Code": "151-0071"
    },
    "RegisteredCountry": {
      "GeoNameID": 1861060,
      "IsInEuropeanUnion": false,
      "IsoCode": "JP",
      "Names": {
        "de": "Japan",
        "en": "Japan",
        "es": "Japón",
        "fr": "Japon",
        "ja": "日本",
        "pt-BR": "Japão",
        "ru": "Япония",
        "zh-CN": "日本"
      }
    },
    "RepresentedCountry": {
      "GeoNameID": 0,
      "IsInEuropeanUnion": false,
      "IsoCode": "",
      "Names": null,
      "Type": ""
    },
    "Subdivisions": [
      {
        "GeoNameID": 1850144,
        "IsoCode": "13",
        "Names": {
          "en": "Tokyo",
          "fr": "Préfecture de Tokyo",
          "ja": "東京都"
        }
      }
    ],
    "Traits": {
      "IsAnonymousProxy": false,
      "IsSatelliteProvider": false
    }
  }
]
```

### location

The `location` tool derives the timezone (location) for one or more addresses by querying an instance of the `server` tool, described below, that has been deployed as a AWS Lambda Function URL.

```
$> ./bin/location -h
Usage of ./bin/location:
  -client-uri string
    	A valid aaronland/go-aws-lambda/functionurl.Client URI in the form of 'functionurl://?credentials={CREDENTIALS}&region={AWS_REGION}' where {CREDENTIALS} is a valid aaronland/go-aws-session credentials strings as described in https://github.com/aaronland/go-aws-session#credentials.
  -uri-template sfomuseum/go-geoip2/app/server
    	A valid URI template where an instance of the sfomuseum/go-geoip2/app/server application that has been deployed as a AWS Lambda Function URL.
```

For example:

```
$> ./bin/location \
	-client-uri 'aws://?credentials={CREDENTIALS}&region={REGION}' \
	-uri-template 'https://FUNCTION_URL_ID.lambda-url.AWS_REGION.on.aws/timezone?address={address}' \
	219.3.47.6
	
Asia/Tokyo
```

_This tool uses the [sfomuseum-go-geoip2/functionurl/location.go](functionurl) package if you need equivalent functionality in your not-command-line code._

### server

`server` exposes a HTTP endpoint where you can query the embedded GeoLite2City database.

```
$> ./bin/server -h
  -database-uri string
    	Valid options are: 'embed://' to use the embedded GeoIP database or the path to an alternate .mmdb database file. (default "embed://")
  -server-uri string
    	A valid aaronland/go-http-server URI string. (default "http://localhost:8080")
```

For example:

```
$> ./bin/server 
2023/06/02 09:59:50 Listening on http://localhost:8080
```

As of this writing it exposes a single `/timezone` endpoint which will return the `tzdata` timezone string for an IP address.

Addresses are inferred from a "?address={IPADDRESS}" parameter or from the IP address of the requestor in that order. For example:

```
$> curl 'http://localhost:8080/timezone?address=219.3.47.6'
Asia/Tokyo
```

#### Lambda

To deploy the `server` tool as an AWS Lambda Function URL first run the `lambda-server` Makefile target:

```
$> make lambda-server
if test -f main; then rm -f main; fi
if test -f server.zip; then rm -f server.zip; fi
GOOS=linux go build -mod readonly -ldflags="-s -w" -o main cmd/server/main.go
zip server.zip main
  adding: main (deflated 52%)
rm -f main
```

Next upload the `server.zip` file to your Lambda function. The function requires no special permissions, policies or roles beyond the ability to execute Lambda functions so you can use the default auto-created role if you want. You will need to assign the following environment variables:

| Name | Value | Notes |
| --- | --- | --- |
| SFOMUSEUM_SERVER_URI | `functionurl://` | Under the hood this is using the [aaronland/go-http-server](https://github.com/aaronland/go-http-server) package to deal with all the Lambda details. |

Access controls for the Lambda Function URL are left to your discretion.

Once complete your function URL will be something like `https://{FUNCTION_URL_ID}.lambda-url.{AWS_REGION}.on.aws` and you can test the server like this:

```
$> curl 'https://{FUNCTION_URL_ID}.lambda-url.{AWS_REGION}.on.aws?address=219.3.47.6'
Asia/Tokyo
```

## GeoLite2 License

This product includes GeoLite2 data created by MaxMind, available from [https://www.maxmind.com](https://www.maxmind.com).

## See also

* https://dev.maxmind.com/geoip/geolite2-free-geolocation-data/#license
* https://github.com/oschwald/geoip2-golang
* https://github.com/aaronland/go-http-server