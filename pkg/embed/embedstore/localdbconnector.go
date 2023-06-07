package embedstore

import "github.com/tzapio/tzap/pkg/types"

type EmbedStore struct {
	types.UnimplementedTGenerator
}

func InitiateLocalDB() (*EmbedStore, error) {
	return &EmbedStore{}, nil
}
