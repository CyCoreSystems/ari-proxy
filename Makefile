
all:
	go build ./...
	go build

docker: all
	docker build -t cycoresystems/ari-proxy ./
	docker push cycoresystems/ari-proxy

test:
	go test ./...

lint:
	gometalinter --skip internal --vendor --deadline=60s ./...

check: all lint test

deps:
	dep ensure

ci: deps check
