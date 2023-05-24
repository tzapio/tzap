package main

import (
	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/types/openai"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/tzapconnect"
	"github.com/tzapio/tzap/pkg/util"
	"github.com/tzapio/tzap/workflows/code/translate"
)

var mission1 string = `
Tzap is a library for Prompts as Code. It provides a toolkit to build, customize, and extend chatbot prompts in a streamlined and extensible manner. 
### General Tzap logic
//file: /home/vscode/go/src/github.com/tzapio/tzap/pkg/tzap/tzap.go
` + util.ReadFileP("/home/vscode/go/src/github.com/tzapio/tzap/pkg/tzap/tzap.go") + `
### Additional types
//file: /home/vscode/go/src/github.com/tzapio/tzap/pkg/tzap/message.go
` + util.ReadFileP("/home/vscode/go/src/github.com/tzapio/tzap/pkg/tzap/message.go")

var mission2 string = `
Tzap is a library for Prompts as Code. It provides a toolkit to build, customize, and extend chatbot prompts in a streamlined and extensible manner. 
### Typescript interface for the Tzap library
//file: /home/vscode/go/src/github.com/tzapio/tzap/ts/src/interface.ts
` + util.ReadFileP("/home/vscode/go/src/github.com/tzapio/tzap/ts/src/interface.ts") + `
### Additional types
//file: /home/vscode/go/src/github.com/tzapio/tzap/pkg/tzap/message.go
` + util.ReadFileP("/home/vscode/go/src/github.com/tzapio/tzap/pkg/tzap/message.go")

var mission3 string = `
Tzap is a library for Prompts as Code. It provides a toolkit to build, customize, and extend chatbot prompts in a streamlined and extensible manner. 
### Typescript interface for the Tzap library
//file: /home/vscode/go/src/github.com/tzapio/tzap/ts/src/interface.ts
` + util.ReadFileP("/home/vscode/go/src/github.com/tzapio/tzap/ts/src/interface.ts") + `
### General Golang Tzap logic
//file: /home/vscode/go/src/github.com/tzapio/tzap/pkg/tzap/tzap.go
` + util.ReadFileP("/home/vscode/go/src/github.com/tzapio/tzap/pkg/tzap/tzap.go") + `
### General Typescript Tzap logic
//file: /home/vscode/go/src/github.com/tzapio/tzap/ts/src/interface.ts
` + util.ReadFileP("/home/vscode/go/src/github.com/tzapio/tzap/ts/src/interface.ts") + `
### Additional types
//file: /home/vscode/go/src/github.com/tzapio/tzap/pkg/tzap/message.go
` + util.ReadFileP("/home/vscode/go/src/github.com/tzapio/tzap/pkg/tzap/message.go")

func main() {

	openai_apikey, err := tzapconnect.LoadOPENAI_APIKEY()
	if err != nil {
		panic(err)
	}
	tzap.
		NewWithConnector(
			tzapconnect.WithConfig(openai_apikey, config.Configuration{
				MD5Rewrites: true,
				OpenAIModel: openai.GPT4,
				EnableLogs:  true})).
		WorkTzap(func(t *tzap.Tzap) {
			t.
				ApplyWorkflow(translate.MakeCodeTSMessage(
					mission1,
					"Generate a typescript interface for the Tzap library.",
					"/home/vscode/go/src/github.com/tzapio/tzap/pkg/types/interfaces.go",
					"/home/vscode/go/src/github.com/tzapio/tzap/ts/src/interface.ts"))
		}).
		WorkTzap(func(t *tzap.Tzap) {
			t.
				ApplyWorkflow(translate.MakeCodeTSMessage(
					mission2,
					"Translate Tzap.go into Tzap.ts.",
					"/home/vscode/go/src/github.com/tzapio/tzap/pkg/tzap/tzap.go",
					"/home/vscode/go/src/github.com/tzapio/tzap/ts/src/tzap.ts"))
		}).
		WorkTzap(func(t *tzap.Tzap) {
			t.LoadFileDir("/home/vscode/go/src/github.com/tzapio/tzap/pkg/util/").
				Map(func(t *tzap.Tzap) *tzap.Tzap {
					return t.
						ApplyWorkflow(
							translate.TranslateCodeFromTo(
								"go",
								"ts",
								"/home/vscode/go/src/github.com/tzapio/tzap/ts/src/",
								mission3+`
### Typescript interface for the Tzap library
//file: /home/vscode/go/src/github.com/tzapio/tzap/ts/src/interface.ts
`+util.ReadFileP("/home/vscode/go/src/github.com/tzapio/tzap/ts/src/interface.ts"),
								"Translate golang library for chatgpt prompting to typescript. You primarely write Typescript. Most files already exists.",
							),
						)
				})
		})

}
