.PHONY: test test-integration build clean

test:
	go test ./...

test-integration:
	OPENEDX_INTEGRATION=1 go test ./integration -v

build:
	go build -o bin/openedx ./cmd/openedx

clean:
	rm -rf bin/
