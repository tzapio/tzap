package config

import (
	"context"

	"github.com/tzapio/tzap/pkg/types/openai"
)

type configKey struct{}

type Configuration struct {
	OpenAIModel    string
	AutoMode       bool
	TruncateLimit  int
	MD5Rewrites    bool
	MD5IncludeList []string
	EnableLogs     bool
	LoggerOutput   string
}

var defaultConfig = Configuration{
	OpenAIModel:    openai.GPT3Dot5Turbo,
	AutoMode:       false,
	TruncateLimit:  0,
	MD5Rewrites:    false,
	MD5IncludeList: []string{""},
	EnableLogs:     false,
	LoggerOutput:   "",
}

func NewContext(ctx context.Context, config Configuration) context.Context {
	return context.WithValue(ctx, configKey{}, withDefaults(config))
}

func FromContext(ctx context.Context) Configuration {
	if config, ok := ctx.Value(configKey{}).(Configuration); ok {
		return config
	}

	return defaultConfig
}

// withDefaults merges the provided configuration with the default configuration,
// and returns a new configuration object with default values for any fields that are not provided.
func withDefaults(userConfig Configuration) Configuration {
	defaults := defaultConfig
	if userConfig.OpenAIModel == "" {
		userConfig.OpenAIModel = defaults.OpenAIModel
	}
	if userConfig.MD5IncludeList == nil {
		userConfig.MD5IncludeList = defaults.MD5IncludeList
	}
	return Configuration{
		OpenAIModel:    userConfig.OpenAIModel,
		AutoMode:       userConfig.AutoMode || defaults.AutoMode,
		TruncateLimit:  userConfig.TruncateLimit,
		MD5Rewrites:    userConfig.MD5Rewrites || defaults.MD5Rewrites,
		MD5IncludeList: userConfig.MD5IncludeList,
		EnableLogs:     userConfig.EnableLogs || defaults.EnableLogs,
		LoggerOutput:   userConfig.LoggerOutput,
	}
}
