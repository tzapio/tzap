package tzap

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/types/openai"
	"github.com/tzapio/tzap/pkg/util"
)

func GetThread(t *Tzap) []types.Message {
	messages := getThread(t)
	if t.InitialSystemContent != "" {
		messages = append([]types.Message{{
			Role:    "system",
			Content: t.InitialSystemContent,
		}}, messages...)
	}
	return messages
}
func getThread(t *Tzap) []types.Message {
	var messages []types.Message

	if t.Parent != nil {
		messages = GetThread(t.Parent)
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
func (t *Tzap) GetThreadAsJSON() (string, error) {
	messages := GetThread(t)
	jsonBytes, err := json.MarshalIndent(messages, "", "  ")
	if err != nil {
		return "", nil
	}
	return string(jsonBytes), err
}
func (t *Tzap) StoreThread(filePath string) *ErrorTzap {
	messages := GetThread(t)
	jsonBytes, err := json.Marshal(messages)
	if err != nil {
		panic(fmt.Errorf("error storing thread: %w", err))
	}

	if err := os.WriteFile(filePath, jsonBytes, 0644); err != nil {
		return t.ErrorTzap(fmt.Errorf("StoreThread: error storing thread: %w", err))
	}

	return t.ErrorTzap(nil)
}
func (t *Tzap) LoadThreadString(content string) *ErrorTzap {
	var messages []types.Message
	err := json.Unmarshal([]byte(content), &messages)
	if err != nil {
		return t.ErrorTzap(fmt.Errorf("error loading thread: %w", err))
	}
	return t.LoadThread(messages).ErrorTzap(nil)
}
func (t *Tzap) LoadThreadFile(filePath string) *ErrorTzap {
	return t.LoadThreadString(util.ReadFileP(filePath))
}

func (t *Tzap) LoadThread(messages []types.Message) *Tzap {
	for _, message := range messages {
		if message.Role == openai.ChatMessageRoleSystem {
			t = t.AddSystemMessage(message.Content)
			continue
		}
		if message.Role == openai.ChatMessageRoleAssistant {
			t = t.AddAssistantMessage(message.Content)
			continue
		}
		if message.Role == openai.ChatMessageRoleUser {
			t = t.AddUserMessage(message.Content)
			continue
		}
	}
	return t
}
func (t *Tzap) storeThread(messages []types.Message) *Tzap {
	for _, message := range messages {
		if message.Role == openai.ChatMessageRoleSystem {
			t.AddSystemMessage(message.Content)
			continue
		}
		if message.Role == openai.ChatMessageRoleAssistant {
			t.AddAssistantMessage(message.Content)
			continue
		}
		if message.Role == openai.ChatMessageRoleUser {
			t.AddUserMessage(message.Content)
			continue
		}
	}
	return t
}

func TruncateToMaxTokens(tg types.TGenerator, messages []types.Message, wordLimit int) []types.Message {
	var result []types.Message
	tokenCount := 0
	if wordLimit < 0 {
		panic(fmt.Sprintf("TruncateToMaxWords wordlimit is %d, set above 1, or 0 to allow unlimited until model fails", wordLimit))
	}
	if wordLimit == 0 {
		return messages
	}

	for i := len(messages) - 1; i >= 0; i-- {
		message := messages[i]
		tokens, err := tg.CountTokens(context.Background(), message.Content)
		if err != nil {
			panic(fmt.Errorf("TruncateToMaxWords: error counting tokens: %w", err))
		}

		if tokenCount+tokens <= wordLimit {
			tokenCount += tokens
			result = append([]types.Message{message}, result...)
		} else {
			break
		}
	}

	return result
}
