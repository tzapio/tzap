package googlevoiceconnector

import (
	"fmt"
	"os"
)

func getGoogleSecretFromEnv() (string, error) {
	apiKey := os.Getenv("GOOGLE_SECRET")

	if apiKey == "" {
		return "", fmt.Errorf("GOOGLE_SECRET environment variable not set")
	}

	return apiKey, nil
}
