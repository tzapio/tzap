package cmd

import (
	"encoding/json"
	"fmt"
	"os"
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
			cmd.Println("Tzap is already initialized.")
			stdin.ConfirmPrompt("Continue anyway?")
		}

		cmd.Println("Initializing Tzap...")
		time.Sleep(time.Millisecond * 500)
		if _, err := os.Stat(".git"); os.IsNotExist(err) {
			if !stdin.ConfirmPrompt("Warning: Trying to find .git in the folder. This command should be run from the root of a project. ") {
				return
			}
		}

		if b := stdin.ConfirmPrompt("Tzap is in Beta. Would you like some general information about Tzap?"); b {
			cmd.Println("\n\nTzap is a code cli tool that is designed to be easy to use.")
			stdin.GetStdinInput("Press enter to continue.")
			cmd.Println("\n\nYou ask tzap to finish a prompt using: tzap prompt \"How do I use X library to enable my backend to do Y\" ")
			stdin.GetStdinInput("Press enter to continue.")
			cmd.Println("\n\nTzap assumes that you are running it from the project root folder. - Tzap attempts to traverse the folder to run from the root folder. During beta, for best results, always run tzap from root folder.")
			stdin.GetStdinInput("Press enter to continue.")
			cmd.Println("\n\nTzap requires an openai apikey. You can get one from https://platform.openai.com/. You need to add a payment method to get started")
			stdin.GetStdinInput("Press enter to continue.")
			cmd.Println("\n\nRegarding costs, embeddings should shows, but it's generally very affordable. https://github.com/twitter/the-algorithm costs around 1.5 USD to embed.")
			stdin.GetStdinInput("Press enter to continue.")
			cmd.Println("\n\nA gpt4 call costs maximum 0.2 dollars and a gpt3.5 (default) costs a fraction of that. https://openai.com/pricing for more info.")
			stdin.GetStdinInput("Press enter to continue.")
			cmd.Println("\n\nYou add your apikey through env variable or .env files. OPENAI_APIKEY=<apikey> for .env file. ")
			stdin.GetStdinInput("Press enter to continue.")
			cmd.Println("\n\nTzap is designed to be used with a .tzapignore file. This file is similar to a .gitignore file, but it is used to ignore files that interfere with search quality. ")
			stdin.GetStdinInput("Press enter to continue.")
			cmd.Println("\n\nTzap is designed to be used with a .tzapinclude file. This file is used to include .")
			stdin.GetStdinInput("Press enter to continue.")
			cmd.Println("\n\nTzap is designed to be used with a .tzapinclude file. This file is used to include .")
			stdin.GetStdinInput("Press enter to continue.")
		}
		// Ask which text editor the user wants to use
		editor := askForEditor()

		// Write editor to the config file
		if err := writeEditorToConfigFile(editor); err != nil {
			panic(err)
		}
		touchTzapignore()
		touchTzapinclude()
		generateViperConfig()
		cmd.Println("Initialization complete.")

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
func askForEditor() string {
	options := []string{"vscode", "code", "vim", "nano", "editor"}
	prompt := fmt.Sprintf("Choose your preferred text editor:\n- %s (alias: code): Edit prompts directly from file and have them automatically open using /usr/bin/code vscode command\n- %s: Opens the file in vim when prompting\n- %s:Opens the file in vim when prompting\n- %s: allows for editing files directly but does not connect to any specific UI.\n- stdin: default, asks for input in CLI.\n\n code is recommended and stdin is the easiest to get started with.", options[0], options[1], options[2], options[3])
	fmt.Print(prompt)

	text := stdin.GetStdinInput("\nEnter your choice (press enter to chose: stdin): ")
	text = strings.TrimSpace(text)
	for _, e := range options {
		if e == text {
			return e
		}
		if e == "" {
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

func generateViperConfig() {

}

func init() {
	RootCmd.AddCommand(initCmd)
}
