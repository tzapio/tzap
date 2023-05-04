
BINARY_NAME=tzap
VERSION=v0.7.8
#VERSION=$(shell git describe --tags --always)
BUILD_FLAGS=-ldflags="-X main.Version=$(VERSION)"

test:
	go test -v ./...

clean:
	rm -rf ../release

tzapWriteRelease:
	if [ -n "$(shell gh release list | grep -E '^$(VERSION)\b')" ]; then \
		tzap ghrelease $(VERSION) \
	else \
		$(error Release $(VERSION) does not exist yet); \
	fi

build: clean
	mkdir -p ../release
	GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) -o ../release/$(BINARY_NAME)-$(VERSION)-darwin-amd64 .
	GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o ../release/$(BINARY_NAME)-$(VERSION)-linux-amd64 .
	GOOS=windows GOARCH=amd64 go build $(BUILD_FLAGS) -o ../release/$(BINARY_NAME)-$(VERSION)-windows-amd64.exe .
	GOOS=darwin GOARCH=arm64 go build $(BUILD_FLAGS) -o ../release/$(BINARY_NAME)-$(VERSION)-darwin-arm64 .
	GOOS=linux GOARCH=arm64 go build $(BUILD_FLAGS) -o ../release/$(BINARY_NAME)-$(VERSION)-linux-arm64 .
	GOOS=windows GOARCH=arm64 go build $(BUILD_FLAGS) -o ../release/$(BINARY_NAME)-$(VERSION)-windows-arm64.exe .


gh-upload:
	if [ -n "$(shell gh release list | grep -E '^$(VERSION)\b')" ]; then \
		gh release upload $(VERSION) ../release/tzap* \
	else \
		$(error Release $(VERSION) does not exist yet); \
	fi
	
	
github-pre-release:
	if [ -z "$(shell gh release list | grep -F pkg/tzapconnect/$(VERSION))" ]; then \
		gh release create pkg/tzapconnect/$(VERSION) --prerelease --title "Release pkg/tzapconnect/$(VERSION)" --notes "Golang specific release"; \
	fi
	if [ -z "$(shell gh release list | grep -F pkg/connectors/openaiconnector/$(VERSION))" ]; then \
		gh release create pkg/connectors/openaiconnector/$(VERSION) --prerelease --title "Release pkg/connectors/openaiconnector/$(VERSION)" --notes "Golang specific release"; \
	fi
	if [ -z "$(shell gh release list | grep -F pkg/connectors/googlevoiceconnector/$(VERSION))" ]; then \
		gh release create pkg/connectors/googlevoiceconnector/$(VERSION) --prerelease --title "Release pkg/connectors/googlevoiceconnector/$(VERSION)" --notes "Golang specific release"; \
	fi
	if [ -z "$(shell gh release list | grep -E '^$(VERSION)\b')" ]; then \
		git tag -d $(git tag -l); \
		gh release create $(VERSION) --prerelease --title "Release $(VERSION)" --notes ""; \
		git fetch; \
	fi

remove-releases:
	gh release delete pkg/tzapconnect/$(VERSION) --yes
	gh release delete pkg/connectors/openaiconnector/$(VERSION) --yes
	gh release delete pkg/connectors/googlevoiceconnector/$(VERSION) --yes
	gh release delete $(VERSION) --yes

	git push --delete origin pkg/tzapconnect/$(VERSION)
	git push --delete origin pkg/connectors/openaiconnector/$(VERSION)
	git push --delete origin pkg/connectors/googlevoiceconnector/$(VERSION)
	git push --delete origin $(VERSION)

	git tag -d pkg/tzapconnect/$(VERSION)
	git tag -d pkg/connectors/openaiconnector/$(VERSION)
	git tag -d pkg/connectors/googlevoiceconnector/$(VERSION)
	git tag -d $(VERSION)

.PHONY: