package types

type ITzap[T any, Z any] interface {
	AddUserMessage(contents ...string) T
	AddAssistantMessage(contents ...string) T
	AddSystemMessage(contents ...string) T
	SetInitialSystemContent(content string) T

	LoadFileDir(dir string) T
	LoadFiles(filepaths []string) T
	ChangeFilepath(filepath string) T

	LoadCompletion(filepath string) T
	LoadCompletionOrRequestCompletionMD5(filePath string) T
	LoadCompletionOrRequestCompletion(filepath string) T

	RequestChatCompletion() T
	StoreCompletion(filePath string) T
	RequestTextToSpeech(language, voice string) T

	//Tzap Primitives
	New() T
	AddTzap(tc T) T
	HijackTzap(tc T) T
	CloneTzap(tc T) T

	//ControlFlow
	WorkTzap(fn func(t T)) T
	IsolatedTzap(fn func(t T)) T
	MutationTzap(fn func(t T) T) T
	Map(func(t T) T) T
	Accumulate(func(t T) T) T
	Exit() T

	ApplyWorkflow(nt NamedWorkflow[T, Z]) T
	ApplyErrorWorkflow(nt NamedWorkflow[T, Z], fn func(t Z) error) T
	ApplyWorkflowFN(nt func(t T) T) T
	ApplyWorkflowP(T) T

	Recursive(func(tzapThatCreatesNewChildren T) T) T
	CheckAndHandleGlobalOccurrences(references int, filename string, noOccurrence, handleOccurrence func(T) T) T
	CheckAndHandleRecurrences(references int, filename string, noReccurance, handleRecurrence func(T) T) T
	FileMustContainHandleGlobalOccurrences(references int, filename string, noOccurrence, handleOccurrence func(T) T) T

	CountTokens(content string) (int, error)
	OffsetTokens(content string, from int, to int) (string, error)
}
type NamedWorkflow[T any, Z any] struct {
	Name     string
	Workflow func(t T) Z // Assuming the function takes no arguments and returns no values
}
