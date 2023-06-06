package embedworkflows

import (
	"fmt"
	"os"

	"github.com/tzapio/tzap/pkg/embed"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/util/stdin"
)

func LoadAndFetchEmbeddings(name types.ProjectName, files []types.FileReader, embedder *embed.Embedder, disableIndex, yes bool) types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap] {
	return types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap]{
		Name: "loadAndFetchEmbeddings",
		Workflow: func(t *tzap.Tzap) *tzap.Tzap {
			return t.MutationTzap(func(t *tzap.Tzap) *tzap.Tzap {
				if !disableIndex {
					println("Checking for file changes. (use -d to disable this check)...\n")
					return t.
						ApplyWorkflow(PrepareEmbedFilesWorkflow(name, files, embedder)).
						ApplyWorkflow(ConfirmEmbeddingSearch(yes)).
						ApplyWorkflow(FetchOrCachedEmbeddingForFilesWorkflow(files)).
						ApplyWorkflow(SaveAndLoadEmbeddingsToDB(name))
				}
				return t
			})
		},
	}
}
func ConfirmEmbeddingSearch(yes bool) types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap] {
	return types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap]{
		Name: "confirmEmbeddingSearch",
		Workflow: func(t *tzap.Tzap) *tzap.Tzap {
			uncachedEmbeddings, ok := t.Data["uncachedEmbeddings"].(types.Embeddings)
			if !ok {
				panic("Loading embeddings went wrong")
			}
			if len(uncachedEmbeddings.Vectors) > 19 {
				price := float64(len(uncachedEmbeddings.Vectors)*400) * 0.0004 / 1000
				if !yes {
					ok := stdin.ConfirmPrompt(fmt.Sprintf(
						"Embeddings - You are about to fetch %d embeddings. Proceed? Estimation tokens: %d. Price is: $0.0004 per 1000 tokens. Estimating %.4f USD",
						len(uncachedEmbeddings.Vectors), len(uncachedEmbeddings.Vectors)*400, price))
					if !ok {
						println("Fetching embeddings aborted by user")
						os.Exit(0)
					}
				}
			}
			return t
		},
	}
}
