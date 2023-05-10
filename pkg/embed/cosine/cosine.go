package cosine

import (
	"math"
	"sort"
)

func dotProduct(a, b []float32) float32 {
	if len(a) != len(b) {
		panic("Vectors must have the same dimensions")
	}

	var dot float32 = 0.0
	for i := range a {
		dot += a[i] * b[i]
	}
	return dot
}

func magnitude(x []float32) float32 {
	var mag float32 = 0.0
	for _, v := range x {
		mag += v * v
	}
	return float32(math.Sqrt(float64(mag)))
}

func CosineSimilarity(a, b []float32) float32 {
	return dotProduct(a, b) / (magnitude(a) * magnitude(b))
}

type Result struct {
	Index      int
	Similarity float32
}

func SearchByCosineSimilarity(results [][]float32, query []float32) []Result {
	similarities := make([]Result, len(results))

	for i, result := range results {
		similarity := CosineSimilarity(result, query)
		similarities[i] = Result{Index: i, Similarity: similarity}
	}

	sort.Slice(similarities, func(i, j int) bool {
		return similarities[i].Similarity > similarities[j].Similarity
	})

	return similarities
}
