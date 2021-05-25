package handlers

import (
	"crypto/md5"
	"encoding/hex"
)

// Generate an MD5 hash of a supplied string, return as a hex string
func md5hex(s string) string {
	hash := md5.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
}