package tzap_test

import (
	"testing"

	"github.com/tzapio/tzap/pkg/types/openai"
	"github.com/tzapio/tzap/pkg/tzap"
)

func Test_ApplyTemplate_RootTzap_ReturnTzapWithAppliedTemplate(t *testing.T) {
	// Arrange
	rootTzap := tzap.InternalNew()
	expectedMessage := "Hello Template!"
	template := tzap.InternalNew().AddUserMessage(
		expectedMessage,
	)

	// Act
	result := rootTzap.ApplyTemplateP(template)

	// Assert
	if result.Message.Role != openai.ChatMessageRoleUser {
		t.Errorf("Expected Role 'user'")
	}
	if result.Message.Content != "Hello Template!" {
		t.Errorf("Expected Content 'Hello Template!'")
	}
}

func Test_ApplyTemplateFN_RootTzap_ReturnTzapWithAppliedTemplate(t *testing.T) {
	// Arrange
	rootTzap := tzap.InternalNew()
	expectedMessage := "Hello Template!"
	templateFunc := func(t *tzap.Tzap) *tzap.Tzap {
		appended := t.AddAssistantMessage(expectedMessage)
		return appended
	}

	// Act
	result := rootTzap.ApplyTemplateFN(templateFunc)

	// Assert
	if result.Message.Role != openai.ChatMessageRoleAssistant {
		t.Errorf("Expected Role '" + openai.ChatMessageRoleAssistant + "'")
	}
	if result.Message.Content != "Hello Template!" {
		t.Errorf("Expected Content 'Hello Template!'")
	}
}
