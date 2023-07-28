package action

import (
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/tzapaction/actionpb"
)

func Implement(t *tzap.Tzap, writeCodeRequest *actionpb.ImplementRequest) (*actionpb.ImplementResponse, error) {
	var fileWrites []*actionpb.FileWrite
	for _, change := range writeCodeRequest.ImplementArgs.Changes {
		response, err := Edit(t, &actionpb.EditRequest{
			EditArgs: &actionpb.EditArgs{
				Mission: writeCodeRequest.ImplementArgs.Mission,
				Plan:    writeCodeRequest.ImplementArgs.Plan,
				Task:    change.Task,
				FileIn:  change.FileIn,
				FileOut: change.FileOut,
			},
		})
		if err != nil {
			continue
		}
		fileWrites = append(fileWrites, response.FileWrites...)
	}
	return &actionpb.ImplementResponse{FileWrites: fileWrites}, nil
}
