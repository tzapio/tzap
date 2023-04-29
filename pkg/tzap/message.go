package tzap

import (
	"fmt"
	"os"
	"strings"

	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/types/openai"
	"github.com/tzapio/tzap/pkg/util/stdin"
)

// AddUserMessage adds a user message to the Tzap
func (t *Tzap) AddUserMessage(contents ...string) *Tzap {
	content := strings.Join(contents, "\n")
	return t.AddTzap(&Tzap{
		Name: "AddUserMessage",
		Message: types.Message{
			Role:    openai.ChatMessageRoleUser,
			Content: content,
		},
	})
}

// LoadUserMessageFromFileOrStdinInput adds a user message from a file or standard input
func (t *Tzap) LoadUserMessageFromFileOrStdinInput(filepath string, task string) *Tzap {
	var content string
	if _, err := os.Stat(filepath); err != nil {
		content = stdin.GetStdinInput(task)
		if err := os.WriteFile(filepath, []byte(content), 0666); err != nil {
			panic(fmt.Errorf("could not write to from %s: %w", filepath, err))
		}
	} else {
		byteContent, err := os.ReadFile(filepath)
		if err != nil {
			panic(fmt.Errorf("could not read from %s: %w", filepath, err))
		}
		content = string(byteContent)
	}
	return t.AddTzap(&Tzap{Name: "LoadTaskFromFileOrCreateTask", Message: types.Message{
		Role:    openai.ChatMessageRoleUser,
		Content: content,
	}})
}

// AddSystemMessage adds a system message to the Tzap
func (t *Tzap) AddSystemMessage(contents ...string) *Tzap {
	content := strings.Join(contents, "\n")
	return t.AddTzap(&Tzap{
		Name: "AddSystemMessage",
		Message: types.Message{
			Role:    openai.ChatMessageRoleSystem,
			Content: content,
		}})
}

// AddAssistantMessage adds an assistant message to the Tzap
func (t *Tzap) AddAssistantMessage(contents ...string) *Tzap {
	content := strings.Join(contents, "\n")

	return t.AddTzap(&Tzap{
		Name: "AddAssistantMessage",
		Message: types.Message{
			Role:    openai.ChatMessageRoleAssistant,
			Content: content,
		}})
}

// AppendMessage appends a message to the current message in the Tzap
func (t *Tzap) AppendMessage(content string) *Tzap {
	lastMessage := t.Message
	newContent := fmt.Sprintf("%s %s", lastMessage.Content, content)
	newMessage := types.Message{
		Role:    lastMessage.Role,
		Content: newContent,
	}

	return t.HijackTzap(&Tzap{
		Message: newMessage,
	})
}

// AppendContent appends content to the current message in the Tzap
func (t *Tzap) AppendContent(sep string, s ...string) *Tzap {
	if t.Message.Content == "" {
		t.Message.Content = strings.Join(s, sep)
		return t
	}
	t.Message.Content = fmt.Sprintf("%s%s%s", t.Message.Content, sep, strings.Join(s, sep))
	return t
}

// PrependContent prepends content to the current message in the Tzap
func (t *Tzap) PrependContent(sep string, s ...string) *Tzap {
	if t.Message.Content == "" {
		t.Message.Content = strings.Join(s, sep)
		return t
	}
	t.Message.Content = fmt.Sprintf("%s%s%s", strings.Join(s, sep), sep, t.Message.Content)
	return t
}

// CombineMessage combines two message functions and creates a new message in the Tzap
func (t *Tzap) CombineMessage(nt1 func(*Tzap) *Tzap, nt2 func(*Tzap) *Tzap) *Tzap {
	types1 := nt1(t).Message
	types2 := nt2(t).Message
	return t.AddTzap(&Tzap{Name: "CombineMessage", Message: types.Message{Role: types1.Role, Content: fmt.Sprintf("%s\n%s", types1.Content, types2.Content)}})
}

func (t *Tzap) SetHeader(header string) *Tzap {
	t.Header = header
	return t
}
