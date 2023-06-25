package main

import (
	"encoding/json"

	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/types/openai"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/tzapconnect"
	"github.com/tzapio/tzap/workflows/code/codegeneration"
)

type RefactoringConfig struct {
	Language         string            `json:"language"`
	Type             string            `json:"type"`
	FileIn           string            `json:"inputs_extension"`
	FileOut          string            `json:"outputs_extension"`
	InspirationFiles map[string]string `json:"files"`
	Mission          string            `json:"mission"`
	Task             string            `json:"task"`
}

type FileConfig struct {
	Description string `json:"description"`
	Filepath    string `json:"filepath"`
}

func main() {
	configJSON := `
        {
                "language": "golang",
				"type":"basic",
				"filein:":"",
				"fileout":"",
                "files": [
                        {
                                "description": "Interfaces for the tzap library",
                                "filepath": "/workspaces/goman/tzaps/interfaces.go"
                        },
                        {
                                "description": "General Tzap logic",
                                "filepath": "/workspaces/goman/tzaps/tzap.go"
                        },
                        {
                                "description": "Additional types",
                                "filepath": "/workspaces/goman/tzaps/msg/message.go"
                        }
                ],
                "mission": "Analyze what can be improved. Refactor the following file to be more readable. Make comments for the functions. Do not add any new public functions, only rewrite.",
                "task": "Make unit tests. Use testify go. If needed create tmp files. Use package tzap_test. Use testnames Test_{function}_{givenCamelCase}_{expectCamelCase}"
        }`
	var configData RefactoringConfig
	if err := json.Unmarshal([]byte(configJSON), &configData); err != nil {
		panic(err)
	}

	openai_apikey, err := tzapconnect.LoadOPENAI_API_KEY()
	if err != nil {
		panic(err)
	}

	tzap.
		NewWithConnector(
			tzapconnect.WithConfig(openai_apikey, config.Configuration{
				MD5Rewrites: true,
				OpenAIModel: openai.GPT4,
				EnableLogs:  true})).
		LoadFileDir("/workspaces/goman/tzaps").
		Map(func(t *tzap.Tzap) *tzap.Tzap {
			return t.
				ApplyWorkflowFN(
					codegeneration.MakeCode(
						codegeneration.BasicRefactoringConfig{}),
				)
		})
}
