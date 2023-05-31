package stubconnector

import (
	"context"

	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/types"
)

type StubConnector struct {
}

func StubWithConfig(conf config.Configuration) types.TzapConnector {
	tg := StubConnector{}

	return func() (types.TGenerator, config.Configuration) {
		return tg, conf
	}
}
func (StubConnector) TextToSpeech(ctx context.Context, content, language, voice string) (*[]byte, error) {
	return &[]byte{}, nil
}
func (StubConnector) SpeechToText(ctx context.Context, audioContent *[]byte, language string) (string, error) {
	return "nil", nil
}
func (StubConnector) GenerateChat(ctx context.Context, messages []types.Message, stream bool) (string, error) {
	return "Hello world", nil
}
func (StubConnector) CountTokens(ctx context.Context, content string) (int, error) {
	return len("Hello world"), nil
}
func (StubConnector) OffsetTokens(ctx context.Context, content string, from int, to int) (string, int, error) {
	return "Hell", 0, nil
}
func (StubConnector) RawTokens(ctx context.Context, content string) ([]string, error) {
	return []string{}, nil
}
func (StubConnector) FetchEmbedding(ctx context.Context, content ...string) ([][1536]float32, error) {
	return [][1536]float32{{0, 1, 2, 3, 4, 5}}, nil
}
func (StubConnector) SearchWithEmbedding(ctx context.Context, embedding types.QueryFilter, k int) (types.SearchResults, error) {
	return types.SearchResults{}, nil
}
func (StubConnector) AddEmbeddingDocument(ctx context.Context, docID string, embedding [1536]float32, metadata map[string]string) error {
	return nil
}
func (StubConnector) GetEmbeddingDocument(ctx context.Context, docID string) (types.Vector, bool, error) {
	return types.Vector{}, false, nil
}
func (StubConnector) DeleteEmbeddingDocument(ctx context.Context, docID string) error {
	return nil
}
func (StubConnector) DeleteEmbeddingDocuments(ctx context.Context, ids []string) error {
	return nil
}
func (StubConnector) ListAllEmbeddingsIds(ctx context.Context) (types.SearchResults, error) {
	return types.SearchResults{}, nil
}
