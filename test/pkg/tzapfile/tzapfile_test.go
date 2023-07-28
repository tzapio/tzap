package tzapfile_test

import (
	"testing"

	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzapfile"
)

func TestDeserializeMessageThread(t *testing.T) {
	content := `
@role:user
Hello!

---


@role:assistant
Hi there!---
\---
---hello
---


@role:user
Can you help me with something?

`

	expected := []types.Message{
		{Role: "user", Content: "Can you help me with something?"},
		{Role: "assistant", Content: "Hi there!---\n---\n---hello"},
		{Role: "user", Content: "Hello!"},
	}

	messages := tzapfile.DeserializeMessageThread(content)

	if len(messages) != len(expected) {
		t.Errorf("Expected %d messages, but got %d", len(expected), len(messages))
	}

	for i, msg := range expected {
		if messages[i].Role != msg.Role {
			t.Errorf("Expected role to be '%s', but got '%s'", msg.Role, messages[i].Role)
		}
		if messages[i].Content != msg.Content {
			t.Errorf("Expected content to be '%s', but got '%s'", msg.Content, messages[i].Content)
		}
	}
}

func TestSerializeMessageThread(t *testing.T) {
	messages := []types.Message{

		{Role: "user", Content: "Hello!"},
		{Role: "assistant", Content: "Hi there!---\n---\nhello---\n---hello"},
		{Role: "user", Content: "Can you help me with something?"},
	}

	expected := `

---
@role:user
Can you help me with something?
---
@role:assistant
Hi there!---
\---
hello---
\---hello
---
@role:user
Hello!
`

	result, err := tzapfile.SerializeMessageThread(messages)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result != expected {
		t.Errorf("Expected serialized message thread to be '%s', but got '%s'", expected, result)
	}
}
