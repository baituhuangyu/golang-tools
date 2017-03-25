package common

import (
	"crypto/sha256"
	"io"
	"fmt"
)

// 计算sha256
func Hash256String(rawString string) string {
	// sha256 计算hashValue
	h := sha256.New()
	io.WriteString(h, rawString)
	return fmt.Sprintf("%X", h.Sum(nil))
}
