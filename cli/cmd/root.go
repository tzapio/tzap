package cmd

import (
	"context"
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
	"github.com/tzapio/tzap/cli/cmd/cmdinstance"
	"github.com/tzapio/tzap/cli/cmd/cmdutil"
	"github.com/tzapio/tzap/internal/logging/tl"
	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/project"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/types/openai"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/tzapconnect"
	"github.com/tzapio/tzap/pkg/tzapconnect/stubconnector"
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
	Stub          bool
	Verbose       bool
	ApiMode       bool
	Yes           bool
	Editor        string
	EmbeddingURL  string
	CompletionURL string
}

var RootCmd = &cobra.Command{
	Use:   "tzap",
	Short: "Tzap Cli",

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		tl.Logger.Println("Cobra CLI Root start")

		if tzapCliSettings.Verbose {
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
			println("Warning: No .tzapinclude file found. Run 'tzap init' Using current directory as root.", err)
		} else {
			os.Chdir(baseDir)
		}
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

	var connector types.TzapConnector
	if tzapCliSettings.Stub {
		connector = stubconnector.StubWithConfig(config)
	} else {
		apikey, err := tzapconnect.LoadOPENAI_API_KEY()
		if err != nil {
			choice := stdin.ConfirmPrompt("Cannot find OPENAI_APIKEY.\n\nWould you like to add it now?")
			if choice {
				apikey = stdin.GetStdinInput("Enter OPENAI_APIKEY to save to .env:\n")
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
			} else {
				println("Aborted, cannot continue without OPENAI_APIKEY.")
				os.Exit(1)
				return nil, err
			}
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

var modelMap map[string]string = map[string]string{
	"gpt35":    openai.GPT3Dot5Turbo,
	"gpt356":   openai.GPT3Dot5Turbo0613,
	"gpt3516":  openai.GPT16,
	"gpt3516k": openai.GPT16,
	"gpt16":    openai.GPT16,
	"gpt4":     openai.GPT4,
}

func init() {
	RootCmd.CompletionOptions.HiddenDefaultCmd = true
	tzapCliSettings.MD5Rewrites = true

	RootCmd.PersistentFlags().StringVarP(&tzapCliSettings.Model, "model", "m", "gpt35", "Define what openai model to use. (Available gpt35 gpt356 (june model) gpt3516 (alias gpt16) gpt4).")
	RootCmd.PersistentFlags().StringVarP(&tzapCliSettings.CompletionURL, "baseurl", "b", "", "Completion URL")
	RootCmd.PersistentFlags().StringVar(&tzapCliSettings.EmbeddingURL, "embeddingbaseurl", "", "Embedding URL")
	//RootCmd.PersistentFlags().BoolVar(&tzapCliSettings.AutoMode, "automode", false, "Some but not all functions prompt if you want to overwrite an existing file. Putting automode to true enaled overwriting for those cases. Setting this to false does not disable anything.")
	//RootCmd.PersistentFlags().IntVar(&tzapCliSettings.TruncateLimit, "truncate", 0, "Truncate limit for the interaction.")
	//RootCmd.PersistentFlags().BoolVar(&tzapCliSettings.MD5Rewrites, "md5rewrites", true, "For some functions, this flag enables overwriting files with the same MD5 hash.")
	//RootCmd.PersistentFlags().BoolVar(&tzapCliSettings.DisableLogs, "disablelogs", false, "Whether to disable logging.")
	RootCmd.PersistentFlags().StringVar(&tzapCliSettings.LoggerOutput, "loggeroutput", ".tzap-data/logs/", "Path and name of the log file.")
	//RootCmd.PersistentFlags().BoolVar(&tzapCliSettings.Stub, "stub", false, "Test non-live mode")
	RootCmd.PersistentFlags().Float32VarP(&tzapCliSettings.Temperature, "temperature", "t", 1.0, "Temperature for the interaction.")
	RootCmd.PersistentFlags().BoolVarP(&tzapCliSettings.Verbose, "verbose", "v", false, "Enable verbose logging")
	RootCmd.PersistentFlags().BoolVar(&tzapCliSettings.ApiMode, "api", false, "ALPHA: Enable clean stdout outputs. Also turns off editor mode.")
	RootCmd.PersistentFlags().BoolVarP(&tzapCliSettings.Yes, "yes", "y", false, "Answer yes on CLI related prompts - cost or similar related questions")
	//RootCmd.PersistentFlags().StringVarP(&tzapCliSettings.Editor, "editor", "e", "vscode", "ALPHA: Select editor mode (stdin, editor, vscode (alias code), vim, nano).")
}
