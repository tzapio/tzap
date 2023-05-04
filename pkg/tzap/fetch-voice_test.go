package tzap_test

import (
	"testing"

	"github.com/tzapio/tzap/pkg/tzap"
)

func Test_RequestTextToSpeech_SynthesizeSpeech(t *testing.T) {
	rootTzap := tzap.InternalNew()
	rootTzap.Data["content"] = "Hello world!"

	// Use the fake implementation of TG
	rootTzap.TG = &mockTG{}

	// Adjust these values to proper sample values for your testing needs.
	sampleLanguage := "en-US"
	sampleVoice := "en-US-Wavenet-A"

	resultTzap := rootTzap.RequestTextToSpeech(sampleLanguage, sampleVoice).HandleError(func(et *tzap.ErrorTzap) *tzap.Tzap {
		t.Errorf("ErrorTzap should not be called")
		return nil
	})
	audioContent, ok := resultTzap.Data["audioContent"].(*[]byte)

	if !ok {
		t.Errorf("Failed to convert 'audioContent' to []byte")
	}
	if audioContent == nil {
		t.Errorf("AudioContent should not be nil")
		return
	}
	if len(*audioContent) == 0 {
		t.Errorf("AudioContent should not be empty")
	}
}

func Test_RequestSpeechToText_TextifySpeech(t *testing.T) {
	rootTzap := tzap.InternalNew()
	sampleAudioContent := []byte("sample audio content")
	rootTzap.Data["audioContent"] = &sampleAudioContent

	// Use the fake implementation of TG
	rootTzap.TG = &mockTG{}

	// Adjust this value to proper sample value for your testing needs.
	sampleLanguage := "en-US"

	resultTzap := rootTzap.RequestTextifySpeech(&sampleAudioContent, sampleLanguage).HandleError(func(et *tzap.ErrorTzap) *tzap.Tzap {
		t.Errorf("ErrorTzap should not be called")
		return nil
	})
	spokenText, ok := resultTzap.Data["spoken"].(string)

	if !ok {
		t.Errorf("Failed to convert 'spoken' to string")
	}
	if spokenText == "" {
		t.Errorf("Spoken text should not be empty string")
	}
}
