package types

type Message struct {
	Role    string
	Content string
}
type CompletionMessage struct {
	Content      string
	FunctionCall *FunctionCall `json:"function_call,omitempty"`
	FinishReason FinishReason  `json:"finish_reason,omitempty"`
}

type FunctionCall struct {
	Name      string `json:"name,omitempty"`
	Arguments string `json:"arguments,omitempty"`
}
type FinishReason string

const (
	FinishReasonStop          FinishReason = "stop"
	FinishReasonLength        FinishReason = "length"
	FinishReasonFunctionCall  FinishReason = "function_call"
	FinishReasonContentFilter FinishReason = "content_filter"
	FinishReasonNull          FinishReason = "null"
)

type MappedInterface map[string]interface{}
