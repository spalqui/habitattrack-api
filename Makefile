.PHONY: build run test clean docker-build docker-run

build:
	go build -o bin/server cmd/server/main.go

run:
	go run cmd/server/main.go

test:
	go test -v ./...

clean:
	rm -rf bin/

docker-build:
	docker build -t property-rental-api .

docker-run:
	docker run -p 8080:8080 --env-file .env property-rental-api

install-deps:
	go mod tidy
	go mod download

lint:
	golangci-lint run

format:
	go fmt ./...