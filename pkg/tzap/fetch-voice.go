package tzap

import (
	"github.com/tzapio/tzap/pkg/types"
)

// RequestTextToSpeech requests synthesized speech using specific (google voices) language and voice.
// It returns a pointer to a new Tzap containing the synthesised speech 'audioContent'.
func (t *Tzap) RequestTextToSpeech(language string, voice string) *Tzap {
	println(t.TG)
	audioContent, err := t.TG.TextToSpeech(t.C, t.Data["content"].(string), language, voice)
	if err != nil {
		panic(err)
	}
	data := types.MappedInterface{
		"audioContent": audioContent,
	}
	withRequestGoogleVoice := t.AddTzap(&Tzap{Name: "withRequestGoogleVoice", Data: data})
	return withRequestGoogleVoice
}
func (t *Tzap) RequestTextifySpeech(audioContent *[]byte, language string) *Tzap {
	text, err := t.TG.SpeechToText(t.C, audioContent, language)
	if err != nil {
		panic(err)
	}
	data := types.MappedInterface{
		"spoken": text,
	}
	withRequestGoogleVoice := t.AddTzap(&Tzap{Name: "withRequestGoogleVoice", Data: data})
	return withRequestGoogleVoice
}
