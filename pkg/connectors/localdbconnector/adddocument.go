package localdbconnector

import (
	"context"

	"github.com/tzapio/tzap/pkg/types"
)

func (idx *LocalembedTGenerator) AddEmbeddingDocument(ctx context.Context, project string, docID string, embedding [1536]float32, metadata types.Metadata) error {
	v := types.Vector{
		ID:        docID,
		TimeStamp: 0,
		Metadata:  metadata,
		Values:    embedding,
	}
	if err := idx.dbs[types.ProjectName(project)].Set(docID, v); err != nil {
		return err
	}
	return nil
}
func (idx *LocalembedTGenerator) AddEmbeddingDocuments(ctx context.Context, project string, vectors []types.Vector) (int, error) {
	pairs := []types.KeyValue[types.Vector]{}
	for _, v := range vectors {
		pairs = append(pairs, types.KeyValue[types.Vector]{Key: v.ID, Value: v})
	}
	wrote, err := idx.dbs[types.ProjectName(project)].BatchSet(pairs)
	if err != nil {
		return wrote, err
	}
	return wrote, nil
}
func (idx *LocalembedTGenerator) DeleteEmbeddingDocument(ctx context.Context, project string, docID string) error {
	if err := idx.dbs[types.ProjectName(project)].Set(docID, types.Vector{}); err != nil {
		return err
	}
	return nil
}
func (idx *LocalembedTGenerator) DeleteEmbeddingDocuments(ctx context.Context, project string, docIDs []string) error {
	var pairs []types.KeyValue[types.Vector]
	for _, docID := range docIDs {
		pairs = append(pairs, types.KeyValue[types.Vector]{Key: docID})
	}
	if _, err := idx.dbs[types.ProjectName(project)].BatchSet(pairs); err != nil {
		return err
	}
	return nil
}
