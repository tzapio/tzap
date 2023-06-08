package embedstore

import "github.com/tzapio/tzap/pkg/types"

type EmbedStore struct {
	types.UnimplementedTGenerator
}

func NewEmbedStore() (*EmbedStore, error) {
	return &EmbedStore{}, nil
}
