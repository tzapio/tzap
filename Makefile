build: gomodtidy
	cd cli && make release
release:
	cd cli && make test
	cd cli && make build
	cd cli && make github-release
releaseZ:
	cd cli && make tag
	cd cli && make gh-upload

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
	cd pkg/connectors/googlevoiceconnector && go mod tidy
	cd examples && go mod tidy
	cd cli && go mod tidy
	go work sync

ts-build:
	cd ts && npm run build