lint:
	go fmt ./...
	golangci-lint run ./... --skip-dirs="(docs|frontend)"

wire:
	gutowire -w ./cmd/inject -p inject

run:
	go run ./cmd/

build:
	go build -o gobana ./cmd/