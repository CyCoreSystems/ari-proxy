
all:
	go build `go list ./... | grep -v /vendor/`
	mkdir -p bin
	go build -o bin/ari-proxy ./cmd/ari-proxy

docker: all
	docker build -t cycoresystems/ari-proxy ./
	docker push cycoresystems/ari-proxy

test:
	go test `go list ./... | grep -v /vendor/`

lint:
	gometalinter --deadline 20s --vendor --fast --skip internal ./...

check: lint test

deps:
	glide cc
	glide i

ci: deps check
