package tzapconnect

import (
	"bufio"
	"errors"
	"os"
	"path"
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
	if apiKey, err := loadEnvFromFile(".env", key); err == nil {
		return apiKey, nil
	}

	if homeDir, err := os.UserHomeDir(); err == nil {
		homeTzapPath := path.Join(homeDir, ".tzap", ".env")
		if apiKey, err := loadEnvFromFile(homeTzapPath, key); err == nil {
			return apiKey, nil
		}
	}

	return "", errors.New("API key not found. Add " + key + "=<key> to environment variable or .env file")
}
func loadEnvFromFile(filePath, key string) (string, error) {
	envData, err := os.ReadFile(filePath)
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
	return "", errors.New("APIkey not found in file: " + filePath)
}
