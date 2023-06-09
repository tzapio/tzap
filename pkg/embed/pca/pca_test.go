package pca_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/tzapio/tzap/pkg/embed/export"
	"github.com/tzapio/tzap/pkg/embed/pca"
)

func TestEmbeddingsTo3D(t *testing.T) {
	embeddings, err := export.GetEmbeddingsFromFile("files.testjson")
	if err != nil {
		t.Fatalf("error getting embeddings from file: %s", err)
	}

	if len(embeddings.Vectors) == 0 {
		t.Fatalf("embeddings file is empty")
	}
	var vectors [][1536]float32
	for _, v := range embeddings.Vectors {
		vectors = append(vectors, v.Values)
	}
	pcaResult := pca.EmbeddingsTo3D(vectors)
	if err != nil {
		t.Fatalf("error computing PCA: %s", err)
	}
	println(pcaResult)

	pcaJSON, err := json.Marshal(pcaResult)
	if err != nil {
		t.Fatalf("error marshalling PCA result: %s", err)
	}

	if err := os.WriteFile("pca.txt", []byte(pcaJSON), 0644); err != nil {
		t.Fatalf("error writing pca.txt: %s", err)
	}
	println("writing pca.txt")

	if len(pcaResult) == 0 {
		t.Fatalf("PCA result has no points")
	}

}
