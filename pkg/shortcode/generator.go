package shortcode

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
)

const (
	base62Chars  = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	shortCodeLen = 7
)

// Generate returns a 7-char Base62 shortcode derived from the SHA-256 of the input.
// Pass a salt (e.g. strconv.Itoa(attempt)) to resolve collisions.
func Generate(longURL string, salt string) string {
	h := sha256.Sum256([]byte(longURL + salt))
	n := binary.BigEndian.Uint64(h[:8])
	return base62Encode(n)[:shortCodeLen]
}

func base62Encode(n uint64) string {
	if n == 0 {
		return string(base62Chars[0])
	}
	result := make([]byte, 0, 11)
	base := uint64(len(base62Chars))
	for n > 0 {
		result = append([]byte{base62Chars[n%base]}, result...)
		n /= base
	}
	// pad to at least shortCodeLen characters
	for len(result) < shortCodeLen {
		result = append([]byte{base62Chars[0]}, result...)
	}
	return fmt.Sprintf("%s", result)
}
