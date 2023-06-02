package cmdutil

import (
	"os"

	ignore "github.com/sabhiram/go-gitignore"
)

func FilterWithExcludePattern(inFiles []string, excludePattern []string) (outFiles []string) {
	object := ignore.CompileIgnoreLines(excludePattern...)

	filteredFiles := make(map[string]struct{}, len(inFiles)/2)
	for _, file := range inFiles {
		if object.MatchesPath(file) {
			filteredFiles[file] = struct{}{}
		}
	}
	for _, file := range inFiles {
		if _, filtered := filteredFiles[file]; !filtered {
			outFiles = append(outFiles, file)
		}
	}
	return outFiles
}
func FilterFileSize(inFiles []string, maxSize int) (outFiles []string) {
	filteredFiles := map[string]struct{}{}
	for _, file := range inFiles {
		if stat, err := os.Stat(file); err == nil && stat.Size() > int64(maxSize) {
			filteredFiles[file] = struct{}{}
		}
	}
	for _, file := range inFiles {
		if _, filtered := filteredFiles[file]; !filtered {
			outFiles = append(outFiles, file)
		}
	}
	return outFiles
}
func FilterWithIncludePattern(inFiles []string, excludePattern []string) (outFiles []string) {
	object := ignore.CompileIgnoreLines(excludePattern...)
	for _, file := range inFiles {
		if object.MatchesPath(file) {
			outFiles = append(outFiles, file)
		}
	}
	return outFiles
}
