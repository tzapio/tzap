package tzap

import (
	"fmt"

	"github.com/tzapio/tzap/internal/filelog"
	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/util/stdin"
)

// FetchTask initializes the chat completion request, fetches the edited content, and sets up the environment.
func (t *Tzap) FetchTask() *Tzap {
	config := config.FromContext(t.C)
	if err := CheckData(t.Data); err != nil {
		panic(err)
	}

	// force print all recent logs.
	Flush()

	filepathValue := t.Data["filepath"].(string)
	fmt.Println("requesting edit file for (" + filepathValue + ")")

	editedContent, err := fetchChatResponse(t, true)
	if err != nil {
		panic(err)
	}

	fmt.Println(editedContent)
	t = applyChanges(t, editedContent)

	if config.AutoMode || stdin.ConfirmToContinue() {
		return t
	}
	panic("Do not continue selected")
}

// RequestOpenAIChat initializes the openai chat completion request and creates a new Tzap with the edited content.
func (t *Tzap) RequestChat() *Tzap {
	output, err := fetchChatResponse(t, true)
	if err != nil {
		panic(err)
	}
	data := types.MappedInterface{
		"content": output,
	}
	withRequestFile := t.AddTzap(&Tzap{Name: "requestChat", Data: data})
	return withRequestFile
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
	result, err := t.TG.GenerateChat(t.C, messages)

	if err != nil {
		filelog.LogData(t.C, err.Error(), filelog.ResponseLog)
		return "", err
	}
	filelog.LogData(t.C, result, filelog.ResponseLog)
	return result, nil
}
