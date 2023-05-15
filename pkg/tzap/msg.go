package tzap

import (
	"fmt"
	"regexp"

	"github.com/tzapio/tzap/pkg/types"
)

func GetMessages(t *Tzap) []types.Message {
	messages := getMessages(t)
	if t.InitialSystemContent != "" {
		messages = append([]types.Message{{
			Role:    "system",
			Content: t.InitialSystemContent,
		}}, messages...)
	}
	return messages
}
func getMessages(t *Tzap) []types.Message {
	var messages []types.Message

	if t.Parent != nil {
		messages = GetMessages(t.Parent)
	}

	if t.Message.Content == "" || t.Message.Role == "" {
		return messages
	}
	key, ok := t.Data["memory"].(string)
	if ok && key != "" {
		mV := Mem[key]
		if mV.Content != "" {
			message := types.Message{
				Role:    mV.Role,
				Content: mV.Content,
			}
			messages = append(messages, message)
		}
	}
	return append(messages, t.Message)
}

// Not accurate way of counting tokens, but an approximation.
func TruncateToMaxWords(messages []types.Message, wordLimit int) []types.Message {
	var result []types.Message
	wordCount := 0
	if wordLimit < 0 {
		panic(fmt.Sprintf("TruncateToMaxWords wordlimit is %d, set above 1, or 0 to allow unlimited until model fails", wordLimit))
	}
	if wordLimit == 0 {
		return messages
	}

	wordSplitter := regexp.MustCompile(`[\s,./()\-_]+`)
	for i := len(messages) - 1; i >= 0; i-- {
		message := messages[i]
		words := wordSplitter.Split(message.Content, -1)
		words = removeEmptyStrings(words)
		if wordCount+len(words) <= wordLimit {

			wordCount += len(words)
			result = append([]types.Message{message}, result...)
		} else {
			break
		}
	}

	return result
}

func removeEmptyStrings(strings []string) []string {
	var result []string
	for _, s := range strings {
		if s != "" {
			result = append(result, s)
		}
	}
	return result
}
