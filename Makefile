GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")

cli:
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/ip2city cmd/ip2city/main.go
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/location cmd/location/main.go
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/server cmd/server/main.go

lambda-server:
	if test -f bootstrap; then rm -f main; fi
	if test -f server.zip; then rm -f server.zip; fi
	GOARCH=arm64 GOOS=linux go build -mod $(GOMOD) -ldflags="-s -w" -tags lambda.norpc -o bootstrap cmd/server/main.go
	zip server.zip bootstrap
	rm -f bootstrap

lambda-server-go1:
	if test -f main; then rm -f main; fi
	if test -f server.zip; then rm -f server.zip; fi
	GOOS=linux go build -mod readonly -ldflags="-s -w" -o main cmd/server/main.go
	zip server.zip main
	rm -f main

