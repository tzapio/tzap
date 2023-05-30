package stdin

import (
	"bufio"
	"os"
	"strings"
)

func ConfirmPrompt(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)

	for tries := 0; tries < 10; tries++ {
		println(prompt + " (y/n):")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if strings.EqualFold(input, "y") {
			return true
		} else if strings.EqualFold(input, "n") {
			return false
		} else {
			println("Invalid input. Please enter 'y' or 'n'.")
		}
	}
	println("no valid input after 10 tries - assuming no")
	return false
}

func GetStdinInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	print(prompt)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}
