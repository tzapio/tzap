package tzap_test

import (
	"context"
	"fmt"

	"github.com/tzapio/tzap/pkg/types"
)

type mockTG struct{}

func (tg *mockTG) TextToSpeech(ctx context.Context, text string, language string, voice string) (*[]byte, error) {
	// Return pre-defined value for testing purposes
	r := []byte("sample audio content")
	return &r, nil
}

func (tg *mockTG) SpeechToText(ctx context.Context, audioContent *[]byte, language string) (string, error) {
	// Return pre-defined value for testing purposes
	return "Hello world!", nil
}
func (tg *mockTG) GenerateChat(ctx context.Context, messages []types.Message) (string, error) {

	if len(messages) == 0 {
		return "", fmt.Errorf("empty messages")
	}
	str := ""
	for i, message := range messages {
		str += fmt.Sprintf("r=%s;c=%s", message.Role, message.Content)
		if i != len(messages)-1 {
			str += "|"
		}
	}
	return str, nil
}
func (tg *mockTG) CountTokens(ctx context.Context, content string) (int, error) {
	// Return pre-defined value for testing purposes
	return 50, nil
}
func (tg *mockTG) OffsetTokens(ctx context.Context, content string, from int, to int) (string, error) {
	// Return pre-defined value for testing purposes
	return "Hello world!", nil

}
