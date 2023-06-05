package stdinworkflows

import (
	"os"
	"os/exec"

	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/util/stdin"
)

func BeforeProceeding(changes string) string {
	file, err := os.CreateTemp("", "tzapchange*")
	if err != nil {
		panic(err)
	}
	defer os.Remove(file.Name())
	file.WriteString(changes)
	file.Close()

	lastChanges := changes
	for {
		println("\n\nFile: ", file.Name())
		println("")
		key := stdin.GetStdinInput("Edit files at file location.\n\n - press c and enter to open in vscode. \n - press v and enter to open in vim. \n - press enter to continue. \n\n")
		if key == "v" {
			// open vim
			cmd := exec.Command("vim", file.Name())
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stderr
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				panic(err)
			}
			key = ""
		}
		if key == "c" {
			// open code
			exec.Command("code", file.Name()).Run()
		}
		if key == "" {
			bytes, err := os.ReadFile(file.Name())
			if err != nil {
				panic(err)
			}
			defer file.Close()
			if string(bytes) == lastChanges {
				return string(bytes)
			} else {
				println("\n\nChanges detected:\n")
				println("---")
				println(string(bytes))
				println("---")
				lastChanges = string(bytes)
				continue
			}
		}
	}

}
func BeforeCompletionWorkflow() types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap] {
	return types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap]{
		Name: "BeforeCompletion",
		Workflow: func(t *tzap.Tzap) *tzap.Tzap {
			priorThread, err := t.Parent.GetThreadAsJSON()
			if err != nil {
				panic(err)
			}
			var outContent string
			t.IsolatedTzap(func(jt *tzap.Tzap) {
				newJson := BeforeProceeding(priorThread)
				jt = jt.LoadThreadString(newJson).
					HandleError(func(et *tzap.ErrorTzap) error {
						return et.Err
					}).RequestChatCompletion()
				outContent = jt.Data["content"].(string)
			})
			return t.AddTzap(&tzap.Tzap{
				Name: "BeforeCompletionWorkflow",
				Data: map[string]interface{}{
					"content": outContent,
				},
			})
		},
	}
}
func BeforeProceedingWorkflow() types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap] {
	return types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap]{
		Name: "ChangeCompletion",
		Workflow: func(t *tzap.Tzap) *tzap.Tzap {
			if t.Message.Content != "" {
				println("Warning - BeforeProceedingWorkflow - ChangeCompletion - Parent has message. This should only apply on Completion parents without message.")
			}
			config := config.FromContext(t.C)
			if config.AutoMode {
				return t
			}
			priorContent := t.Data["content"].(string)
			t.Data["content"] = BeforeProceeding(priorContent)
			return t
		},
	}
}
