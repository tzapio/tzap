package tzap_test

import (
	"testing"

	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
)

func TestGetMessagesWithinLimit(t *testing.T) {
	messages := []types.Message{
		{Role: "", Content: "a The truncated messages should match..."},     // 3 words
		{Role: "", Content: "b The truncated messages should match..."},     // 3 words
		{Role: "", Content: "c The truncated messages should match..."},     // 3 words
		{Role: "", Content: "d The truncated messages should match qeq..."}, // 3 words
		{Role: "", Content: "The truncated messages should match"},          // 3 words
		{Role: "", Content: "f The truncated messages should match"},        // 3 words
	}

	expectedResult := []types.Message{
		{Role: "", Content: "The truncated messages should match"},   // 3 words
		{Role: "", Content: "f The truncated messages should match"}, // 3 words
	}

	wordLimit := 12

	result := tzap.TruncateToMaxWords(messages, wordLimit)
	if len(result) != len(expectedResult) {
		t.Errorf("Expected %d messages, but got %d", len(expectedResult), len(result))
		return
	}
	for i := range result {
		if result[i].Content != expectedResult[i].Content {
			t.Errorf("Expected message content '%s', but got '%s'", expectedResult[i].Content, result[i].Content)
		}
	}
}
