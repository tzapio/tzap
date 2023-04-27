package types

import (
	"context"

	"github.com/tzapio/tzap/pkg/config"
)

type TGenerator interface {
	TextToSpeech(ctx context.Context, content, language, voice string) (*[]byte, error)
	SpeechToText(ctx context.Context, audioContent *[]byte, language string) (string, error)
	GenerateChat(ctx context.Context, messages []Message) (string, error)
	CountTokens(ctx context.Context, content string) (int, error)
	OffsetTokens(ctx context.Context, content string, from int, to int) (string, error)
}
type UnimplementedTGenerator struct {
	TGenerator
}

func (UnimplementedTGenerator) TextToSpeech(ctx context.Context, content, language, voice string) (*[]byte, error) {
	panic("UnimplementedTGenerator - Probably not supported yet.")
}
func (UnimplementedTGenerator) SpeechToText(ctx context.Context, audioContent *[]byte, language string) (string, error) {
	panic("UnimplementedTGenerator - Probably not supported yet.")
}
func (UnimplementedTGenerator) GenerateChat(ctx context.Context, messages []Message) (string, error) {
	panic("UnimplementedTGenerator - Probably not supported yet.")
}

// TODO
type UnimplementedTGenerator2 struct {
	TextToSpeech
	SpeechToText
	GenerateChat
}
type TzapConnector func() (TGenerator, config.Configuration)
type TextToSpeech func(ctx context.Context, content, language, voice string) (*[]byte, error)
type SpeechToText func(ctx context.Context, audioContent *[]byte, language string) (string, error)
type GenerateChat func(ctx context.Context, messages []Message) (string, error)
