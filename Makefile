build: gomodtidy
	go test ./...
	make -C cli build
release:
	make gomodtidy
	go test ./...
	git pull
	git push
	make -C cli build
	make -C cli tzapPrepareRelease
	make -C cli tzapWriteRelease
	make -C cli github-upload
releaseOther:
	make -C cli github-otherpkgs-release



exGithubDoc:
	go run examples/githubdoc/main.go
exMadebygpt:
	go run examples/madebygpt/main.go
exRefactoring:
	go run examples/refactoring/main.go
exTesting:
	go run examples/testing/main.go

gomodtidy:
	go mod tidy
	cd pkg/connectors/openaiconnector && go mod tidy
	cd pkg/tzapconnect && go mod tidy
	cd pkg/connectors/redisembeddbconnector && go mod tidy
	cd pkg/connectors/googlevoiceconnector && go mod tidy
	cd examples && go mod tidy
	cd cli && go mod tidy
	go work sync

ts-build:
	cd ts && npm run build

.PHONY: release
