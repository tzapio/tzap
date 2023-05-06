// file: templates/code/files/findRelevantFile.go
package files

import (
	"path/filepath"
	"strings"

	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/util"
)

func FindRelevantFile(workdir, language, prompt, selectionCriteria string) types.NamedTemplate[*tzap.Tzap, *tzap.Tzap] {
	return types.NamedTemplate[*tzap.Tzap, *tzap.Tzap]{
		Name: "findRelevantFile",
		Template: func(t *tzap.Tzap) *tzap.Tzap {

			// List files in the workdir
			files := util.ListFilesInDir(workdir)

			// Filter files based on the provided language
			filteredFiles := filterFilesByLanguage(files, language)

			listPrompt := "### " + prompt + ":\n\n"
			for _, file := range filteredFiles {
				listPrompt += "[" + filepath.Base(file) + "]\n"
			}

			t = t.
				AddUserMessage(listPrompt). // Add message to prompt for relevant file
				AddUserMessage(             // Prompt GPT with the selection criteria to determine the relevant file/files
					strings.TrimSpace(strings.ReplaceAll(prompt, "\n", "")) + "\n" + selectionCriteria)

			// Add the selected file/files as user messages in the Tzap instance
			for _, file := range filteredFiles {
				t = t.AddUserMessage(file)
			}
			return t
		},
	}
}

// filterFilesByLanguage filters a list of files based on given language string
func filterFilesByLanguage(files []string, language string) []string {
	var filteredFiles []string

	fileExtension := languageToFileExtension(language)
	for _, file := range files {
		if filepath.Ext(file) == fileExtension {
			filteredFiles = append(filteredFiles, file)
		}
	}
	return filteredFiles
}

// languageToFileExtension returns the appropriate file extension for a given language string
func languageToFileExtension(language string) string {
	switch language {
	case "Go", "go", "Golang", "golang":
		return ".go"
	case "Python", "python", "py":
		return ".py"
	case "JavaScript", "javascript", "js":
		return ".js"
	case "TypeScript", "typescript", "ts":
		return ".ts"
	case "Kotlin", "kotlin", "kt":
		return ".kt"
	default:
		return ""
	}
}
