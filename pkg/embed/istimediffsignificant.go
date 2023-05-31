package embed

import "time"

func isTimePassedSignificant(a, b int64) bool {
	return absInt64(a-b) > int64(time.Second)
}
func absInt64(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}
