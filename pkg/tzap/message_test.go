package tzap_test

import (
	"os"
	"testing"

	"github.com/tzapio/tzap/pkg/types/openai"
	"github.com/tzapio/tzap/pkg/tzap"
)

func Test_AddUserMessage_SingleString_AddedContent(t *testing.T) {
	original := tzap.InternalNew()
	result := original.AddUserMessage("Hello world!")

	if result == nil {
		t.Errorf("Expected a non-nil result, but got nil")
		return
	}

	if result.Message.Content != "Hello world!" {
		t.Errorf("Expected message content 'Hello world!', but got %v", result.Message.Content)
	}
	if result.Message.Role != openai.ChatMessageRoleUser {
		t.Errorf("Expected message role 'user', but got %v", result.Message.Role)
	}
}

func Test_AddUserMessage_MultipleStrings_AddedContent(t *testing.T) {
	original := tzap.InternalNew()
	result := original.AddUserMessage("Hello world!", "How are you?")

	if result == nil {
		t.Errorf("Expected a non-nil result, but got nil")
		return
	}
	if result.Message.Content != "Hello world!\nHow are you?" {
		t.Errorf("Expected message content 'Hello world!\nHow are you?', but got %v", result.Message.Content)
	}
	if result.Message.Role != openai.ChatMessageRoleUser {
		t.Errorf("Expected message role 'user', but got %v", result.Message.Role)
	}
}

func Test_LoadUserMessageFromFileOrStdinInput_FileExists_LoadedContent(t *testing.T) {
	filepath := "temp_testfile.txt"
	testContent := "Example file content."

	err := os.WriteFile(filepath, []byte(testContent), 0600)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(filepath)

	original := tzap.InternalNew()
	result := original.LoadUserMessageFromFileOrStdinInput(filepath, "Example file content task")

	if result == nil {
		t.Errorf("Expected a non-nil result, but got nil")
		return
	}
	if result.Message.Content != testContent {
		t.Errorf("Expected message content '%v', but got %v", testContent, result.Message.Content)
	}
	if result.Message.Role != openai.ChatMessageRoleUser {
		t.Errorf("Expected message role 'user', but got %v", result.Message.Role)
	}
}

func Test_AddSystemMessage_SingleString_AddedContent(t *testing.T) {
	original := tzap.InternalNew()
	result := original.AddSystemMessage("System message example")

	if result == nil {
		t.Errorf("Expected a non-nil result, but got nil")
		return
	}
	if result.Message.Content != "System message example" {
		t.Errorf("Expected message content 'System message example', but got %v", result.Message.Content)
	}
	if result.Message.Role != openai.ChatMessageRoleSystem {
		t.Errorf("Expected message role 'system', but got %v", result.Message.Role)
	}
}

func Test_AddSystemMessage_MultipleStrings_AddedContent(t *testing.T) {
	original := tzap.InternalNew()
	result := original.AddSystemMessage("System message example", "Another message")

	if result == nil {
		t.Errorf("Expected a non-nil result, but got nil")
		return
	}
	if result.Message.Content != "System message example\nAnother message" {
		t.Errorf("Expected message content 'System message example\nAnother message', but got %v", result.Message.Content)
	}
	if result.Message.Role != openai.ChatMessageRoleSystem {
		t.Errorf("Expected message role 'system', but got %v", result.Message.Role)
	}
}

func Test_AddAssistantMessage_SingleString_AddedContent(t *testing.T) {
	original := tzap.InternalNew()
	result := original.AddAssistantMessage("Assistant message example")

	if result == nil {
		t.Errorf("Expected a non-nil result, but got nil")
		return
	}
	if result.Message.Content != "Assistant message example" {
		t.Errorf("Expected message content 'Assistant message example', but got %v", result.Message.Content)
	}
	if result.Message.Role != openai.ChatMessageRoleAssistant {
		t.Errorf("Expected message role 'assistant', but got %v", result.Message.Role)
	}
}

func Test_AddAssistantMessage_MultipleStrings_AddedContent(t *testing.T) {
	original := tzap.InternalNew()
	result := original.AddAssistantMessage("Assistant message example", "Another message")

	if result == nil {
		t.Errorf("Expected a non-nil result, but got nil")
		return
	}
	if result.Message.Content != "Assistant message example\nAnother message" {
		t.Errorf("Expected message content 'Assistant message example\nAnother message', but got %v", result.Message.Content)
	}
	if result.Message.Role != openai.ChatMessageRoleAssistant {
		t.Errorf("Expected message role 'assistant', but got %v", result.Message.Role)
	}
}

func Test_AppendMessage_Content_Appended(t *testing.T) {
	original := tzap.InternalNew().AddUserMessage("Hello world!")
	result := original.AppendMessage("Testing append")

	if result == nil {
		t.Errorf("Expected a non-nil result, but got nil")
		return
	}
	if result.Message.Content != "Hello world! Testing append" {
		t.Errorf("Expected message content 'Hello world! Testing append', but got %v", result.Message.Content)
	}
}

func Test_AppendContent_SeparatorAndStrings_ContentAppended(t *testing.T) {
	original := tzap.InternalNew()
	separator := ", "
	strings := []string{"foo", "bar", "baz"}

	for _, s := range strings {
		original.AppendContent(separator, s)
	}

	if original.Message.Content != "foo, bar, baz" {
		t.Errorf("Expected message content 'foo, bar, baz', but got %v", original.Message.Content)
	}
}

func Test_PrependContent_SeparatorAndStrings_ContentPrepended(t *testing.T) {
	original := tzap.InternalNew()
	separator := ", "
	strings := []string{"foo", "bar", "baz"}

	for _, s := range strings {
		original.PrependContent(separator, s)
	}

	if original.Message.Content != "baz, bar, foo" {
		t.Errorf("Expected message content 'baz, bar, foo', but got %v", original.Message.Content)
	}
}

func Test_CombineMessage_TwoMessageFunctions_CombinedContent(t *testing.T) {
	original := tzap.InternalNew()
	fn1 := func(t *tzap.Tzap) *tzap.Tzap { return t.AddUserMessage("Hello world!") }
	fn2 := func(t *tzap.Tzap) *tzap.Tzap { return t.AddAssistantMessage("How can I help you?") }

	result := original.CombineMessage(fn1, fn2)

	if result == nil {
		t.Errorf("Expected a non-nil result, but got nil")
		return
	}
	if result.Message.Content != "Hello world!\nHow can I help you?" {
		t.Errorf("Expected message content 'Hello world!\nHow can I help you?', but got %v", result.Message.Content)
	}
	if result.Message.Role != openai.ChatMessageRoleUser {
		t.Errorf("Expected message role 'user', but got %v", result.Message.Role)
	}
}
