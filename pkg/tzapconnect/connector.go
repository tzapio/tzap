package tzapconnect

import (
	"context"
	"os"

	"github.com/tzapio/tzap/internal/logging/tl"
	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/connectors/localdbconnector"
	"github.com/tzapio/tzap/pkg/connectors/openaiconnector"
	"github.com/tzapio/tzap/pkg/types"
)

func WithConfig(openai_apikey string, libs types.ProjectDB, conf config.Configuration) types.TzapConnector {
	newLibs := make(types.ProjectDB)
	newLibs["@LOCAL"] = "./.tzap-data/"
	for k, v := range libs {
		newLibs[k] = v
	}

	tg, err := newBaseconnector(openai_apikey, newLibs)
	if err != nil {
		println(err)
		os.Exit(1)
	}
	return func() (types.TGenerator, config.Configuration) {
		return tg, conf
	}
}

func newBaseconnector(openai_apikey string, libs types.ProjectDB) (types.TGenerator, error) {
	tl.Logger.Println("Initializing tzapConnect")
	openaiC := openaiconnector.InitiateOpenaiClient(openai_apikey)

	tl.Logger.Println("Open AI Client Initialized")
	embeddingC, err := localdbconnector.InitiateLocalDB(libs)
	if err != nil {
		return nil, err
	}
	tl.Logger.Println("Local DB Client Initialized")
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
func (pc PartialComposite) FetchEmbedding(ctx context.Context, content ...string) ([][1536]float32, error) {
	return pc.OpenaiTgenerator.FetchEmbedding(ctx, content...)
}
func (pc PartialComposite) CountTokens(ctx context.Context, content string) (int, error) {
	return pc.OpenaiTgenerator.CountTokens(content)
}
func (pc PartialComposite) OffsetTokens(ctx context.Context, content string, from int, to int) (string, int, error) {
	return pc.OpenaiTgenerator.OffsetTokens(content, from, to)
}
func (pc PartialComposite) SearchWithEmbedding(ctx context.Context, project string, embedding types.QueryFilter, k int) (types.SearchResults, error) {
	return pc.EmbeddingGenerator.SearchWithEmbedding(ctx, project, embedding, k)
}
func (pc PartialComposite) AddEmbeddingDocument(ctx context.Context, project string, docID string, embedding [1536]float32, metadata types.Metadata) error {
	return pc.EmbeddingGenerator.AddEmbeddingDocument(ctx, project, docID, embedding, metadata)
}
func (pc PartialComposite) GetEmbeddingDocument(ctx context.Context, project string, docID string) (types.Vector, bool, error) {
	return pc.EmbeddingGenerator.GetEmbeddingDocument(ctx, project, docID)
}
func (pc PartialComposite) DeleteEmbeddingDocument(ctx context.Context, project string, docID string) error {
	return pc.EmbeddingGenerator.DeleteEmbeddingDocument(ctx, project, docID)
}
func (pc PartialComposite) DeleteEmbeddingDocuments(ctx context.Context, project string, ids []string) error {
	return pc.EmbeddingGenerator.DeleteEmbeddingDocuments(ctx, project, ids)
}
func (pc PartialComposite) ListAllEmbeddingsIds(ctx context.Context, project string) (types.SearchResults, error) {
	return pc.EmbeddingGenerator.ListAllEmbeddingsIds(ctx, project)
}
