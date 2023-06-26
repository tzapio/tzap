test: gomodtidy
	go test ./...
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

wasm: 
	cd cli/wasm && GOOS=js GOARCH=wasm go build -o tzap.wasm
wasml: 
	npx live-server cli/wasm

.PHONY: release


proto:
	go install \
		google.golang.org/protobuf/cmd/protoc-gen-go@v1.28

	protoc \
	-I cli/proto \
	--go_out=cli/ --go_opt=paths=import \
	--go-grpc_out=cli/ --go-grpc_opt=paths=import \
	tzap.proto prompt.proto search.proto

docu:
	cd documentation && npm start
