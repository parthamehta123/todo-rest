run:
	go run ./cmd/server

test:
	go test ./...

build:
	go build -o bin/todo ./cmd/server
