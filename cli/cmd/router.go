package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tzapio/tzap/cli/cmd/cmdui"
	"github.com/tzapio/tzap/cli/cmd/cmdutil"
	"github.com/tzapio/tzap/cli/cmd/resolver"
	"github.com/tzapio/tzap/internal/logging/tl"
	"github.com/tzapio/tzap/pkg/tzapaction/action"
	"github.com/tzapio/tzap/pkg/tzapaction/actionpb"

	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/types/openai"
	"github.com/tzapio/tzap/pkg/tzap"
)

func init() {
	RootCmd.AddCommand(routerCmd)
}

var routerCmd = &cobra.Command{
	Aliases: []string{"r"},
	Use:     "router <prompt>",
	Short:   "Generate code by combining prompt and code-search",
	Long: `The 'router' command generates content based on code-searching existing files.
	This enables GPT to be able to generate code with depth. To add breadth, the user can recommend 
	needed Inspiration files like interfaces and types to enhance GPTs general understanding.
	The inspiration files should be a comma-separated list of file paths.`,
	Run: routerFunc,
}

func routerFunc(cmd *cobra.Command, args []string) {
	tl.EnableUICompletionLogger()
	t := cmdutil.GetTzapFromContext(cmd.Context())

	embedsCount := embedsCountFlag

	if promptFile == "-" {
		promptFile = ""
	}
	if tzapCliSettings.ApiMode {
		tzapCliSettings.Editor = "api"
	}
	cmdUI := cmdui.NewCMDUI(promptFile, tzapCliSettings.Editor)
	messageThread := cmdui.NewMessageThread()
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
				SearchArgss: []*actionpb.SearchArgs{
					{
						SearchQuery: searchQuery,
						Lib:         lib,
						EmbedsCount: embedsCount,
					},
				},
				Thread: action.ToPBMessage(truncThread),
			}

			println(cmdutil.Bold("\nSearch query: "), cmdutil.Yellow(searchQuery))
			t.WorkTzap(func(t *tzap.Tzap) {
				t = t.
					ApplyWorkflow(action.RouterWorkflow(promptWorkflowArgs)).
					IfFunctionCall(
						func(tzapFunc *tzap.Tzap) *tzap.Tzap {
							fc := tzapFunc.Data["content"].(types.CompletionMessage).FunctionCall
							if fc != nil {
								tzapFunc.Message.Role = "function"
								tzapFunc.Message.Content = fc.Name + "\n" + fc.Arguments
								messageThread.Append(types.Message{Role: "function", Content: fc.Name + "\n" + fc.Arguments})
								println("---")
								result, err := resolver.LocalRun("/"+fc.Name, fc.Arguments)
								if err != nil {
									panic(err)
								}
								println("---")
								if tzapCliSettings.ApiMode {
									fmt.Print(result)
									os.Exit(0)
									return tzapFunc
								}
								type FileWriter struct {
									FileWrites []*actionpb.FileWrite `json:"fileWrites"`
								}
								var fWriter *FileWriter = &FileWriter{}

								if err := json.Unmarshal([]byte(result), fWriter); err == nil {
									for _, fileWrite := range fWriter.FileWrites {
										cmdUI.EditFile(fileWrite)
									}
								}
							}
							os.Exit(0)
							return tzapFunc
						},
						func(notTzapFunc *tzap.Tzap) *tzap.Tzap {
							notTzapFunc = notTzapFunc.AsAssistantMessage()
							messageThread.Append(notTzapFunc.Message)
							cmdUI.SaveMessageThreadToFile(messageThread.GetMessages())
							return notTzapFunc
						})

				if tzapCliSettings.ApiMode {
					fmt.Print(t.Message.Content)
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
