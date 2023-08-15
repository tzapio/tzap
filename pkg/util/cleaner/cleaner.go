package cleaner

import (
	"regexp"
	"strings"
)

func FileWriteClean(content string) string {
	//Strip first line if it starts with: ####
	if len(content) > 4 && content[0:4] == "####" {
		content = content[strings.Index(content, "\n")+1:]
	}
	pattern := "```[^\\n]*\\n?([^!```]*)```"
	re := regexp.MustCompile(pattern)

	matches := re.FindAllStringSubmatch(content, -1)
	if len(matches) > 1 {
		return content
	}

	for _, match := range matches {
		if len(match) >= 2 {
			return (match[1])
		}

	}
	return content
}
