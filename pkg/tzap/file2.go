package tzap

import (
	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/types/openai"
	"github.com/tzapio/tzap/pkg/util/stdin"
)

func applyChanges(t *Tzap, editedContent string) *Tzap {
	config := config.FromContext(t.C)
	filepathValue := t.Data["filepath"].(string)

	autoMode := config.AutoMode
	makeChange := autoMode
	if !autoMode {
		if err := stdin.ConfirmAndApplyChanges(filepathValue, editedContent, stdin.ApplyChanges); err == nil {
			makeChange = true
		}
	}
	if makeChange {
		stdin.ApplyChanges(filepathValue, editedContent)
		writeMessageMD5(filepathValue, t)
		data := types.MappedInterface{
			"filepath": filepathValue,
			"content":  editedContent,
		}
		withEditFile := t.HijackTzap(&Tzap{Name: "withEditFile", Message: types.Message{
			Role:    openai.ChatMessageRoleAssistant,
			Content: editedContent,
		}, Data: data})

		return withEditFile
	}
	return t
}
