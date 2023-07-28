package tzapconnect

import (
	"context"
	"os"

	"github.com/tzapio/tzap/internal/logging/tl"
	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/embed/embedstore"

	"github.com/tzapio/tzap/pkg/connectors/openaiconnector"
	"github.com/tzapio/tzap/pkg/types"
)

func WithConfig(openai_apikey string, conf config.Configuration) types.TzapConnector {
	tg, err := newBaseconnector(openai_apikey, conf)
	if err != nil {
		println(err)
		os.Exit(1)
	}
	return func() (types.TGenerator, config.Configuration) {
		return tg, conf
	}
}

func newBaseconnector(openai_apikey string, conf config.Configuration) (types.TGenerator, error) {
	tl.Logger.Println("Initializing tzapConnect")
	openaiC := openaiconnector.InitiateOpenaiClient(openai_apikey, conf)

	tl.Logger.Println("Open AI Client Initialized")

	tl.Logger.Println("Local DB Client Initialized")
	partialComposite := PartialComposite{OpenaiTgenerator: openaiC}
	var myInterface types.TGenerator = partialComposite
	return myInterface, nil
}

type PartialComposite struct {
	*types.UnimplementedTGenerator
	OpenaiTgenerator *openaiconnector.OpenaiTgenerator
	VoiceGenerator   types.TGenerator
}

func (pc PartialComposite) TextToSpeech(ctx context.Context, content, language, voice string) (*[]byte, error) {
	return pc.VoiceGenerator.TextToSpeech(ctx, content, language, voice)
}
func (pc PartialComposite) GenerateChat(ctx context.Context, messages []types.Message, stream bool, functions string) (types.CompletionMessage, error) {
	return pc.OpenaiTgenerator.GenerateChat(ctx, messages, stream, functions)
}
func (pc PartialComposite) FetchEmbedding(ctx context.Context, content ...string) ([][1536]float32, error) {
	return pc.OpenaiTgenerator.FetchEmbedding(ctx, content...)
}
func (pc PartialComposite) CountTokens(ctx context.Context, content string) (int, error) {
	return pc.OpenaiTgenerator.CountTokens(content)
}
func (pc PartialComposite) OffsetTokens(ctx context.Context, content string, from int, to int) (string, int, error) {
	return pc.OpenaiTgenerator.OffsetTokens(content, from, to)
}
func (pc PartialComposite) SearchWithEmbedding(ctx context.Context, embedding types.QueryFilter, k int) (types.SearchResults, error) {
	return embedstore.EmbedStore.SearchWithEmbedding(ctx, embedding, k)
}
func (pc PartialComposite) AddEmbeddingDocument(ctx context.Context, docID string, embedding [1536]float32, metadata types.Metadata) error {
	return embedstore.EmbedStore.AddEmbeddingDocument(ctx, docID, embedding, metadata)
}
func (pc PartialComposite) GetEmbeddingDocument(ctx context.Context, docID string) (types.Vector, bool, error) {
	return embedstore.EmbedStore.GetEmbeddingDocument(ctx, docID)
}
func (pc PartialComposite) DeleteEmbeddingDocument(ctx context.Context, docID string) error {
	return embedstore.EmbedStore.DeleteEmbeddingDocument(ctx, docID)
}
func (pc PartialComposite) DeleteEmbeddingDocuments(ctx context.Context, ids []string) error {
	return embedstore.EmbedStore.DeleteEmbeddingDocuments(ctx, ids)
}
func (pc PartialComposite) ListAllEmbeddingsIds(ctx context.Context) (types.SearchResults, error) {
	return embedstore.EmbedStore.ListAllEmbeddingsIds(ctx)
}
