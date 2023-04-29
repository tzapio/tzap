BINARY_NAME=tzap
VERSION=$(shell git describe --tags --always --dirty="-dev")
BUILD_FLAGS=-ldflags="-X main.Version=$(VERSION)"

build:
	go build $(BUILD_FLAGS) -o $(BINARY_NAME) .

test:
	go test -v ./...

clean:
	rm -f $(BINARY_NAME)

release: clean
	mkdir -p ../release
	GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) -o ../release/$(BINARY_NAME)-$(VERSION)-darwin-amd64 .
	GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o ../release/$(BINARY_NAME)-$(VERSION)-linux-amd64 .
	GOOS=windows GOARCH=amd64 go build $(BUILD_FLAGS) -o ../release/$(BINARY_NAME)-$(VERSION)-windows-amd64.exe .

github-release:
	github-release upload --user aleksandarfaraj --repo tzapio/tzap --tag $(VERSION) --name "$(VERSION)" --file ../release/$(BINARY_NAME)-$(VERSION)-darwin-amd64
	github-release upload --user aleksandarfaraj --repo tzapio/tzap --tag $(VERSION) --name "$(VERSION)" --file ../release/$(BINARY_NAME)-$(VERSION)-linux-amd64
	github-release upload --user aleksandarfaraj --repo tzapio/tzap --tag $(VERSION) --name "$(VERSION)" --file ../release/$(BINARY_NAME)-$(VERSION)-windows-amd64.exe

.PHONY: build test clean release github-release