package main

import (
	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/types/openai"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/tzapconnect"

	"github.com/tzapio/tzap/pkg/util"
	"github.com/tzapio/tzap/templates/code/codegeneration"
)

func main() {
	tzap.NewWithConnector(tzapconnect.WithConfig(
		config.Configuration{
			AutoMode:    true,
			OpenAIModel: openai.GPT4,
			MD5Rewrites: true,
		})).
		LoadFileDir("/workspaces/goman/tzaps", "*.go").
		Map(func(t *tzap.Tzap) *tzap.Tzap {
			return t.
				ApplyTemplateFN(
					codegeneration.MakeCodeGO(`
You are helping the user writing a library for chatgpt prompting. You primarely write Golang. Most files already exists. Do not create new data structures.
### Current interface: 
//file: /workspaces/goman/tzaps/interfaces.go
`+util.ReadFileP("/workspaces/goman/tzaps/interfaces.go")+`
### General Tzap logic
//file: /workspaces/goman/tzaps/tzap.go
`+util.ReadFileP("/workspaces/goman/tzaps/tzap.go")+`
### Additional types
//file: /workspaces/goman/tzaps/msg/message.go
`+util.ReadFileP("/workspaces/goman/tzaps/msg/message.go"),
						//	"Make unit tests. Use testify go. If needed create tmp files. Use package tzap_test. Use testnames Test_{function}_{givenCamelCase}_{expectCamelCase}."),
						"Analyze what can be improved. Refactor the following file to be more readable. Make comments for the functions. Do not add any new public functions, only rewrite."),
				)
		})
}
