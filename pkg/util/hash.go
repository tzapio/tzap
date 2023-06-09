package util

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
func MD5HashByte(text []byte) string {
	hash := md5.Sum(text)
	return hex.EncodeToString(hash[:])
}
