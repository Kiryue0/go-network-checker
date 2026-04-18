.PHONY: build test run

build:
	go build -o netcheck ./cmd/netcheck

test:
	go test -race ./...

run:
	go run ./cmd/netcheck

clean:
	rm -f netcheck