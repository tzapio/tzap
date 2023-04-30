package tzapconnect

import (
	"context"
	"fmt"
	"os"

	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/connectors/openaiconnector"
	"github.com/tzapio/tzap/pkg/types"
)

func WithConfig(conf config.Configuration) types.TzapConnector {
	tg, err := newBaseconnector()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return func() (types.TGenerator, config.Configuration) {
		return tg, conf
	}
}

func newBaseconnector() (types.TGenerator, error) {
	apiKey, err := getOpenAIAPIKeyFromEnv()
	if err != nil {
		return nil, err
	}
	openaiC := openaiconnector.InitiateOpenaiClient(apiKey)

	partialComposite := PartialComposite{OpenaiTgenerator: openaiC}
	var myInterface types.TGenerator = partialComposite
	return myInterface, nil
}

type PartialComposite struct {
	types.UnimplementedTGenerator
	OpenaiTgenerator *openaiconnector.OpenaiTgenerator
	VoiceGenerator   types.TGenerator
}

func (pc PartialComposite) TextToSpeech(ctx context.Context, content, language, voice string) (*[]byte, error) {
	return pc.VoiceGenerator.TextToSpeech(ctx, content, language, voice)
}
func (pc PartialComposite) GenerateChat(ctx context.Context, messages []types.Message, stream bool) (string, error) {
	return pc.OpenaiTgenerator.GenerateChat(ctx, messages, stream)
}
func (pc PartialComposite) CountTokens(ctx context.Context, content string) (int, error) {
	return pc.OpenaiTgenerator.CountTokens(content)
}
func (pc PartialComposite) OffsetTokens(ctx context.Context, content string, from int, to int) (string, error) {
	return pc.OpenaiTgenerator.OffsetTokens(content, from, to)
}
func getOpenAIAPIKeyFromEnv() (string, error) {
	apiKey := os.Getenv("OPENAI_APIKEY")
	if apiKey == "" {
		return "", fmt.Errorf("OPENAI_APIKEY environment variable not set\n\n\t\texport OPENAI_APIKEY=<apikey>")
	}

	return apiKey, nil
}
