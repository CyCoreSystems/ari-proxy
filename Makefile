
all:
	mkdir -p bin
	go build -o bin/ari-proxy ./cmd/ari-proxy

docker: all
	docker build -t cycoresystems/ari-proxy ./
	docker push cycoresystems/ari-proxy


