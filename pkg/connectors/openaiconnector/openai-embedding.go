package openaiconnector

import (
	"context"

	"github.com/sashabaranov/go-openai"
	"github.com/tzapio/tzap/pkg/connectors/openaiconnector/tokenizer"
)

func (ot OpenaiTgenerator) FetchEmbedding(ctx context.Context, content ...string) ([][]float32, error) {
	request := openai.EmbeddingRequest{
		Model: openai.AdaEmbeddingV2,
		Input: content,
	}
	response, err := ot.client.CreateEmbeddings(ctx, request)
	if err != nil {
		return [][]float32{}, err
	}
	embeddings := [][]float32{}
	for _, embedding := range response.Data {
		embeddings = append(embeddings, embedding.Embedding)
	}
	return embeddings, nil
}

func (ot OpenaiTgenerator) CountTokens(content string) (int, error) {
	return tokenizer.CountTokens(content)
}

func (ot OpenaiTgenerator) OffsetTokens(content string, from int, to int) (string, int, error) {
	return tokenizer.OffsetTokens(content, from, to)
}

func (ot OpenaiTgenerator) RawTokens(content string) ([]string, error) {
	return tokenizer.RawTokens(content)
}
