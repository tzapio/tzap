package cmd

import (
	"encoding/json"

	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteEditorToConfigFile(t *testing.T) {
	// Set up
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	dir, err := os.MkdirTemp("", "tzap-test")
	if err != nil {
		t.Fatal(err)
	}

	defer os.RemoveAll(dir)

	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(cwd)
	// Test function
	editor := "vim"
	err = writeEditorToConfigFile(editor)
	if err != nil {
		t.Fatal(err)
	}

	// Check that written data is correct
	data, err := os.ReadFile(dir + "/.tzap-data/config.json")
	assert.NoError(t, err)

	var cfg map[string]interface{}
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		t.Fatal(err)
	}
	if editor != cfg["editor"] {
		t.Fatalf("expected %s, got %s", editor, cfg["editor"])
	}
}
