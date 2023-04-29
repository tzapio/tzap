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

var magicCmd = &cobra.Command{
	Use:   "magic <file> <prompt>",
	Short: "Attach a file and add a prompt to a Tzap message",
	Long:  `Attach a file and add a prompt to a Tzap message.`,
	Args:  cobra.ExactArgs(2),
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
	rootCmd.AddCommand(magicCmd)
}
