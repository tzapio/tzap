package localdbconnector

import (
	"context"
	"encoding/json"

	"github.com/tzapio/tzap/pkg/types"
)

func (idx *LocalembedTGenerator) AddEmbeddingDocument(ctx context.Context, docID string, embedding []float32, metadata map[string]string) error {
	v := types.Vector{
		ID:        docID,
		TimeStamp: 0,
		Metadata:  metadata,
		Values:    embedding,
	}
	vectorBytes, err := json.Marshal(v)
	if err != nil {
		return err
	}
	if err := idx.db.Set(docID, string(vectorBytes)); err != nil {
		return err
	}
	return nil
}
func (idx *LocalembedTGenerator) AddEmbeddingDocuments(ctx context.Context, vectors []types.Vector) (int, error) {
	pairs := []types.KeyValue{}
	for _, v := range vectors {
		vectorBytes, err := json.Marshal(v)
		if err != nil {
			return 0, err
		}
		pairs = append(pairs, types.KeyValue{Key: v.ID, Value: string(vectorBytes)})
	}
	wrote, err := idx.db.BatchSet(pairs)
	if err != nil {
		return wrote, err
	}
	return wrote, nil
}
func (idx *LocalembedTGenerator) DeleteEmbeddingDocument(ctx context.Context, docID string) error {
	if err := idx.db.Set(docID, ""); err != nil {
		return err
	}
	return nil
}
