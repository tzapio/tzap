package openaiconnector

import (
	"context"
	"errors"
	"time"

	"github.com/sashabaranov/go-openai"
	"github.com/tzapio/tzap/internal/logging/tl"
)

func (ot *OpenaiTgenerator) FetchEmbedding(ctx context.Context, content ...string) ([][1536]float32, error) {
	tl.Logger.Println("Fetching embeddings for", len(content), "strings")
	request := openai.EmbeddingRequest{
		Model: openai.AdaEmbeddingV2,
		Input: content,
	}
	retries := 3
	for i := 0; i < retries; i++ {
		response, err := ot.client.CreateEmbeddings(ctx, request)
		if err != nil {
			e := &openai.APIError{}
			if errors.As(err, &e) {
				switch e.HTTPStatusCode {
				case 401:
					panic("invalid open ai api key")
				case 429:
					println(err.Error())
				case 500:
					// openai server error (retry)
				default:
					println("unknown error", e.HTTPStatusCode, err.Error())
					return nil, errors.New("embedding failed")
				}
			}
			println("Fetching embedding failed. Retry in 2 seconds.", err.Error())
			time.Sleep(2 * time.Second)
			continue
		}
		embeddings := [][1536]float32{}
		for _, embedding := range response.Data {
			embeddings = append(embeddings, [1536]float32(embedding.Embedding))
		}
		return embeddings, nil
	}
	return nil, errors.New("embedding failed")
}
