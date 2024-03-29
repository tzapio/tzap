
BINARY_NAME=tzap
PREVIOUS_VERSION=v0.9.3
VERSION=v0.9.4
#VERSION=$(shell git describe --tags --always)
BUILD_FLAGS=-ldflags="-X main.Version=$(VERSION) -s -w" -trimpath

test:
	go test -v ./...

clean:
	rm -rf ../release

build: clean
	mkdir -p ../release
# GOOS=js GOARCH=wasm go build $(BUILD_FLAGS) -o ../release/$(BINARY_NAME).wasm ./wasm/main.go
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build $(BUILD_FLAGS) -o ../release/$(BINARY_NAME)-darwin-amd64 .
#upx ../release/$(BINARY_NAME)-darwin-amd64
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(BUILD_FLAGS) -o ../release/$(BINARY_NAME)-linux-amd64 .
	upx ../release/$(BINARY_NAME)-linux-amd64
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build $(BUILD_FLAGS) -o ../release/$(BINARY_NAME)-windows-amd64.exe .
#	upx ../release/$(BINARY_NAME)-windows-amd64.exe
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build $(BUILD_FLAGS) -o ../release/$(BINARY_NAME)-darwin-arm64 .
#	upx ../release/$(BINARY_NAME)-darwin-arm64
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build $(BUILD_FLAGS) -o ../release/$(BINARY_NAME)-linux-arm64 .
	upx ../release/$(BINARY_NAME)-linux-arm64
	GOOS=windows GOARCH=arm64 CGO_ENABLED=0 go build $(BUILD_FLAGS) -o ../release/$(BINARY_NAME)-windows-arm64.exe .
#upx ../release/$(BINARY_NAME)-windows-arm64.exe
tzapPrepareRelease:
	git tag -d $(git tag -l);
	gh release create $(VERSION) --prerelease --title "Release $(VERSION)" --notes "";
	git fetch;
	
tzapWriteRelease:
	cd ../ && go run ./cli/main.go ghrelease $(PREVIOUS_VERSION) $(VERSION)

github-upload:
	gh release upload $(VERSION) ../release/tzap*
	
	
github-otherpkgs-release:
	if [ -z "$(shell gh release list | grep -F pkg/tzapconnect/$(VERSION))" ]; then \
		gh release create pkg/tzapconnect/$(VERSION) --prerelease --title "Release pkg/tzapconnect/$(VERSION)" --notes "Golang specific release"; \
	fi
	if [ -z "$(shell gh release list | grep -F pkg/tzapaction/$(VERSION))" ]; then \
		gh release create pkg/tzapaction/$(VERSION) --prerelease --title "Release pkg/tzapaction/$(VERSION)" --notes "Golang specific release"; \
	fi
	if [ -z "$(shell gh release list | grep -F pkg/connectors/openaiconnector/$(VERSION))" ]; then \
		gh release create pkg/connectors/openaiconnector/$(VERSION) --prerelease --title "Release pkg/connectors/openaiconnector/$(VERSION)" --notes "Golang specific release"; \
	fi
	if [ -z "$(shell gh release list | grep -F cli/$(VERSION))" ]; then \
		gh release create cli/$(VERSION) --prerelease --title "Release cli/$(VERSION)" --notes "Golang specific release"; \
	fi


remove-releases:
	-gh release delete pkg/tzapconnect/$(VERSION) --yes
	-gh release delete pkg/tzapaction/$(VERSION) --yes
	-gh release delete pkg/connectors/openaiconnector/$(VERSION) --yes
	-gh release delete $(VERSION) --yes

	-git push --delete origin pkg/tzapconnect/$(VERSION)
	-git push --delete origin pkg/tzapaction/$(VERSION)
	-git push --delete origin pkg/connectors/openaiconnector/$(VERSION)
	-git push --delete origin cli/$(VERSION)
	-git push --delete origin $(VERSION)

	-git tag -d pkg/tzapconnect/$(VERSION)
	-git tag -d pkg/tzapaction/$(VERSION)
	-git tag -d pkg/connectors/openaiconnector/$(VERSION)
	-git tag -d cli/$(VERSION)
	-git tag -d $(VERSION)

.PHONY: