
BINARY_NAME=tzap
VERSION=v0.7.3
#VERSION=$(shell git describe --tags --always)
BUILD_FLAGS=-ldflags="-X main.Version=$(VERSION)"

build:
	go build $(BUILD_FLAGS) -o $(BINARY_NAME) .

test:
	go test -v ./...

clean:
	rm -f ../release/tzap*

release: clean
	mkdir -p ../release
	GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) -o ../release/$(BINARY_NAME)-$(VERSION)-darwin-amd64 .
	GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o ../release/$(BINARY_NAME)-$(VERSION)-linux-amd64 .
	GOOS=windows GOARCH=amd64 go build $(BUILD_FLAGS) -o ../release/$(BINARY_NAME)-$(VERSION)-windows-amd64.exe .
	GOOS=darwin GOARCH=arm64 go build $(BUILD_FLAGS) -o ../release/$(BINARY_NAME)-$(VERSION)-darwin-arm64 .
	GOOS=linux GOARCH=arm64 go build $(BUILD_FLAGS) -o ../release/$(BINARY_NAME)-$(VERSION)-linux-arm64 .
	GOOS=windows GOARCH=arm64 go build $(BUILD_FLAGS) -o ../release/$(BINARY_NAME)-$(VERSION)-windows-arm64.exe .

gh-upload:
	gh release upload $(VERSION) ../release/tzap*
	
github-release:
	gh release create v0.7.3 --tag $(VERSION) --name "$(VERSION)" --file ../release/$(BINARY_NAME)-$(VERSION)-darwin-amd64
	gh release create v0.7.3 --tag $(VERSION) --name "$(VERSION)" --file ../release/$(BINARY_NAME)-$(VERSION)-linux-amd64
	gh release create v0.7.3 --tag $(VERSION) --name "$(VERSION)" --file ../release/$(BINARY_NAME)-$(VERSION)-windows-amd64.exe

.PHONY: build test clean release github-release