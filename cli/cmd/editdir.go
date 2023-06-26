// editdirCmd refactor source code using the given template in the given directory.
// templates available: Add documentation to code (documentation), Refactor code (refactor), Add Unit Tests (unittest), Mermaid (mermaid)".
// The command takes two arguments. The first argument is the template name, and the second argument is the directory path.
// It lists all files in the directory, sets some configurations, and runs the refactoring command using the provided template.
// It also includes some user interface operations such as showing the updated files, handling panics, and logging.
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tzapio/tzap/cli/action"
	"github.com/tzapio/tzap/cli/cmd/cliworkflows"
	"github.com/tzapio/tzap/cli/cmd/cmdutil"
	"github.com/tzapio/tzap/internal/logging/tl"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/util"
	"github.com/tzapio/tzap/workflows/code/codegeneration"
	"github.com/tzapio/tzap/workflows/stdinworkflows"
)

var editdirCmd = &cobra.Command{
	Use:    "editdir [template] [directory]",
	Short:  "Refactor source code using the given template in the given directory",
	Hidden: true,
	Args:   cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		// store the template and directory path
		template := args[0]
		dir := args[1]

		tl.EnableUICompletionLogger()

		files, err := util.ListGlob(dir)
		if err != nil {
			return fmt.Errorf("failed to list files in directory: %w", err)
		}

		err = tzap.HandlePanic(func() {
			t := cmdutil.GetTzapFromContext(cmd.Context())
			defer t.HandleShutdown()

			for _, file := range files {
				config := getConfig(template, file)

				cmd.Println("Editing file: ", file, "template", template)
				promptWorkflowArgs := action.PromptWorkflowArgs{
					InspirationFiles: []string{config.FileIn},
					SearchQuery:      util.ReadFileP(config.FileIn),
					EmbedsCount:      60,
					NCount:           60,
					DisableIndex:     false,
					Yes:              tzapCliSettings.Yes,
					MessageThread:    []types.Message{},
				}
				t.
					ApplyWorkflow(action.PromptWorkflow(promptWorkflowArgs)).
					ApplyWorkflowFN(codegeneration.MakeCode(config)).
					ApplyWorkflow(stdinworkflows.BeforeProceedingWorkflow()).
					ApplyWorkflow(cliworkflows.PrintFileDiff(config.FileOut)).
					StoreCompletion(config.FileOut)
			}
		})

		if err != nil {
			return fmt.Errorf("failed to refactor directory: %w", err)
		}

		return nil
	},
}

func getConfig(template string, file string) codegeneration.BasicRefactoringConfig {
	if template == "refactor" {
		refactorConfig := codegeneration.BasicRefactoringConfig{
			FileIn:           file,
			FileOut:          file,
			Mission:          "Improve code readability and maintainability",
			Task:             "Refactor code to use better variable names and remove duplication. Refactor the following file to be more readable. Make comments for the functions. Do not add any new public functions, only rewrite.",
			Plan:             "Do not write any text because this file will be saved directly to " + file,
			OutputFormat:     "{golang code}",
			Example:          "package something\nfunc doSomething(w http.ResponseWriter, r *http.Request, db *sql.DB) error {\n          // function logic\n}",
			InspirationFiles: []string{},
		}
		return refactorConfig
	}
	if template == "documentation" {
		fileOut := util.ReplaceExt(file, ".md")
		refactorConfig := codegeneration.BasicRefactoringConfig{
			FileIn:           file,
			FileOut:          fileOut,
			Mission:          "Document code in markdown format",
			Task:             "Add documentation to the following file. Add comments to the functions. Add a markdown file with the same name as the file and add the documentation to it.",
			Plan:             "Do not write any text because this file will be saved directly to " + fileOut,
			OutputFormat:     "{markdown documentation}",
			InspirationFiles: []string{},
		}
		return refactorConfig
	}
	if template == "test" {
		fileOut := util.ReplaceExt(file, "_test.go")
		refactorConfig := codegeneration.BasicRefactoringConfig{
			FileIn:           file,
			FileOut:          fileOut,
			Mission:          "Add golang unit tests to the code",
			Task:             "Add unit tests to the following file. Add a test file with the same name as the file and add the unit tests to it.",
			Plan:             "Do not write any text because this file will be saved directly to " + fileOut,
			OutputFormat:     "{golang unit test}",
			InspirationFiles: []string{},
		}
		return refactorConfig
	}
	return codegeneration.BasicRefactoringConfig{}
}
func init() {
	RootCmd.AddCommand(editdirCmd)
}
