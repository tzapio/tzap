package tzapfile

import (
	"fmt"
	"strings"

	"github.com/tzapio/tzap/pkg/types"
)

func DeserializeMessageThread(content string) []types.Message {
	var messages []types.Message
	messageLines := strings.Split(content, "---")
	for _, messageLine := range messageLines {
		message := types.Message{}
		lines := strings.Split(strings.TrimSpace(messageLine), "\n")
		if len(lines) > 0 {
			if strings.HasPrefix(lines[0], "@role:") {
				message.Role = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(lines[0]), "@role:"))
				message.Content = strings.TrimSpace(strings.Join(lines[1:], "\n"))
			} else {
				message.Role = "user"
				message.Content = strings.TrimSpace(strings.Join(lines, "\n"))
			}
		}
		messages = append(messages, message)
	}
	reverseMessages := []types.Message{}
	for i := len(messages) - 1; i >= 0; i-- {
		if strings.TrimSpace(messages[i].Content) == "" || strings.TrimSpace(messages[i].Role) == "" {
			continue
		}
		reverseMessages = append(reverseMessages, messages[i])
	}
	return reverseMessages
}
func SerializeMessageThread(messages []types.Message) (string, error) {
	reversedMessages := []types.Message{}
	for i := len(messages) - 1; i >= 0; i-- {
		if strings.TrimSpace(messages[i].Content) == "" || strings.TrimSpace(messages[i].Role) == "" {
			continue
		}
		reversedMessages = append(reversedMessages, messages[i])
	}
	s := strings.Builder{}
	if len(reversedMessages) > 0 {
		s.WriteString("\n\n")
	}

	for _, msg := range reversedMessages {
		if _, err := s.WriteString(fmt.Sprintf("---\n@role:%s\n%s\n", msg.Role, msg.Content)); err != nil {
			return "", err
		}
	}
	return s.String(), nil
}
