package documents

import (
	"os"

	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/templates/code/files"
)

func ReadmeGithub(libraryDescription string, inspirationFiles []string, outputFile string, extraInstruction string) func(*tzap.Tzap) *tzap.Tzap {
	return func(t *tzap.Tzap) *tzap.Tzap {
		wd, _ := os.Getwd()
		println("Hello", wd)
		return t.
			AddSystemMessage(
				"Library Description: "+libraryDescription,
				"Output: Write github README.md presentation about the project.").
			ApplyTemplate(files.InspirationTemplate(inspirationFiles)).
			AddUserMessage(extraInstruction).
			LoadTaskOrRequestNewTask(outputFile)
	}
}
