package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/types/openai"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/tzapconnect"
)

var magicWandCmd = &cobra.Command{
	Use:   "magicwand <file> <prompt>",
	Short: "Use the magicwand to perform quick edits on code or documents, including adding functions, comments, and clarifying error messages.",
	Long: `The 'magicwand' command allows you to use a magicwand to make quick edits to existing code or documents, 
including adding functions, comments, and clarifying error messages. The command takes two arguments, the 
file name and your prompt. The file must exist in the current directory or you should provide the absolute
path for the file. If the file is not found, an error message will be printed on the console. Use this 
command to make quick edits and seek help from your teammates versus an AI model.`,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			absoluteFilePath, _ := filepath.Abs(filename)
			fmt.Println("Could not find file:", absoluteFilePath, " err:", err)
			os.Exit(1)
		}

		content := strings.Join(args[1:], " ")
		tzap.
			NewWithConnector(tzapconnect.WithConfig(config.Configuration{SupressLogs: true, MD5Rewrites: true})).
			AddSystemMessage(fmt.Sprintf("Task: %s", content)).
			LoadTask(filename).
			MutationTzap(func(t *tzap.Tzap) *tzap.Tzap {
				t.Message.Role = openai.ChatMessageRoleUser
				return t
			}).
			LoadTaskOrRequestNewTaskMD5(filename)
	},
}

func init() {
	rootCmd.AddCommand(magicWandCmd)
}
