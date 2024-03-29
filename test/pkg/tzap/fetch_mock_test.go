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
func (tg *mockTG) GenerateChat(ctx context.Context, messages []types.Message, stream bool, functions string) (types.CompletionMessage, error) {
	if len(messages) == 0 {
		return types.CompletionMessage{}, fmt.Errorf("empty messages")
	}
	str := ""
	for i, message := range messages {
		str += fmt.Sprintf("r=%s;c=%s", message.Role, message.Content)
		if i != len(messages)-1 {
			str += "|"
		}
	}
	return types.CompletionMessage{Content: str}, nil
}
func (tg *mockTG) CountTokens(ctx context.Context, content string) (int, error) {
	// Return pre-defined value for testing purposes
	return 50, nil
}
func (tg *mockTG) OffsetTokens(ctx context.Context, content string, from int, to int) (string, int, error) {
	// Return pre-defined value for testing purposes
	return "Hello world!", 0, nil
}
func (tg *mockTG) RawTokens(ctx context.Context, content string) ([]string, error) {
	// Return pre-defined value for testing purposes
	return []string{}, nil
}
func (tg *mockTG) FetchEmbedding(ctx context.Context, content ...string) ([][1536]float32, error) {
	return [][1536]float32{}, nil
}
func (tg *mockTG) SearchWithEmbedding(ctx context.Context, embedding types.QueryFilter, k int) (types.SearchResults, error) {
	return types.SearchResults{}, nil
}
func (tg *mockTG) AddEmbeddingDocument(ctx context.Context, docID string, embedding [1536]float32, metadata types.Metadata) error {
	return nil
}
func (tg *mockTG) GetEmbeddingDocument(ctx context.Context, docID string) (types.Vector, bool, error) {
	return types.Vector{}, false, nil
}
func (tg *mockTG) DeleteEmbeddingDocument(ctx context.Context, docID string) error {
	return nil
}
func (tg *mockTG) DeleteEmbeddingDocuments(ctx context.Context, ids []string) error {
	return nil
}
func (tg *mockTG) ListAllEmbeddingsIds(ctx context.Context) (types.SearchResults, error) {
	return types.SearchResults{}, nil
}
