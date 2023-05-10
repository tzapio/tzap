package tzap

import (
	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/types/openai"
	"github.com/tzapio/tzap/pkg/util/stdin"
)

func applyChanges(t *Tzap, filePath, editedContent string) *Tzap {
	config := config.FromContext(t.C)

	autoMode := config.AutoMode
	makeChange := autoMode
	if !autoMode {
		if err := stdin.ConfirmAndApplyChanges(filePath, editedContent, stdin.ApplyChanges); err == nil {
			makeChange = true
		}
	}
	if makeChange {
		stdin.ApplyChanges(filePath, editedContent)
		writeMessageMD5(filePath, t)
		data := types.MappedInterface{
			"filepath": filePath,
			"content":  editedContent,
		}
		withEditFile := t.AddTzap(&Tzap{Name: "withEditFile", Message: types.Message{
			Role:    openai.ChatMessageRoleAssistant,
			Content: editedContent,
		}, Data: data})

		return withEditFile
	}
	return t
}
