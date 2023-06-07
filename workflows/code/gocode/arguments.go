package gocode

import (
	"strings"

	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
)

func DeserializedArguments(dataname string, args []string) types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap] {
	return types.NamedWorkflow[*tzap.Tzap, *tzap.Tzap]{
		Name: "deserializedArguments",
		Workflow: func(t *tzap.Tzap) *tzap.Tzap {
			t.Data = make(types.MappedInterface)
			t.Data[dataname] = strings.Join(args, " ")
			return t
		}}
}
