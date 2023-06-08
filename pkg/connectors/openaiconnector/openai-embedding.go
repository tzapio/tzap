package openaiconnector

import (
	"context"

	"github.com/sashabaranov/go-openai"
	"github.com/tzapio/tzap/internal/logging/tl"
)

func (ot OpenaiTgenerator) FetchEmbedding(ctx context.Context, content ...string) ([][1536]float32, error) {
	tl.Logger.Println("Fetching embeddings for", len(content), "strings")
	request := openai.EmbeddingRequest{
		Model: openai.AdaEmbeddingV2,
		Input: content,
	}
	response, err := ot.client.CreateEmbeddings(ctx, request)
	if err != nil {
		return [][1536]float32{}, err
	}
	embeddings := [][1536]float32{}
	for _, embedding := range response.Data {
		embeddings = append(embeddings, [1536]float32(embedding.Embedding))
	}
	return embeddings, nil
}

func (ot OpenaiTgenerator) CountTokens(content string) (int, error) {
	return ot.tokenizer.CountTokens(content)
}

func (ot OpenaiTgenerator) OffsetTokens(content string, from int, to int) (string, int, error) {
	return ot.tokenizer.OffsetTokens(content, from, to)
}

func (ot OpenaiTgenerator) RawTokens(content string) ([]string, error) {
	return ot.tokenizer.RawTokens(content)
}
