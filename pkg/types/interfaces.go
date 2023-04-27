package types

type ITzap[Tzap any] interface {
	//
	AddUserMessage(contents ...string) Tzap
	AddAssistantMessage(contents ...string) Tzap
	AddSystemMessage(contents ...string) Tzap
	SetHeader(header string) Tzap

	LoadFileDir(dir string, ext string) Tzap
	ChangeFilepath(filepath string) Tzap
	FetchTask() Tzap
	LoadTask(filepath string) Tzap

	LoadTaskOrRequestNewTask(filepath string) Tzap
	LoadTaskOrRequestNewTaskMD5(filepath string) Tzap

	RequestChat() Tzap
	RequestTextToSpeech(language, voice string) Tzap

	//Tzap Primitives
	New() Tzap
	AddTzap(tc Tzap) Tzap
	HijackTzap(tc Tzap) Tzap

	//ControlFlow
	BubbleTzap(fn func(t Tzap)) Tzap
	IsolatedTzap(fn func(t Tzap)) Tzap
	MutationTzap(fn func(t Tzap) Tzap) Tzap
	Map(func(Tzap) Tzap) Tzap
	Accumulate(func(Tzap) Tzap) Tzap
	Exit() Tzap

	ApplyTemplate(template Tzap) Tzap
	ApplyTemplateFN(nt func(Tzap) Tzap) Tzap

	Recursive(func(tzapThatCreatesNewChildren Tzap) Tzap) Tzap
	CheckAndHandleGlobalOccurrences(references int, filename string, noOccurrence, handleOccurrence func(Tzap) Tzap) Tzap
	CheckAndHandleRecurrences(references int, filename string, noReccurance, handleRecurrence func(Tzap) Tzap) Tzap
	FileMustContainHandleGlobalOccurrences(references int, filename string, noOccurrence, handleOccurrence func(Tzap) Tzap) Tzap
}
