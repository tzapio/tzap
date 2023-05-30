package cmdutil

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	ignore "github.com/sabhiram/go-gitignore"
	"github.com/tzapio/tzap/internal/logging/tl"
)

type FileInDirEvaluator struct {
	includeMatcher *ignore.GitIgnore
	excludeMatcher *ignore.GitIgnore
}

func NewFileInDirEvaluator() (*FileInDirEvaluator, error) {
	excludePatterns := []string{".git", ".DS_Store", "desktop.ini"}
	ignoreFiles := []string{".tzapignore"}
	if _, err := os.Stat(".gitignore"); err == nil {
		ignoreFiles = append(ignoreFiles, ".gitignore")
	}
	patterns, err := ReadFilterPatternFiles(ignoreFiles...)
	if err != nil {
		return nil, err
	}
	excludePatterns = append(excludePatterns, patterns...)

	includePatterns, err := ReadFilterPatternFiles(".tzapinclude")
	if err != nil {
		return nil, err
	}

	return NewFileInDirEvaluatorWithPatterns(excludePatterns, includePatterns), nil
}
func NewFileInDirEvaluatorWithPatterns(excludePatterns []string, includePatterns []string) *FileInDirEvaluator {
	excludeMatcher := ignore.CompileIgnoreLines(excludePatterns...)
	includeMatcher := ignore.CompileIgnoreLines(includePatterns...)
	return &FileInDirEvaluator{excludeMatcher: excludeMatcher, includeMatcher: includeMatcher}
}
func (e *FileInDirEvaluator) ShouldKeepPath(path string) bool {
	exclude := e.excludeMatcher.MatchesPath(path)
	include := e.includeMatcher.MatchesPath(path)
	return include && !exclude
}
func (e *FileInDirEvaluator) ShouldTraverseDir(path string) bool {
	return !e.excludeMatcher.MatchesPath(path)
}
func (e *FileInDirEvaluator) WalkDir(dir string) ([]string, error) {
	var list []string
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		tl.Logger.Println("WALKDIR", dir, path, d.Name())
		if err != nil {
			return err
		}
		if d.IsDir() {
			if e.ShouldTraverseDir(path) || d.Name() == "." {
				tl.Logger.Println("KEEPDIR", path)
				return nil
			} else {
				tl.Logger.Println("SKIPDIR", path)
				return filepath.SkipDir
			}
		} else if e.ShouldKeepPath(path) {
			tl.Logger.Println("KEEPFILE", path)
			absPath, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			relPath, err := filepath.Rel(cwd, absPath)
			if err != nil {
				return err
			}
			relPath = strings.TrimPrefix(relPath, "./")
			list = append(list, filepath.ToSlash(relPath))
			return nil
		}
		tl.Logger.Println("SKIPFILE", path)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return list, nil
}
