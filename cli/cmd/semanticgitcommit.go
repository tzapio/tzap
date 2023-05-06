package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/tzapconnect"
	"github.com/tzapio/tzap/templates/code/git"
	"github.com/tzapio/tzap/templates/code/gocode"
)

func RequestChat() types.NamedTemplate[*tzap.Tzap, *tzap.ErrorTzap] {
	return types.NamedTemplate[*tzap.Tzap, *tzap.ErrorTzap]{
		Name: "RequestChat",
		Template: func(t *tzap.Tzap) *tzap.ErrorTzap {
			extraPrompt := t.Data["extraPrompt"].(string)
			diff := t.Data["git-diff"].(string)
			t = t.SetHeader(`Write one commit using semantic commit specification. \n\n` + CV100)
			if extraPrompt != "" {
				t = t.AddUserMessage(extraPrompt)
			}
			t = t.AddUserMessage(diff)

			return t.RequestChat().ErrorTzap(nil)
		}}
}

var semanticGitcommitCmd = &cobra.Command{
	Use:   "semantic:gitcommit [clarifying prompt]",
	Short: "Generate a git commit message using ChatGPT",
	Long:  `Prompts ChatGPT to generate a commit message and commits it to the current git repo. The generated commit message is based on the diff of the currently staged files.`,
	Run: func(cmd *cobra.Command, args []string) {
		tzapConnector := tzapconnect.WithConfig(config.Configuration{SupressLogs: true, OpenAIModel: modelMap[settings.Model]})
		err := tzap.HandlePanic(func() {
			tzap.NewWithConnector(tzapConnector).
				ApplyTemplate(gocode.DeserializedArguments("extraPrompt", args)).
				ApplyErrorTemplate(git.GitDiff(), func(et *tzap.ErrorTzap) error {
					return et.Err
				}).
				ApplyErrorTemplate(git.ValidateDiff(), func(et *tzap.ErrorTzap) error {
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
				ApplyErrorTemplate(git.GitCommit(), func(et *tzap.ErrorTzap) error {
					return et.Err
				})
		})
		println(err)
	},
}

func init() {
	rootCmd.AddCommand(semanticGitcommitCmd)
}
