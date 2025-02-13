package utils

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateRandomSecret 生成随机的密钥
func GenerateRandomSecret() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
