package cmdui

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/util/stdin"
)

type CMDUI struct {
	filePath string
	editor   string
}

func NewCMDUI(promptFile string, editor string) *CMDUI {
	if editor == "stdin" {
		if promptFile == "" {
			return &CMDUI{filePath: promptFile, editor: editor}
		}
		println("stdin editor does not support promptFile. change .tzap-data/config.json to editor: editor, vscode, vim, nano. {\"editor\":\"choice\"}")
		os.Exit(1)
	}
	if promptFile != "" {
		editorUI := CMDUI{filePath: promptFile, editor: editor}
		return &editorUI
	}

	err := os.MkdirAll(".tzap-data/chats", 0755)
	if err != nil {
		panic(err)
	}
	file, err := os.CreateTemp(".tzap-data/chats", "tzapchange*.txt")
	if err != nil {
		panic(err)
	}
	file.Close()
	return &CMDUI{filePath: file.Name(), editor: editor}
}
func (ui *CMDUI) Init() {
	if ui.editor == "vscode" || ui.editor == "code" {
		exec.Command("code", "-r", "--goto", ui.filePath+":-1").Run()
		return
	}
}
func (ui *CMDUI) RunEditor() {
	if ui.editor == "vscode" || ui.editor == "code" {
		println("Write at top of file and hit save 3 times in 1 seconds to trigger re-prompt. ")
		exec.Command("code", "-r", "--goto", ui.filePath+":-1").Run()
		ui.WatchSavesToFile(time.Second*1, 3)
		return
	}
	if ui.editor == "editor" {
		println("Write at top of file and hit save 3 times in 1 seconds to trigger re-prompt. ")
		ui.WatchSavesToFile(time.Second*1, 3)
		return
	}
	if ui.editor == "vim" {
		println("Write at top of file.")
		cmd := exec.Command("vim", ui.filePath)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			panic(err)
		}
		return
	}
	if ui.editor == "nano" {
		println("Write at top of file.")
		cmd := exec.Command("nano", ui.filePath)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			panic(err)
		}
		return
	}
	if ui.editor == "api" {
		return
	}
	if ui.editor == "stdin" || ui.editor != "api" {
		return
	}
}
func (ui *CMDUI) AddPromptTextWithStdinUI(thread []types.Message) []types.Message {
	if ui.editor == "stdin" {
		prompt := stdin.GetStdinInput("Follow up prompt: ")
		return append(thread, types.Message{Role: "user", Content: prompt})
	}
	println("\n\nFile: ", ui.filePath)
	println("")
	err := ui.SaveMessageThreadToFile(thread)
	if err != nil {
		panic(err)
	}
	ui.RunEditor()
	return ui.ReadMessagesFromFile()
}

func (ui *CMDUI) ReadMessagesFromFile() []types.Message {
	bytes, err := os.ReadFile(ui.filePath)
	if err != nil {
		panic(err)
	}
	content := string(bytes)
	return ui.DeserializeMessageThread(content)
}
func (ui *CMDUI) DeserializeMessageThread(content string) []types.Message {
	var messages []types.Message
	messageLines := strings.Split(content, "---")
	for _, messageLine := range messageLines {
		message := types.Message{}
		lines := strings.Split(strings.TrimSpace(messageLine), "\n")
		if len(lines) > 0 {
			if strings.HasPrefix(lines[0], "@role:") {
				message.Role = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(lines[0]), "@role:"))
				message.Content = strings.TrimSpace(strings.Join(lines[1:], "\n"))
			} else {
				message.Role = "user"
				message.Content = strings.TrimSpace(strings.Join(lines, "\n"))
			}
		}
		messages = append(messages, message)
	}
	reverseMessages := []types.Message{}
	for i := len(messages) - 1; i >= 0; i-- {
		if strings.TrimSpace(messages[i].Content) == "" || strings.TrimSpace(messages[i].Role) == "" {
			continue
		}
		reverseMessages = append(reverseMessages, messages[i])
	}
	return reverseMessages
}

// SaveMessageThreadToFile saves the thread to a file
// The file is saved in reverse order, so that the last message is the first line in the file
// The first line may skip role and defaults to user
func (ui *CMDUI) SaveMessageThreadToFile(messages []types.Message) error {
	text, err := ui.SerializeMessageThread(messages)
	if err != nil {
		return err
	}

	if err := os.WriteFile(ui.filePath, []byte(text), 0644); err != nil {
		return err
	}
	return nil
}

func (ui *CMDUI) SerializeMessageThread(messages []types.Message) (string, error) {
	reversedMessages := []types.Message{}
	for i := len(messages) - 1; i >= 0; i-- {
		if strings.TrimSpace(messages[i].Content) == "" || strings.TrimSpace(messages[i].Role) == "" {
			continue
		}
		reversedMessages = append(reversedMessages, messages[i])
	}
	s := strings.Builder{}
	if len(reversedMessages) > 0 {
		s.WriteString("\n\n")
	}

	for _, msg := range reversedMessages {
		if _, err := s.WriteString(fmt.Sprintf("---\n@role:%s\n%s\n", msg.Role, msg.Content)); err != nil {
			return "", err
		}
	}
	return s.String(), nil
}
