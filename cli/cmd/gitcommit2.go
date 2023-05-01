package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/tzapconnect"
	"github.com/tzapio/tzap/pkg/util/stdin"
)

const (
	maxTokensForGPT4    = 8000
	maxTokensForDefault = 4000
)

var gitcommit2Cmd = &cobra.Command{
	Use:   "gitcommit2",
	Short: "Generate a git commit message using ChatGPT",
	Long:  `Prompts ChatGPT to generate a commit message and commits it to the current git repo. The generated commit message is based on the diff of the currently staged files.`,
	Run: func(cmd *cobra.Command, args []string) {
		diff := exec.Command("git", "diff",
			"--staged",
			"--patch-with-raw",
			"--unified=2",
			"--color=never",
			"--no-renames",
			"--ignore-space-change",
			"--ignore-all-space",
			"--ignore-blank-lines",
		)
		out, err := diff.CombinedOutput()
		if err != nil {
			fmt.Println("Could not get diff:", err)
			return
		}
		fmt.Println(string(out))

		var contextSize int
		if settings.Model == "gpt4" {
			contextSize = maxTokensForGPT4
		} else {
			contextSize = maxTokensForDefault
		}

		tzapConnector := tzapconnect.WithConfig(config.Configuration{SupressLogs: true, OpenAIModel: modelMap[settings.Model]})
		t := tzap.NewWithConnector(tzapConnector).
			SetHeader(`Write using semantic commits specification. \n\n` + CV100).
			AddUserMessage(string(out))

		headerCount, err := t.CountTokens(t.Parent.Header)
		if err != nil {
			fmt.Println("Could not count tokens:", err)
			return
		}
		max := contextSize - headerCount - 500
		c, err := t.CountTokens(t.Message.Content)
		if err != nil {
			fmt.Println("Could not count tokens:", err)
			return
		}

		if c >= max {
			fmt.Printf("WARNING: diff is too long. TRUNCATING TO %d of %d estimated tokens\n", max, c)
		}
		ok := stdin.ConfirmToContinue()
		if !ok {
			return
		}

		offsetStart := 0
		offsetEnd := 0 + max
		t.Message.Content, err = t.OffsetTokens(t.Message.Content, offsetStart, offsetEnd)
		if err != nil {
			fmt.Println("Could not offset tokens:", err)
			return
		}

		content := t.RequestChat().Data["content"].(string)
		fmt.Println("\n", content)
		ok = stdin.ConfirmToContinue()
		if !ok {
			return
		}

		cmd2 := exec.Command("git", "commit", "-m", content)
		if err := cmd2.Run(); err != nil {
			fmt.Printf("Could not git commit. Content: %s, Error: %s\n", content, err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(gitcommit2Cmd)
}