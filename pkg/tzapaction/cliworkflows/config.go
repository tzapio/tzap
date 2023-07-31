package cliworkflows

import "context"

var cliConfigContextKey = struct{ cliConfigContextKey string }{}

type CLIWorkflowConfig struct {
	Usd          float64
	Yes          bool
	DisableIndex bool
}

func SetCLIWorkflowConfigInContext(ctx context.Context, config *CLIWorkflowConfig) context.Context {
	return context.WithValue(ctx, cliConfigContextKey, config)
}
func GetCLIWorkflowConfigFromContext(ctx context.Context) *CLIWorkflowConfig {
	value := ctx.Value(cliConfigContextKey)
	if value == nil {
		return nil
	}
	return value.(*CLIWorkflowConfig)
}
