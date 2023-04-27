package stdin

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ApplyChanges writes the edited content to the file
func ApplyChanges(filename, content string) error {
	dir := filepath.Dir(filename)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}
	return os.WriteFile(filename, []byte(content), 0644)
}
func ConfirmToContinue() bool {
	reader := bufio.NewReader(os.Stdin)
	var input string = ""
	for input != "y" && input != "n" && input != "Y" && input != "N" {
		fmt.Print("Continue? (y/n): ")
		raw, _ := reader.ReadString('\n')
		input = strings.TrimSpace(raw)
	}
	return input == "y" || input == "Y"
}
func ConfirmAndApplyChanges(filename string, editedContent string, cb func(filename string, content string) error) error {
	reader := bufio.NewReader(os.Stdin)
	var err error

	for tries := 0; tries < 10; tries++ {

		fmt.Print("Apply changes? (" + filename + ") (y/n): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if strings.EqualFold(input, "y") {
			err = cb(filename, editedContent)
			if err != nil {
				panic(fmt.Errorf("error applying changes: %w", err))
			} else {
				fmt.Println("Changes applied.")
				return nil
			}
		} else if strings.EqualFold(input, "n") {
			fmt.Println("Changes not applied.")
			return fmt.Errorf("changes not applied")
		} else {
			fmt.Println("Invalid input. Please enter 'y' or 'n'.")
		}
	}
	panic("Crashing after 5 invalid inputs. Probably something wrong.")
}

func GetStdinInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}
