module github.com/tzapio/tzap/cli

go 1.20

replace github.com/tzapio/tzap => ../

replace github.com/tzapio/tzap/pkg/tzapconnect => ../pkg/tzapconnect

replace github.com/tzapio/tzap/pkg/connectors/openaiconnector => ../pkg/connectors/openaiconnector

require (
	github.com/spf13/cobra v1.7.0
	github.com/tzapio/tzap v0.0.0-00010101000000-000000000000
	github.com/tzapio/tzap/pkg/tzapconnect v0.0.0-00010101000000-000000000000
)

require (
	github.com/dlclark/regexp2 v1.9.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/sashabaranov/go-openai v1.9.3 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/tiktoken-go/tokenizer v0.1.0 // indirect
	github.com/tzapio/tzap/pkg/connectors/openaiconnector v0.0.0-00010101000000-000000000000 // indirect
)
