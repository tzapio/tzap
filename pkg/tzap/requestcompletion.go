package tzap

import (
	"fmt"

	"github.com/tzapio/tzap/internal/logging/filelog"
	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/types/openai"
	"github.com/tzapio/tzap/pkg/util"
	"github.com/tzapio/tzap/pkg/util/stdin"
)

func (t *Tzap) StoreCompletion(filePath string) *Tzap {
	config := config.FromContext(t.C)
	editedContent := t.Data["content"].(string)

	autoMode := config.AutoMode
	makeChange := autoMode
	if !autoMode {
		if ok := stdin.ConfirmPrompt("Overwrite file: " + filePath); ok {
			makeChange = true
		}
	}
	if makeChange {
		err := util.MkdirPAndWriteFile(filePath, editedContent)
		if err != nil {
			panic(fmt.Errorf("error applying changes: %w", err))
		}
		writeMessageMD5(filePath, t)
		data := types.MappedInterface{
			"filepath": filePath,
			"content":  editedContent,
		}
		withEditFile := t.CloneTzap(&Tzap{Name: "storeCompletion", Message: types.Message{
			Role:    openai.ChatMessageRoleAssistant,
			Content: editedContent,
		}, Data: data})
		return withEditFile
	}

	if config.AutoMode || stdin.ConfirmPrompt("Continue on?") {
		return t
	}
	panic("Do not continue selected")
}

// RequestChatCompletion initializes the openai chat completion request and creates a new Tzap with the edited content.
func (t *Tzap) RequestChatCompletion() *Tzap {

	output, err := fetchChatResponse(t, true)
	if err != nil {
		panic(err)
	}
	data := types.MappedInterface{
		"content": output,
	}
	requestChat := t.AddTzap(&Tzap{Name: "requestChat", Data: data})

	return requestChat
}

// RequestOpenAIChat initializes the openai chat completion request and creates a new Tzap with the edited content.
func (t *Tzap) CountTokens(content string) (int, error) {
	return t.TG.CountTokens(t.C, content)
}

// RequestOpenAIChat initializes the openai chat completion request and creates a new Tzap with the edited content.
func (t *Tzap) OffsetTokens(content string, from int, to int) (string, error) {
	return t.TG.OffsetTokens(t.C, content, from, to)
}

// fetchChatResponse requests openai-chat completion for the given Tzap and returns the modified content.
func fetchChatResponse(t *Tzap, stream bool) (string, error) {
	config := config.FromContext(t.C)

	messages := TruncateToMaxWords(GetMessages(t), config.TruncateLimit)

	filelog.LogData(t.C, t, filelog.TzapLog)
	GenerateGraphvizDotFile(t, FillGraphVizGraph())
	filelog.LogData(t.C, messages, filelog.RequestLog)
	println("\n--- Completion:")
	result, err := t.TG.GenerateChat(t.C, messages, stream)

	if err != nil {
		filelog.LogData(t.C, err.Error(), filelog.ResponseLog)
		return "", err
	}
	println("\n---")
	getMessagesGraphViz(t)
	GenerateGraphvizDotFile(t, FillGraphVizGraph())
	filelog.LogData(t.C, result, filelog.ResponseLog)

	return result, nil
}
