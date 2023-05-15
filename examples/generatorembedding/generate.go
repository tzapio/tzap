package main

import (
	"os"
	"strings"

	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/types/openai"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/tzapconnect"
	"github.com/tzapio/tzap/workflows/code/embed"
)

func main() {
	filename := os.Args[1]
	println(filename)
	content := strings.Join(os.Args[2:], " ")
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
		ApplyWorkflow(embed.EmbeddingInspirationWorkflow(
			content,
			[]string{
				"./pkg/types/interfaces.go",
				"./pkg/tzap/tzap.go",
				//	"./workflows/code/files/embeddingInspirationWorkflow.go",
			}, 10, 15)).
		AddUserMessage(content).
		LoadCompletionOrRequestCompletionMD5(filename)
}
