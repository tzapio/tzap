package action_test

import (
	"testing"

	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/tzapaction/action"
	"github.com/tzapio/tzap/pkg/tzapaction/actionpb"
	mock_types "github.com/tzapio/tzap/test/mocks/pkg/types"
	"go.uber.org/mock/gomock"
)

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	tg := mock_types.NewMockTGenerator(ctrl)
	tzap := tzap.InjectNew(tg, config.DefaultConfig)

	expectedContent := "hello world"

	expectedFile := "pkg/tzapaction/action/create.go"
	tg.EXPECT().
		CountTokens(gomock.Any(), gomock.Any()).
		Return(1, nil).
		AnyTimes()

	tg.EXPECT().
		GenerateChat(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(types.CompletionMessage{
			Content:      expectedContent,
			FinishReason: "stop",
		}, nil)

	request := &actionpb.EditRequest{
		EditArgs: &actionpb.EditArgs{
			Mission: "Write to file",
			Task:    "Add tests that cover different scenarios and edge cases for the Create function",
			Plan:    "Add the tests using the testing package and cover both success cases and failure cases",
			FileIn:  expectedFile,
			Code: `
				// Write your test code here
			`,
		},
	}

	response, err := action.Create(tzap, request)
	if err != nil {
		t.Errorf("Error creating file: %v", err)
	}
	if response.FileWrites[0].Fileout != expectedFile {
		t.Errorf("Expected file %s, but got %s", expectedFile, response.FileWrites[0].Fileout)
	}

	if response.FileWrites[0].Contentout != expectedContent {
		t.Errorf("Expected content\n%s\nbut got\n%s", expectedContent, response.FileWrites[0].Contentout)
	}
}
