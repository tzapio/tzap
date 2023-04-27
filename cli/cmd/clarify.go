package cmd

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/types/openai"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/tzapconnect"
)

var magicCmd = &cobra.Command{
	Use:   "magic",
	Short: "attaches your file and adds your prompt",
	Long:  `tbd`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("magic called")
		println(args[0])
		if len(args) < 2 {
			fmt.Println("tzap magic <file> <prompt...>")
			os.Exit(1)
		}
		filename := args[0]
		if _, err := os.Open(filename); err != nil {
			absoluteFilePath := path.Join(os.Getenv("PWD"), filename)
			fmt.Println("Could not find file:", absoluteFilePath, " err:", err)
			os.Exit(1)
		}
		content := strings.Join(args[1:], " ")
		tzap.
			NewWithConnector(
				tzapconnect.WithConfig(config.Configuration{SupressLogs: true, MD5Rewrites: true}),
			).
			SetHeader(`Task: ` + content).
			LoadTask(filename). //add as message
			MutationTzap(func(t *tzap.Tzap) *tzap.Tzap {
				t.Message.Role = openai.ChatMessageRoleUser
				return t
			}).
			LoadTaskOrRequestNewTaskMD5(filename) // ask to rewrite.

	},
}

func init() {
	rootCmd.AddCommand(magicCmd)
}
