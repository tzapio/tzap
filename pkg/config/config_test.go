package config_test

import (
	"context"
	"testing"

	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/types/openai"
)

func TestNewContext(t *testing.T) {
	// create a sample context
	ctx := context.TODO()

	// create a sample configuration
	cfg := config.Configuration{
		OpenAIModel:   openai.GPT16,
		AutoMode:      true,
		TruncateLimit: 100,
		EnableLogs:    true,
	}

	// create a new context with the configuration
	newCtx := config.NewContext(ctx, cfg)

	// retrieve the configuration from the new context
	newCfg := config.FromContext(newCtx)

	// check that the retrieved configuration is equal to the original configuration
	if newCfg != cfg {
		t.Fatalf("Expected %v, got %v", cfg, newCfg)
	}
}

func TestNewContext_DefaultConfiguration(t *testing.T) {
	// create a sample context
	ctx := context.TODO()

	// create a new context with the default configuration
	newCtx := config.NewContext(ctx, config.Configuration{Temperature: 1})

	// retrieve the configuration from the new context
	newCfg := config.FromContext(newCtx)

	// check that the retrieved configuration is equal to the default configuration
	if newCfg != config.DefaultConfig {
		t.Fatalf("Expected %v, got %v", config.DefaultConfig, newCfg)
	}
}

func TestFromContext(t *testing.T) {
	// create a sample context
	ctx := context.TODO()

	// create a sample configuration
	cfg := config.Configuration{
		OpenAIModel:   openai.GPT16,
		AutoMode:      true,
		TruncateLimit: 100,
		EnableLogs:    true,
	}

	// create a new context with the configuration
	newCtx := config.NewContext(ctx, cfg)

	// retrieve the configuration from the new context using FromContext
	newCfg := config.FromContext(newCtx)

	// check that the retrieved configuration is equal to the original configuration
	if newCfg != cfg {
		t.Fatalf("Expected %v, got %v", cfg, newCfg)
	}
}

func TestFromContext_DefaultConfiguration(t *testing.T) {
	// create a sample context
	ctx := context.TODO()

	// retrieve the configuration from the context using FromContext
	cfg := config.FromContext(ctx)

	// check that the retrieved configuration is equal to the default configuration
	if cfg != config.DefaultConfig {
		t.Fatalf("Expected %v, got %v", config.DefaultConfig, cfg)
	}
}

func TestWithDefaults(t *testing.T) {
	// create a sample configuration with some defaults
	cfg := config.Configuration{
		OpenAIModel:  openai.GPT3Dot5Turbo,
		AutoMode:     false,
		EnableLogs:   false,
		LoggerOutput: "",
		Temperature:  0,
	}

	// set some values to an empty configuration
	emptyCfg := config.Configuration{}
	emptyCfg.OpenAIModel = ""

	// add the defaults to the empty configuration
	defaultCfg := config.WithDefaults(emptyCfg)

	// check that the default configuration is equal to the original configuration with defaults
	if defaultCfg != cfg {
		t.Fatalf("Expected %v, got %v", cfg, defaultCfg)
	}
}
