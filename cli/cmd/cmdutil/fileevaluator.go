package cmdutil

import (
	"os"

	ignore "github.com/sabhiram/go-gitignore"
)

type FileEvaluator struct {
	includeMatcher *ignore.GitIgnore
	excludeMatcher *ignore.GitIgnore
}

func NewFileEvaluator() (*FileEvaluator, error) {
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

	return NewFileEvaluatorWithPatterns(excludePatterns, includePatterns), nil
}
func NewFileEvaluatorWithPatterns(excludePatterns []string, includePatterns []string) *FileEvaluator {
	excludeMatcher := ignore.CompileIgnoreLines(excludePatterns...)
	includeMatcher := ignore.CompileIgnoreLines(includePatterns...)
	return &FileEvaluator{excludeMatcher: excludeMatcher, includeMatcher: includeMatcher}
}
func (e *FileEvaluator) ShouldKeepPath(path string) bool {
	exclude := e.excludeMatcher.MatchesPath(path)
	include := e.includeMatcher.MatchesPath(path)
	return include && !exclude
}
