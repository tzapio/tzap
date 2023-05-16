package tzapconnect

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

func LoadOPENAI_APIKEY() (string, error) {
	return loadAPIKey("OPENAI_APIKEY")
}

func loadAPIKey(key string) (string, error) {
	// Try to get API key from environment variable.
	apiKey := os.Getenv(key)
	if apiKey != "" {
		return apiKey, nil
	}

	// Try to get API key from .env file.
	envData, err := os.ReadFile(".env")
	if err == nil {
		scanner := bufio.NewScanner(strings.NewReader(string(envData)))
		for scanner.Scan() {
			line := scanner.Text()
			splits := strings.SplitN(line, "=", 2)
			if len(splits) == 2 && splits[0] == key {
				return strings.TrimSpace(splits[1]), nil
			}
		}
	}
	return "", errors.New("API key not found. Add " + key + "=<key> to environment variable or .env file")
}
