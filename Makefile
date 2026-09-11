.PHONY : build run
build:
	CGO_ENABLED=0,GOOS=linux,GOARCH =amd64 go build -trimpath -ldflags="-s -w" -o bin/api ./cmd/api

run: build
	@./bin/api
