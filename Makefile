VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "0.0.0-dev")

.PHONY: test test-integration build clean dist

test:
	go test ./...

test-integration:
	OPENEDX_INTEGRATION=1 go test ./integration -v

build:
	go build -o bin/openedx ./cmd/openedx

clean:
	rm -rf bin/ dist/

dist:
	@mkdir -p dist
	GOOS=darwin  GOARCH=arm64 go build -ldflags "-s -w -X main.version=$(VERSION)" -o dist/openedx-$(VERSION)-darwin-arm64    ./cmd/openedx
	GOOS=darwin  GOARCH=amd64 go build -ldflags "-s -w -X main.version=$(VERSION)" -o dist/openedx-$(VERSION)-darwin-amd64    ./cmd/openedx
	GOOS=linux   GOARCH=amd64 go build -ldflags "-s -w -X main.version=$(VERSION)" -o dist/openedx-$(VERSION)-linux-amd64     ./cmd/openedx
	GOOS=linux   GOARCH=arm64 go build -ldflags "-s -w -X main.version=$(VERSION)" -o dist/openedx-$(VERSION)-linux-arm64     ./cmd/openedx
	GOOS=windows GOARCH=amd64 go build -ldflags "-s -w -X main.version=$(VERSION)" -o dist/openedx-$(VERSION)-windows-amd64.exe ./cmd/openedx
