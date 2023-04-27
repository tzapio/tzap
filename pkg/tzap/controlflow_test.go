package tzap_test

import (
	"testing"

	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
)

func Test_MutationTzap_givenTzap_expectMutated(t *testing.T) {
	tt := tzap.InternalNew()
	tt.Message = types.Message{Role: "userqqqqqqqqq", Content: "Hello"}
	mutatedTzap := tt.MutationTzap(func(t *tzap.Tzap) *tzap.Tzap {
		t.Message.Content = "Mutated"

		return t
	})

	if mutatedTzap.Message.Content != "Mutated" {
		t.Errorf("Expected content to be 'Mutated', but got '%s'", mutatedTzap.Message.Content)
	}
	if mutatedTzap.Message.Role != "userqqqqqqqqq" {
		t.Errorf("Expected role to be 'userqqqqqqqqq', but got '%s'", mutatedTzap.Message.Role)
	}

}

func Test_BubbleTzap_givenTzap_expectUnchanged(t *testing.T) {
	tt := tzap.InternalNew()
	tt.Message = (types.Message{Role: "user", Content: "Hello"})
	bubbledTzap := tt.BubbleTzap(func(t *tzap.Tzap) {
		t.Message.Content = "Changed2" // should this usecase be allowed
	})

	if bubbledTzap.Message.Content != "Changed2" {
		t.Errorf("Expected content to be 'Changed2', but got '%s'", bubbledTzap.Message.Content)
	}
}

func Test_IsolatedTzap_givenTzap_expectOriginalUnchanged(t *testing.T) {
	tt := tzap.InternalNew()
	tt.Message = types.Message{Role: "user", Content: "Hello"}
	afterIsolation := tt.IsolatedTzap(func(t *tzap.Tzap) {
		t.Message.Content = "Changed"
	})

	if tt.Message.Content != "Hello" {
		t.Errorf("Expected original content to be 'Hello', but got '%s'", tt.Message.Content)
	}
	if afterIsolation.Message.Content != tt.Message.Content {
		t.Errorf("Expected content to be '%s', but got '%s'", tt.Message.Content, afterIsolation.Message.Content)
	}
}

func Test_Exit_givenTzap_expectPanic(t *testing.T) {
	tt := tzap.InternalNew()
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected a panic, but did not get one")
			}
		}()
		tt.Exit()
	}()
}

func Test_Map_givenChildrenTzaps_expectMappedChildren(t *testing.T) {
	tt := tzap.InternalNew()
	children := []*tzap.Tzap{
		tzap.InternalNew().AddUserMessage("Child1"),
		tzap.InternalNew().AddUserMessage("Child2"),
		tzap.InternalNew().AddUserMessage("Child3"),
	}
	expected := []string{"Child1", "Child2", "Child3"}
	tt.Data["children"] = children

	mapped := tt.Map(func(t *tzap.Tzap) *tzap.Tzap {
		return t.AddUserMessage(t.Message.Content + " Mapped")
	})

	mappedChildren := mapped.Data["children"].([]*tzap.Tzap)
	for i, child := range mappedChildren {
		expectedContent := expected[i] + " Mapped"
		if child.Message.Content != expectedContent {
			t.Errorf("Expected content to be '%s', but got '%s'", expectedContent, child.Message.Content)
		}
	}
}
