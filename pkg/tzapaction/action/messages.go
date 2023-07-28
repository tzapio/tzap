package action

import (
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzapaction/actionpb"
)

// ToPBMessage converts Tzap messages to protobuf messages
func ToPBMessage(tzapMessages []types.Message) []*actionpb.Message {
	var messages []*actionpb.Message
	for _, m := range tzapMessages {
		messages = append(messages, &actionpb.Message{
			Role:    m.Role,
			Content: m.Content,
		})
	}
	return messages
}

// ToTzapMessage converts protobuf messages to Tzap messages
func ToTzapMessage(pbMessages []*actionpb.Message) []types.Message {
	var messages []types.Message
	for _, m := range pbMessages {
		messages = append(messages, types.Message{
			Role:    m.Role,
			Content: m.Content,
		})
	}
	return messages
}
