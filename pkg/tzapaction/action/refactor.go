package action

import (
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/tzapaction/actionpb"
	"github.com/tzapio/tzap/pkg/tzapaction/cliworkflows/codegeneration"
	"github.com/tzapio/tzap/pkg/util"
	"github.com/tzapio/tzap/pkg/util/cleaner"
)

func Refactor(t *tzap.Tzap, refactorRequest *actionpb.RefactorRequest) (*actionpb.RefactorResponse, error) {
	refactorArgs := refactorRequest.RefactorArgs
	if refactorArgs.FileOut == "" {
		refactorArgs.FileOut = refactorArgs.FileIn
	}

	t = t.
		ApplyWorkflow(SearchWorkflow(&actionpb.SearchArgs{
			ExcludeFiles: append([]string{refactorArgs.FileIn, refactorArgs.FileOut}, refactorArgs.InspirationFiles...),
			SearchQuery:  refactorArgs.Mission + " " + refactorArgs.Task + " " + refactorArgs.Plan,
			EmbedsCount:  20,
		})).
		ApplyWorkflowFN(codegeneration.MakeCode(refactorArgs))

	contentIn := util.ReadFileP(refactorArgs.FileIn)
	contentOut := t.Data["content"].(types.CompletionMessage).Content

	return &actionpb.RefactorResponse{
		FileWrites: []*actionpb.FileWrite{
			{
				Filein:     refactorArgs.FileIn,
				Contentin:  contentIn,
				Fileout:    refactorArgs.FileOut,
				Contentout: cleaner.FileWriteClean(contentOut),
			},
		},
	}, nil
}
