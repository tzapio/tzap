module github.com/tzapio/tzap/pkg/tzapconnect

go 1.20

replace github.com/tzapio/tzap => ../../

replace github.com/tzapio/tzap/pkg/connectors/openaiconnector => ../connectors/openaiconnector

require (
	github.com/tzapio/tzap v0.0.0-00010101000000-000000000000
	github.com/tzapio/tzap/pkg/connectors/openaiconnector v0.0.0-00010101000000-000000000000
)

require (
	github.com/dlclark/regexp2 v1.9.0 // indirect
	github.com/sashabaranov/go-openai v1.12.0 // indirect
	github.com/tzapio/tokenizer v0.0.4 // indirect
)
