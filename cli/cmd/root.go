package cmd

import (
	"context"
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
	"github.com/tzapio/tzap/cli/cmd/cmdutil"
	"github.com/tzapio/tzap/cli/cmd/cmdutil/fileevaluator/cmdinstance"
	"github.com/tzapio/tzap/internal/logging/tl"
	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/project"
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
	Editor        string
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
		if cmd.Name() == "init" || cmd.Name() == "help" || cmd.Name() == "install" {
			return nil
		}

		baseDir, err := cmdutil.SearchForTzapincludeAndGetRootDir()
		if err != nil {
			return err
		}
		os.Chdir(baseDir)
		data, err := os.ReadFile(".tzap-data/config.json")
		if err == nil {
			var cfg map[string]interface{}
			if err := json.Unmarshal(data, &cfg); err == nil {
				if editor, ok := cfg["editor"].(string); ok {
					settings.Editor = editor
				}
			}
		} else {
			tl.Logger.Println("No config.json found")
			os.WriteFile(".tzap-data/config.json", []byte(`{"editor":"code"}`), 0644)
		}
		tl.Logger.Println("Current working directory:", baseDir)
		t, err := initializeTzap()
		if err != nil {
			return err
		}

		tl.Logger.Println("Tzap initialized")

		var projectP project.Project
		if lib != "" {
			var name project.ProjectName = project.ProjectName(lib)
			libProject, err := cmdinstance.NewLocalLibProject(baseDir, name)
			if err != nil {
				return err
			}
			tl.Logger.Println("Loaded lib ProjectDB:", name, libProject)
			projectP = libProject
		} else {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			localProject, err := cmdinstance.NewLocalProject(cwd)
			if err != nil {
				return err
			}
			projectP = localProject
		}
		t = t.
			AddContextChange(func(c context.Context) context.Context {
				return project.SetProjectInContext(c, projectP)
			})
		cmd.SetContext(cmdutil.SetTzapInContext(cmd.Context(), t))
		return nil
	},
}

func initializeTzap() (*tzap.Tzap, error) {

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
			return nil, err
		}
		connector = tzapconnect.WithConfig(apikey, config)
	}
	t := tzap.NewWithConnector(connector)

	return t, nil
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
	RootCmd.PersistentFlags().BoolVar(&settings.AutoMode, "automode", false, "Some but not all functions prompt if you want to overwrite an existing file. Putting automode to true enaled overwriting for those cases. Setting this to false does not disable anything.")
	RootCmd.PersistentFlags().IntVar(&settings.TruncateLimit, "truncate", 0, "Truncate limit for the interaction.")
	RootCmd.PersistentFlags().BoolVar(&settings.MD5Rewrites, "md5rewrites", true, "For some functions, this flag enables overwriting files with the same MD5 hash.")
	RootCmd.PersistentFlags().StringVar(&settings.IncludeList, "include", "", "Files include MD5 matching pattern.")
	RootCmd.PersistentFlags().BoolVar(&settings.DisableLogs, "disablelogs", false, "Whether to disable logging.")
	RootCmd.PersistentFlags().StringVar(&settings.LoggerOutput, "loggeroutput", "./.tzap-data/logs", "Path and name of the log file.")
	RootCmd.PersistentFlags().BoolVar(&settings.Stub, "stub", false, "Test non-live mode")
	RootCmd.PersistentFlags().Float32VarP(&settings.Temperature, "temperature", "t", 1.0, "Temperature for the interaction.")
	RootCmd.PersistentFlags().BoolVarP(&settings.Verbose, "verbose", "v", false, "Enable verbose logging")
	RootCmd.PersistentFlags().BoolVar(&settings.ApiMode, "api", false, "ALPHA: Enable clean stdout outputs. Also turns off editor mode.")
	RootCmd.PersistentFlags().BoolVarP(&settings.Yes, "yes", "y", false, "Answer yes on CLI related prompts - cost or similar related questions")
	RootCmd.PersistentFlags().StringVarP(&settings.Editor, "editor", "e", "vscode", "ALPHA: Select editor mode (stdin, editor, vscode, vim, nano).")
}
