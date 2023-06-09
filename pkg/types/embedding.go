package types

type Vector struct {
	ID        string        `json:"id"`
	TimeStamp int           `json:"timestamp"`
	Metadata  Metadata      `json:"metadata"`
	Values    [1536]float32 `json:"values"`
}
type Metadata struct {
	ID            string `json:"id"`
	Filename      string `json:"filename"`
	Start         int    `json:"start"`
	End           int    `json:"end"`
	LineStart     int    `json:"lineStart"`
	TruncatedEnd  int    `json:"truncatedEnd"`
	SplitPart     string `json:"splitPart"`
	RealSplitPart string `json:"realSplitPart"`
}
type Embeddings struct {
	Vectors []*Vector `json:"vectors"`
}

type Query struct {
	Values [1536]float32 `json:"values"`
}
type QueryJson struct {
	Queries []Query `json:"queries"`
}

type QueryFilter struct {
	Filter map[string]string `json:"filter"`
	Values [1536]float32     `json:"values"`
}
type QueryRequest struct {
	TopK            int           `json:"topK"`
	IncludeMetadata bool          `json:"includeMetadata"`
	Namespace       string        `json:"namespace"`
	Queries         []QueryFilter `json:"queries"`
}
type SearchResult struct {
	Vector     Vector    `json:"vector"`
	PCA        []float32 `json:"pca"`
	Similarity float32   `json:"score"`
}
type SearchResults struct {
	Results []SearchResult
}
type KeyValue[T any] struct {
	Key   string `json:"key"`
	Value T      `json:"value"`
}
