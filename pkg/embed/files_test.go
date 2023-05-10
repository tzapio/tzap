//file: pkg/embed/files_test.go

package embed_test

/*
TODO requires refactor to automate

func Test_GetGoFilesInDir_givenEmptyDir_expectNoFiles(t *testing.T) {
	dir := "testdata/empty"
	expected := 0

	files := embed.GetGoFilesInDir(dir)

	if expected != len(files) {
		t.Errorf("Expected %d files, got %d", expected, len(files))
	}

}

func Test_GetGoFilesInDir_givenGoDir_expectGoFiles(t *testing.T) {
	dir := "testdata/go_files"
	expected := 2

	files := embed.GetGoFilesInDir(dir)
	if expected != len(files) {
		t.Errorf("Expected %d files, got %d", expected, len(files))
	}

}

func Test_ProcessFile_givenGoFile_expectTokensAndLines(t *testing.T) {
	filename := "testdata/go_files/file1.go"
	tzap := tzap.InternalNew()

	fileTokens, lines, _, err := embed.ProcessFile(filename, tzap)
	if err != nil {
		t.Errorf("Error processing file: %s", err)
	}

	if fileTokens != 14 {
		t.Errorf("Expected %d tokens, got %d", 14, fileTokens)
	}
	if lines == 9 {
		t.Errorf("Expected %d lines, got %d", 9, lines)
	}
}

func Test_ProcessFile_givenInvalidFile_expectError(t *testing.T) {
	filename := "testdata/invalid/file.go"
	tzap := tzap.InternalNew()

	_, _, _, err := embed.ProcessFile(filename, tzap)
	if err == nil {
		t.Errorf("Expected error processing file: %s", filename)
	}
}

func Test_ProcessOffset_givenGoFileAndOffset_expectVector(t *testing.T) {
	filename := "testdata/go_files/file1.go"
	content := `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}
`
	start := 0
	end := 10
	fileTokens := 14
	tzap := tzap.InternalNew()

	vector, err := embed.ProcessOffset(tzap, filename, content, start, end, fileTokens)

	if err != nil {
		t.Errorf("Error processing offset: %s", err)
	}
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if vector.ID == "" {
		t.Error("Expected vector.id to be set")
	}

	if len(vector.Metadata) == 0 {
		t.Error("Expected vector Metadata not to be empty")
	}
}
*/
