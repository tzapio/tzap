package stdin

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ConfirmPrompt(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)

	for tries := 0; tries < 10; tries++ {
		fmt.Printf("%s (y/n):", prompt)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if strings.EqualFold(input, "y") {
			return true
		} else if strings.EqualFold(input, "n") {
			return false
		} else {
			fmt.Println("Invalid input. Please enter 'y' or 'n'.")
		}
	}
	fmt.Println("no valid input after 10 tries - assuming no")
	return false
}

func GetStdinInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}
