package openaiconnector

import (
	"context"

	"github.com/sashabaranov/go-openai"
	"github.com/tzapio/tzap/pkg/types"
)

type openaiconnectorKey struct{}

type OpenaiTgenerator struct {
	client *openai.Client
}

func NewContext(ctx context.Context, openaiTgenerator OpenaiTgenerator) context.Context {
	return context.WithValue(ctx, openaiconnectorKey{}, openaiTgenerator)
}

func FromContext(ctx context.Context) types.TGenerator {
	if config, ok := ctx.Value(openaiconnectorKey{}).(types.TGenerator); ok {
		return config
	}
	panic("OpenaiTgenerator is not configured")
}
