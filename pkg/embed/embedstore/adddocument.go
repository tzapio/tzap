package embedstore

import (
	"context"

	"github.com/tzapio/tzap/pkg/project"
	"github.com/tzapio/tzap/pkg/types"
)

func (idx *embedStore) AddEmbeddingDocument(ctx context.Context, docID string, embedding [1536]float32, metadata types.Metadata) error {
	v := types.Vector{
		ID:        docID,
		TimeStamp: 0,
		Metadata:  metadata,
		Values:    embedding,
	}
	embeddingCollection := project.GetProjectFromContext(ctx).GetEmbeddingCollection()
	if err := embeddingCollection.Set(docID, v); err != nil {
		return err
	}
	return nil
}
func (idx *embedStore) AddEmbeddingDocuments(ctx context.Context, vectors []types.Vector) (int, error) {
	pairs := []types.KeyValue[types.Vector]{}
	for _, v := range vectors {
		pairs = append(pairs, types.KeyValue[types.Vector]{Key: v.ID, Value: v})
	}
	embeddingCollection := project.GetProjectFromContext(ctx).GetEmbeddingCollection()
	wrote, err := embeddingCollection.BatchSet(pairs)
	if err != nil {
		return wrote, err
	}
	return wrote, nil
}
func (idx *embedStore) DeleteEmbeddingDocument(ctx context.Context, docID string) error {
	embeddingCollection := project.GetProjectFromContext(ctx).GetEmbeddingCollection()
	if err := embeddingCollection.Set(docID, types.Vector{}); err != nil {
		return err
	}
	return nil
}
func (idx *embedStore) DeleteEmbeddingDocuments(ctx context.Context, docIDs []string) error {
	var pairs []types.KeyValue[types.Vector]
	for _, docID := range docIDs {
		pairs = append(pairs, types.KeyValue[types.Vector]{Key: docID})
	}
	embeddingCollection := project.GetProjectFromContext(ctx).GetEmbeddingCollection()
	if _, err := embeddingCollection.BatchSet(pairs); err != nil {
		return err
	}
	return nil
}
