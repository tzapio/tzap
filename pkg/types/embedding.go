package types

type Vector struct {
	ID        string            `json:"id"`
	TimeStamp int               `json:"timestamp"`
	Metadata  map[string]string `json:"metadata"`
	Values    []float32         `json:"values"`
}

type Embeddings struct {
	Vectors []Vector `json:"vectors"`
}

type Query struct {
	Values []float32 `json:"values"`
}
type QueryJson struct {
	Queries []Query `json:"queries"`
}

type QueryFilter struct {
	Filter map[string]string `json:"filter"`
	Values []float32         `json:"values"`
}
type QueryRequest struct {
	TopK            int           `json:"topK"`
	IncludeMetadata bool          `json:"includeMetadata"`
	Namespace       string        `json:"namespace"`
	Queries         []QueryFilter `json:"queries"`
}
type SearchResults struct {
	Results []Vector
}
type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
