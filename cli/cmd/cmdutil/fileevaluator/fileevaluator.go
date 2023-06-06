package fileevaluator

import (
	"os"

	ignore "github.com/sabhiram/go-gitignore"
	"github.com/tzapio/tzap/pkg/types"
)

type FileEvaluator struct {
	includeMatcher *ignore.GitIgnore
	excludeMatcher *ignore.GitIgnore
	name           types.ProjectName
}

var baseExcludePatterns = []string{".git", ".DS_Store", "desktop.ini"}

func New(name types.ProjectName) (*FileEvaluator, error) {
	ignoreFiles := []string{".tzapignore"}
	if _, err := os.Stat(".gitignore"); err == nil {
		ignoreFiles = append(ignoreFiles, ".gitignore")
	}
	excludePatternsFromFile, err := ReadFilterPatternFiles(ignoreFiles...)
	if err != nil {
		return nil, err
	}
	excludePatterns := append(baseExcludePatterns, excludePatternsFromFile...)

	includePatternsFromFile, err := ReadFilterPatternFiles(".tzapinclude")
	if err != nil {
		return nil, err
	}

	return NewWithPatterns(name, excludePatterns, includePatternsFromFile), nil
}
func NewWithPatterns(name types.ProjectName, excludePatterns []string, includePatterns []string) *FileEvaluator {
	excludeMatcher := ignore.CompileIgnoreLines(excludePatterns...)
	includeMatcher := ignore.CompileIgnoreLines(includePatterns...)
	return &FileEvaluator{name: name, excludeMatcher: excludeMatcher, includeMatcher: includeMatcher}
}
func (e *FileEvaluator) ShouldKeepPath(path string) bool {
	exclude := e.excludeMatcher.MatchesPath(path)
	include := e.includeMatcher.MatchesPath(path)
	return include && !exclude
}
