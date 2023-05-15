package cmdutil

import (
	"bufio"
	"os"
	"strings"
)

func ReadFilterPatternFiles(patternFilePath ...string) ([]string, error) {
	ignorePatternss := [][]string{}
	for _, path := range patternFilePath {
		ignorePatterns, err := readPatternFile(path)
		if err != nil {
			return nil, err
		}
		ignorePatternss = append(ignorePatternss, ignorePatterns)
	}
	mergedPatterns := mergeFilterPatterns(ignorePatternss...)
	return mergedPatterns, nil
}
func readPatternFile(patternFilePath string) ([]string, error) {
	file, err := os.Open(patternFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	filterPatterns := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if len(line) > 0 && !strings.HasPrefix(line, "#") {
			filterPatterns = append(filterPatterns, line)
		}
	}

	return filterPatterns, scanner.Err()
}

func mergeFilterPatterns(patternGroups ...[]string) []string {
	mergedPatterns := make([]string, 0)
	existing := make(map[string]bool)

	for _, patterns := range patternGroups {
		for _, pattern := range patterns {
			if !existing[pattern] {
				mergedPatterns = append(mergedPatterns, pattern)
				existing[pattern] = true
			}
		}
	}

	return mergedPatterns
}
