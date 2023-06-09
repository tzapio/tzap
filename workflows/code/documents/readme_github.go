package documents

import (
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/workflows/code/fileworkflows"
)

func ReadmeGithub(libraryDescription string, inspirationFiles []string, outputFile string, extraInstruction string) func(*tzap.Tzap) *tzap.Tzap {
	return func(t *tzap.Tzap) *tzap.Tzap {
		return t.
			AddSystemMessage(
				"Library Description: "+libraryDescription,
				"Output: Write github README.md presentation about the project. Use the files only as inspiration.").
			ApplyWorkflow(fileworkflows.InspirationWorkflow(inspirationFiles)).
			AddUserMessage(extraInstruction).
			LoadCompletionOrRequestCompletionMD5(outputFile)
	}
}
