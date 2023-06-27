package tzapconnect

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

func LoadOPENAI_API_KEY() (string, error) {
	key, err := loadAPIKey("OPENAI_APIKEY")
	if err != nil {
		return loadAPIKey("OPENAI_API_KEY")
	}
	return key, err
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
			if len(splits) == 2 && strings.EqualFold(splits[0], key) {
				return strings.ReplaceAll(strings.TrimSpace(splits[1]), "\"", ""), nil
			}
		}
	}
	return "", errors.New("API key not found. Add " + key + "=<key> to environment variable or .env file")
}
