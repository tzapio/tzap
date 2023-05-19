package tzap_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
)

func Test_ErrorTzap_ErrorNotHandled_ReturnsError(t *testing.T) {
	// Given
	tzapObj := tzap.InternalNew()
	tzapObj.TG = &mockTG{}
	err := tzap.HandlePanic(func() {
		tzapObj = tzapObj.
			AddTzap(&tzap.Tzap{
				InitialSystemContent: "validHeader",
				Name:                 "Mock",
				Message:              types.Message{Role: "mocked", Content: "Hello!"},
				Data: types.MappedInterface{
					"filepath": "validFilePath",
				},
			}).
			ErrorTzap(fmt.Errorf("MOCK ERROR")).
			HandleError(func(et *tzap.ErrorTzap) error {
				return et.Err
			}).
			AddSystemMessage("Hello!")
	})

	if tzapObj.Message.Content == "Hello!" {
		t.Errorf("Expected tzapObj Message Content not to be 'Hello!'")
	}
	if err == nil {
		t.Errorf("Expected err to be not nil, but got nil")
	}
	if err != nil && !strings.HasSuffix(err.Error(), "MOCK ERROR") {
		t.Errorf("Expected err to be 'MOCK ERROR', but got '%s'", err.Error())
	}
}
