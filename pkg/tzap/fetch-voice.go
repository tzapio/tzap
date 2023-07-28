package tzap

import (
	"github.com/tzapio/tzap/pkg/types"
)

// RequestTextToSpeech requests synthesized speech using specific (google voices) language and voice.
// It returns a pointer to a new Tzap containing the synthesised speech 'audioContent'.
func (t *Tzap) RequestTextToSpeech(language string, voice string) *ErrorTzap {
	completionMessage := t.Data["content"].(types.CompletionMessage)
	audioContent, err := t.TG.TextToSpeech(t.C, completionMessage.Content, language, voice)
	if err != nil {
		return t.ErrorTzap(err)
	}
	data := types.MappedInterface{
		"audioContent": audioContent,
	}
	withRequestGoogleVoice := t.AddTzap(&Tzap{Name: "withRequestGoogleVoice", Data: data})
	return withRequestGoogleVoice.ErrorTzap(nil)
}
func (t *Tzap) RequestTextifySpeech(audioContent *[]byte, language string) *ErrorTzap {
	text, err := t.TG.SpeechToText(t.C, audioContent, language)
	if err != nil {
		return t.ErrorTzap(err)
	}
	data := types.MappedInterface{
		"spoken": text,
	}
	withRequestGoogleVoice := t.AddTzap(&Tzap{Name: "withRequestGoogleVoice", Data: data})
	return withRequestGoogleVoice.ErrorTzap(nil)
}
