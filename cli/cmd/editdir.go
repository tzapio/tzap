package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tzapio/tzap/cli/cmd/cmdutil"
	"github.com/tzapio/tzap/internal/logging/tl"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/workflows/code/codegeneration"
	"github.com/tzapio/tzap/workflows/stdinworkflows"
)

var refactordirCmd = &cobra.Command{
	Use:    "editdir [template] [directory]",
	Short:  "Perform the refactoring operation on a directory. Templates available: Add documentation to code (documentation), Refactor code (refactor), Add Unit Tests (unittest), Mermaid (mermaid)",
	Hidden: true,
	Args:   cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := args[0]
		tl.EnableUICompletionLogger()

		refactorConfig := codegeneration.BasicRefactoringConfig{
			FileIn: dir,
		}

		// Customize the refactorConfig based on your requirements

		err := tzap.HandlePanic(func() {
			t := cmdutil.GetTzapFromContext(cmd.Context())
			defer t.HandleShutdown()
			t.
				ApplyWorkflowFN(codegeneration.MakeCode(refactorConfig)).
				ApplyWorkflow(stdinworkflows.BeforeProceedingWorkflow()).
				StoreCompletion(refactorConfig.FileOut)

		})

		if err != nil {
			return fmt.Errorf("failed to refactor directory: %w", err)
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(refactordirCmd)
}
