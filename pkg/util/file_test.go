package util_test

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/tzapio/tzap/pkg/util"
)

func TestListFilesInDir(t *testing.T) {
	// Create a temp directory
	tempDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %s", err)
	}
	defer os.RemoveAll(tempDir)
	// Create files in the temp directory
	dummyData := []byte("file")
	file1 := filepath.Join(tempDir, "file1.txt")
	err = os.WriteFile(file1, dummyData, 0666)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %s", err)
	}
	file2 := filepath.Join(tempDir, "file2.txt")
	err = os.WriteFile(file2, dummyData, 0666)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %s", err)
	}

	// Test the ListFilesInDir function
	files, err := util.ListFilesInDir(tempDir)
	if err != nil {
		t.Fatalf("ListFilesInDir returned an error: %s", err)
	}

	// Sort the files slice to ensure consistent ordering for test
	sort.Strings(files)

	// Check if the returned files are correct
	expectedFiles := []string{"file1.txt", "file2.txt"}
	if len(files) != len(expectedFiles) {
		t.Fatalf("Expected %d files, got %d", len(expectedFiles), len(files))
	}
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %s", err)
	}
	for i, file := range files {
		tmpFilePath := filepath.Join(tempDir, file)
		relativePath, err := filepath.Rel(cwd, tmpFilePath)
		if err != nil {
			t.Fatalf("Failed to get relative path: %s", err)
		}
		if file != relativePath {
			t.Errorf("Expected file at index %d to be %s, got %s", i, relativePath, file)
		}
	}

	// Test the ListFilesInDir function with an invalid directory
	_, err = util.ListFilesInDir("/nonexistent")
	if err == nil {
		t.Error("Expected error when calling ListFilesInDir with nonexistent directory, got nil")
	}
}
