package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/tzapio/tzap/cli/cmd/cmdutil/fileevaluator"
	"github.com/tzapio/tzap/pkg/util/stdin"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Setup .tzapignore, .tzapinclude, and configuration files",
	Long:  `This command initiates the tzap configuration by creating .tzapignore, .tzapinclude, and configuration files.`,
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(".tzapconfig"); os.IsExist(err) {
			println("Tzap is already initialized.")
			stdin.ConfirmPrompt("Continue anyway?")
		}

		println("Initializing Tzap...")
		if !tzapCliSettings.Yes {
			time.Sleep(time.Millisecond * 500)
			if _, err := os.Stat(".git"); os.IsNotExist(err) {
				if !stdin.ConfirmPrompt("\n\nWarning: Trying to find .git in the folder. This command should be run from the root of a project.\n\nIt's safe to continue, this is just incase you did not stand in your root folder. Continue anyway?") {
					return
				}
			}

			if b := stdin.ConfirmPrompt("\n\nTzap is in Beta. Would you like some general information about Tzap?"); b {
				println("\n\nTzap is a code cli tool that is designed to be easy to use.\n")
				stdin.GetStdinInput("Press enter to continue.")
				println("\n\nYou ask tzap to finish a prompt using: tzap prompt \"How do I use X library to enable my backend to do Y\"\n")
				stdin.GetStdinInput("Press enter to continue.")
				println("\n\nTzap assumes that you are running it from the project root folder. - Tzap attempts to traverse the folder to run from the root folder. During beta, for best results, always run tzap from root folder.\n")
				stdin.GetStdinInput("Press enter to continue.")
				println("\n\nTzap requires an openai apikey. You can get one from https://platform.openai.com/. You need to add a payment method to get started\n")
				stdin.GetStdinInput("Press enter to continue.")
				println("\n\nRegarding costs, embeddings should shows, but it's generally very affordable. A huge project like https://github.com/twitter/the-algorithm costs max 1 USD to embed. Usually less.\n")
				stdin.GetStdinInput("Press enter to continue.")
				println("\n\nA gpt4 call costs maximum 0.2 dollars and a gpt3.5 (default) costs a fraction of that. https://openai.com/pricing for more info.\n")
				stdin.GetStdinInput("Press enter to continue.")
				println("\n\nYou add your apikey through env variable or .env files. OPENAI_APIKEY=<apikey> for .env file.\n")
				stdin.GetStdinInput("Press enter to continue.")
				println("\n\nTzap is designed to be used with a .tzapignore file. This file is similar to a .gitignore file, but it is used to ignore files that interfere with search quality.\n")
				stdin.GetStdinInput("Press enter to continue.")
				println("\n\nTzap is designed to be used with a .tzapinclude file. This file is used to include.\n")
				stdin.GetStdinInput("Press enter to continue.")
				println("\n\nTzap is designed to be used with a .tzapinclude file. This file is used to include.\n")
				stdin.GetStdinInput("Press enter to continue.")
			}
		}
		if !tzapCliSettings.Yes {
			initializeTzap()
		}
		// Ask which text editor the user wants to use
		editor := askForEditor()

		// Write editor to the config file
		if err := writeEditorToConfigFile(editor); err != nil {
			panic(err)
		}
		touchTzapignore()
		touchTzapinclude()

		println("Initialization complete.")

		println("\nHere are some commands to try: ")
		println("\n   $ tzap change user component to contain an email field")
		println("\n   $ tzap refactor my user component")
		println("\n   $ tzap document my user component")
		println("\n   $ tzap test my user component")
	},
}

func writeEditorToConfigFile(selected string) error {
	cfg := map[string]string{
		"editor": selected,
	}

	cfgJSON, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("error parsing json for config: %w", err)
	}
	if err := os.MkdirAll(".tzap-data", 0755); err != nil {
		return fmt.Errorf("error creating .tzap-data folder: %w", err)
	}

	if err := os.WriteFile(".tzap-data/config.json", cfgJSON, 0644); err != nil {
		return fmt.Errorf("error writing to config.json: %w", err)
	}

	return nil
}
func checkVSCodeInstalled() bool {
	cmd := exec.Command("code", "--help")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}
func askForEditor() string {
	if tzapCliSettings.Yes {
		return "stdin"
	}
	options := []string{"vscode", "code", "vim", "nano", "editor", "stdin"}
	if checkVSCodeInstalled() {
		time.Sleep(time.Millisecond * 1000)
		runWithVscode := stdin.ConfirmPrompt(`

Would you like to run tzap with VSCode? Tzap will open VSCode when prompting for input and show diffs in VSCode.`)

		if runWithVscode {
			return "vscode"
		}
	}

	prompt := fmt.Sprintf(`
	
Choose your preferred text editor:
	- %s: Edit prompts directly from file, get diffs and more
	- %s: Opens the file in vim when prompting
	- %s: Opens the file in nano when prompting
	- %s: allows for editing files directly but does not connect to any specific UI.
	- stdin (default): asks for input in CLI.
`, options[0], options[2], options[3], options[4])
	fmt.Print(prompt)

	text := stdin.GetStdinInput("\nEnter your choice (press enter for default choice): ")
	text = strings.TrimSpace(text)
	for _, e := range options {
		if text == e {
			return e
		}
		if text == "" {
			return "stdin"
		}
	}

	fmt.Println("Invalid input. Please choose a valid text editor.")
	return askForEditor()
}
func touchTzapignore() {
	if _, err := os.Stat(".tzapignore"); err == nil {
		println("Warning: .tzapignore already exists.")
		time.Sleep(time.Millisecond * 500)
		return
	}
	var gitignoreContent string

	if _, err := os.Stat(".gitignore"); os.IsExist(err) {
		println("Warning: .gitignore does not exist.")
		content, _ := os.ReadFile(".gitignore")
		if len(content) > 0 {
			gitignoreContent = string(content)
		}
	} else if err != nil {
		println("Warning: did not copy for .tzapignore. .gitignore error: ", err)
		time.Sleep(time.Millisecond * 500)
	}
	tzapIgnoreContent := fileevaluator.BaseTzapIgnore + "\n" + gitignoreContent

	if err := os.WriteFile(".tzapignore", []byte(tzapIgnoreContent), 0644); err != nil {
		println("Error:", err)
	}
	println("Created file .tzapignore")
	time.Sleep(time.Millisecond * 500)
}

func touchTzapinclude() {
	//if not exist, copy .gitignore to .tzapignore
	if _, err := os.Stat(".tzapinclude"); err == nil {
		println("Warning: .tzapinclude already exists.")
		time.Sleep(time.Millisecond * 500)
		return
	}

	commonLanguage := fileevaluator.BaseTzapInclude

	if err := os.WriteFile(".tzapinclude", []byte(commonLanguage), 0644); err != nil {
		println("Error:", err)
	}
	println("Created file .tzapinclude")
	time.Sleep(time.Millisecond * 500)
}

func init() {
	RootCmd.AddCommand(initCmd)
	initCmd.PersistentFlags().BoolVarP(&tzapCliSettings.Yes, "yes", "y", false, "Skip all prompts and use default values")
}
