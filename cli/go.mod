module github.com/tzapio/tzap/cli

go 1.20

replace github.com/tzapio/tzap => ../

replace github.com/tzapio/tzap/pkg/tzapconnect => ../pkg/tzapconnect

replace github.com/tzapio/tzap/pkg/connectors/openaiconnector => ../pkg/connectors/openaiconnector

require (
	github.com/AlecAivazis/survey/v2 v2.3.6
	github.com/fatih/color v1.15.0
	github.com/sabhiram/go-gitignore v0.0.0-20210923224102-525f6e181f06
	github.com/spf13/cobra v1.7.0
	github.com/tzapio/tzap v0.0.0-00010101000000-000000000000
	github.com/tzapio/tzap/pkg/tzapconnect v0.0.0-00010101000000-000000000000
)

require (
	github.com/dlclark/regexp2 v1.9.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/mgutz/ansi v0.0.0-20170206155736-9520e82c474b // indirect
	github.com/sashabaranov/go-openai v1.9.3 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/testify v1.8.1 // indirect
	github.com/tiktoken-go/tokenizer v0.1.0 // indirect
	github.com/tzapio/tzap/pkg/connectors/openaiconnector v0.0.0-00010101000000-000000000000 // indirect
	golang.org/x/sys v0.7.0 // indirect
	golang.org/x/term v0.7.0 // indirect
	golang.org/x/text v0.9.0 // indirect
)
