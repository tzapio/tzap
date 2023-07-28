package action

import (
	"context"
	"os"

	"github.com/tzapio/tzap/cli/cmd/cmdinstance"
	"github.com/tzapio/tzap/internal/logging/tl"
	"github.com/tzapio/tzap/pkg/embed"
	"github.com/tzapio/tzap/pkg/embed/embedstore"
	"github.com/tzapio/tzap/pkg/project"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/tzapaction/actionpb"
	"github.com/tzapio/tzap/pkg/tzapaction/cliworkflows"
	"github.com/tzapio/tzap/pkg/tzapaction/cliworkflows/threadworkflows"
	"github.com/tzapio/tzap/pkg/util/singlewait"
	"github.com/tzapio/tzap/workflows/code/embedworkflows"
)

func Search(t *tzap.Tzap, request *actionpb.SearchRequest) (*actionpb.SearchResponse, error) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	projectP, err := cmdinstance.LoadProject(cwd, request.SearchArgs.Lib)
	if err != nil {
		return nil, err
	}

	resultT := t.AddContextChange(func(c context.Context) context.Context {
		return project.SetProjectInContext(c, projectP)
	}).
		ApplyWorkflow(loadAndSearchEmbeddingsWorkflow(request.SearchArgs))
	searchResult := embedstore.TightenSearchResults(resultT.Data["searchResults"].(types.SearchResults).Results)
	embeddings := make([]*actionpb.Embedding, 0, len(searchResult.Results))
	for _, result := range searchResult.Results {
		embeddings = append(embeddings, &actionpb.Embedding{Content: result.Vector.Metadata.SplitPart,
			File:      result.Vector.Metadata.Filename,
			LineStart: int32(result.Vector.Metadata.LineStart)})
	}
	return &actionpb.SearchResponse{Embeddings: embeddings}, nil
}

func loadAndSearchEmbeddingsWorkflow(args *actionpb.SearchArgs) types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap] {
	return types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap]{
		Name: "loadAndSearchEmbeddings",
		Workflow: func(t *tzap.Tzap) *tzap.Tzap {
			queryWait := singlewait.New(func() types.QueryRequest {
				tl.Logger.Println("loadAndSearchEmbeddings: Getting query")
				query, err := embed.NewQuery(t, args.SearchQuery)
				if err != nil {
					panic(err)
				}
				tl.Logger.Println("loadAndSearchEmbeddings: Query received")
				return query
			})

			return t.
				ApplyWorkflow(cliworkflows.IndexFilesAndEmbeddings()).
				ApplyWorkflow(embedworkflows.SearchFilesWorkflow(queryWait.GetData(), args.ExcludeFiles, int(args.EmbedsCount)))
		},
	}
}
func SearchWorkflow(args *actionpb.SearchArgs) types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap] {
	return types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap]{
		Workflow: func(t *tzap.Tzap) *tzap.Tzap {
			result, err := Search(t, &actionpb.SearchRequest{SearchArgs: args})
			if err != nil {
				panic(err)
			}
			return t.
				ApplyWorkflow(cliworkflows.PrintEmbeddings(result.Embeddings)).
				ApplyWorkflow(threadworkflows.EmbedWorkflow(result.Embeddings))
		},
	}
}
