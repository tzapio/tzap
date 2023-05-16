package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tzapio/tzap/cli/cmd/cmdutil"
	"github.com/tzapio/tzap/pkg/types/openai"
	"github.com/tzapio/tzap/pkg/tzap"
)

var magicWandCmd = &cobra.Command{
	Use:   "magicwand <file> <prompt>",
	Short: "Use magicwand to quickly edit code or documents",
	Long: `The 'magicwand' command allows you to use a magicwand to make quick edits to existing code or document,
including adding functions, comments, and clarifying error messages. 

This command takes two arguments: the file name and your prompt. The file must be relative to the current directory, 
or you should provide the absolute path to the file. If the file is not found, an error message will be displayed.`,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			absoluteFilePath, _ := filepath.Abs(filename)
			fmt.Println("Could not find file:", absoluteFilePath, " err:", err)
			os.Exit(1)
		}

		content := strings.Join(args[1:], " ")
		cmdutil.GetTzapFromContext(cmd.Context()).
			AddSystemMessage(fmt.Sprintf("Task: %s", content)).
			LoadCompletion(filename).
			MutationTzap(func(t *tzap.Tzap) *tzap.Tzap {
				t.Message.Role = openai.ChatMessageRoleUser
				return t
			}).
			LoadCompletionOrRequestCompletionMD5(filename)
	},
}

func init() {
	rootCmd.AddCommand(magicWandCmd)
}
