package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tzapio/tzap/cli/cmd/cmdutil"
	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/tl"
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
}

var rootCmd = &cobra.Command{
	Use:     "tzap",
	Short:   "Tzap Cli!",
	Long:    `tbd`,
	Version: "v0.7.19",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if settings.Verbose {
			tl.EnableLogger()
		}
		//check subcommand if init or help
		if cmd.Name() == "init" || cmd.Name() == "help" {
			return nil
		}

		err := cmdutil.SearchForTzapincludeAndChangeDir()
		if err != nil {
			return err
		}
		// print cwd
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		tl.Logger.Println("Current working directory:", cwd)
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
		cmd.SetContext(cmdutil.SetTzapInContext(cmd.Context(), t))
		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var modelMap map[string]string = map[string]string{"gpt35": openai.GPT3Dot5Turbo, "gpt4": openai.GPT4}

func init() {
	rootCmd.PersistentFlags().StringVarP(&settings.Model, "model", "m", "gpt35", "Define what openai model to use. (Available gpt35 gpt4).")
	rootCmd.PersistentFlags().BoolVar(&settings.AutoMode, "automode", false, "Whether to press yes on continue prompts.")
	rootCmd.PersistentFlags().IntVar(&settings.TruncateLimit, "truncate", 0, "Truncate limit for the interaction.")
	rootCmd.PersistentFlags().BoolVar(&settings.MD5Rewrites, "md5rewrites", true, "For some functions, this flag enables overwriting files with the same MD5 hash.")
	rootCmd.PersistentFlags().StringVar(&settings.IncludeList, "include", "", "Files include MD5 matching pattern.")
	rootCmd.PersistentFlags().BoolVar(&settings.DisableLogs, "enablelogs", true, "Whether to enable logging.")
	rootCmd.PersistentFlags().StringVar(&settings.LoggerOutput, "loggeroutput", "./.tzap-data/logs", "Path and name of the log file.")
	rootCmd.PersistentFlags().BoolVar(&settings.Stub, "stub", false, "Test non-live mode")
	rootCmd.PersistentFlags().Float32VarP(&settings.Temperature, "temperature", "t", 1.0, "Temperature for the interaction.")
	rootCmd.PersistentFlags().BoolVarP(&settings.Verbose, "verbose", "v", false, "Enable verbose logging")
}
