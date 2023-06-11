package embedstore

import (
	"errors"
	"sort"
	"strings"

	"github.com/tzapio/tzap/pkg/embed"
	"github.com/tzapio/tzap/pkg/types"
)

// Sorts search results in a way that preserves original order based on filenames
func TightenSearchResults(searchResults []types.SearchResult) []types.SearchResult {
	// group search results by filename
	searchResultsByFilename := make(map[string][]types.SearchResult)
	for _, sr := range searchResults {
		filename := sr.Vector.Metadata.Filename
		searchResultsByFilename[filename] = append(searchResultsByFilename[filename], sr)
	}

	// sort each group of search results by start position
	var sortedResults []types.SearchResult
	for _, results := range searchResultsByFilename {
		sort.SliceStable(results, func(i, j int) bool {
			return results[i].Vector.Metadata.Start < results[j].Vector.Metadata.Start
		})
		consecutiveResults := groupConsecutiveMetadata(results)
		sortedResults = append(sortedResults, consecutiveResults...)
	}

	return sortedResults
}

// Groups consecutive metadata together from sorted search results
func groupConsecutiveMetadata(searchResults []types.SearchResult) []types.SearchResult {
	var resultsWithConsecutive []types.SearchResult
	var currentGroup []types.SearchResult

	for i := range searchResults {
		if len(currentGroup) == 0 {
			currentGroup = append(currentGroup, searchResults[i])
		} else {
			if searchResults[i].Vector.Metadata.Start == currentGroup[len(currentGroup)-1].Vector.Metadata.Start+200 {
				currentGroup = append(currentGroup, searchResults[i])
			} else {
				resultWithConcatenatedMetadata := concatenateConsecutiveMetadata(currentGroup)
				resultsWithConsecutive = append(resultsWithConsecutive, resultWithConcatenatedMetadata)
				currentGroup = []types.SearchResult{searchResults[i]}
			}
		}
	}

	// concatenate last group of consecutive metadata, if any
	if len(currentGroup) > 0 {
		resultWithConcatenatedMetadata := concatenateConsecutiveMetadata(currentGroup)
		resultsWithConsecutive = append(resultsWithConsecutive, resultWithConcatenatedMetadata)
	}

	return resultsWithConsecutive
}

// Concatenates consecutive metadata from a group of search results
func concatenateConsecutiveMetadata(searchResults []types.SearchResult) types.SearchResult {
	if len(searchResults) == 1 {
		return searchResults[0]
	}
	if len(searchResults) == 0 {
		panic(errors.New("searchResult may not be empty"))
	}
	first := searchResults[0].Vector.Metadata
	last := searchResults[len(searchResults)-1].Vector.Metadata
	filename := searchResults[0].Vector.Metadata.Filename
	return types.SearchResult{
		Vector: types.Vector{
			Metadata: types.Metadata{
				Filename:     filename,
				Start:        first.Start,
				LineStart:    first.LineStart,
				End:          last.End,
				TruncatedEnd: last.TruncatedEnd,
				SplitPart:    concatSplitPart(filename, searchResults),
			},
		},
	}
}
func concatSplitPart(filename string, searchResults []types.SearchResult) string {
	var splitPart strings.Builder = strings.Builder{}
	if len(searchResults) <= 1 {
		panic(errors.New("searchResult may not be empty or only be single element"))
	}
	for i, sr := range searchResults {
		if i == len(searchResults)-1 {
			splitPart.WriteString(embed.StripEmbedHeader(sr.Vector.Metadata.SplitPart))
		} else {
			splitPart.WriteString(sr.Vector.Metadata.RealSplitPart)
		}
	}
	return embed.AddEmbedHeader(filename, splitPart.String())
}

// Returns true if there are any consecutive metadata in the search results
func hasConsecutiveMetadata(searchResults []types.SearchResult) bool {
	var start int
	var hasConsecutive bool

	for _, sr := range searchResults {
		if start == 0 {
			start = sr.Vector.Metadata.Start
		} else if sr.Vector.Metadata.Start == start+200 {
			hasConsecutive = true
			break
		}
	}

	return hasConsecutive
}
