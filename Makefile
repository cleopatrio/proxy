ENV=test
GIT_COMMIT:=$(shell git rev-parse HEAD)

format:
	@go fmt ./...

test: format
	go test -cover -v ./...

run: test
	go run .
