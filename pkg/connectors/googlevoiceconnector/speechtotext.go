package googlevoiceconnector

import (
	"context"
	"fmt"
	"log"

	speech "cloud.google.com/go/speech/apiv1"
	"cloud.google.com/go/speech/apiv1/speechpb"
	"google.golang.org/api/option"
)

func getSpeechClient(secret string) *speech.Client {
	ctx := context.Background()
	//
	client, err := speech.NewClient(ctx, option.WithTelemetryDisabled(), option.WithCredentialsJSON([]byte(secret)))
	if err != nil {
		log.Fatalf("Failed to create texttospeech client: %v", err)
		panic(err)
	}
	return client
}
func (g *GoogleTgenerator) SpeechToText(ctx context.Context, audioContent *[]byte, language string) (string, error) {
	req := &speechpb.RecognizeRequest{
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{
				Content: *audioContent,
			},
		},
		Config: &speechpb.RecognitionConfig{LanguageCode: language},
	}
	resp, err := g.speechtotextClient.Recognize(ctx, req)
	if err != nil {
		return "nil", err
	}
	_ = resp
	fmt.Printf("%+v\n", resp)
	println(resp.Results[0].Alternatives[0].Transcript)
	return resp.Results[0].Alternatives[0].Transcript, nil
}
