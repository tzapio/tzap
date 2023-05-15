package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tzapio/tzap/cli/cmd/cmdutil"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/workflows/code/git"
	"github.com/tzapio/tzap/workflows/code/gocode"
	"github.com/tzapio/tzap/workflows/truncate"
)

var semanticGitcommitCmd = &cobra.Command{
	Aliases: []string{"c"},
	Use:     "semantic:gitcommit [clarifying prompt]",
	Short:   "Generate a git commit message using ChatGPT",
	Long:    `Prompts ChatGPT to generate a commit message and commits it to the current git repo. The generated commit message is based on the diff of the currently staged files.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := tzap.HandlePanic(func() {
			t := cmdutil.GetTzapFromContext(cmd.Context())
			defer t.HandleShutdown()
			t.
				ApplyWorkflow(gocode.DeserializedArguments("extraPrompt", args)).
				ApplyErrorWorkflow(git.GitDiff(), func(et *tzap.ErrorTzap) error {
					return et.Err
				}).
				ApplyErrorWorkflow(git.ValidateDiff(), func(et *tzap.ErrorTzap) error {
					return et.Err
				}).
				ApplyWorkflow(truncate.SetContextSize()).
				ApplyErrorWorkflow(truncate.CountTokens(), func(et *tzap.ErrorTzap) error {
					return et.Err
				}).
				ApplyErrorWorkflow(truncate.TruncateTokens(), func(et *tzap.ErrorTzap) error {
					return et.Err
				}).
				ApplyErrorWorkflow(RequestChat(), func(et *tzap.ErrorTzap) error {
					return et.Err
				}).
				ApplyErrorWorkflow(gocode.DisplayAndConfirm(), func(et *tzap.ErrorTzap) error {
					return et.Err
				}).
				ApplyErrorWorkflow(git.GitCommit(), func(et *tzap.ErrorTzap) error {
					return et.Err
				})
		})
		if err != nil {
			println(err.Error())
		}
	},
}

// RequestChat is a workflow that requests a chat from ChatGPT.
func RequestChat() types.NamedWorkflow[*tzap.Tzap, *tzap.ErrorTzap] {
	return types.NamedWorkflow[*tzap.Tzap, *tzap.ErrorTzap]{
		Name: "RequestChat",
		Workflow: func(t *tzap.Tzap) *tzap.ErrorTzap {
			extraPrompt := t.Data["extraPrompt"].(string)
			diff := t.Data["git-diff"].(string)
			t = t.SetInitialSystemContent(`Write one commit using semantic commit specification. \n\n` + CV100)
			if extraPrompt != "" {
				t = t.AddUserMessage(extraPrompt)
			}
			t = t.AddUserMessage(diff)

			return t.RequestChatCompletion().ErrorTzap(nil)
		}}
}

func init() {
	rootCmd.AddCommand(semanticGitcommitCmd)
}
