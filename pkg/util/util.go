package util

func CreateSpaces(n int) string {
	spaces := ""
	for i := 0; i < n; i++ {
		spaces += " "
	}
	return spaces
}
