package cosine_test

import (
	"encoding/json"
	"math"
	"os"
	"testing"

	"github.com/tzapio/tzap/pkg/embed/cosine"
	"github.com/tzapio/tzap/pkg/types"
)

func TestCosineDistance(t *testing.T) {
	var epsilon = 0.000001
	filecontent, err := os.ReadFile("test-files-1.json")
	if err != nil {
		panic(err)
	}
	var data types.Embeddings
	if err := json.Unmarshal(filecontent, &data); err != nil {
		panic(err)
	}
	queryContent, err := os.ReadFile("test-query.json")
	if err != nil {
		panic(err)
	}
	var query types.QueryJson
	if err := json.Unmarshal(queryContent, &query); err != nil {
		panic(err)
	}

	tests := []struct {
		name     string
		a        [1536]float32
		b        [1536]float32
		expected float32
	}{
		{
			name:     "identical vectors",
			a:        data.Vectors[0].Values,
			b:        query.Queries[0].Values,
			expected: 0.644763,
			//epsilon 0.000001
		},
		{
			name:     "orthogonal vectors",
			a:        data.Vectors[1].Values,
			b:        query.Queries[0].Values,
			expected: 0.626373,
			//epsilon 0.000001
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			distance := cosine.CosineSimilarity(test.a, test.b)
			if math.Abs(float64(distance-test.expected)) > epsilon {
				t.Errorf("Expected distance: %f, got: %f", test.expected, distance)
			}
		})
	}

	/*
		var results = [][1536]float{}

		for _, vector := range data.Vectors {
			results = append(results, vector.Values)
		}
		// SearchByCosineSimilarity
		res := cosine.SearchByCosineSimilarity(results, query.Queries[0].Values)
		for _, r := range res {
			println(data.Vectors[r.Index].ID, fmt.Sprintf("%f", r.Similarity), data.Vectors[r.Index].Metadata["filename"])
		}*/
}
