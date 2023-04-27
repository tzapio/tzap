package splitter

import (
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
)

// GenericOutputter is a structure that defines a generic transformer and a generic callback.
type GenericOutputter[T interface{}] struct {
	genericTransformer func(*tzap.Tzap) []T
	genericCallback    func(int, *tzap.Tzap, T) *tzap.Tzap
}

// GenericSplitter takes a tzap.Tzap and applies the transformer and callback provided from the GenericOutputter.
// It returns an updated tzap.Tzap with data containing the children.
func (y GenericOutputter[T]) GenericSplitter(t *tzap.Tzap) *tzap.Tzap {
	children := make([]*tzap.Tzap, 0)
	for i, payload := range y.genericTransformer(t) {
		task := y.genericCallback(i, t, payload)
		children = append(children, task)
	}

	data := types.MappedInterface{"children": children}
	return t.HijackTzap(&tzap.Tzap{Name: "GenericSplitter", Data: data})
}

// NewGenericOutputter creates a new GenericOutputter with the provided parameters.
func NewGenericOutputter[T interface{}](
	genericSplitter func(*tzap.Tzap) []T,
	genericCallback func(int, *tzap.Tzap, T) *tzap.Tzap,
) GenericOutputter[T] {
	return GenericOutputter[T]{genericSplitter, genericCallback}
}
