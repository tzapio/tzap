module github.com/tzapio/tzap/pkg/example

go 1.20

replace github.com/tzapio/tzap => ../

replace github.com/tzapio/tzap/pkg/connectors/googlevoiceconnector => ../pkg/connectors/googlevoiceconnector

replace github.com/tzapio/tzap/pkg/connectors/openaiconnector => ../pkg/connectors/openaiconnector

replace github.com/tzapio/tzap/pkg/tzapconnect => ../pkg/tzapconnect

require (
	github.com/tzapio/tzap v0.7.10
	github.com/tzapio/tzap/pkg/tzapconnect v0.7.10
)

require (
	github.com/dlclark/regexp2 v1.9.0 // indirect
	github.com/sashabaranov/go-openai v1.9.0 // indirect
	github.com/tiktoken-go/tokenizer v0.1.0 // indirect
	github.com/tzapio/tzap/pkg/connectors/openaiconnector v0.7.10 // indirect
)
