package tzap

import (
	"fmt"
	"strings"

	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/types/openai"
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

func (t *Tzap) SetInitialSystemContent(content string) *Tzap {
	t.InitialSystemContent = content
	return t
}
