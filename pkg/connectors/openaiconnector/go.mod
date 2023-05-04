module github.com/tzapio/tzap/pkg/connectors/openaiconnector

go 1.20

replace github.com/tzapio/tzap => ../../../

require (
	github.com/sashabaranov/go-openai v1.9.0
	github.com/tiktoken-go/tokenizer v0.1.0
	github.com/tzapio/tzap v.0.7.9
)

require github.com/dlclark/regexp2 v1.9.0 // indirect
