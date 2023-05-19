package tzap_test

import (
	"testing"

	"github.com/tzapio/tzap/pkg/tzap"
)

func Test_ChangeFilepath_givenPath_expectChangedFilepath(t *testing.T) {
	tt := tzap.InternalNew()
	newFilepath := "/new/file/path"

	changedTzap := tt.ChangeFilepath(newFilepath)

	if changedTzap.Data["filepath"] != newFilepath {
		t.Errorf("Expected filepath to be changed to '%s', but got '%s'", newFilepath, changedTzap.Data["filepath"])
	}
}
