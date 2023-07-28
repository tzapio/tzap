package embedworkflows

import (
	"github.com/tzapio/tzap/internal/logging/tl"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
)

// k is amount of embeddings to be included.
// When using inspiration files, embeddings are likely to be duplicated and as such are filtered out. n is used to increase how many embeddings are fetched but are trimmed to only contain top K after filtering.
func SearchFilesWorkflow(query types.QueryRequest, excludeFiles []string, k int) types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap] {
	return types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap]{
		Name: "searchFilesWorkflow",
		Workflow: func(t *tzap.Tzap) *tzap.Tzap {
			tl.Logger.Println("searchFilesWorkflow")
			if len(query.Queries) == 0 {
				panic("empty embeddings")
			}

			if len(query.Queries) > 1 {
				panic("should only return one embedding")
			}
			embedding := query.Queries[0]
			searchResults, err := t.TG.SearchWithEmbedding(t.C, embedding, k)
			if err != nil {
				panic(err)
			}
			filteredResults := filterSearchResults(searchResults, excludeFiles, k)

			data := types.MappedInterface{
				"queryResult":   query,
				"searchResults": filteredResults,
			}
			tl.Logger.Println("searchFilesWorkflow ending")
			return t.AddTzap(&tzap.Tzap{Name: "searchResults", Data: data})
		},
	}
}

func filterSearchResults(searchResults types.SearchResults, excludedFiles []string, k int) types.SearchResults {
	filteredResults := []types.SearchResult{}
	for _, result := range searchResults.Results {
		if len(filteredResults) >= k && k > -1 {
			break
		}
		fileName := result.Vector.Metadata.Filename
		isExcluded := false
		for _, excludedFile := range excludedFiles {
			if fileName == excludedFile {
				isExcluded = true
				break
			}
		}
		if !isExcluded {
			filteredResults = append(filteredResults, result)
		}

	}
	return types.SearchResults{Results: filteredResults}
}
