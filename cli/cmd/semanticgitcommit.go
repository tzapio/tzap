package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/tzapconnect"
	"github.com/tzapio/tzap/pkg/util/stdin"
	"github.com/tzapio/tzap/templates/code/gocode"
)

func RequestChat() types.NamedTemplate[*tzap.Tzap, *tzap.ErrorTzap] {
	return types.NamedTemplate[*tzap.Tzap, *tzap.ErrorTzap]{
		Name: "RequestChat",
		Template: func(t *tzap.Tzap) *tzap.ErrorTzap {
			extraPrompt := t.Data["extraPrompt"].(string)
			diff := t.Data["git-diff"].(string)
			t = t.SetHeader(`Write using semantic commits specification. \n\n` + CV100)
			if extraPrompt != "" {
				t = t.AddUserMessage(extraPrompt)
			}
			t = t.AddUserMessage(diff)

			return t.RequestChat().ErrorTzap(nil)
		}}
}

var gitcommit3Cmd = &cobra.Command{
	Use:   "semantic:gitcommit [clarifying prompt]",
	Short: "Generate a git commit message using ChatGPT",
	Long:  `Prompts ChatGPT to generate a commit message and commits it to the current git repo. The generated commit message is based on the diff of the currently staged files.`,
	Run: func(cmd *cobra.Command, args []string) {
		tzapConnector := tzapconnect.WithConfig(config.Configuration{SupressLogs: true, OpenAIModel: modelMap[settings.Model]})
		tzap.NewWithConnector(tzapConnector).
			ApplyTemplate(gocode.DeserializedArguments("extraPrompt", args)).
			ApplyErrorTemplate(gocode.GitDiff(), func(et *tzap.ErrorTzap) error {
				return et.Err
			}).
			ApplyErrorTemplate(gocode.ValidateDiff(), func(et *tzap.ErrorTzap) error {
				return et.Err
			}).
			ApplyTemplate(gocode.SetContextSize()).
			ApplyErrorTemplate(gocode.CountTokens(), func(et *tzap.ErrorTzap) error {
				return et.Err
			}).
			ApplyErrorTemplate(gocode.TruncateTokens(), func(et *tzap.ErrorTzap) error {
				return et.Err
			}).
			ApplyErrorTemplate(RequestChat(), func(et *tzap.ErrorTzap) error {
				return et.Err
			}).
			ApplyErrorTemplate(gocode.DisplayAndConfirm(), func(et *tzap.ErrorTzap) error {
				return et.Err
			}).
			ApplyErrorTemplate(gocode.GitCommit(), func(et *tzap.ErrorTzap) error {
				return et.Err
			})
	},
}
var gitcommit2Cmd = &cobra.Command{
	Use:   "semantic:gitcommit [clarifying prompt]",
	Short: "Generate a git commit message using ChatGPT",
	Long:  `Prompts ChatGPT to generate a commit message and commits it to the current git repo. The generated commit message is based on the diff of the currently staged files.`,
	Run: func(cmd *cobra.Command, args []string) {
		extraPrompt := strings.Join(args, " ")
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
			contextSize = 8000
		} else {
			contextSize = 4000
		}

		tzapConnector := tzapconnect.WithConfig(config.Configuration{SupressLogs: true, OpenAIModel: modelMap[settings.Model]})
		t := tzap.NewWithConnector(tzapConnector).
			SetHeader(`Write using semantic commits specification. \n\n` + CV100)
		if extraPrompt != "" {
			t = t.AddUserMessage(extraPrompt)
		}
		t = t.AddUserMessage(string(out))

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
		if c == 0 {
			fmt.Printf("Diff is empty. Stage files to continue.\n")
			return
		}
		println("Estimated tokens:", c)
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
	rootCmd.AddCommand(gitcommit3Cmd)
}
