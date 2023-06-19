package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tzapio/tzap/cli/action"
	"github.com/tzapio/tzap/cli/cmd/cmdutil"
	"github.com/tzapio/tzap/internal/logging/tl"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/util"
	"github.com/tzapio/tzap/workflows/code/codegeneration"
	"github.com/tzapio/tzap/workflows/stdinworkflows"
)

var refactorCmd = &cobra.Command{
	Use:   "refactor [filein] [fileout]",
	Short: "Refactors code",
	Long: `The refactor command enables you to refactor code using either command-line flags or a configuration file. 
It is used to generate refactor and document code or generate documentation files.`,
	Run: func(cmd *cobra.Command, args []string) {
		tl.EnableUICompletionLogger()

		if len(args) > 0 {
			fileIn := args[0]
			fileOut := fileIn
			if len(args) > 1 {
				fileOut = args[1]
			}
			runRefactor(cmd, &codegeneration.BasicRefactoringConfig{FileIn: fileIn, FileOut: fileOut})
		} else if configFile != "" {
			runRefactorWithConfigFile(cmd, configFile)
		} else {
			fmt.Fprintf(os.Stderr, "error: either filein and fileout or refactorconfig flag is required\n")
			cmd.Println(refactorJSONExample)
			os.Exit(1)
		}
	},
}

var configFile string

func init() {
	RootCmd.AddCommand(refactorCmd)

	refactorCmd.Flags().StringVar(&configFile, "refactorconfig", "", "the path to the refactor file")
	refactorCmd.Flags().StringSliceVarP(&inspirationFiles, "inspiration", "i", []string{}, "Optional comma-separated list of inspiration files or multiple -i flags.")
}

func runRefactor(cmd *cobra.Command, basicConfig *codegeneration.BasicRefactoringConfig) {
	if basicConfig.FileOut == "" {
		basicConfig.FileOut = basicConfig.FileIn
	}

	err := tzap.HandlePanic(func() {
		t := cmdutil.GetTzapFromContext(cmd.Context())
		defer t.HandleShutdown()

		t.ApplyWorkflow(action.SearchWorkflow(action.PromptWorkflowArgs{
			InspirationFiles: []string{basicConfig.FileIn},
			SearchQuery:      util.ReadFileP(basicConfig.FileIn),
			EmbedsCount:      40,
			NCount:           40,
			DisableIndex:     false,
			Yes:              true,
			MessageThread:    []types.Message{},
		})).
			ApplyWorkflowFN(codegeneration.MakeCode(codegeneration.BasicRefactoringConfig{
				FileIn:  basicConfig.FileIn,
				FileOut: basicConfig.FileOut,
			})).
			ApplyWorkflow(stdinworkflows.BeforeProceedingWorkflow()).
			StoreCompletion(basicConfig.FileOut)
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func runRefactorWithConfigFile(cmd *cobra.Command, configFile string) {
	//#### you should refactor this part such that it uses codegeneration.BasicRefactoringConfig
	config, err := loadConfig(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading config file: %v\n", err)
		cmd.Println(refactorJSONExample)
		os.Exit(1)
	}

	runRefactor(cmd, config)
}

func loadConfig(filePath string) (*codegeneration.BasicRefactoringConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var configData codegeneration.BasicRefactoringConfig
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
    "searchQuery": "BETA: if set; append embeddings using searchQuery as key",
    "example": "func doSomething(w http.ResponseWriter, r *http.Request, db *sql.DB) error {\n          // function logic\n    }",
    "inspirationfiles": [
        "/path/to/inspiration/file1.go",
        "/path/to/inspiration/file2.go"
    ]
}`
