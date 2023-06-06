package types

import (
	"context"

	"github.com/tzapio/tzap/pkg/config"
)

type TGenerator interface {
	TextToSpeech(ctx context.Context, content, language, voice string) (*[]byte, error)
	SpeechToText(ctx context.Context, audioContent *[]byte, language string) (string, error)
	FetchEmbedding(ctx context.Context, content ...string) ([][1536]float32, error)
	AddEmbeddingDocument(ctx context.Context, project string, id string, embedding [1536]float32, metadata Metadata) error
	GetEmbeddingDocument(ctx context.Context, project string, id string) (Vector, bool, error)
	DeleteEmbeddingDocument(ctx context.Context, project string, id string) error
	DeleteEmbeddingDocuments(ctx context.Context, project string, ids []string) error
	SearchWithEmbedding(ctx context.Context, project string, embedding QueryFilter, k int) (SearchResults, error)
	ListAllEmbeddingsIds(ctx context.Context, project string) (SearchResults, error)
	GenerateChat(ctx context.Context, messages []Message, stream bool) (string, error)
	CountTokens(ctx context.Context, content string) (int, error)
	OffsetTokens(ctx context.Context, content string, from int, to int) (string, int, error)
	RawTokens(ctx context.Context, content string) ([]string, error)
}
type TzapConnector func() (TGenerator, config.Configuration)
type ProjectDir string
type ProjectName string

type ProjectDB map[ProjectName]ProjectDir

const LOCALPROJECTNAME ProjectName = "@LOCAL"
