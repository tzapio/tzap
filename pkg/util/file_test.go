package util_test

import (
	"os"
	"path"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/tzapio/tzap/pkg/util"
)

func TestListFilesInDir(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test_list_files")
	if err != nil {
		t.Fatalf("Error creating temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test files and subdirectories
	filesAndDirs := []string{
		"a.txt",
		"b.txt",
		"subdir/c.txt",
		"subdir/d.txt",
	}
	for _, fd := range filesAndDirs {
		path := filepath.Join(tempDir, fd)
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Error creating subdirectory: %v", err)
		}
		if _, err := os.Create(path); err != nil {
			t.Fatalf("Error creating test file: %v", err)
		}
	}

	expected := []string{
		path.Join(tempDir, "./a.txt"),
		path.Join(tempDir, "./b.txt"),
		path.Join(tempDir, "./subdir/c.txt"),
		path.Join(tempDir, "./subdir/d.txt"),
	}

	result := util.ListFilesInDir(tempDir)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ListFilesInDir returned incorrect result. Expected: %v, Got: %v", expected, result)
	}
}
