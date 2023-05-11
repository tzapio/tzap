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
	println(filename)
	content := strings.Join(os.Args[2:], " ")
	tzap.
		NewWithConnector(
			tzapconnect.WithConfig(
				config.Configuration{
					MD5Rewrites: true,
					OpenAIModel: openai.GPT4,
					EnableLogs:  true,
				})).
		ApplyTemplate(embed.EmbeddingInspirationTemplate(
			content,
			[]string{
				"./pkg/types/interfaces.go",
				"./pkg/tzap/tzap.go",
				//	"./templates/code/files/embeddingInspirationTemplate.go",
			},
		)).
		AddUserMessage(content).
		LoadTaskOrRequestNewTaskMD5(filename)
}
