package main

import (
	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/types/openai"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/tzapconnect"
	"github.com/tzapio/tzap/workflows/code/documents"
)

func main() {
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
		ApplyWorkflowFN(documents.ReadmeGithub(
			"Tzap is a library for Prompts as Code.",
			[]string{
				".tzap/usecases",
				"pkg/types/interfaces.go",
				".tzap/examples",
				"cli/cmd/semanticgitcommit.go",
				"pkg/tzap/tzap.go",
				"examples/githubdoc/main.go",
				"README.md",
			},
			"README.md",
			"",
		))
}
