package cliworkflows

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/tzapio/tzap/cli/cmd/cmdutil"
	"github.com/tzapio/tzap/internal/logging/tl"
	"github.com/tzapio/tzap/pkg/embed"
	"github.com/tzapio/tzap/pkg/project"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/tzapaction/actionpb"
	"github.com/tzapio/tzap/pkg/util"
	"github.com/tzapio/tzap/pkg/util/stdin"
	"github.com/tzapio/tzap/workflows/code/embedworkflows"
)

func LoadAndFetchEmbeddings(files []types.FileReader, embedder *embed.Embedder, yes bool, usd float64) types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap] {
	return types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap]{
		Name: "loadAndFetchEmbeddings",
		Workflow: func(t *tzap.Tzap) *tzap.Tzap {
			return t.MutationTzap(func(t *tzap.Tzap) *tzap.Tzap {
				return t.
					ApplyWorkflow(embedworkflows.PrepareEmbedFilesWorkflow(files, embedder)).
					ApplyWorkflow(ConfirmEmbeddingSearch(yes, usd)).
					ApplyWorkflow(embedworkflows.FetchOrCachedEmbeddingForFilesWorkflow(files)).
					ApplyWorkflow(embedworkflows.SaveAndLoadEmbeddingsToDB())
			})
		},
	}
}

func ConfirmEmbeddingSearch(yes bool, usd float64) types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap] {
	return types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap]{
		Name: "confirmEmbeddingSearch",
		Workflow: func(t *tzap.Tzap) *tzap.Tzap {
			uncachedEmbeddings, ok := t.Data["uncachedEmbeddings"].(*types.Embeddings)
			if !ok {
				panic("Loading embeddings went wrong")
			}
			if len(uncachedEmbeddings.Vectors) > 19 {
				price := float64(len(uncachedEmbeddings.Vectors)*400) * 0.0001 / 1000
				if !yes {
					ok := stdin.ConfirmPrompt(fmt.Sprintf(
						"Embeddings - You are about to fetch %d embeddings. Proceed? Estimation tokens: %d. Price is: $0.0004 per 1000 tokens. Estimating %.4f USD",
						len(uncachedEmbeddings.Vectors),
						len(uncachedEmbeddings.Vectors)*400,
						price))
					if !ok {
						fmt.Print("Aborting...")
						os.Exit(1)
					}
				} else {
					if price > usd {
						fmt.Print("Price too high", price, ">", usd)
						os.Exit(1)
					}
				}
			}
			return t
		},
	}
}

func IndexZipFilesAndEmbeddings(name project.ProjectName, projectDir project.ProjectDir, zipURL string, disableIndex bool, yes bool, usd float64) types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap] {
	return types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap]{
		Name: "indexFilesAndEmbeddings",
		Workflow: func(t *tzap.Tzap) *tzap.Tzap {
			if disableIndex {
				return t
			}
			projectP := project.GetProjectFromContext(t.C)
			projectP.GetTimestampCache().StartInit()
			projectP.GetEmbeddingsCache().StartInit()
			files, err := projectP.GetFiles()
			if err != nil {
				panic(err)
			}
			projectP.GetEmbeddingCollection().StartInit()
			embedder := embed.NewEmbedder(projectP.GetEmbeddingsCache(), projectP.GetTimestampCache())
			tl.Logger.Println("Indexing files...")

			return t.ApplyWorkflow(LoadAndFetchEmbeddings(files, embedder, yes, usd))
		},
	}
}

func IndexFilesAndEmbeddings() types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap] {
	return types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap]{
		Name: "indexFilesAndEmbeddings",
		Workflow: func(t *tzap.Tzap) *tzap.Tzap {
			config := GetCLIWorkflowConfigFromContext(t.C)
			var disableIndex bool
			var usd float64
			var yes bool
			if config != nil {
				disableIndex = config.DisableIndex
				usd = config.Usd
				yes = config.Yes
			} else {
				tl.Logger.Println("Warning: No CLIWorkflowConfig found in context")
				disableIndex = false
				usd = 0.001
				yes = false
			}

			projectP := project.GetProjectFromContext(t.C)
			if disableIndex || !projectP.CanIndex() {
				projectP.GetEmbeddingCollection().StartInit()
				return t
			}

			filesStampsDB := projectP.GetTimestampCache()
			filesStampsDB.StartInit()

			embeddingCacheDB := projectP.GetEmbeddingsCache()
			embeddingCacheDB.StartInit()

			projectP.GetEmbeddingCollection().StartInit()
			tl.Logger.Println("Indexing files...")
			files, err := projectP.GetFiles()
			if err != nil {
				panic(err)
			}

			embedder := embed.NewEmbedder(embeddingCacheDB, filesStampsDB)
			tl.Logger.Println("Finished index files...")
			if !disableIndex {
				println("Checking for file changes. " + cmdutil.Black("(use -d to disable this check)...\n"))
				return t.ApplyWorkflow(LoadAndFetchEmbeddings(files, embedder, yes, usd))
			}
			return t

		},
	}
}

func PrintInspirationFiles(inspirationFiles []string) types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap] {
	return types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap]{
		Name: "listInspirationFiles",
		Workflow: func(t *tzap.Tzap) *tzap.Tzap {
			return t.WorkTzap(func(t *tzap.Tzap) {
				if len(inspirationFiles) == 0 {
					println("\nInspiration files: None" + cmdutil.Black(" (use --inspiration to add more)"))
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

func PrintEmbeddings(embeddings []*actionpb.Embedding) types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap] {
	return types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap]{
		Name: "listInspirationFiles",
		Workflow: func(t *tzap.Tzap) *tzap.Tzap {
			return t.WorkTzap(func(t *tzap.Tzap) {
				tl.Logger.Println("Showing results")
				println("\nSearch result embeddings:")
				for _, embedding := range embeddings {
					tokens, err := t.CountTokens(embedding.Content)
					if err != nil {
						panic(err)
					}
					fmt.Fprintf(os.Stderr, "\t"+cmdutil.Black("t:%d")+"\t%s\n", tokens, cmdutil.Cyan(cmdutil.FormatVectorToClickable(embedding)))
				}
				println()
			})
		},
	}
}

func PrintFileDiff(compareToFile string) types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap] {
	return types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap]{
		Name: "printFileDiff",
		Workflow: func(t *tzap.Tzap) *tzap.Tzap {
			return t.WorkTzap(func(t *tzap.Tzap) {
				oldContent := ""
				if _, err := os.Stat(compareToFile); err == nil {
					oldContent = util.ReadFileP(compareToFile)
				}
				dmp := diffmatchpatch.New()
				completionMessage := t.Data["content"].(types.CompletionMessage)
				diffs := dmp.DiffPrettyText(dmp.DiffMain(oldContent, completionMessage.Content, false))
				println(diffs)
			})
		},
	}
}

func Print2VSCode(compareToFile *actionpb.FileWrite) types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap] {
	return types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap]{
		Name: "printFileDiff",
		Workflow: func(t *tzap.Tzap) *tzap.Tzap {
			return t.WorkTzap(func(t *tzap.Tzap) {
				temp, err := os.CreateTemp("", "tzapcontent*")
				if err != nil {
					panic(err)
				}

				if _, err := temp.Write([]byte(compareToFile.Contentout)); err != nil {
					panic(err)
				}

				if _, err := os.Stat(compareToFile.Fileout); err == nil {
					exec.Command("code", "-d", compareToFile.Fileout, temp.Name()).Run()
				} else {
					exec.Command("code", temp.Name()).Run()
				}
			})
		},
	}
}
