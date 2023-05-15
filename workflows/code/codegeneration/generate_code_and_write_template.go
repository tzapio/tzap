package codegeneration

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/workflows/gptasfunction"
)

type CodeGeneration struct {
	Code     string
	FilePath string
	Type     string
}

func GenerateCodeAndApplyWorkflow() types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap] {
	return types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap]{
		Name: "GenerateCodeAndApplyWorkflow",
		Workflow: func(t *tzap.Tzap) *tzap.Tzap {
			return t.
				WorkTzap(func(t *tzap.Tzap) {
					// Load the file content
					t.IsolatedTzap(func(ti *tzap.Tzap) {
						codeContent, ok := t.Data["content"].(string)
						if !ok {
							panic("Could not extract content")
						}
						println("code:" + codeContent)
						ti.
							AddSystemMessage("You are rewriting code into JSON. It is not allowed to use newline in the json.",
								`Template:
- {"code":"{code}","filePath":"{filePath}", "type": "{full code OR partial code}"}
`).
							AddAssistantMessage(codeContent).
							AddUserMessage("Extract as JSON").
							RequestChatCompletion(). // Run the completion
							WorkTzap(func(t *tzap.Tzap) {
								completion, ok := t.Data["content"].(string)
								if !ok {
									panic("Could not extract content")
								}
								codeSegments := strings.Split(completion, "\n")

								var jsonObjects []CodeGeneration
								for _, codeSeg := range codeSegments {
									obj := CodeGeneration{}
									codeIndex := strings.Index(codeSeg, "{")
									if codeIndex == -1 {
										continue
									}
									if err := json.Unmarshal([]byte(codeSeg[codeIndex:]), &obj); err != nil {
										continue
									}
									jsonObjects = append(jsonObjects, obj)
								}
								for _, jsonObject := range jsonObjects {
									// Use the GPTAsFunction worfklow to transform the JSON object
									if _, err := os.Stat(jsonObject.FilePath); os.IsNotExist(err) {
										os.WriteFile(jsonObject.FilePath, []byte(jsonObject.Type), 0644)
										return
									}
									oldFileContent, err := os.ReadFile(jsonObject.FilePath)
									if err != nil {
										panic(err)
									}
									sysPrompt := "Transfer the changes onto the user response. You are now editing:" +
										jsonObject.FilePath + "\n\nchanges:\n" + jsonObject.Code
									t.ApplyWorkflow(gptasfunction.GPTAsFunction(sysPrompt,
										string(oldFileContent))).
										StoreCompletion(jsonObject.FilePath)
								}
							})
					})
				})

		},
	}
}
