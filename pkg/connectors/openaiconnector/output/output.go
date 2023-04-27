package output

import (
	"github.com/sashabaranov/go-openai"
	"github.com/tzapio/tzap/pkg/types"
)

func GetOpenAICompletionMessage(messages []types.Message) []openai.ChatCompletionMessage {
	requestMessages := make([]openai.ChatCompletionMessage, len(messages))
	for i, message := range messages {
		requestMessages[i] = openai.ChatCompletionMessage{
			Role:    message.Role,
			Content: message.Content,
		}
	}
	return requestMessages
}
