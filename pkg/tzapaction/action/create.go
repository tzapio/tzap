package action

import (
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/tzapaction/actionpb"
	"github.com/tzapio/tzap/pkg/util/cleaner"
)

func Create(t *tzap.Tzap, request *actionpb.EditRequest) (*actionpb.EditResponse, error) {
	t = t.
		ApplyWorkflow(SearchWorkflow(&actionpb.SearchArgs{
			ExcludeFiles: []string{request.EditArgs.FileIn, request.EditArgs.FileOut},
			SearchQuery:  request.EditArgs.Mission + " " + request.EditArgs.Task + " " + request.EditArgs.Plan,
			EmbedsCount:  20,
		})).
		AddSystemMessage(
			"The mission is: "+request.EditArgs.Mission,
			"The plan is to: "+request.EditArgs.Plan,
			"The current task is: "+request.EditArgs.Task,
			"You will take the following file and apply the users changes on it and output this whole file. You need to think twice on the implementation, you need to fulfil the suggestion, plan, and current prompt, over just taking changes for what they are. Your output will be written to file so it is important to ONLY write the filecontents.",
			"You are creating the file:", request.EditArgs.FileIn).
		AddUserMessage("Add the following:" + request.EditArgs.Code).
		RequestChatCompletion()

	contentOut := t.Data["content"].(types.CompletionMessage).Content

	return &actionpb.EditResponse{
		FileWrites: []*actionpb.FileWrite{{
			Fileout:    request.EditArgs.FileIn,
			Contentout: cleaner.FileWriteClean(contentOut),
		}}}, nil
}
