
all:
	go build ./...
	go build

docker: all
	docker build -t cycoresystems/ari-proxy ./
	docker push cycoresystems/ari-proxy

test:
	go test ./...

lint:
	golangci-lint run

check: all lint test

deps:
	dep ensure

ci: deps check
