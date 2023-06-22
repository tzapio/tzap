package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/tzapio/tzap/example/tdev/prompts"

	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/template"
	"github.com/tzapio/tzap/pkg/types/openai"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/tzapconnect"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: tdev.go <prompt>")
		return
	}
	projectName := os.Args[1]
	prompt := os.Args[2]
	generated := path.Join("generated", projectName)

	if err := os.MkdirAll(generated, 0755); err != nil {
		panic(err)
	}
	openai_apikey, err := tzapconnect.LoadOPENAI_API_KEY()
	if err != nil {
		panic(err)
	}

	t := tzap.
		NewWithConnector(
			tzapconnect.WithConfig("", openai_apikey, config.Configuration{
				OpenAIModel: openai.GPT16,
				MD5Rewrites: true,
				AutoMode:    true,
				EnableLogs:  true}))

	t.
		AddSystemMessage(prompts.First).
		AddUserMessage(prompt).
		LoadCompletionOrRequestCompletion(path.Join(generated, "filepaths.json")).
		WorkTzap(func(t *tzap.Tzap) {
			filePaths := t.Data["content"].(string)
			var fileP []string
			json.Unmarshal([]byte(filePaths), &fileP)
			fmt.Println(filePaths)

			secondStep := template.NewWorkflowStep("second", prompts.Second)
			output, err := secondStep.Execute(map[string]interface {
			}{
				"prompt":           prompt,
				"filepaths_string": filePaths,
			})
			if err != nil {
				panic(err)
			}
			t.IsolatedTzap(func(t *tzap.Tzap) {

				sharedDependenciesPath := path.Join(generated, "shared_dependencies.md")
				q := t.
					AddSystemMessage(output).
					AddUserMessage(prompt).
					LoadCompletionOrRequestCompletion(sharedDependenciesPath)
				sharedDependencies := q.Data["content"].(string)

				thirdStep := template.NewWorkflowStep("third", prompts.Third)
				thirdOutput, err := thirdStep.Execute(map[string]interface {
				}{
					"prompt":              prompt,
					"filepaths_string":    filePaths,
					"shared_dependencies": sharedDependencies,
				})
				forth := template.NewWorkflowStep("forth", prompts.Forth)
				if err != nil {
					panic(err)
				}
				for _, file := range fileP {
					/*if fileutil.IsNotFolderAndHasExtension(file) {*/
					forthOutput, _ := forth.Execute(map[string]interface {
					}{
						"prompt":              prompt,
						"filename":            file,
						"filepaths_string":    filePaths,
						"shared_dependencies": sharedDependencies,
					})

					t.
						AddSystemMessage(thirdOutput).
						AddUserMessage(forthOutput).
						RequestChatCompletion().
						StoreCompletion(path.Join(generated, file))

					//}
				}
			})
		})
}
