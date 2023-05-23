package tzap

import (
	"fmt"
	"os"

	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/types/openai"
	"github.com/tzapio/tzap/pkg/util"
)

// LoadFileDir exposes an array of Tzaps in the previous elements .Data["children"]. Each child is a .LoadTask(file)
func (t *Tzap) LoadFileDir(dir string) *Tzap {
	_, err := os.ReadDir(dir)
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		panic(err)
	}
	files, err := util.ListFilesInDir(dir)
	if err != nil {
		fmt.Printf("Error listing files in directory: %v\n", err)
		panic(err)
	}
	return t.LoadFiles(files)
}

// LoadFiles exposes an array of Tzaps in with .Data["children"]. Each child is a .LoadTask(file)
func (t *Tzap) LoadFiles(filepaths []string) *Tzap {
	var ts []*Tzap
	t = t.AddTzap(&Tzap{
		Name: "LoadFiles",
		Data: types.MappedInterface{},
	})
	for _, file := range filepaths {
		// Check if the file is a regular file and its name contains "test" if test is true.
		if info, err := os.Stat(file); err == nil && !info.IsDir() {
			// Load the file content and create a Tzap with the file content as the message content.
			ts = append(ts, t.LoadFileAsCompletion(file))
		} else if err != nil {
			panic(err)
		}
	}
	t.Data["children"] = ts
	return t
}

// LoadFileAsCompletion loads a file and returns a Tzap with the file's content
func (t *Tzap) LoadFileAsCompletion(filePath string) *Tzap {
	Log(t, "Adding file", filePath)
	originalContent, err := os.ReadFile(filePath)
	if err != nil {
		panic(fmt.Errorf("cannot add file: %w", err))
	}

	data := types.MappedInterface{
		"filepath": filePath,
		"content":  string(originalContent),
	}
	loadTaskFromFile := t.AddTzap(
		&Tzap{
			Name: "LoadCompletion",
			Message: types.Message{
				Role:    openai.ChatMessageRoleAssistant,
				Content: string(originalContent),
			},
			Data: data,
		})

	return loadTaskFromFile
}

// LoadTaskOrRequestNewTask loads a file if it exists, otherwise requests a new file content from OpenAI and applies the changes to the original file
func (t *Tzap) LoadCompletionOrRequestCompletion(filePath string) *Tzap {
	Log(t, "Opening", filePath)
	t = t.AddTzap(&Tzap{
		Name: "LoadCompletionOrRequestCompletion"})

	var out *Tzap
	if _, err := os.Stat(filePath); err != nil {
		out = t.
			RequestChatCompletion().
			StoreCompletion(filePath)
	} else {
		out = t.LoadFileAsCompletion(filePath)
	}

	return out
}
