package types

import (
	"context"

	"github.com/tzapio/tzap/pkg/config"
)

type TGenerator interface {
	TextToSpeech(ctx context.Context, content, language, voice string) (*[]byte, error)
	SpeechToText(ctx context.Context, audioContent *[]byte, language string) (string, error)
	GenerateChat(ctx context.Context, messages []Message, stream bool) (string, error)
	CountTokens(ctx context.Context, content string) (int, error)
	OffsetTokens(ctx context.Context, content string, from int, to int) (string, error)
}
type TzapConnector func() (TGenerator, config.Configuration)
