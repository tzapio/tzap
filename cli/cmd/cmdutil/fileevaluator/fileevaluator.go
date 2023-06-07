package fileevaluator

import (
	"os"
	"path"

	ignore "github.com/sabhiram/go-gitignore"
	"github.com/tzapio/tzap/internal/logging/tl"
)

type FileEvaluator struct {
	includeMatcher *ignore.GitIgnore
	excludeMatcher *ignore.GitIgnore
}

var baseExcludePatterns = []string{".git", ".DS_Store", "desktop.ini"}

func New(baseDir string) (*FileEvaluator, error) {
	gitIgnorePath := path.Join(baseDir, ".gitignore")
	tzapIgnorePath := path.Join(baseDir, ".tzapignore")
	tzapIncludePath := path.Join(baseDir, ".tzapinclude")
	ignoreFiles := []string{tzapIgnorePath}
	tl.Logger.Println("gitIgnorePath", gitIgnorePath)
	tl.Logger.Println("tzapIgnorePath", tzapIgnorePath)
	tl.Logger.Println("tzapIncludePath", tzapIncludePath)
	if _, err := os.Stat(gitIgnorePath); err == nil {
		ignoreFiles = append(ignoreFiles, gitIgnorePath)
	}
	excludePatternsFromFile, err := ReadFilterPatternFiles(ignoreFiles...)
	if err != nil {
		return nil, err
	}
	excludePatterns := append(baseExcludePatterns, excludePatternsFromFile...)

	includePatternsFromFile, err := ReadFilterPatternFiles(tzapIncludePath)
	if err != nil {
		return nil, err
	}

	return NewWithPatterns(excludePatterns, includePatternsFromFile), nil
}
func NewWithPatterns(excludePatterns []string, includePatterns []string) *FileEvaluator {
	excludeMatcher := ignore.CompileIgnoreLines(excludePatterns...)
	includeMatcher := ignore.CompileIgnoreLines(includePatterns...)
	return &FileEvaluator{excludeMatcher: excludeMatcher, includeMatcher: includeMatcher}
}
func (e *FileEvaluator) ShouldKeepPath(path string) bool {
	exclude := e.excludeMatcher.MatchesPath(path)
	include := e.includeMatcher.MatchesPath(path)
	return include && !exclude
}
func (e *FileEvaluator) ShouldTraverseDir(path string) bool {
	return !e.excludeMatcher.MatchesPath(path)
}
