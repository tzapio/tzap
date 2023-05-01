
BINARY_NAME=tzap
VERSION=v0.7.7
#VERSION=$(shell git describe --tags --always)
BUILD_FLAGS=-ldflags="-X main.Version=$(VERSION)"

test:
	go test -v ./...

clean:
	rm -f ../release/tzap*

tag: release github-release
	go run ./main.go ghrelease $(VERSION)

release: clean 
	mkdir -p ../release
	GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) -o ../release/$(BINARY_NAME)-$(VERSION)-darwin-amd64 .
	GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o ../release/$(BINARY_NAME)-$(VERSION)-linux-amd64 .
	GOOS=windows GOARCH=amd64 go build $(BUILD_FLAGS) -o ../release/$(BINARY_NAME)-$(VERSION)-windows-amd64.exe .
	GOOS=darwin GOARCH=arm64 go build $(BUILD_FLAGS) -o ../release/$(BINARY_NAME)-$(VERSION)-darwin-arm64 .
	GOOS=linux GOARCH=arm64 go build $(BUILD_FLAGS) -o ../release/$(BINARY_NAME)-$(VERSION)-linux-arm64 .
	GOOS=windows GOARCH=arm64 go build $(BUILD_FLAGS) -o ../release/$(BINARY_NAME)-$(VERSION)-windows-arm64.exe .


gh-upload: release github-release
	gh release upload $(VERSION) ../release/tzap*
	
github-release:
	gh release create pkg/tzapconnect/$(VERSION) --prerelease --title "Release pkg/tzapconnect/$(VERSION)" --notes "Golang specific library release"
	gh release create pkg/connectors/openaiconnector/$(VERSION) --prerelease --title "Release pkg/connectors/openaiconnector/$(VERSION)" --notes "Golang specific library release"
	gh release create pkg/connectors/googlevoiceconnector/$(VERSION) --prerelease --title "Release pkg/connectors/googlevoiceconnector/$(VERSION)" --notes "Golang specific library release"

.PHONY: build test clean release github-release