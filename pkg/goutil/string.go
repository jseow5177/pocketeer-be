package goutil

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
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
