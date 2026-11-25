package crypto

import (
	"crypto/md5"
	"encoding/hex"
)

// MD5 computes the MD5 hash of a string and returns its hexadecimal representation.
// Note that MD5 is retained only for generating cache keys and must not be used for security purposes.
func MD5(s string) string {
	sum := md5.Sum([]byte(s))
	return hex.EncodeToString(sum[:])
}
