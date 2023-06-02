package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tzapio/tzap/cli/cmd/cmdutil"
	"github.com/tzapio/tzap/internal/logging/tl"
	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/types/openai"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/tzapconnect"
	"github.com/tzapio/tzap/pkg/tzapconnect/stubconnector"
)

var settings struct {
	Model         string
	AutoMode      bool
	TruncateLimit int
	ConfigPath    string
	MD5Rewrites   bool
	IncludeList   string
	DisableLogs   bool
	LoggerOutput  string
	Stub          bool
	Temperature   float32
	Verbose       bool
	ApiMode       bool
	Yes           bool
}

var RootCmd = &cobra.Command{
	Use:   "tzap",
	Short: "Tzap Cli",

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		tl.Logger.Println("Cobra CLI Root start")
		if settings.Verbose {
			tl.EnableLogger()
			tl.EnableUICompletionLogger()
			tl.EnableUILogger()
		}
		//check subcommand if init or help
		if cmd.Name() == "init" || cmd.Name() == "help" {
			return nil
		}

		root, err := cmdutil.SearchForTzapincludeAndGetRootDir()
		if err != nil {
			return err
		}
		tl.Logger.Println("Current working directory:", root)
		os.Chdir(root)
		config := config.Configuration{
			OpenAIModel:   modelMap[settings.Model],
			AutoMode:      settings.AutoMode,
			TruncateLimit: settings.TruncateLimit,
			MD5Rewrites:   settings.MD5Rewrites,
			EnableLogs:    !settings.DisableLogs,
			LoggerOutput:  settings.LoggerOutput,
			Temperature:   settings.Temperature,
		}
		var connector types.TzapConnector
		if settings.Stub {
			connector = stubconnector.StubWithConfig(config)
		} else {
			apikey, err := tzapconnect.LoadOPENAI_APIKEY()
			if err != nil {
				return err
			}
			connector = tzapconnect.WithConfig(apikey, config)
		}

		t := tzap.NewWithConnector(connector)
		tl.Logger.Println("Tzap initialized")
		cmd.SetContext(cmdutil.SetTzapInContext(cmd.Context(), t))
		return nil
	},
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var modelMap map[string]string = map[string]string{"gpt35": openai.GPT3Dot5Turbo, "gpt4": openai.GPT4}

func init() {
	RootCmd.PersistentFlags().StringVarP(&settings.Model, "model", "m", "gpt35", "Define what openai model to use. (Available gpt35 gpt4).")
	RootCmd.PersistentFlags().BoolVar(&settings.AutoMode, "automode", false, "Whether to press yes on continue prompts.")
	RootCmd.PersistentFlags().IntVar(&settings.TruncateLimit, "truncate", 0, "Truncate limit for the interaction.")
	RootCmd.PersistentFlags().BoolVar(&settings.MD5Rewrites, "md5rewrites", true, "For some functions, this flag enables overwriting files with the same MD5 hash.")
	RootCmd.PersistentFlags().StringVar(&settings.IncludeList, "include", "", "Files include MD5 matching pattern.")
	RootCmd.PersistentFlags().BoolVar(&settings.DisableLogs, "disablelogs", false, "Whether to disable logging.")
	RootCmd.PersistentFlags().StringVar(&settings.LoggerOutput, "loggeroutput", "./.tzap-data/logs", "Path and name of the log file.")
	RootCmd.PersistentFlags().BoolVar(&settings.Stub, "stub", false, "Test non-live mode")
	RootCmd.PersistentFlags().Float32VarP(&settings.Temperature, "temperature", "t", 1.0, "Temperature for the interaction.")
	RootCmd.PersistentFlags().BoolVarP(&settings.Verbose, "verbose", "v", false, "Enable verbose logging")
	RootCmd.PersistentFlags().BoolVar(&settings.ApiMode, "api", false, "Enable clean stdout outputs")
	RootCmd.PersistentFlags().BoolVarP(&settings.Yes, "yes", "y", false, "Answer yes on prompts")
}
