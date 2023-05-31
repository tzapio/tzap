package types

import (
	"context"

	"github.com/tzapio/tzap/pkg/config"
)

type TGenerator interface {
	TextToSpeech(ctx context.Context, content, language, voice string) (*[]byte, error)
	SpeechToText(ctx context.Context, audioContent *[]byte, language string) (string, error)
	FetchEmbedding(ctx context.Context, content ...string) ([][1536]float32, error)
	AddEmbeddingDocument(ctx context.Context, id string, embedding [1536]float32, metadata map[string]string) error
	GetEmbeddingDocument(ctx context.Context, id string) (Vector, bool, error)
	DeleteEmbeddingDocument(ctx context.Context, id string) error
	DeleteEmbeddingDocuments(ctx context.Context, ids []string) error
	SearchWithEmbedding(ctx context.Context, embedding QueryFilter, k int) (SearchResults, error)
	ListAllEmbeddingsIds(ctx context.Context) (SearchResults, error)
	GenerateChat(ctx context.Context, messages []Message, stream bool) (string, error)
	CountTokens(ctx context.Context, content string) (int, error)
	OffsetTokens(ctx context.Context, content string, from int, to int) (string, int, error)
	RawTokens(ctx context.Context, content string) ([]string, error)
}
type TzapConnector func() (TGenerator, config.Configuration)
