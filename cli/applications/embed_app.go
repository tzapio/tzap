package applications

import (
	"github.com/tzapio/tzap/cli/cmd/cliworkflows"
	"github.com/tzapio/tzap/pkg/embed"
	"github.com/tzapio/tzap/pkg/embed/localdb/singlewait"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/workflows/code/embedworkflows"
)

func LoadAndSearchEmbeddings(excludeFiles []string, searchQuery string, k int, n int, disableIndex, yes bool) types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap] {
	return types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap]{
		Name: "loadAndSearchEmbeddings",
		Workflow: func(t *tzap.Tzap) *tzap.Tzap {
			queryWait := singlewait.New(func() types.QueryRequest {
				query, err := embed.GetQuery(t, searchQuery)
				if err != nil {
					panic(err)
				}
				return query
			})

			return t.
				ApplyWorkflow(cliworkflows.IndexFilesAndEmbeddings("./", disableIndex, yes)).
				ApplyWorkflow(embedworkflows.EmbeddingInspirationWorkflow(queryWait.GetData(), excludeFiles, k, n))
		},
	}
}
