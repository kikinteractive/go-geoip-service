.PHONY: build test

all: build

build: clean
	go build

build-linux: clean
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build

test:
	go list ./... | grep -v vendor | xargs -I{} go test -v '{}' -check.v

clean:
	rm -f go-geoip-service
