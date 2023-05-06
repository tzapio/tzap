package types

type ITzap[T any, Z any] interface {
	//
	AddUserMessage(contents ...string) T
	AddAssistantMessage(contents ...string) T
	AddSystemMessage(contents ...string) T
	SetHeader(header string) T

	LoadFileDir(dir string, ext string) T
	ChangeFilepath(filepath string) T
	FetchTask() T

	LoadTask(filepath string) T
	LoadTaskOrRequestNewTask(filepath string) T
	LoadTaskOrRequestNewTaskMD5(filepath string) T

	RequestChat() T
	RequestTextToSpeech(language, voice string) T

	//Tzap Primitives
	New() T
	AddTzap(tc T) T
	HijackTzap(tc T) T

	//ControlFlow
	BubbleTzap(fn func(t T)) T
	IsolatedTzap(fn func(t T)) T
	MutationTzap(fn func(t T) T) T
	Map(func(T) T) T
	Accumulate(func(T) T) T
	Exit() T

	ApplyTemplate(nt NamedTemplate[T, Z]) T
	ApplyTemplateFN(nt func(T) T) T
	ApplyTemplateP(T) T

	Recursive(func(tzapThatCreatesNewChildren T) T) T
	CheckAndHandleGlobalOccurrences(references int, filename string, noOccurrence, handleOccurrence func(T) T) T
	CheckAndHandleRecurrences(references int, filename string, noReccurance, handleRecurrence func(T) T) T
	FileMustContainHandleGlobalOccurrences(references int, filename string, noOccurrence, handleOccurrence func(T) T) T
}
type NamedTemplate[T any, Z any] struct {
	Name     string
	Template func(T) Z // Assuming the function takes no arguments and returns no values
}
