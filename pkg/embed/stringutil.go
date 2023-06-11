package embed

import "strings"

func AddEmbedHeader(filename string, splitPart string) string {
	return "####embedding from file: " + filename + "\n" + splitPart
}

func StripEmbedHeader(splitPart string) string {
	return splitPart[strings.Index(splitPart, "\n")+1:]
}
