package cliworkflows

import (
	"fmt"
	"os"
	"strings"

	"github.com/tzapio/tzap/cli/cmd/cmdutil"
	"github.com/tzapio/tzap/internal/logging/tl"
	"github.com/tzapio/tzap/pkg/embed"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/util"
	"github.com/tzapio/tzap/pkg/util/stdin"
	"github.com/tzapio/tzap/workflows/code/embedworkflows"
)

func IndexFilesAndEmbeddings(dir string, disableIndex, yes bool) types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap] {
	return types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap]{
		Name: "indexFilesAndEmbeddings",
		Workflow: func(t *tzap.Tzap) *tzap.Tzap {
			if disableIndex {
				return t
			}
			fileInDirEvaluator, err := cmdutil.NewFileEvaluator()
			if err != nil {
				panic(err)
			}
			embedder := embed.NewEmbedder(t)
			tl.Logger.Println("Indexing files...")
			files, err := fileInDirEvaluator.WalkDir("./")
			if err != nil {
				panic(err)
			}
			tl.Logger.Println("Finished index files...")

			return t.ApplyWorkflow(LoadAndFetchEmbeddings(files, embedder, disableIndex, yes))
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
func LoadAndFetchEmbeddings(files []types.FileReader, embedder *embed.Embedder, disableIndex, yes bool) types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap] {
	return types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap]{
		Name: "loadAndFetchEmbeddings",
		Workflow: func(t *tzap.Tzap) *tzap.Tzap {
			return t.MutationTzap(func(t *tzap.Tzap) *tzap.Tzap {
				if !disableIndex {
					println("Checking for file changes. (use -d to disable this check)...\n")
					return t.
						ApplyWorkflow(embedworkflows.PrepareEmbedFilesWorkflow(files, embedder)).
						ApplyWorkflow(ConfirmEmbeddingSearch(yes)).
						ApplyWorkflow(embedworkflows.FetchOrCachedEmbeddingForFilesWorkflow()).
						ApplyWorkflow(embedworkflows.SaveAndLoadEmbeddingsToDB())
				}
				return t
			})
		},
	}
}
func PrintInspirationFiles(inspirationFiles []string) types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap] {
	return types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap]{
		Name: "listInspirationFiles",
		Workflow: func(t *tzap.Tzap) *tzap.Tzap {
			return t.WorkTzap(func(t *tzap.Tzap) {
				if len(inspirationFiles) == 0 {
					println(cmdutil.Bold("Inspiration files: None (use --inspiration to add more)"))
					return
				}
				println(cmdutil.Bold("\nInspiration files:"))
				for _, inspirationFile := range inspirationFiles {
					inspirationFile = strings.TrimSpace(inspirationFile)
					tokens, err := t.CountTokens(util.ReadFileP(inspirationFile))
					if err != nil {
						panic(err)
					}
					fmt.Fprintf(os.Stderr, "\tt:%d\t%s\n", tokens, cmdutil.Cyan(inspirationFile))

				}
			})
		},
	}
}
