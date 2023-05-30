module github.com/tzapio/tzap/pkg/connectors/openaiconnector

go 1.20

replace github.com/tzapio/tzap => ../../../

require (
	github.com/sashabaranov/go-openai v1.9.5
	github.com/tzapio/tokenizer v0.0.3
	github.com/tzapio/tzap v0.0.0-00010101000000-000000000000
)

require github.com/dlclark/regexp2 v1.9.0 // indirect
