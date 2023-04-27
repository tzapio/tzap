package googlevoiceconnector

import (
	"context"
	"fmt"
	"log"

	speech "cloud.google.com/go/speech/apiv1"
	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
	"github.com/tzapio/tzap/pkg/types"
	"google.golang.org/api/option"
)

type GoogleTgenerator struct {
	*types.UnimplementedTGenerator
	texttospeechClient *texttospeech.Client
	speechtotextClient *speech.Client
}

func InitiateGoogleClient(jsonSecret string) *GoogleTgenerator {
	return &GoogleTgenerator{texttospeechClient: getClient(jsonSecret), speechtotextClient: getSpeechClient(jsonSecret)}
}

func getClient(secret string) *texttospeech.Client {
	ctx := context.Background()

	client, err := texttospeech.NewClient(ctx, option.WithTelemetryDisabled(),
		option.WithCredentialsJSON([]byte(secret)))
	if err != nil {
		log.Fatalf("Failed to create texttospeech client: %v", err)
		panic(err)
	}
	return client
}
func (g *GoogleTgenerator) SynthesizeSpeech(outputText, language, voice string) (*[]byte, error) {
	// Set up the client
	ctx := context.Background()

	// Configure the request
	request := &texttospeechpb.SynthesizeSpeechRequest{
		Input: &texttospeechpb.SynthesisInput{

			InputSource: &texttospeechpb.SynthesisInput_Ssml{
				Ssml: outputText,
			},
		},
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: language,
			Name:         voice,
		},
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_MP3,
			SpeakingRate:  1.1,
		},
	}
	fmt.Printf("%+v\n", request)
	// Call the API
	response, err := g.texttospeechClient.SynthesizeSpeech(ctx, request)
	if err != nil {
		fmt.Print(err.Error())
		return nil, fmt.Errorf("failed to synthesize text %s", err)
	}
	return &response.AudioContent, nil

}
