package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

// GenerateSalt 生成随机盐值?
func GenerateSalt() (string, error) {
	bytes := make([]byte, 16) // 16 字节的随机盐值?
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// HashPassword 使用 SHA256 和盐值加密密码？
func HashPassword(password, salt string) string {
	hash := sha256.Sum256([]byte(password + salt))
	return hex.EncodeToString(hash[:])
}

// VerifyPassword 验证密码是否正确
func VerifyPassword(password, hashedPassword, salt string) bool {
	// 使用相同的盐和哈希算法计算输入密码的哈希值
	computedHash := HashPassword(password, salt)
	// 比较计算出的哈希值与存储的哈希值
	return computedHash == hashedPassword
}
