package tzap

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/types/openai"
	"github.com/tzapio/tzap/pkg/util"
)

// LoadFileDir exposes an array of Tzaps in the previous elements .Data["children"]. Each child is a .LoadTask(file)
func (t *Tzap) LoadFileDir(dir string, match string) *Tzap {
	pattern := filepath.Join(dir, match)
	files, err := filepath.Glob(pattern)
	if err != nil {
		panic(fmt.Errorf("error reading directory: %w", err))
	}
	return t.LoadFiles(files)
}

// LoadFiles exposes an array of Tzaps in with .Data["children"]. Each child is a .LoadTask(file)
func (t *Tzap) LoadFiles(filepaths []string) *Tzap {
	var ts []*Tzap
	for _, file := range filepaths {
		// Check if the file is a regular file and its name contains "test" if test is true.
		if info, err := os.Stat(file); err == nil && !info.IsDir() {
			// Load the file content and create a Tzap with the file content as the message content.
			ts = append(ts, t.LoadTask(file))
		} else if err != nil {
			panic(err)
		}
	}

	return t.AddTzap(&Tzap{
		Name: "WithAddFile",
		Data: types.MappedInterface{"children": ts},
	})
}

// LoadTask loads a file and returns a Tzap with the file's content
func (t *Tzap) LoadTask(filepath string) *Tzap {
	Log(t, "Adding file", filepath)
	originalContent, err := util.ReadFile(filepath)
	if err != nil {
		panic(fmt.Errorf("cannot add file: %w", err))
	}

	data := types.MappedInterface{
		"filepath": filepath,
		"content":  originalContent,
	}
	withAddFile := t.AddTzap(
		&Tzap{
			Name: "WithAddFile",
			Message: types.Message{
				Role:    openai.ChatMessageRoleUser,
				Content: originalContent,
			},
			Data: data,
		})
	return withAddFile
}

// PrepareOutputTask creates a Tzap with an empty file content to be used for outputting to a file
func (t *Tzap) PrepareOutputTask(filepath string) *Tzap {
	data := types.MappedInterface{
		"filepath": filepath,
		"content":  "",
	}
	withOutputTask := t.AddTzap(&Tzap{
		Name: "withOutputTask",
		Data: data,
	})
	return withOutputTask
}

// LoadTaskOrRequestNewTask loads a file if it exists, otherwise requests a new file content from OpenAI and applies the changes to the original file
func (t *Tzap) LoadTaskOrRequestNewTask(filepath string) *Tzap {
	Log(t, "Opening", filepath)
	if _, err := os.Stat(filepath); err != nil {
		return t.
			PrepareOutputTask(filepath).
			FetchTask()
	} else {
		return t.LoadTask(filepath)
	}
}
