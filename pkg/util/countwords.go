package util

import "regexp"

func CountWords(text string) int {
	wordSplitter := regexp.MustCompile(`(?:[\s,./()\-_]+|\w{3,}|[A-Z]+)`)
	words := wordSplitter.Split(text, -1)
	return len(words)
}

func removeEmptyStrings(strings []string) []string {
	var result []string
	for _, str := range strings {
		if str != "" {
			result = append(result, str)
		}
	}
	return result
}
