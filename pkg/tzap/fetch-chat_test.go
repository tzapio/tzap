package tzap_test

import (
	"os"
	"path"
	"testing"

	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
)

func Test_RequestChat_ValidFetch_OpenAIChatRequested(t *testing.T) {
	// Given
	tzapObj := tzap.InternalNew()
	tzapObj.TG = &mockTG{}
	tzapObj = tzapObj.AddTzap(&tzap.Tzap{
		Header:  "validHeader",
		Name:    "Mock",
		Message: types.Message{Role: "mocked", Content: "Hello!"},
		Data: types.MappedInterface{
			"filepath": "validFilePath",
		},
	})

	// When
	responseTzap := tzapObj.RequestChat()

	// Expect
	if responseTzap.Parent.Parent.Header == "validHeader" {
		t.Errorf("Expected header to be 'validHeader', but got '%s'", responseTzap.Parent.Parent.Header)
	}
	if responseTzap.Data["content"] != "r=system;c=validHeader|r=mocked;c=Hello!" {
		t.Errorf("Expected content to be 'r=system;c=validHeader|r=mocked;c=Hello!', but got '%s'", responseTzap.Data["content"])
	}
	if responseTzap.Parent.Data["filepath"] != "validFilePath" {
		t.Errorf("Expected filepath to be 'validFilePath!', but got '%s'", responseTzap.Parent.Data["validFilePath"])
	}
}

func Test_FetchTask_ValidFetch_ChangesApplied(t *testing.T) {
	// Given
	f, err := os.MkdirTemp("", "example")
	if err != nil {
		panic(err)
	}
	filename := path.Join(f, "validFilePath.txt")

	tzapObj := tzap.NewWithConnector(func() (types.TGenerator, config.Configuration) { return nil, config.Configuration{AutoMode: true} })

	tzapObj.TG = &mockTG{}
	tzapObj = tzapObj.AddTzap(&tzap.Tzap{
		Header:  "validHeader",
		Name:    "Mock",
		Message: types.Message{Role: "mocked", Content: "Hello!"},
		Data: types.MappedInterface{
			"filepath": filename,
			"content":  "someOldFile",
		},
	})

	// When
	responseTzap := tzapObj.FetchTask()

	// Expect
	if responseTzap.Header != "" {
		t.Errorf("Expected header to be '', but got '%s'", responseTzap.Header)
	}
	if responseTzap.Message.Content == "" {
		t.Errorf("Expected content to not be empty string")
	}
}
