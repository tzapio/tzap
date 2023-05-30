package embed

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/tzapio/tzap/pkg/types"
)

func SaveEmbeddingToFile(e types.Embeddings) {
	if err := BatchEmbeddings(e); err != nil {
		panic(err)
	}
	if err := SaveVectorsToFile(e, "./.tzap-data/files.json"); err != nil {
		panic(err)
	}
}
func SaveVectorsToFile(e types.Embeddings, filePath string) error {
	embeddingJSON, err := json.Marshal(e)
	if err != nil {
		return err
	}

	if err := os.WriteFile(filePath, embeddingJSON, 0644); err != nil {
		return err
	}
	//fmt.Printf("Upserted files (count: %d): %s\n", len(e.Vectors), filePath)
	return nil
}

func BatchEmbeddings(e types.Embeddings) error {
	batchSize := 100
	batchNumber := 1
	var batch []types.Vector

	for i, vector := range e.Vectors {
		batch = append(batch, vector)
		deletePreviousBatch()
		if (i+1)%batchSize == 0 || i == len(e.Vectors)-1 {
			filePath := fmt.Sprintf("./.tzap-data/files-%d.json", batchNumber)
			batchEmbeddingJson := types.Embeddings{
				Vectors: batch,
			}
			err := SaveVectorsToFile(batchEmbeddingJson, filePath)
			if err != nil {
				return err
			}
			batch = nil
			batchNumber++
		}
	}
	return nil
}
func deletePreviousBatch() error {
	// Remove previous embedding files
	for i := 1; ; i++ {
		filePath := fmt.Sprintf("./.tzap-data/files-%d.json", i)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			break
		}
		if err := os.Remove(filePath); err != nil {
			return err
		}
	}
	return nil
}
