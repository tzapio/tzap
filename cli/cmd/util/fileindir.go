package util

func GetNonExcludedFiles(files []string) []string {
	baseIgnore := []string{".git/", "node_modules", ".DS_Store", "desktop.ini"}
	excludePatterns, err := ReadFilterPatternFiles(".tzapignore", ".gitignore")
	if err != nil {
		panic(err)
	}
	println("total files filter: ", len(files))
	excludePatterns = append(excludePatterns, baseIgnore...)
	files = FilterWithExcludePattern(files, excludePatterns)
	println("after exclude filter: ", len(files))
	includePatterns, err := ReadFilterPatternFiles(".tzapinclude")
	if err != nil {
		panic(err)
	}
	files = FilterWithIncludePattern(files, includePatterns)
	println("after include filter: ", len(files))
	var kilobyte = 1024
	files = FilterFileSize(files, 20*kilobyte)
	println("after size filter: ", len(files))
	return files
}
