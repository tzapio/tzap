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
				"pkg/types/interfaces.go",
				"pkg/types/structs.go",
				"pkg/tzap/templates.go",
				"pkg/tzap/file.go",
				"pkg/tzap/fetch-chat.go",
				"pkg/tzap/tzap.go",
				"examples/githubdoc/main.go",
				"examples/refactoring/main.go",
			},
			"README.md",
			"",
		))
}
