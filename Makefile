.PHONY: build halo-server simulator-http test clean

build: halo-server simulator-http

halo-server:
	go build -o halo-server ./cmd/halo-server

simulator-http:
	go build -o simulator-http ./cmd/simulator-http

test:
	go test ./...

clean:
	rm -f halo-server simulator-http
