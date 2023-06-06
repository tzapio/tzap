package cliworkflows

import (
	"fmt"
	"os"
	"strings"

	"github.com/tzapio/tzap/cli/cmd/cmdutil"
	"github.com/tzapio/tzap/cli/cmd/cmdutil/fileevaluator"
	"github.com/tzapio/tzap/internal/logging/tl"
	"github.com/tzapio/tzap/pkg/embed"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/util"
	"github.com/tzapio/tzap/workflows/code/embedworkflows"
)

func IndexZipFilesAndEmbeddings(name types.ProjectName, projectDir types.ProjectDir, zipURL string, disableIndex, yes bool) types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap] {
	return types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap]{
		Name: "indexFilesAndEmbeddings",
		Workflow: func(t *tzap.Tzap) *tzap.Tzap {
			if disableIndex {
				return t
			}
			fileInDirEvaluator, err := fileevaluator.New(name)
			if err != nil {
				panic(err)
			}
			fileStampsDB, err := embed.NewFilestampsDB(projectDir)
			if err != nil {
				panic(err)
			}
			files, err := fileInDirEvaluator.WalkDirFromURL(zipURL)
			if err != nil {
				panic(err)
			}
			embedder := embed.NewEmbedder(t, name, fileStampsDB, files)
			tl.Logger.Println("Indexing files...")
			tl.Logger.Println("Finished index files...")
			return t.ApplyWorkflow(embedworkflows.LoadAndFetchEmbeddings(name, files, embedder, disableIndex, yes))
		},
	}
}
func IndexFilesAndEmbeddings(dir string, disableIndex, yes bool) types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap] {
	return types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap]{
		Name: "indexFilesAndEmbeddings",
		Workflow: func(t *tzap.Tzap) *tzap.Tzap {
			if disableIndex {
				return t
			}
			fileEvaluator, err := fileevaluator.New(types.LOCALPROJECTNAME)
			if err != nil {
				panic(err)
			}
			filesStampsDB, err := embed.NewFilestampsDB("./.tzap-data")

			tl.Logger.Println("Indexing files...")
			files, err := fileEvaluator.WalkDir("./")
			if err != nil {
				panic(err)
			}
			embedder := embed.NewEmbedder(t, types.LOCALPROJECTNAME, filesStampsDB, files)
			tl.Logger.Println("Finished index files...")

			return t.ApplyWorkflow(embedworkflows.LoadAndFetchEmbeddings(types.LOCALPROJECTNAME, files, embedder, disableIndex, yes))
		},
	}
}
func PrintInspirationFiles(inspirationFiles []string) types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap] {
	return types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap]{
		Name: "listInspirationFiles",
		Workflow: func(t *tzap.Tzap) *tzap.Tzap {
			return t.WorkTzap(func(t *tzap.Tzap) {
				if len(inspirationFiles) == 0 {
					println(cmdutil.Bold("\nInspiration files: None (use --inspiration to add more)"))
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
func PrintSearchResults(searchResults types.SearchResults) types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap] {
	return types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap]{
		Name: "listInspirationFiles",
		Workflow: func(t *tzap.Tzap) *tzap.Tzap {
			return t.WorkTzap(func(t *tzap.Tzap) {
				tl.Logger.Println("Showing results")
				println("\nSearch result embeddings:")
				for _, result := range searchResults.Results {
					tokens, err := t.CountTokens(result.Vector.Metadata.SplitPart)
					if err != nil {
						panic(err)
					}
					fmt.Fprintf(os.Stderr, "\tt:%d\t%s\n", tokens, cmdutil.Cyan(cmdutil.FormatVectorToClickable(result.Vector)))
				}
				println()
			})
		},
	}
}
