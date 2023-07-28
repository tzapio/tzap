module github.com/tzapio/tzap/pkg/tzapaction

go 1.20

replace github.com/tzapio/tzap => ../../

replace github.com/tzapio/tzap/cli => ../../cli

replace github.com/tzapio/tzap/pkg/tzapconnect => ../tzapconnect

replace github.com/tzapio/tzap/pkg/connectors/openaiconnector => ../connectors/openaiconnector

require (
	github.com/sashabaranov/go-openai v1.14.1
	github.com/sergi/go-diff v1.3.1
	github.com/tzapio/tzap v0.8.9
	github.com/tzapio/tzap/cli v0.0.0-20230720082106-8d6dc0712682
	go.uber.org/mock v0.2.0
	google.golang.org/protobuf v1.31.0
)

require (
	github.com/fatih/color v1.15.0 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/sabhiram/go-gitignore v0.0.0-20210923224102-525f6e181f06 // indirect
	golang.org/x/sys v0.10.0 // indirect
)
