package main

import (
	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/types/openai"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/tzapconnect"
	"github.com/tzapio/tzap/templates/code/documents"
)

func main() {
	tzap.NewWithConnector(tzapconnect.WithConfig(
		config.Configuration{
			AutoMode:    true,
			OpenAIModel: openai.GPT4,
			MD5Rewrites: true,
		})).
		ApplyTemplateFN(documents.ReadmeGithub(
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
