package goutil

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"html/template"

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

func Base64Encode(b []byte, padding rune) string {
	return base64.StdEncoding.WithPadding(padding).EncodeToString(b)
}

func Base64Decode(s string, padding rune) ([]byte, error) {
	return base64.StdEncoding.WithPadding(padding).DecodeString(s)
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

func ParseTemplate(t *template.Template, data interface{}) (string, error) {
	b := new(bytes.Buffer)
	if err := t.Execute(b, data); err != nil {
		return "", err
	}

	return b.String(), nil
}
