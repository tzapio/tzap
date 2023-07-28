package tzap

import "github.com/tzapio/tzap/pkg/types"

func (t *Tzap) IfFunctionCall(tzapFunc func(*Tzap) *Tzap, notTzapFunc func(*Tzap) *Tzap) *Tzap {
	if t.Data["content"].(types.CompletionMessage).FunctionCall != nil {
		return tzapFunc(t)
	}
	return notTzapFunc(t)
}
