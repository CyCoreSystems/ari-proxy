

all:
	mkdir -p bin
	go build -o bin/ari-proxy ./cmd/ari-proxy
