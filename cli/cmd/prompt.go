package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/tzapio/tzap/cli/cmd/cmdui"
	"github.com/tzapio/tzap/cli/cmd/cmdutil"
	"github.com/tzapio/tzap/internal/logging/tl"
	"github.com/tzapio/tzap/pkg/tzapaction/action"
	"github.com/tzapio/tzap/pkg/tzapaction/actionpb"

	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/types/openai"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/tzapfile"
)

var inspirationFiles []string
var embedsCountFlag int32
var promptFile string
var searchQuery string

func init() {
	RootCmd.AddCommand(promptCmd)
	promptCmd.Flags().StringSliceVarP(&inspirationFiles,
		"inspiration", "i", []string{}, "Comma-separated list of inspiration files or multiple -i flags.")
	promptCmd.Flags().Int32VarP(&embedsCountFlag, "embeds", "k", 30,
		"Number of embeddings to use for the prompt generation")
	promptCmd.Flags().StringVarP(&promptFile, "promptfile", "f", "", "Read from file instead of prompt")
	promptCmd.Flags().StringVarP(&lib, "lib", "l", "", "BETA: select library to search.")
}

var promptCmd = &cobra.Command{
	Aliases: []string{"p", "embeddingprompt"},
	Use:     "prompt <prompt>",
	Short:   "Generate code by combining prompt and code-search",
	Long: `The 'prompt' command generates content based on code-searching existing files.
	This enables GPT to be able to generate code with depth. To add breadth, the user can recommend 
	needed Inspiration files like interfaces and types to enhance GPTs general understanding.
	The inspiration files should be a comma-separated list of file paths.`,
	Run: promptFunc,
}

func promptFunc(cmd *cobra.Command, args []string) {
	tl.EnableUICompletionLogger()
	t := cmdutil.GetTzapFromContext(cmd.Context())

	embedsCount := embedsCountFlag

	if promptFile == "-" {
		promptFile = ""
	}
	cmdUI := cmdui.NewCMDUI(promptFile, tzapCliSettings.Editor)
	messageThread := cmdui.NewMessageThread()
	if tzapCliSettings.ApiMode {
		tzapCliSettings.Editor = "api"
	}

	if promptFile != "" {
		messageThread.SetMessages(cmdUI.ReadMessagesFromFile())
	}

	if len(args) > 0 {
		userMessage := types.Message{
			Content: strings.Join(args[0:], " "),
			Role:    openai.ChatMessageRoleUser,
		}
		messageThread.Append(userMessage)
	}
	cmdUI.Init()
	err := tzap.HandlePanic(func() {
		defer t.HandleShutdown()

		for {
			if !messageThread.IsLastMessageFromUser() {
				messageThread.SetMessages(
					cmdUI.AddPromptTextWithStdinUI(
						messageThread.GetMessages(),
					),
				)
				continue
			}
			searchQuery = messageThread.LastMessage().Content
			truncThread := tzap.TruncateToMaxTokens(t.TG, messageThread.GetMessages(), 4000)

			promptWorkflowArgs := &actionpb.PromptArgs{
				InspirationFiles: inspirationFiles,
				SearchQuery:      searchQuery,
				EmbedsCount:      embedsCount,
				Thread:           action.ToPBMessage(truncThread),
			}

			cmd.Println(cmdutil.Bold("\nSearch query: "), cmdutil.Yellow(searchQuery))
			t.WorkTzap(func(t *tzap.Tzap) {
				t = t.
					ApplyWorkflow(action.PromptWorkflow(promptWorkflowArgs)).
					AsAssistantMessage()
				messageThread.Append(t.Message)

				if tzapCliSettings.ApiMode {
					cmdUI.SaveMessageThreadToFile(messageThread.GetMessages())
					fmt.Print(t.Message.Content)
					os.Exit(0)
					threadText, err := tzapfile.SerializeMessageThread(messageThread.GetMessages())
					if err != nil {
						panic(err)
					}
					fmt.Print(threadText)
					os.Exit(0)
					return
				}
			})
		}
	})
	if err != nil {
		panic(err)
	}

}
