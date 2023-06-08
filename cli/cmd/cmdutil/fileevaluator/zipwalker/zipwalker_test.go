package zipwalker_test

import (
	"archive/zip"
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tzapio/tzap/cli/cmd/cmdutil/fileevaluator"
	"github.com/tzapio/tzap/cli/cmd/cmdutil/fileevaluator/zipwalker"
)

func TestWalkDirFromURL(t *testing.T) {
	// Create a test ZIP file
	testZipContent, err := createTestZipFile()
	if err != nil {
		t.Fatalf("failed to create test ZIP file: %v", err)
	}

	// Create a test HTTP server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/zip")
		w.Write(testZipContent)
	}))
	defer ts.Close()

	// Instantiate a FileInDirEvaluator
	evaluator := fileevaluator.NewWithPatterns([]string{}, []string{"*.txt"})
	zipwalker := zipwalker.New(evaluator, "/", ts.URL)
	// Call WalkDirFromURL with the test server URL
	result, err := zipwalker.GetFiles()
	if err != nil {
		t.Fatalf("error walking directory from URL: %v", err)
	}

	expected := []string{"file1.txt", "file2.txt"}
	for i, file := range expected {
		assert.Equal(t, file, result[i].FilePath())
	}

}

func createTestZipFile() ([]byte, error) {
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	files := []struct {
		Name    string
		Content string
	}{
		{"file1.txt", "File 1 Content"},
		{"file2.txt", "File 2 Content"},
		{"file3.ignored", "File 3 Ignored"},
	}

	for _, file := range files {
		fWriter, err := zipWriter.Create(file.Name)
		if err != nil {
			return nil, err
		}
		_, err = fWriter.Write([]byte(file.Content))
		if err != nil {
			return nil, err
		}
	}

	err := zipWriter.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
