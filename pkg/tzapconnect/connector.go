package tzapconnect

import (
	"context"
	"os"

	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/connectors/localdbconnector"
	"github.com/tzapio/tzap/pkg/connectors/openaiconnector"
	"github.com/tzapio/tzap/pkg/types"
)

func WithConfig(openai_apikey string, conf config.Configuration) types.TzapConnector {
	tg, err := newBaseconnector(openai_apikey)
	if err != nil {
		println(err)
		os.Exit(1)
	}
	return func() (types.TGenerator, config.Configuration) {
		return tg, conf
	}
}

func newBaseconnector(openai_apikey string) (types.TGenerator, error) {
	openaiC := openaiconnector.InitiateOpenaiClient(openai_apikey)
	embeddingC, err := localdbconnector.InitiateLocalDB("./.tzap-data/fileembeddings.db")
	if err != nil {
		return nil, err
	}
	partialComposite := PartialComposite{OpenaiTgenerator: openaiC, EmbeddingGenerator: embeddingC}
	var myInterface types.TGenerator = partialComposite
	return myInterface, nil
}

type PartialComposite struct {
	*types.UnimplementedTGenerator
	OpenaiTgenerator   *openaiconnector.OpenaiTgenerator
	VoiceGenerator     types.TGenerator
	EmbeddingGenerator types.TGenerator
}

func (pc PartialComposite) TextToSpeech(ctx context.Context, content, language, voice string) (*[]byte, error) {
	return pc.VoiceGenerator.TextToSpeech(ctx, content, language, voice)
}
func (pc PartialComposite) GenerateChat(ctx context.Context, messages []types.Message, stream bool) (string, error) {
	return pc.OpenaiTgenerator.GenerateChat(ctx, messages, stream)
}
func (pc PartialComposite) FetchEmbedding(ctx context.Context, content ...string) ([][]float32, error) {
	return pc.OpenaiTgenerator.FetchEmbedding(ctx, content...)
}
func (pc PartialComposite) CountTokens(ctx context.Context, content string) (int, error) {
	return pc.OpenaiTgenerator.CountTokens(content)
}
func (pc PartialComposite) OffsetTokens(ctx context.Context, content string, from int, to int) (string, int, error) {
	return pc.OpenaiTgenerator.OffsetTokens(content, from, to)
}
func (pc PartialComposite) SearchWithEmbedding(ctx context.Context, embedding types.QueryFilter, k int) (types.SearchResults, error) {
	return pc.EmbeddingGenerator.SearchWithEmbedding(ctx, embedding, k)
}
func (pc PartialComposite) AddEmbeddingDocument(ctx context.Context, docID string, embedding []float32, metadata map[string]string) error {
	return pc.EmbeddingGenerator.AddEmbeddingDocument(ctx, docID, embedding, metadata)
}
func (pc PartialComposite) GetEmbeddingDocument(ctx context.Context, docID string) (types.Vector, bool, error) {
	return pc.EmbeddingGenerator.GetEmbeddingDocument(ctx, docID)
}
func (pc PartialComposite) DeleteEmbeddingDocument(ctx context.Context, docID string) error {
	return pc.EmbeddingGenerator.DeleteEmbeddingDocument(ctx, docID)
}
func (pc PartialComposite) DeleteEmbeddingDocuments(ctx context.Context, ids []string) error {
	return pc.EmbeddingGenerator.DeleteEmbeddingDocuments(ctx, ids)
}
func (pc PartialComposite) ListAllEmbeddingsIds(ctx context.Context) (types.SearchResults, error) {
	return pc.EmbeddingGenerator.ListAllEmbeddingsIds(ctx)
}
