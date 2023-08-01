package fileevaluator

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ReadFilterPatternFiles(patternFilePath ...string) ([]string, error) {
	ignorePatternss := [][]string{}
	var errs []error
	for _, path := range patternFilePath {
		ignorePatterns, err := readPatternFile(path)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		ignorePatternss = append(ignorePatternss, ignorePatterns)
	}
	mergedPatterns := mergeFilterPatterns(ignorePatternss...)
	if len(errs) > 0 {
		err := fmt.Errorf("failed to read filter patterns: %v", errs)
		return mergedPatterns, err
	}
	return mergedPatterns, nil
}
func readPatternFile(patternFilePath string) ([]string, error) {
	content, err := os.ReadFile(patternFilePath)
	if err != nil {
		return nil, err
	}
	return ReadPatternString(string(content))
}
func ReadPatternString(content string) ([]string, error) {
	filterPatterns := []string{}
	scanner := bufio.NewScanner(strings.NewReader(content))
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
