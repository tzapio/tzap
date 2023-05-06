package documents

import (
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/templates/code/files"
)

func ReadmeGithub(libraryDescription string, inspirationFiles []string, outputFile string, extraInstruction string) func(*tzap.Tzap) *tzap.Tzap {
	return func(t *tzap.Tzap) *tzap.Tzap {
		return t.
			AddSystemMessage(
				"Library Description: "+libraryDescription,
				"Output: Write github README.md presentation about the project. Use the files only as inspiration.").
			ApplyTemplate(files.InspirationTemplate(inspirationFiles)).
			AddUserMessage(extraInstruction).
			LoadTaskOrRequestNewTaskMD5(outputFile)
	}
}
