package googlevoiceconnector

import (
	"fmt"
	"os"

	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzapconnect"
)

// AddToTGenerator adds GoogleTgenerator to an existing types.TGenerator and returns a new types.TGenerator.
// How to use: tzapconnector.NewTzap()
func AddToTGenerator(tg types.TGenerator) (types.TGenerator, error) {
	googleSecret, err := getGoogleSecretFromEnv()
	if err != nil {
		return nil, err
	}
	googleT := InitiateGoogleClient(googleSecret)
	partialComposite, ok := tg.(tzapconnect.PartialComposite)
	if !ok {
		return nil, fmt.Errorf("failed to cast TGenerator to PartialComposite")
	}

	partialComposite.VoiceGenerator = googleT
	return partialComposite, nil
}

func getGoogleSecretFromEnv() (string, error) {
	apiKey := os.Getenv("GOOGLE_SECRET")

	if apiKey == "" {
		return "", fmt.Errorf("GOOGLE_SECRET environment variable not set")
	}

	return apiKey, nil
}
