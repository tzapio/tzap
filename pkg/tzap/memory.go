package tzap

import "github.com/tzapio/tzap/pkg/types"

// Mem represents a global memory storage for messages.
var Mem = map[string]*types.Message{}

// Memory adds a new Tzap tied to a memory key.
func (t *Tzap) Memory(role, key string) *Tzap {
	data := map[string]interface{}{}
	data["memory"] = key
	Mem[key] = &types.Message{Role: role, Content: ""}
	return t.AddTzap(&Tzap{Name: "memoryTzap", Data: data})
}

// Memorize stores the current Tzap message content under a given key.
func (t *Tzap) Memorize(key string) *Tzap {
	oldMemory := Mem[key]
	Mem[key] = &types.Message{Role: oldMemory.Role, Content: t.Message.Content}
	println("Memorized ", Mem[key].Content)
	return t
}

// MemorizeReq stores the Chat message content under a given key.
func (t *Tzap) MemorizeReq(key string) *Tzap {
	oldMemory := Mem[key]
	Mem[key] = &types.Message{Role: oldMemory.Role, Content: t.RequestChatCompletion().Data["content"].(types.CompletionMessage).Content}
	println("Memorized ", Mem[key].Content)
	return t
}

// GetMemory returns the content of a memory key, or an empty string if the key does not exist.
func GetMemory(key string) string {
	memory := Mem[key]
	if memory == nil {
		return ""
	}
	return memory.Content
}
