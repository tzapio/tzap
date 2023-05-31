package localdbconnector

import (
	"context"

	"github.com/tzapio/tzap/pkg/types"
)

func (idx *LocalembedTGenerator) AddEmbeddingDocument(ctx context.Context, docID string, embedding [1536]float32, metadata map[string]string) error {
	v := types.Vector{
		ID:        docID,
		TimeStamp: 0,
		Metadata:  metadata,
		Values:    embedding,
	}
	if err := idx.db.Set(docID, v); err != nil {
		return err
	}
	return nil
}
func (idx *LocalembedTGenerator) AddEmbeddingDocuments(ctx context.Context, vectors []types.Vector) (int, error) {
	pairs := []types.KeyValue[types.Vector]{}
	for _, v := range vectors {
		pairs = append(pairs, types.KeyValue[types.Vector]{Key: v.ID, Value: v})
	}
	wrote, err := idx.db.BatchSet(pairs)
	if err != nil {
		return wrote, err
	}
	return wrote, nil
}
func (idx *LocalembedTGenerator) DeleteEmbeddingDocument(ctx context.Context, docID string) error {
	if err := idx.db.Set(docID, types.Vector{}); err != nil {
		return err
	}
	return nil
}
func (idx *LocalembedTGenerator) DeleteEmbeddingDocuments(ctx context.Context, docIDs []string) error {
	var pairs []types.KeyValue[types.Vector]
	for _, docID := range docIDs {
		pairs = append(pairs, types.KeyValue[types.Vector]{Key: docID})
	}
	if _, err := idx.db.BatchSet(pairs); err != nil {
		return err
	}
	return nil
}
