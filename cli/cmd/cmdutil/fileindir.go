package cmdutil

func GetNonExcludedFiles(files []string) []string {
	baseIgnore := []string{".git/", ".DS_Store", "desktop.ini"}
	excludePatterns, err := ReadFilterPatternFiles(".tzapignore", ".gitignore")
	if err != nil {
		panic(err)
	}

	excludePatterns = append(excludePatterns, baseIgnore...)
	files = FilterWithExcludePattern(files, excludePatterns)

	includePatterns, err := ReadFilterPatternFiles(".tzapinclude")
	if err != nil {
		panic(err)
	}
	files = FilterWithIncludePattern(files, includePatterns)

	var kilobyte = 1024
	files = FilterFileSize(files, 10000*kilobyte)

	return files
}
