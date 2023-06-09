package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/tzapio/tzap/cli/action"
	"github.com/tzapio/tzap/cli/cmd/cliworkflows"
	"github.com/tzapio/tzap/cli/cmd/cmdui"
	"github.com/tzapio/tzap/cli/cmd/cmdutil"
	"github.com/tzapio/tzap/internal/logging/tl"
	"github.com/tzapio/tzap/workflows/code/fileworkflows"

	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/types/openai"
	"github.com/tzapio/tzap/pkg/tzap"
)

var inspirationFiles []string
var embedsCountFlag int
var nCountFlag int
var promptFile string
var disableIndex bool
var searchQuery string

func init() {
	RootCmd.AddCommand(promptCmd)
	promptCmd.Flags().StringSliceVarP(&inspirationFiles,
		"inspiration", "i", []string{}, "Comma-separated list of inspiration files or multiple -i flags.")
	promptCmd.Flags().IntVarP(&embedsCountFlag, "embeds", "k", 10,
		"Number of embeddings to use for the prompt generation")
	promptCmd.Flags().IntVarP(&nCountFlag, "searchsize", "n", 15,
		"Number of embeddings to include in the search space before filtering out the matches with inspiration files.")
	promptCmd.Flags().BoolVarP(&disableIndex, "disableindex", "d", false,
		"For large projects disabling indexing speeds up the process.")
	promptCmd.Flags().StringVarP(&promptFile, "promptfile", "f", "", "Read from file instead of prompt")
	promptCmd.Flags().StringVarP(&lib, "lib", "l", "", "BETA: select library to search.")
}

var promptCmd = &cobra.Command{
	Aliases: []string{"p", "prompt"},
	Use:     "prompt <prompt>",
	Short:   "Generate code by combining prompt and code-search",
	Long: `The 'prompt' command generates content based on code-searching existing files.
	This enables GPT to be able to generate code with depth. To add breadth, the user can recommend 
	needed Inspiration files like interfaces and types to enhance GPTs general understanding.
	The inspiration files should be a comma-separated list of file paths.`,
	Run: func(cmd *cobra.Command, args []string) {
		tl.EnableUICompletionLogger()
		t := cmdutil.GetTzapFromContext(cmd.Context())

		embedsCount := embedsCountFlag
		nCount := nCountFlag
		if embedsCountFlag > nCountFlag {
			nCount = embedsCountFlag + 5
		}

		err := tzap.HandlePanic(func() {
			defer t.HandleShutdown()
			if tzapCliSettings.ApiMode {
				tzapCliSettings.Editor = "api"
			}
			cmdUI := cmdui.NewCMDUI(promptFile, tzapCliSettings.Editor)
			messageThread := cmdui.MessageThread{}

			if promptFile != "" {
				messageThread.SetMessageThread(cmdUI.ReadMessageThreadFromFile())
			}
			if len(args) > 0 {
				userMessage := types.Message{
					Content: strings.Join(args[0:], " "),
					Role:    openai.ChatMessageRoleUser,
				}
				messageThread.Append(userMessage)
			}
			cmdUI.Init()
			for {
				if !messageThread.IsLastMessageFromUser() {
					messageThread.SetMessageThread(
						cmdUI.AddPromptTextWithStdinUI(
							messageThread.GetMessageThread(),
						),
					)
					continue
				}
				searchQuery = messageThread.LastMessage().Content

				cmd.Println(cmdutil.Bold("\nSearch query: "), cmdutil.Yellow(searchQuery))

				output := action.LoadAndSearchEmbeddings(t, action.LoadAndSearchEmbeddingsArgs{
					ExcludeFiles: inspirationFiles,
					SearchQuery:  searchQuery,
					K:            embedsCount,
					N:            nCount,
					DisableIndex: disableIndex,
					Yes:          tzapCliSettings.Yes,
				})

				t.
					ApplyWorkflow(cliworkflows.PrintInspirationFiles(inspirationFiles)).
					ApplyWorkflow(cliworkflows.PrintSearchResults(output.SearchResults)).
					ApplyWorkflow(fileworkflows.InspirationWorkflow(inspirationFiles)).
					MutationTzap(func(t *tzap.Tzap) *tzap.Tzap {
						searchResults := output.SearchResults
						if len(searchResults.Results) > 0 {
							t = t.AddSystemMessage(
								"The following file contents are embeddings for the user input:",
							)
							for _, result := range searchResults.Results {
								t = t.AddSystemMessage(result.Vector.Metadata.SplitPart)
							}
						}
						return t
					}).
					MutationTzap(func(t *tzap.Tzap) *tzap.Tzap {
						cmd.Println(cmdutil.Bold("--- Completion"))
						truncThread := tzap.TruncateToMaxTokens(t.TG, messageThread.GetMessageThread(), 1000)
						t = t.LoadThread(truncThread).RequestChatCompletion().AsAssistantMessage()
						cmd.Println(cmdutil.Bold("\n---"))
						messageThread.Append(t.Message)

						if tzapCliSettings.ApiMode {
							cmdUI.SaveMessageThreadToFile(messageThread.GetMessageThread())
							threadText, err := cmdUI.SerializeMessageThread(messageThread.GetMessageThread())
							if err != nil {
								panic(err)
							}
							fmt.Print(threadText)
							os.Exit(0)
							return t
						}

						return t
					})
			}
		})
		if err != nil {
			panic(err)
		}

	},
}