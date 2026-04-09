.PHONY: test build clean

test:
	go test ./...

build:
	go build -o bin/openedx ./cmd/openedx

clean:
	rm -rf bin/
