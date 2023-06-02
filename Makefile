GOMOD=readonly

cli:
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/ip2city cmd/ip2city/main.go
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/server cmd/server/main.go

lambda-server:
	if test -f main; then rm -f main; fi
	if test -f server.zip; then rm -f server.zip; fi
	GOOS=linux go build -mod $(GOMOD) -ldflags="-s -w" -o main cmd/server/main.go
	zip server.zip main
	rm -f main
