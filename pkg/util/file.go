package util

import (
	"fmt"
	"os"
	"path/filepath"
)

func ReadFile(filepath string) (string, error) {
	b, err := os.ReadFile(filepath)
	if err != nil {
		return "", fmt.Errorf("cannot read file %s", filepath)
	}
	return string(b), err
}
func ReadFileP(filepath string) string {
	b, err := ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func ListFilesInDir(dir string) []string {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			panic(err)
		}
		if !info.IsDir() {
			absPath, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			files = append(files, filepath.ToSlash(absPath))
		}
		return nil
	})

	if err != nil {
		panic(err)

	}

	return files
}
