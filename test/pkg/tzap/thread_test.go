package tzap_test

import (
	"testing"

	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
)

func TestGetThreadWithinLimit(t *testing.T) {
	thread := []types.Message{
		{Role: "", Content: "a The truncated thread should match..."},     // 50 words mock
		{Role: "", Content: "b The truncated thread should match..."},     // 50 words mock
		{Role: "", Content: "c The truncated thread should match..."},     // 50 words mock
		{Role: "", Content: "d The truncated thread should match qeq..."}, // 50 words mock
		{Role: "", Content: "The truncated thread should match"},          // 50 words mock
		{Role: "", Content: "f The truncated thread should match"},        // 50 words mock
	}

	expectedResult := []types.Message{
		{Role: "", Content: "The truncated thread should match"},   // 3 words mock
		{Role: "", Content: "f The truncated thread should match"}, // 3 words mock
	}

	wordLimit := 120

	result := tzap.TruncateToMaxTokens(&mockTG{}, thread, wordLimit)
	if len(result) != len(expectedResult) {
		t.Errorf("Expected %d thread, but got %d", len(expectedResult), len(result))
		return
	}
	for i := range result {
		if result[i].Content != expectedResult[i].Content {
			t.Errorf("Expected message content '%s', but got '%s'", expectedResult[i].Content, result[i].Content)
		}
	}
}
