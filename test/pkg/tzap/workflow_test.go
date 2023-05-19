package tzap_test

import (
	"testing"

	"github.com/tzapio/tzap/pkg/types/openai"
	"github.com/tzapio/tzap/pkg/tzap"
)

func Test_ApplyWorkflow_RootTzap_ReturnTzapWithAppliedWorkflow(t *testing.T) {
	// Arrange
	rootTzap := tzap.InternalNew()
	expectedMessage := "Hello Workflow!"
	workflow := tzap.InternalNew().AddUserMessage(
		expectedMessage,
	)

	// Act
	result := rootTzap.ApplyWorkflowP(workflow)

	// Assert
	if result.Message.Role != openai.ChatMessageRoleUser {
		t.Errorf("Expected Role 'user'")
	}
	if result.Message.Content != "Hello Workflow!" {
		t.Errorf("Expected Content 'Hello Workflow!'")
	}
}

func Test_ApplyWorkflowFN_RootTzap_ReturnTzapWithAppliedWorkflow(t *testing.T) {
	// Arrange
	rootTzap := tzap.InternalNew()
	expectedMessage := "Hello Workflow!"
	workflowFunc := func(t *tzap.Tzap) *tzap.Tzap {
		appended := t.AddAssistantMessage(expectedMessage)
		return appended
	}

	// Act
	result := rootTzap.ApplyWorkflowFN(workflowFunc)

	// Assert
	if result.Message.Role != openai.ChatMessageRoleAssistant {
		t.Errorf("Expected Role '" + openai.ChatMessageRoleAssistant + "'")
	}
	if result.Message.Content != "Hello Workflow!" {
		t.Errorf("Expected Content 'Hello Workflow!'")
	}
}
