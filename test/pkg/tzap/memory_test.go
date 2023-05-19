package tzap_test

import (
	"testing"

	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
)

func Test_Memory_Memorize(t *testing.T) {
	tt := tzap.InternalNew()
	memKey := "testKey"
	tt.Memory("root", memKey)
	if tzap.GetMemory(memKey) != "" {
		t.Errorf("Expected empty string for memory key '%s'", memKey)
	}
	tt.Message.Content = "test data"
	tt.Memorize(memKey)
	if tzap.GetMemory(memKey) != "test data" {
		t.Errorf("Expected 'test data' for memory key '%s'", memKey)
	}
}

func Test_Memory_MemorizeReq(t *testing.T) {
	tt := tzap.NewWithConnector(func() (types.TGenerator, config.Configuration) { return nil, config.Configuration{AutoMode: true} })
	tt.TG = &mockTG{}
	memKey := "testReqKey"
	tt.Memory("root", memKey)
	if tzap.GetMemory(memKey) != "" {
		t.Errorf("Expected empty string for memory key '%s'", memKey)
	}

	tt.
		AddAssistantMessage("requested data").MemorizeReq(memKey)

	if tzap.GetMemory(memKey) != "r=assistant;c=requested data" {
		t.Errorf("Expected 'r=assistant;c=requested data' for memory key '%s'", tzap.GetMemory(memKey))
	}
}
