package output

import "github.com/tzapio/tzap/pkg/types"

func GetText(messages []types.Message) string {
	txt := ""
	for _, message := range messages {
		txt += "\n###" + message.Role + "\n" + message.Content
	}
	return txt
}
