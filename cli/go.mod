module github.com/tzapio/tzap/cli

go 1.20

replace github.com/tzapio/tzap => ../

replace github.com/tzapio/tzap/pkg/tzapaction => ../pkg/tzapaction

replace github.com/tzapio/tzap/pkg/tzapconnect => ../pkg/tzapconnect

replace github.com/tzapio/tzap/pkg/connectors/openaiconnector => ../pkg/connectors/openaiconnector

require (
	github.com/fatih/color v1.15.0
	github.com/fsnotify/fsnotify v1.6.0
	github.com/sabhiram/go-gitignore v0.0.0-20210923224102-525f6e181f06
	github.com/spf13/cobra v1.7.0
	github.com/stretchr/testify v1.8.4
	github.com/tzapio/tzap v0.9.4
	github.com/tzapio/tzap/pkg/tzapaction v0.9.4
	github.com/tzapio/tzap/pkg/tzapconnect v0.9.4
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dlclark/regexp2 v1.10.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/labstack/echo/v4 v4.11.1 // indirect
	github.com/labstack/gommon v0.4.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sashabaranov/go-openai v1.14.1 // indirect
	github.com/sergi/go-diff v1.3.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/tzapio/tokenizer v0.0.4 // indirect
	github.com/tzapio/tzap/pkg/connectors/openaiconnector v0.9.4 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	golang.org/x/crypto v0.11.0 // indirect
	golang.org/x/net v0.12.0 // indirect
	golang.org/x/sys v0.10.0 // indirect
	golang.org/x/text v0.11.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
