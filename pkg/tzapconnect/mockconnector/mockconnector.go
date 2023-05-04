package mockconnector

import (
	"context"

	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/types"
)

type MockConnector struct {
}

func MockWithConfig(conf config.Configuration) types.TzapConnector {
	tg := MockConnector{}

	return func() (types.TGenerator, config.Configuration) {
		return tg, conf
	}
}
func (MockConnector) TextToSpeech(ctx context.Context, content, language, voice string) (*[]byte, error) {
	return &[]byte{}, nil
}
func (MockConnector) SpeechToText(ctx context.Context, audioContent *[]byte, language string) (string, error) {
	return "nil", nil
}
func (MockConnector) GenerateChat(ctx context.Context, messages []types.Message, stream bool) (string, error) {
	return "Hello world", nil
}
func (MockConnector) CountTokens(ctx context.Context, content string) (int, error) {
	return len("Hello world"), nil
}
func (MockConnector) OffsetTokens(ctx context.Context, content string, from int, to int) (string, error) {
	return "Hell", nil
}
