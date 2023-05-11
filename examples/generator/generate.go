package main

import (
	"os"
	"strings"

	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/types/openai"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/tzapconnect"
	"github.com/tzapio/tzap/templates/code/embed"
)

func main() {
	filename := os.Args[1]
	content := strings.Join(os.Args[2:], " ")
	tzap.
		NewWithConnector(
			tzapconnect.WithConfig(
				config.Configuration{
					MD5Rewrites: true,
					OpenAIModel: openai.GPT4,
				})).
		ApplyTemplate(embed.InspirationTemplate(
			[]string{
				"README.md",
				"cli/cmd/semanticgitcommit.go",
				"pkg/types/structs.go",
				"pkg/tzap/templates.go",
				"pkg/tzap/tzap.go",
				"templates/code/gocode/arguments.go",
			},
		)).
		AddUserMessage(content).
		LoadTaskOrRequestNewTask(filename)
}
