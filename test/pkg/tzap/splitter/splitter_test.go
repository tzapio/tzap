package splitter_test

import (
	"testing"

	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/tzap/splitter"
)

func Test_GenericSplitter_TransformedItems_ReturnsExpectedData(t *testing.T) {
	root := tzap.InternalNew()

	transformerFn := func(t *tzap.Tzap) []string {
		return []string{"a", "b", "c"}
	}
	callbackFn := func(i int, t *tzap.Tzap, s string) *tzap.Tzap {
		tc := t.AddSystemMessage()
		tc.InitialSystemContent = s
		return tc
	}

	goObj := splitter.NewGenericOutputter(transformerFn, callbackFn)
	output := goObj.GenericSplitter(root)

	if output == nil {
		t.Error("should return a Tzap object")
		return
	}
	data, ok := output.Data["children"].([]*tzap.Tzap)
	if !ok {
		t.Error("should have children of type []*Tzap")
	}
	if len(data) != 3 {
		t.Errorf("should have 3 children, got %d", len(data))
	}
	if data[0].InitialSystemContent != "a" {
		t.Errorf("First child should have header a, got %s", data[0].InitialSystemContent)
	}
	if data[1].InitialSystemContent != "b" {
		t.Errorf("Second child should have header b, got %s", data[1].InitialSystemContent)
	}
	if data[2].InitialSystemContent != "c" {
		t.Errorf("Third child should have header c, got %s", data[2].InitialSystemContent)
	}
}
