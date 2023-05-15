package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/tzapio/tzap/cli/cmd/cmdutil"
	"github.com/tzapio/tzap/pkg/util/stdin"
)

var gitcommitCmd = &cobra.Command{
	Use:   "gitcommit",
	Short: "Prompts ChatGPT to generate a commit message and commits it to the current git repo",
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

		t := cmdutil.GetTzapFromContext(cmd.Context()).
			AddSystemMessage(`Write a git commit message maximum 30 words.
			
Workflow:
{brief git commit message}`).
			AddUserMessage(string(out))

		c, err := t.CountTokens(t.Message.Content)
		if err != nil {
			fmt.Println("Could not count tokens:", err)
			return
		}
		if c >= 3900 {
			fmt.Printf("WARNING: diff is too long. TRUNCATING TO 3900 of %d estimated tokens\n", c)
		}
		fmt.Printf("Summarizing %d estimated tokens\n", c)
		if !stdin.ConfirmPrompt("Continue?") {
			return
		}

		offsetStart := 0
		offsetEnd := 0 + 3900
		t.Message.Content, err = t.OffsetTokens(t.Message.Content, offsetStart, offsetEnd)
		if err != nil {
			fmt.Println("Could not offset tokens:", err)
			return
		}

		content := t.RequestChatCompletion().Data["content"].(string)
		fmt.Println("\n", content)
		if !stdin.ConfirmPrompt("Make the git commit?") {
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
	rootCmd.AddCommand(gitcommitCmd)
}
