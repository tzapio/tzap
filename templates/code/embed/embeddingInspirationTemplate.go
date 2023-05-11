package embed

import (
	"github.com/tzapio/tzap/pkg/embed"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
)

func EmbeddingInspirationTemplate(input string, inspirationFiles []string) types.NamedTemplate[*tzap.Tzap, *tzap.Tzap] {
	return types.NamedTemplate[*tzap.Tzap, *tzap.Tzap]{
		Name: "embeddingInspirationTemplate",
		Template: func(t *tzap.Tzap) *tzap.Tzap {
			return t.
				ApplyTemplate(InspirationTemplate(inspirationFiles)).
				ApplyTemplate(SearchFilesTemplate(input, inspirationFiles))

		},
	}
}

func SearchFilesTemplate(input string, excludeFiles []string) types.NamedTemplate[*tzap.Tzap, *tzap.Tzap] {
	return types.NamedTemplate[*tzap.Tzap, *tzap.Tzap]{
		Name: "searchFilesTemplate",
		Template: func(t *tzap.Tzap) *tzap.Tzap {
			query, err := embed.GetQuery(t, input)
			if err != nil {
				panic(err)
			}
			if len(query.Queries) == 0 {
				panic("empty embeddings")
			}

			if len(query.Queries) > 1 {
				panic("should only return one embedding")
			}
			embedding := query.Queries[0]
			searchResults, err := t.TG.SearchWithEmbedding(t.C, embedding, 20)
			if err != nil {
				panic(err)
			}
			filteredResults := filterSearchResults(searchResults, excludeFiles, 15)

			t = t.AddSystemMessage(
				"The following file contents are embeddings for the user input:",
			)
			println("Using embeddings:")
			for _, result := range filteredResults.Results {
				t = t.AddSystemMessage(result.Metadata["splitPart"])
				println(result.ID)
			}

			return t
		},
	}
}

func filterSearchResults(searchResults types.SearchResults, excludedFiles []string, k int) types.SearchResults {
	filteredResults := []types.Vector{}
	for _, result := range searchResults.Results {
		fileName := result.Metadata["filename"]
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
		if len(filteredResults) >= k {
			break
		}
	}
	return types.SearchResults{Results: filteredResults}
}
