package cmd

import (
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
	"github.com/tzapio/tzap/cli/cmd/cmdutil"
	"github.com/tzapio/tzap/internal/logging/tl"
	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/types/openai"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/tzapaction/cliworkflows"
	"github.com/tzapio/tzap/pkg/tzapconnect"

	"github.com/tzapio/tzap/pkg/util/stdin"
)

var tzapCliSettings struct {
	Model       string
	Temperature float32

	AutoMode      bool
	TruncateLimit int
	ConfigPath    string
	MD5Rewrites   bool
	DisableLogs   bool
	LoggerOutput  string
	Verbose       bool
	ApiMode       bool
	Yes           bool
	DisableIndex  bool
	Editor        string
	EmbeddingURL  string
	CompletionURL string
}

var RootCmd = &cobra.Command{
	Use:   "tzap",
	Short: "Tzap Cli",

	PersistentPreRunE: preRun,
}

func preRun(cmd *cobra.Command, args []string) error {
	tl.Logger.Println("Cobra CLI Root start")

	if tzapCliSettings.Verbose {
		tl.EnableLogger()
		tl.EnableUICompletionLogger()
		tl.EnableUILogger()
	}

	if isInitOrHelp(cmd.Name()) {
		return nil
	}

	baseDir, err := cmdutil.SearchForTzapincludeAndGetRootDir()
	if err != nil {
		println("Warning: No .tzapinclude file found. Run 'tzap init'. Using current directory as root.", err)
	} else {
		os.Chdir(baseDir)
	}

	err = initializeConfig()
	if err != nil {
		return err
	}

	tl.Logger.Println("Current working directory:", baseDir)
	t, err := initializeTzap()
	if err != nil {
		return err
	}
	cmd.SetContext(cliworkflows.SetCLIWorkflowConfigInContext(cmd.Context(), &cliworkflows.CLIWorkflowConfig{
		DisableIndex: tzapCliSettings.DisableIndex,
		Yes:          tzapCliSettings.Yes,
	}))
	tl.Logger.Println("Tzap initialized")
	cmd.SetContext(cmdutil.SetTzapInContext(cmd.Context(), t))
	return nil
}

func isInitOrHelp(command string) bool {
	return command == "init" || command == "help" || command == "install"
}

func initializeConfig() error {
	data, err := os.ReadFile(".tzap-data/config.json")
	if err == nil {
		var cfg map[string]interface{}
		if err := json.Unmarshal(data, &cfg); err == nil {
			if editor, ok := cfg["editor"].(string); ok {
				tzapCliSettings.Editor = editor
			}
		}
	} else {
		tl.Logger.Println("No config.json found")
		os.WriteFile(".tzap-data/config.json", []byte(`{"editor":"stdin"}`), 0644)
		tzapCliSettings.Editor = "stdin"
	}
	return nil
}

func initializeTzap() (*tzap.Tzap, error) {
	cfg := createConfigFromSettings()

	apikey, err := tzapconnect.LoadOPENAI_API_KEY()
	if err != nil {
		choice := stdin.ConfirmPrompt("Cannot find OPENAI_APIKEY.\n\nWould you like to add it now?")
		if choice {
			apikey = stdin.GetStdinInput("Enter OPENAI_APIKEY to save to .env:\n")
			saveAPIKey(apikey)
		} else {
			println("Aborted, cannot continue without OPENAI_APIKEY.")
			os.Exit(1)
			return nil, err
		}
	}
	connector := tzapconnect.WithConfig(apikey, cfg)

	t := tzap.NewWithConnector(connector)
	return t, nil
}

func createConfigFromSettings() config.Configuration {
	return config.Configuration{
		OpenAIModel:   modelMap[tzapCliSettings.Model],
		AutoMode:      tzapCliSettings.Yes, // automode == yes
		TruncateLimit: tzapCliSettings.TruncateLimit,
		MD5Rewrites:   tzapCliSettings.MD5Rewrites,
		EnableLogs:    !tzapCliSettings.DisableLogs,
		LoggerOutput:  tzapCliSettings.LoggerOutput,
		Temperature:   tzapCliSettings.Temperature,
		EmbeddingURL:  tzapCliSettings.EmbeddingURL,
		CompletionURL: tzapCliSettings.CompletionURL,
	}
}

func saveAPIKey(apikey string) {
	f, err := os.OpenFile(".env", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	if _, err := f.Write([]byte("\nOPENAI_APIKEY=" + apikey)); err != nil {
		panic(err)
	}
	if err := f.Close(); err != nil {
		panic(err)
	}
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var modelMap map[string]string = map[string]string{
	"gpt35":    openai.GPT3Dot5Turbo,
	"gpt356":   openai.GPT3Dot5Turbo0613,
	"gpt3516":  openai.GPT16,
	"gpt3516k": openai.GPT16,
	"gpt16":    openai.GPT16,
	"gpt4":     openai.GPT4,
	"gpt32":    openai.GPT432K,
}

func init() {
	RootCmd.CompletionOptions.HiddenDefaultCmd = true
	tzapCliSettings.MD5Rewrites = true

	RootCmd.PersistentFlags().StringVarP(&tzapCliSettings.Model, "model", "m", "gpt16", "Define what openai model to use. (Available gpt35 gpt356 (june model) gpt3516 (alias gpt16) gpt4).")
	RootCmd.PersistentFlags().StringVarP(&tzapCliSettings.CompletionURL, "baseurl", "b", "", "Completion URL")
	RootCmd.PersistentFlags().StringVar(&tzapCliSettings.EmbeddingURL, "embeddingbaseurl", "", "Embedding URL")
	RootCmd.PersistentFlags().StringVar(&tzapCliSettings.LoggerOutput, "loggeroutput", ".tzap-data/logs/", "Path and name of the log file.")
	RootCmd.PersistentFlags().Float32VarP(&tzapCliSettings.Temperature, "temperature", "t", 1.0, "Temperature for the interaction.")
	RootCmd.PersistentFlags().BoolVarP(&tzapCliSettings.Verbose, "verbose", "v", false, "Enable verbose logging")
	RootCmd.PersistentFlags().BoolVar(&tzapCliSettings.ApiMode, "api", false, "ALPHA: Enable clean stdout outputs. Also turns off editor mode.")
	RootCmd.PersistentFlags().BoolVarP(&tzapCliSettings.Yes, "yes", "y", false, "Answer yes on CLI related prompts - cost or similar related questions")
	RootCmd.Flags().BoolVarP(&tzapCliSettings.DisableIndex, "disableindex", "d", false,
		"For large projects disabling indexing speeds up the process.")
}
