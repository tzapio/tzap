package main

import (
	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/types/openai"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/tzapconnect"
	"github.com/tzapio/tzap/workflows/code/codegeneration"
)

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
		LoadFileDir("./pkg/embed").
		Map(func(t *tzap.Tzap) *tzap.Tzap {
			return t.
				ApplyWorkflowFN(
					codegeneration.MakeCode(
						codegeneration.BasicRefactoringConfig{},
						/*"golang", ".go", ".go", `
						You are helping the user writing a library for chatgpt prompting. You primarely write Golang. Most files already exists. Do not create new data structures.

						###file: ./pkg/types/interfaces.go
						`+util.ReadFileP("./pkg/types/interfaces.go")+`

						###file: ./pkg/tzap/tzap.go
						`+util.ReadFileP("./pkg/tzap/tzap.go")+`

						###file: ./pkg/types/embedding.go
						`+util.ReadFileP("./pkg/types/embedding.go"),
												"Make unit tests. Use testify go. If needed create tmp files. Use package tzap_test. Use testnames Test_{function}_{givenCamelCase}_{expectCamelCase}.", map[string]string{}),
						*/),
				)
		})
}
