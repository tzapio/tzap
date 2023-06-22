package main

import (
	"github.com/tzapio/tzap/internal/logging/mermaid"
	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/types/openai"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/tzapconnect"
)

func main() {
	openai_apikey, err := tzapconnect.LoadOPENAI_API_KEY()
	if err != nil {
		panic(err)
	}

	t := tzap.
		NewWithConnector(
			tzapconnect.WithConfig("", openai_apikey, config.Configuration{
				MD5Rewrites: true,
				OpenAIModel: openai.GPT4,
				EnableLogs:  true}))

	// Assume there's a Tzap hierarchy created like t1 -> t2 -> t3
	t1 := t.CopyConnection().AddAssistantMessage("t2", "t2 message")

	// Fill the Mermaid graph with the Tzap hierarchy
	graph := mermaid.FillMermaidGraph(t1)

	// Generate the Mermaid markup file based on the created graph

	if err := mermaid.GenerateMermaidMarkupFile("out/mermaid.mmd", graph); err != nil {
		panic(err)
	}
}
