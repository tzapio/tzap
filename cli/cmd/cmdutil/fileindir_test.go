package cmdutil_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tzapio/tzap/cli/cmd/cmdutil"
)

func TestWalkDir(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}
	dir, err := ioutil.TempDir("", "testdir")
	if err != nil {
		t.Fatalf("failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(dir)
	defer os.Chdir(cwd)
	err = os.Chdir(dir)
	if err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}
	err = ioutil.WriteFile(filepath.Join(dir, "testfile.txt"), []byte("test"), 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	err = ioutil.WriteFile(filepath.Join(dir, "qe.go"), []byte("test"), 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	err = ioutil.WriteFile(filepath.Join(dir, ".gitignore"), []byte("testfile.txt"), 0644)
	if err != nil {
		t.Fatalf("failed to create test .gitignore file: %v", err)
	}
	err = ioutil.WriteFile(filepath.Join(dir, ".tzapignore"), []byte("testfile.txt"), 0644)
	if err != nil {
		t.Fatalf("failed to create test .gitignore file: %v", err)
	}
	err = ioutil.WriteFile(filepath.Join(dir, ".tzapinclude"), []byte("*.go"), 0644)
	if err != nil {
		t.Fatalf("failed to create test .gitignore file: %v", err)
	}
	evaluator, err := cmdutil.NewFileInDirEvaluator()
	if err != nil {
		t.Fatalf("error creating FileInDirEvaluator: %v", err)
	}

	expected := []string{"qe.go"}

	result, err := evaluator.WalkDir(dir)
	if err != nil {
		t.Fatalf("error walking directory: %v", err)
	}

	if len(result) != len(expected) {
		t.Fatalf("expected %d entries, but got %d", len(expected), len(result))
	}

	for i, file := range expected {
		if file != result[i] {
			t.Fatalf("expected %s, but got %s", file, result[i])
		}
	}
}

func Test_shouldTraverseDir(t *testing.T) {
	evaluator := cmdutil.NewFileInDirEvaluatorWithPatterns([]string{"exclude"}, []string{"*.txt"})
	tests := map[string]bool{
		"exclude":         false,
		"include":         true,
		"exclude/exc1":    false,
		"include/inc2":    true,
		"exclude2/inc2":   true,
		"exclude/include": false,
	}
	for path, expected := range tests {
		actual := evaluator.ShouldTraverseDir(path)
		assert.Equal(t, expected, actual, "Path '%s' should have been %t", path, expected)
	}

}
func Test_shouldKeepPath(t *testing.T) {
	evaluator := cmdutil.NewFileInDirEvaluatorWithPatterns([]string{"exclude"}, []string{"*.txt"})
	tests := map[string]bool{
		"exclude/file.txt":  false,
		"include/file.txt":  true,
		"exclude/exc1":      false,
		"include2/inc2.txt": true,
		"exclude/include":   false,
	}
	for path, expected := range tests {
		actual := evaluator.ShouldKeepPath(path)
		assert.Equal(t, expected, actual, "Path '%s' should have been %t", path, expected)
	}

}
