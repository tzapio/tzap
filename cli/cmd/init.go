package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Setup .tzapignore, .tzapinclude, and configuration files",
	Long:  `This command initiates the tzap configuration by creating .tzapignore, .tzapinclude, and configuration files.`,
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(".tzapconfig"); os.IsExist(err) {
			fmt.Println("Tzap is already initialized.")
			return
		}

		fmt.Println("Initializing Tzap...")
		time.Sleep(time.Millisecond * 500)
		if _, err := os.Stat(".git"); os.IsNotExist(err) {
			if !surveyConfirm("Warning: This command should be run from the root of a project") {
				return
			}
		}

		touchTzapignore()
		touchTzapinclude()
		generateViperConfig()
		fmt.Println("Initialization complete.")

	},
}

func surveyConfirm(prompt string) bool {
	confirm := false
	promptError := survey.AskOne(&survey.Confirm{
		Message: prompt + "\nAre you sure you want to continue?",
	}, &confirm)

	if promptError != nil {
		fmt.Println("Error:", promptError)
		return false
	}

	return confirm
}
func touchTzapignore() {
	//if not exist, copy .gitignore to .tzapignore
	if _, err := os.Stat(".tzapignore"); err == nil {
		fmt.Println("Warning: .tzapignore already exists.")
		time.Sleep(time.Millisecond * 500)
		return
	}
	var gitignoreContent string

	if _, err := os.Stat(".gitignore"); os.IsExist(err) {
		fmt.Println("Warning: .gitignore does not exist.")
		content, err := os.ReadFile(".gitignore")
		if err != nil {
			fmt.Println("Error:", err)
		}
		gitignoreContent = string(content)
	} else if err != nil {
		fmt.Println("Warning: did not copy for .tzapignore. .gitignore error: ", err)
		time.Sleep(time.Millisecond * 500)
	}
	tzapIgnoreContent := `# Tzap ignore file. Add extra files like test folders, or other files that interfere with search (embeddings) quality. 
node_modules
\n\n	
# copied from .gitignore
` + gitignoreContent
	if err := os.WriteFile(".tzapignore", []byte(tzapIgnoreContent), 0644); err != nil {
		fmt.Println("Error:", err)
	}
	time.Sleep(time.Millisecond * 500)
}
func touchTzapinclude() {

	//if not exist, copy .gitignore to .tzapignore
	if _, err := os.Stat(".tzapinclude"); err == nil {
		fmt.Println("Warning: .tzapinclude already exists.")
		time.Sleep(time.Millisecond * 500)
		return
	}

	commonLanguage := `# Common languages. Example, remove .js if .js files are only compiled bundles.
*.js
*.tsx
*.ts
*.py
*.go
*.java
*.c
*.cpp
*.h
*.hpp
*.rb
*.php
`

	if err := os.WriteFile(".tzapinclude", []byte(commonLanguage), 0644); err != nil {
		fmt.Println("Error:", err)
	}
	println("Created file .tzapinclude")
	time.Sleep(time.Millisecond * 500)
}

func generateViperConfig() {

}

func init() {
	rootCmd.AddCommand(initCmd)
}
