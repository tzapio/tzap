package tzap_test

import (
	"os"
	"testing"

	"github.com/tzapio/tzap/pkg/tzap"
)

func Test_LoadTask_given_existing_file_path_expect_tzap_with_file_contents(t *testing.T) {
	// Prepare a temp file with content
	content := "Test content"
	tempFile, err := os.CreateTemp(os.TempDir(), "Test_LoadTask_given_existing_file_path_expect_tzap_with_file_contents")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	_, err = tempFile.WriteString(content)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	// Load the file with Tzap
	tt := tzap.InternalNew()
	tt.TG = &mockTG{}
	loadedTzap := tt.LoadTask(tempFile.Name())

	// Check if the loadedTzap contains the correct content
	if loadedTzap.Message.Content != content {
		t.Errorf("Expected content to be '%s', but got '%s'", content, loadedTzap.Message.Content)
	}
}

func Test_LoadFileDir_given_path_and_match_expect_tzaps_with_file_contents(t *testing.T) {
	// Prepare temp files with content
	content1 := "Test content 1"
	content2 := "Test content 2"
	tempFile1, err := os.CreateTemp(os.TempDir(), "Test_LoadFileDir_given_path_and_match_expect_tzaps_with_file_contents_1")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile1.Name())
	_, err = tempFile1.WriteString(content1)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile2, err := os.CreateTemp(os.TempDir(), "Test_LoadFileDir_given_path_and_match_expect_tzaps_with_file_contents_2")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile2.Name())
	_, err = tempFile2.WriteString(content2)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	// Load the files with Tzap
	path := os.TempDir()
	match := "Test_LoadFileDir_given_path_and_match_expect_tzaps_with_file_contents_*"
	tt := tzap.InternalNew()
	tt.TG = &mockTG{}
	loadedTzap := tt.LoadFileDir(path, match)

	// Check if the loadedTzap contains two tzaps with the correct content
	children, ok := loadedTzap.Data["children"].([]*tzap.Tzap)
	if !ok {
		t.Fatalf("Failed to convert children to []*tzap.Tzap")
	}
	if len(children) != 2 {
		t.Fatalf("Expected 2 children, but got %d", len(children))
	}

	// Check if the contents are correct
	foundContent1 := false
	foundContent2 := false
	for _, child := range children {
		if child.Message.Content == content1 {
			foundContent1 = true
		}
		if child.Message.Content == content2 {
			foundContent2 = true
		}
	}
	if !foundContent1 {
		t.Errorf("Expected to find content1 '%s', but did not", content1)
	}
	if !foundContent2 {
		t.Errorf("Expected to find content2 '%s', but did not", content2)
	}
}

func Test_LoadTaskOrRequestNewTask_given_existing_file_path_expect_tzap_with_file_contents(t *testing.T) {
	// Prepare a temp file with content
	content := "Test content"
	tempFile, err := os.CreateTemp(os.TempDir(), "Test_LoadTaskOrRequestNewTask_given_existing_file_path_expect_tzap_with_file_contents")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	_, err = tempFile.WriteString(content)
	if err != nil {
		t.Fatalf("Failed to write content to temp file: %v", err)
	}

	// Load the file with Tzap
	tt := tzap.InternalNew()

	tt.TG = &mockTG{}
	loadedTzap := tt.LoadTaskOrRequestNewTask(tempFile.Name())

	// Check if the loadedTzap contains the correct content
	if loadedTzap.Message.Content != content {
		t.Errorf("Expected content to be %s but got %s", content, loadedTzap.Message.Content)
	}
}
