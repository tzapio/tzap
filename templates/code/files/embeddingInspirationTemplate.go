package files

import (
	"github.com/tzapio/tzap/pkg/embed"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
)

func EmbeddingInspirationTemplate(input string, excludedFiles []string) types.NamedTemplate[*tzap.Tzap, *tzap.Tzap] {
	return types.NamedTemplate[*tzap.Tzap, *tzap.Tzap]{
		Name: "embeddingInspirationTemplate",
		Template: func(t *tzap.Tzap) *tzap.Tzap {
			return t.
				IsolatedTzap(func(t *tzap.Tzap) {
					//check if index exists
				}).
				ApplyTemplate(InspirationTemplate(excludedFiles)).
				ApplyTemplate(SearchFilesTemplate(input, excludedFiles))

		},
	}
}

func SearchFilesTemplate(input string, excludedFiles []string) types.NamedTemplate[*tzap.Tzap, *tzap.Tzap] {
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
			searchResults, err := t.TG.SearchWithEmbedding(t.C, embedding, 10)
			if err != nil {
				panic(err)
			}
			filteredResults := filterSearchResults(searchResults, excludedFiles)

			s := ""

			for i, result := range filteredResults.Results {
				if i > 0 {
					s += "\n"
				}
				s += "###file: " + result.Metadata["filename"] + "\n"
				s += result.Metadata["splitPart"] + "\n"
			}

			return t.AddSystemMessage(
				"The following file contents are embeddings for the user input:",
				s,
			)
		},
	}
}

func filterSearchResults(searchResults types.SearchResults, excludedFiles []string) types.SearchResults {
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
	}
	return types.SearchResults{Results: filteredResults}
}
