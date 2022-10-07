format:
	@go fmt ./...

test:
	@go test -cover -v ./...

run:
	@go run main.go

build:
	@go build -v ./...

install:
	@go install -v ./...
