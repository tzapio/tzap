package action_test

import (
	"testing"

	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/tzapaction/action"
	"github.com/tzapio/tzap/pkg/tzapaction/actionpb"
	"github.com/tzapio/tzap/pkg/util"
	mock_types "github.com/tzapio/tzap/test/mocks/pkg/types"
	"go.uber.org/mock/gomock"
)

func TestRefactor(t *testing.T) {
	ctrl := gomock.NewController(t)
	tg := mock_types.NewMockTGenerator(ctrl)
	tzap := tzap.InjectNew(tg, config.DefaultConfig)
	tg.EXPECT().
		CountTokens(gomock.Any(), gomock.Any()).
		Return(1, nil).
		AnyTimes()
	tg.EXPECT().
		OffsetTokens(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return("", 50, nil).
		AnyTimes()
	expectedContentOut := `package action`
	tg.EXPECT().
		GenerateChat(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(types.CompletionMessage{
			Content:      expectedContentOut,
			FinishReason: "stop",
		}, nil)

	expectedFileIn := "pkg/tzapaction/action/refactor.go"
	expectedFileOut := expectedFileIn

	request := &actionpb.RefactorRequest{
		RefactorArgs: &actionpb.RefactorArgs{
			FileIn:       expectedFileIn,
			FileOut:      expectedFileOut,
			Mission:      "Improve code readability and maintainability",
			Task:         "Refactor code to use better variable names and remove duplication. Refactor the following file to be more readable. Make comments for the functions. Do not add any new public functions, only rewrite.",
			Plan:         "Do not write any text because this file will be saved directly to " + expectedFileOut,
			OutputFormat: "golang",
			Example: `package something
func doSomething(w http.ResponseWriter, r *http.Request, db *sql.DB) error {
	// function logic
}`,
			InspirationFiles: []string{},
		},
	}

	response, err := action.Refactor(tzap, request)
	if err != nil {
		t.Errorf("Error refactoring file: %v", err)
	}
	if len(response.FileWrites) != 1 {
		t.Errorf("Expected 1 file write, but got %d", len(response.FileWrites))
	}

	fileWrite := response.FileWrites[0]
	if fileWrite.Filein != expectedFileIn {
		t.Errorf("Expected file in %s, but got %s", expectedFileIn, fileWrite.Filein)
	}

	if fileWrite.Fileout != expectedFileOut {
		t.Errorf("Expected file out %s, but got %s", expectedFileOut, fileWrite.Fileout)
	}

	expectedContentIn := util.ReadFileP(expectedFileIn)

	if fileWrite.Contentin != expectedContentIn {
		t.Errorf("Expected content in\n%s\nbut got\n%s", expectedContentIn, fileWrite.Contentin)
	}

	if fileWrite.Contentout != expectedContentOut {
		t.Errorf("Expected content out\n%s\nbut got\n%s", expectedContentOut, fileWrite.Contentout)
	}
}
