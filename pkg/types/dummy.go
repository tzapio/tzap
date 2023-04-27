package types

import "context"

type UnimplementedTGenerator struct {
	TGenerator
}

func (UnimplementedTGenerator) TextToSpeech(ctx context.Context, content, language, voice string) (*[]byte, error) {
	panic("UnimplementedTGenerator - Probably not supported yet.")
}
func (UnimplementedTGenerator) SpeechToText(ctx context.Context, audioContent *[]byte, language string) (string, error) {
	panic("UnimplementedTGenerator - Probably not supported yet.")
}
func (UnimplementedTGenerator) GenerateChat(ctx context.Context, messages []Message, stream bool) (string, error) {
	panic("UnimplementedTGenerator - Probably not supported yet.")
}
