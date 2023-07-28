package action

import (
	"os"

	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/tzapaction/actionpb"
	"github.com/tzapio/tzap/pkg/util"
)

func Edit(t *tzap.Tzap, request *actionpb.EditRequest) (*actionpb.EditResponse, error) {
	if _, err := os.Stat(request.EditArgs.FileIn); os.IsNotExist(err) {
		return Create(t, request)
	}
	if request.EditArgs.FileOut == "" {
		request.EditArgs.FileOut = request.EditArgs.FileIn
	}
	contentIn := util.ReadFileP(request.EditArgs.FileIn)
	t = t.
		AddSystemMessage(
			"The mission is: "+request.EditArgs.Mission,
			"The plan is to: "+request.EditArgs.Plan,
			"The current task is: "+request.EditArgs.Task,
			"You will take the following file and apply the users changes on it and output this whole file. You need to think twice on the implementation, you need to fulfil the suggestion, plan, and current prompt, over just taking changes for what they are. Your output will be written to file so it is important to ONLY write the filecontents.",
			"You are editing the file:", request.EditArgs.FileIn,
			"####file: "+request.EditArgs.FileIn+"\n"+contentIn).
		AddUserMessage("Make the following change:" + request.EditArgs.Code).
		RequestChatCompletion()
	contentOut := t.Data["content"].(types.CompletionMessage).Content

	return &actionpb.EditResponse{
		FileWrites: []*actionpb.FileWrite{{
			Filein:     request.EditArgs.FileIn,
			Contentin:  contentIn,
			Fileout:    request.EditArgs.FileOut,
			Contentout: contentOut,
		}}}, nil
}
