package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func RandomStr(length int) string {
	l := length

	if length <= 0 {
		l = 10
	}

	b := make([]byte, l/2)
	_, _ = rand.Read(b)
	s := hex.EncodeToString(b)
	return s
}
