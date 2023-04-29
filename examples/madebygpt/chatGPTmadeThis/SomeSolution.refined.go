package main

import (
	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/types/openai"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/tzapconnect"
)

func main() {
	// Initialize Tzap with configuration options.
	t := tzap.NewWithConnector(tzapconnect.WithConfig(
		config.Configuration{
			AutoMode:    true,
			OpenAIModel: openai.GPT4,
			MD5Rewrites: true,
		}))

	// Load all Go files in the directory.
	t.
		LoadFileDir("/path/to/your/codebase", "*.go").
		// Apply improvements using GPT-3.5-turbo.
		Map(func(t *tzap.Tzap) *tzap.Tzap {
			return t.ApplyTemplateFN(
				// Customize the prompt template to optimize the codebase.
				func(t *tzap.Tzap) *tzap.Tzap {
					// Read file content.
					content := t.Data["content"].(string)
					filepath := t.Data["filepath"].(string)
					return t.
						// Set up the prompt template.
						AddUserMessage("You are an AI code reviewer powered by OpenAI's GPT-3.5-turbo. Your goal is to review and improve the Go code provided. Analyze the following code and suggest improvements to it. Make sure you also provide clear explanations for your suggestions.").
						AddAssistantMessage(
							"//file: "+filepath,
							content).
						// Fetch suggestions from GPT-3.5-turbo.
						// NOTE: this does not save results anywhere.
						RequestChat()
				},
			)
		}).
		// Output the improvement report to the console.
		Each(func(t *tzap.Tzap) {
			filename := t.Data["filepath"].(string)
			fileImprovements := t.Message.Content

			println("File:", filename)
			if fileImprovements != "" {
				println(fileImprovements)
				println("-------------------------")
			} else {
				println("No improvements suggested.")
				println("-------------------------")
			}
		})
}
