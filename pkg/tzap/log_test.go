package tzap_test

import (
	"context"
	"testing"
	"time"

	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/tzap"
)

func Test_Log_singleMessage_messageBuffered(t *testing.T) {
	ctx := context.Background()
	ctx = config.NewContext(ctx, config.Configuration{
		EnableLogs: true,
	})
	tt := &tzap.Tzap{C: ctx}
	tt.Name = "TestTzap"

	tzap.ResetFlush()
	tzap.Logf(tt, "Test message %d", 1)
	tzap.Logf(tt, "Test message %d", 1)
	tzap.Logf(tt, "Test message %d", 1)
	messageBufferLen := len(tzap.MessageBuffer)

	if messageBufferLen != 2 {
		t.Errorf("Expected message buffer length to be 1, got %d", messageBufferLen)
	}
}

func Test_Log_multipleMessages_messageBufferOverflow(t *testing.T) {
	ctx := context.Background()
	ctx = config.NewContext(ctx, config.Configuration{
		EnableLogs: true,
	})
	tt := &tzap.Tzap{C: ctx}
	tt.Name = "TestTzap"

	for i := 0; i < 20; i++ {
		tzap.Logf(tt, "Test message %d", i+1)
	}

	messageBufferLen := len(tzap.MessageBuffer)

	if messageBufferLen != 15 {
		t.Errorf("Expected message buffer length to be 15, got %d", messageBufferLen)
	}
}

func Test_Log_flush_messageBufferEmpty(t *testing.T) {
	tt := tzap.InternalNew()
	tt.Name = "TestTzap"

	for i := 0; i < 10; i++ {
		tzap.Logf(tt, "Test message %d", i+1)
	}

	tzap.Flush()

	messageBufferLen := len(tzap.MessageBuffer)

	if messageBufferLen != 0 {
		t.Errorf("Expected message buffer length to be 0, got %d", messageBufferLen)
	}
}

func Test_LastExecutionTime_limitedLogging(t *testing.T) {
	tt := tzap.InternalNew()
	tt.Name = "TestTzap"

	tzap.Log(tt, "Initial log")
	initialExecutionTime := tzap.LastExecutionTime

	time.Sleep(500 * time.Millisecond)
	tzap.Log(tt, "Second log within rate limit")
	secondExecutionTime := tzap.LastExecutionTime

	if initialExecutionTime != secondExecutionTime {
		t.Errorf("Expected initial execution time to be equal to second execution time, but they were not equal")
	}
}

func Test_LastExecutionTime_unlimitedLogging(t *testing.T) {
	ctx := context.Background()
	ctx = config.NewContext(ctx, config.Configuration{
		EnableLogs: true,
	})
	tt := &tzap.Tzap{C: ctx}
	tt.Name = "TestTzap"

	tzap.Log(tt, "Initial log")
	initialExecutionTime := tzap.LastExecutionTime

	time.Sleep(1500 * time.Millisecond)
	tzap.Log(tt, "Second log outside rate limit")
	secondExecutionTime := tzap.LastExecutionTime

	if initialExecutionTime == secondExecutionTime {
		t.Errorf("Expected initial execution time to be different from second execution time, but they were equal")
	}
}
