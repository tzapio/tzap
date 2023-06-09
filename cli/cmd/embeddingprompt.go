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
	RootCmd.AddCommand(embeddingPromptCmd)
	embeddingPromptCmd.Flags().StringSliceVarP(&inspirationFiles,
		"inspiration", "i", []string{}, "Comma-separated list of inspiration files or multiple -i flags.")
	embeddingPromptCmd.Flags().IntVarP(&embedsCountFlag, "embeds", "k", 10,
		"Number of embeddings to use for the prompt generation")
	embeddingPromptCmd.Flags().IntVarP(&nCountFlag, "searchsize", "n", 15,
		"Number of embeddings to include in the search space before filtering out the matches with inspiration files.")
	embeddingPromptCmd.Flags().BoolVarP(&disableIndex, "disableindex", "d", false,
		"For large projects disabling indexing speeds up the process.")
	embeddingPromptCmd.Flags().StringVarP(&promptFile, "promptfile", "f", "", "Read from file instead of prompt")
	embeddingPromptCmd.Flags().StringVarP(&lib, "lib", "l", "", "BETA: select library to search.")
}

var embeddingPromptCmd = &cobra.Command{
	Aliases: []string{"p", "prompt"},
	Use:     "embeddingprompt <prompt>",
	Short:   "Generate code or document content using code-search",
	Long: `The 'embeddingprompt' command generates content based on code-searching existing files.
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

			cmdUI := cmdui.NewCMDUI(promptFile, settings.Editor)
			thread := []types.Message{}
			var lastUserMessage *types.Message
			if promptFile != "" {
				thread = cmdUI.DeserializeThread()
				if len(thread) > 0 {
					lastMessage := thread[len(thread)-1]
					if lastMessage.Role == openai.ChatMessageRoleUser {
						lastUserMessage = &lastMessage
					}
				}
				if len(args) > 0 {
					content := strings.Join(args[0:], " ")
					userMessage := types.Message{
						Content: content,
						Role:    openai.ChatMessageRoleUser,
					}
					thread = append(thread, userMessage)
					lastUserMessage = &userMessage
				}
			} else {
				if len(args) > 0 {
					content := strings.Join(args[0:], " ")
					userMessage := types.Message{
						Content: content,
						Role:    openai.ChatMessageRoleUser,
					}
					thread = append(thread, userMessage)
					lastUserMessage = &userMessage
				}
			}
			cmdUI.Init()
			for {
				if lastUserMessage == nil {
					thread = cmdUI.AddPromptTextWithStdinUI(thread)
					currentMessage := &thread[len(thread)-1]
					if currentMessage.Content == "" || currentMessage.Role != openai.ChatMessageRoleUser {
						continue
					}
					lastUserMessage = currentMessage
				}
				searchQuery = lastUserMessage.Content
				cmd.Println(cmdutil.Bold("\nSearch query: "), cmdutil.Yellow(searchQuery))

				output := action.LoadAndSearchEmbeddings(t, action.LoadAndSearchEmbeddingsArgs{
					ExcludeFiles: inspirationFiles,
					SearchQuery:  searchQuery,
					K:            embedsCount,
					N:            nCount,
					DisableIndex: disableIndex,
					Yes:          settings.Yes,
				})

				t.
					ApplyWorkflow(cliworkflows.PrintInspirationFiles(inspirationFiles)).
					ApplyWorkflow(cliworkflows.PrintSearchResults(output.SearchResults)).
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
						truncThread := tzap.TruncateToMaxTokens(t.TG, thread, 1000)
						t = t.LoadThread(truncThread).RequestChatCompletion().AsAssistantMessage()
						cmd.Println(cmdutil.Bold("\n---"))
						thread = append(thread, types.Message{
							Content: t.Message.Content,
							Role:    t.Message.Role,
						})
						if settings.ApiMode {
							cmdUI.SaveThread(thread)
							threadText, err := cmdUI.SerializeThread(thread)
							if err != nil {
								panic(err)
							}
							fmt.Print(threadText)
							os.Exit(0)
							return t
						}
						lastUserMessage = nil
						return t
					})
			}
		})
		if err != nil {
			panic(err)
		}

	},
}
