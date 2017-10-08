
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
	gometalinter --disable=gotype --disable=errcheck client/... server/... proxy/...

check: all lint test

deps:
	dep ensure

ci: deps check
