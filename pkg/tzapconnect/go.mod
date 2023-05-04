module github.com/tzapio/tzap/pkg/tzapconnect

go 1.20

replace github.com/tzapio/tzap => ../../

replace github.com/tzapio/tzap/pkg/connectors/openaiconnector => ../connectors/openaiconnector

require (
	github.com/tzapio/tzap v0.7.8
	github.com/tzapio/tzap/pkg/connectors/openaiconnector v0.7.8
)

require (
	github.com/dlclark/regexp2 v1.9.0 // indirect
	github.com/sashabaranov/go-openai v1.9.0 // indirect
	github.com/tiktoken-go/tokenizer v0.1.0 // indirect
)
