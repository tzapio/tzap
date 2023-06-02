package main

import (
	"os"
	"strings"

	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/tzapconnect"
	"github.com/tzapio/tzap/workflows/code/embedworkflows"
)

func main() {
	filename := os.Args[1]
	content := strings.Join(os.Args[2:], " ")
	openai_apikey, err := tzapconnect.LoadOPENAI_APIKEY()
	if err != nil {
		panic(err)
	}
	tzap.
		NewWithConnector(
			tzapconnect.WithConfig(openai_apikey, config.Configuration{MD5Rewrites: true})).
		ApplyWorkflow(embedworkflows.InspirationWorkflow(
			[]string{
				"README.md",
				"cli/cmd/semanticgitcommit.go",
				"pkg/types/structs.go",
				"pkg/tzap/workflows.go",
				"pkg/tzap/tzap.go",
				"workflows/code/gocode/arguments.go",
			},
		)).
		AddUserMessage(content).
		LoadCompletionOrRequestCompletion(filename)
}
