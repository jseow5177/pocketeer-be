package goutil

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"

	"github.com/rs/xid"
)

func RandByte(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func HMACSha256(str string, secret []byte) ([]byte, error) {
	h := hmac.New(sha256.New, secret)
	_, err := h.Write([]byte(str))
	if err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

func Base64Encode(b []byte) string {
	return base64.StdEncoding.WithPadding(base64.NoPadding).EncodeToString(b)
}

func Base64Decode(s string) ([]byte, error) {
	return base64.StdEncoding.WithPadding(base64.NoPadding).DecodeString(s)
}

func NextXID() string {
	return xid.New().String()
}

func ContainString(slice []string, target string) bool {
	for _, v := range slice {
		if target == v {
			return true
		}
	}

	return false
}

func RemoveDuplicateString(slice []string) []string {
	var (
		m   = make(map[string]bool)
		arr = make([]string, 0)
	)
	for _, s := range slice {
		if _, ok := m[s]; !ok {
			m[s] = true
			arr = append(arr, s)
		}
	}
	return arr
}
