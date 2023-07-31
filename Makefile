ENV=test
LOG_LEVEL=5
GIT_COMMIT:=$(shell git rev-parse HEAD)

format:
	@go fmt ./...

test: format
	go test -cover -v ./...

run:
	go run .

build-for-docker:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' -o proxy

build-dockerfile:
	docker build -t proxy . --build-arg GIT_COMMIT=$(GIT_COMMIT)

examples:
	@echo "=> [GET] request to <http://viacep.com.br> requesting information about a given CEP:"
	@curl --connect-to viacep.com.br:80:localhost:5000 http://viacep.com.br/ws/29100010/json
	@echo "\n"

	@echo "=> [GET] request to <http://jsonplaceholder.typicode.com> which returns a single post:"
	@curl --connect-to jsonplaceholder.typicode.com:80:localhost:5000 http://jsonplaceholder.typicode.com/posts/1
	@echo "\n"

	@echo "=> [POST] request to <http://jsonplaceholder.typicode.com> which simulates the creation of a new post:"
	@curl --connect-to jsonplaceholder.typicode.com:80:localhost:5000 http://jsonplaceholder.typicode.com/posts -X POST
	@echo "\n"

	@echo "=> [GET] forward <http://example.com/people/*> request"
	@curl --connect-to example.com:80:localhost:5000 http://example.com/people/1
	@echo "\n"

	@echo "=> [GET] forward <http://example.com/friends> request"
	@curl --connect-to example.com:80:localhost:5000 http://example.com/friends
	@echo "\n"
