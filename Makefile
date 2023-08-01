test: gomodtidy
	go test ./...
build: gomodtidy
	go test ./...
	make -C cli build
release:
	go test ./...
	make -C cli build
	make -C cli github-upload
release-local:
	make gomodtidy
	go test ./...
	git pull
	git push
	make -C cli build
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
install:
	go install ./cli 
gomodtidy:
	go mod tidy
	cd pkg/tzapaction && go mod tidy
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
	go install \
		google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
	protoc \
	-I pkg/tzapaction/proto \
	--go_out=pkg/tzapaction/ --go_opt=paths=import \
	tzap.proto prompt.proto search.proto refactor.proto edit.proto implement.proto common.proto 

docu:
	cd documentation && npm start

dockerbuild:
	docker build -t tzapio/tzap .
dockerpush:
	docker push tzapio/tzap