package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tzapio/tzap/cli/cmd/cmdui"
	"github.com/tzapio/tzap/cli/cmd/cmdutil"
	"github.com/tzapio/tzap/internal/logging/tl"
	"github.com/tzapio/tzap/pkg/tzapaction/action"
	"github.com/tzapio/tzap/pkg/tzapaction/actionpb"
)

var refactorCmd = &cobra.Command{
	Use:   "trefactor [filein] [fileout] \nOR\n tzap refactor --refactorconfig refactorconfig.json \nOR\n tzap refactor --filein filein [--fileout fileout] [see for all params: tzap refactor --help] \n\n Json example: \n" + refactorJSONExample,
	Short: "Refactors code",
	Long: `The refactor command enables you to refactor code using either command-line flags or a configuration file. 
It is used to generate refactor and document code or generate documentation files.`,
	Run: func(cmd *cobra.Command, args []string) {
		tl.EnableUICompletionLogger()
		if len(args) > 0 {
			basicConfig.FileIn = args[0]
			if len(args) > 1 {
				basicConfig.FileOut = args[1]
			}
		} else if configFile != "" {
			config, err := loadConfig(configFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error loading config file: %v\n", err)
				println(refactorJSONExample)
				os.Exit(1)
			}
			basicConfig = config
		}

		if basicConfig.FileOut == "" {
			basicConfig.FileOut = basicConfig.FileIn
		}
		if (basicConfig.FileIn == "") || (basicConfig.Task == "") {
			fmt.Fprintf(os.Stderr, "error: filein and task are required\n")
			println(refactorJSONExample)
			os.Exit(1)
		}
		if tzapCliSettings.ApiMode {
			tzapCliSettings.Editor = "api"
		}
		cmdUI := cmdui.NewCMDUI("", tzapCliSettings.Editor)

		t := cmdutil.GetTzapFromContext(cmd.Context())
		result, err := action.Refactor(t, &actionpb.RefactorRequest{RefactorArgs: basicConfig})
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		for _, fileWrite := range result.FileWrites {
			cmdUI.EditFile(fileWrite)
		}
	},
}

var configFile string
var basicConfig *actionpb.RefactorArgs = &actionpb.RefactorArgs{}

func init() {
	RootCmd.AddCommand(refactorCmd)

	refactorCmd.Flags().StringVar(&configFile, "refactorconfig", "", "the path to the refactor file")
	refactorCmd.Flags().StringVar(&basicConfig.FileIn, "filein", "", "required - the file to refactor")
	refactorCmd.Flags().StringVar(&basicConfig.FileOut, "fileout", "", "optional - the output file (default filein)")
	refactorCmd.Flags().StringVar(&basicConfig.Mission, "mission", "", "optional - a description of the overall mission for the project")
	refactorCmd.Flags().StringVar(&basicConfig.Task, "task", "Analyze what can be improved. Refactor code to use better variable names and remove duplication. Refactor the following file to be more readable. Add documentation. Do not add any new public functions, only rewrite.", "required - a description of the refactoring task")
	refactorCmd.Flags().StringVar(&basicConfig.Plan, "plan", "", "a description to steer output. Recommended if you generate something else than code")
	refactorCmd.Flags().StringVar(&basicConfig.OutputFormat, "outputformat", "", "recommended - e.g. \"golang\")")
	refactorCmd.Flags().StringVar(&basicConfig.Example, "example", "", "optional an example of the refactoring task, {typescript code}")
	refactorCmd.Flags().StringSliceVarP(&basicConfig.InspirationFiles, "inspiration", "i", []string{}, "Optional comma-separated list of inspiration files or multiple -i flags.")
}

func loadConfig(filePath string) (*actionpb.RefactorArgs, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var configData actionpb.RefactorArgs
	err = json.Unmarshal(data, &configData)
	if err != nil {
		return nil, err
	}

	return &configData, nil
}

const refactorJSONExample = `{
    "filein": "input.go",
    "fileout": "output.go",
    "mission": "Improve code readability and maintainability",
    "task": "Analyze what can be improved. Refactor code to use better variable names and remove duplication. Refactor the following file to be more readable. Make comments for the functions. Do not add any new public functions, only rewrite.",
    "plan": "Do not write any text because this file will be saved directly to output.go",
    "outputformat": "golang",
    "example": "func doSomething(w http.ResponseWriter, r *http.Request, db *sql.DB) error {\n          // function logic\n    }",
    "inspirationfiles": [
        "/path/to/inspiration/file1.go",
        "/path/to/inspiration/file2.go"
    ]
}`
