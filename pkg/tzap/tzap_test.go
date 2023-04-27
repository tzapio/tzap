package tzap_test

import (
	"testing"

	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
)

func Test_NewTzap_Defaults_RootTzapWithEmptyMessage(t *testing.T) {
	tt := tzap.InternalNew()

	if tt.Name != "RootTzap" {
		t.Errorf("expected name to be RootTzap but got %s", tt.Name)
	}
	if tt.Header != "" {
		t.Errorf("expected header to be empty but got %s", tt.Header)
	}
	if tt.Message.Role != "" {
		t.Errorf("expected message role to be empty but got %s", tt.Message.Role)
	}
	if tt.Message.Content != "" {
		t.Errorf("expected message content to be empty but got %s", tt.Message.Content)
	}
	if tt.Data == nil {
		t.Errorf("expected data to be not nil")
	}
	if tt.C == nil {
		t.Errorf("expected context to be not nil")
	}
}

func Test_AddTzap_EmptyTzap_NewTzapAsChild(t *testing.T) {
	parent := tzap.InternalNew()
	child := &tzap.Tzap{}

	result := parent.AddTzap(child)

	if result != child {
		t.Errorf("expected returned value to be child but got %v", result)
	}
	if result.Parent != parent {
		t.Errorf("expected child parent to be parent but got %v", result.Parent)
	}
	if result.C == nil {
		t.Errorf("expected context to be not nil")
	}
}

func Test_CloneTzap_Defaults_ClonedTzapWithInitialProperties(t *testing.T) {
	parent := tzap.InternalNew()
	child := &tzap.Tzap{
		Name:   "ChildTzap",
		Header: "ChildHeader",
		Message: types.Message{
			Role:    "ChildRole",
			Content: "ChildContent",
		},
		Data: types.MappedInterface{
			"key": "value",
		},
	}

	clonedChild := parent.CloneTzap(child)

	if clonedChild.Name != child.Name {
		t.Errorf("expected cloned name to be %s but got %s", child.Name, clonedChild.Name)
	}
	if clonedChild.Header != child.Header {
		t.Errorf("expected cloned header to be %s but got %s", child.Header, clonedChild.Header)
	}
	if clonedChild.Message.Role != child.Message.Role {
		t.Errorf("expected cloned message role to be %s but got %s", child.Message.Role, clonedChild.Message.Role)
	}
	if clonedChild.Message.Content != child.Message.Content {
		t.Errorf("expected cloned message content to be %s but got %s", child.Message.Content, clonedChild.Message.Content)
	}
	if clonedChild.Data["key"] != child.Data["key"] {
		t.Errorf("expected cloned data to be %v but got %v", child.Data, clonedChild.Data)
	}
	if clonedChild.C == nil {
		t.Errorf("expected context to be not nil")
	}
}

func Test_HijackTzap_Defaults_HijackedTzapWithParentProperties(t *testing.T) {
	parent := tzap.InternalNew()
	child := &tzap.Tzap{
		Name:   "ChildTzap",
		Header: "ChildHeader",
		Message: types.Message{
			Role:    "ChildRole",
			Content: "ChildContent",
		},
		Data: types.MappedInterface{
			"key": "value",
		},
	}

	hijackedChild := parent.HijackTzap(child)

	if hijackedChild != child {
		t.Errorf("expected hijacked child to be child but got %v", hijackedChild)
	}
	if hijackedChild.C == nil {
		t.Errorf("expected context to be not nil")
	}
}
