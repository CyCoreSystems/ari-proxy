
all:
	go build `go list ./... | grep -v /vendor/`
	mkdir -p bin
	go build

docker: all
	docker build -t cycoresystems/ari-proxy ./
	docker push cycoresystems/ari-proxy

test:
	go test `go list ./... | grep -v /vendor/`

lint:
	gometalinter --disable=gotype --disable=errcheck ./... --skip internal --vendor

check: all lint test

deps:
	dep ensure

ci: deps check
