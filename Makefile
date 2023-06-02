ENV=test
GIT_COMMIT:=$(shell git rev-parse HEAD)

format:
	@go fmt ./...

test: format
	go test -p 1 ./...

run: test
	go run .
