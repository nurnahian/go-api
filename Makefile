.PHONY: build run docker-build docker-run

build:
	go build -o bin/go-api ./cmd/go-api

run:
	go run ./cmd/go-api

docker-build:
	docker build -t go-api .

docker-run:
	docker run -p 8080:8080 go-api