package gocode

import (
	"fmt"

	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/util/stdin"
)

func DisplayAndConfirm() types.NamedWorkflow[*tzap.Tzap, *tzap.ErrorTzap] {
	return types.NamedWorkflow[*tzap.Tzap, *tzap.ErrorTzap]{
		Name: "DisplayAndConfirm",
		Workflow: func(t *tzap.Tzap) *tzap.ErrorTzap {

			ok := stdin.ConfirmPrompt("\nContinue with this commit?")

			if !ok {
				return t.ErrorTzap(fmt.Errorf("commit aborted by user"))
			}

			return t.ErrorTzap(nil)
		}}
}
