package config

import (
	"context"

	"github.com/tzapio/tzap/pkg/types/openai"
)

type configKey struct{}

type Configuration struct {
	OpenAIModel   string
	EmbedModel    string
	CompletionURL string
	EmbeddingURL  string
	AutoMode      bool
	TruncateLimit int
	MD5Rewrites   bool
	EnableLogs    bool
	LoggerOutput  string
	Temperature   float32
}

var DefaultConfig = Configuration{
	OpenAIModel:   openai.GPT3Dot5Turbo,
	AutoMode:      false,
	TruncateLimit: 0,
	MD5Rewrites:   false,
	EnableLogs:    false,
	LoggerOutput:  "",
	Temperature:   1.0,
}

func NewContext(ctx context.Context, config Configuration) context.Context {
	return context.WithValue(ctx, configKey{}, WithDefaults(config))
}

func FromContext(ctx context.Context) Configuration {
	if config, ok := ctx.Value(configKey{}).(Configuration); ok {
		return config
	}
	return DefaultConfig
}

// withDefaults merges the provided configuration with the default configuration,
// and returns a new configuration object with default values for any fields that are not provided.
func WithDefaults(userConfig Configuration) Configuration {
	defaults := DefaultConfig
	if userConfig.OpenAIModel == "" {
		userConfig.OpenAIModel = defaults.OpenAIModel
	}
	return Configuration{
		OpenAIModel:   userConfig.OpenAIModel,
		AutoMode:      userConfig.AutoMode || defaults.AutoMode,
		TruncateLimit: userConfig.TruncateLimit,
		MD5Rewrites:   userConfig.MD5Rewrites || defaults.MD5Rewrites,
		EnableLogs:    userConfig.EnableLogs || defaults.EnableLogs,
		LoggerOutput:  userConfig.LoggerOutput,
		Temperature:   userConfig.Temperature,
	}
}
