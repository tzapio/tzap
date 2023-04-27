package documents

import (
	"os"

	"github.com/tzapio/tzap/pkg/tzap"
)

func ReadmeGithub(libraryDescription string, inspirationFiles []string, outputFile string, extraInstruction string) func(*tzap.Tzap) *tzap.Tzap {
	return func(t *tzap.Tzap) *tzap.Tzap {
		wd, _ := os.Getwd()
		println("Hello", wd)
		return t.
			AddSystemMessage(
				"Library Description: "+libraryDescription,
				"Output: Write github README.md presentation about the project.").
			LoadFiles(inspirationFiles).
			Accumulate(func(t *tzap.Tzap) *tzap.Tzap {
				return t.AddUserMessage(
					"//file: "+t.Data["filepath"].(string),
					t.Data["content"].(string))
			}).
			AddUserMessage(extraInstruction).
			LoadTaskOrRequestNewTaskMD5(outputFile)
	}
}
