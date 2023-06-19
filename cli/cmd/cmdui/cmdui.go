package cmdui

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/types/openai"
	"github.com/tzapio/tzap/pkg/tzapfile"
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
	// Start the file number from 10000
	fileNumber := 10000
	var filename string

	for {
		fileNumber--
		filename = ".tzap-data/chats/tzapchange" + strconv.Itoa(fileNumber) + ".md"
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			break
		}

	}

	return &CMDUI{filePath: filename, editor: editor}
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
		messages := ui.ReadMessagesFromFile()
		if len(messages) == 0 {
			println("No messages found in file.")
			os.Exit(1)
			return
		}
		if strings.TrimSpace(messages[len(messages)-1].Role) == openai.ChatMessageRoleAssistant {
			println("No messages found in file. ", messages[len(messages)-1].Content, messages[len(messages)-1].Role)
			os.Exit(1)
			return
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
		messages := ui.ReadMessagesFromFile()
		println(len(messages))
		if len(messages) == 0 {
			println("No messages found in file.")
			os.Exit(1)
			return
		}
		if strings.TrimSpace(messages[len(messages)-1].Content) == "" {
			println("No user message found in file.")
			os.Exit(1)
			return
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
	return tzapfile.DeserializeMessageThread(content)
}

// SaveMessageThreadToFile saves the thread to a file
// The file is saved in reverse order, so that the last message is the first line in the file
// The first line may skip role and defaults to user
func (ui *CMDUI) SaveMessageThreadToFile(messages []types.Message) error {
	text, err := tzapfile.SerializeMessageThread(messages)
	if err != nil {
		return err
	}

	if err := os.WriteFile(ui.filePath, []byte(text), 0644); err != nil {
		return err
	}
	return nil
}
