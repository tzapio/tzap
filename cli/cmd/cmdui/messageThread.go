package cmdui

import (
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/types/openai"
)

type MessageThread struct {
	messageThread []types.Message
}

func NewMessageThread() *MessageThread {
	return &MessageThread{}
}
func (m *MessageThread) Append(message types.Message) {
	m.messageThread = append(m.messageThread, message)
}
func (m *MessageThread) IsLastMessageFromUser() bool {
	if len(m.messageThread) > 0 {
		lastMessage := m.LastMessage()
		if lastMessage.Role == openai.ChatMessageRoleUser && lastMessage.Content != "" {
			return true
		}
	}
	return false
}
func (m *MessageThread) LastMessage() *types.Message {
	if len(m.messageThread) > 0 {
		lastMessage := m.messageThread[len(m.messageThread)-1]
		if lastMessage.Content != "" {
			return &lastMessage
		}
	}
	return nil
}
func (m *MessageThread) SetMessages(messageThread []types.Message) {
	m.messageThread = messageThread
}
func (m *MessageThread) GetMessages() []types.Message {
	return m.messageThread

}
