package encrypte

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
)

func Md5Encode(inputStr string) string {
	md5Hash := md5.New()
	md5Hash.Write([]byte(inputStr))
	md5Bytes := md5Hash.Sum(nil)
	return hex.EncodeToString(md5Bytes)
}

func Md5Verify(inputStr string, hashStr string) bool {
	return Md5Encode(inputStr) == hashStr
}

// Sha256Encode 使用 SHA-256 算法对输入字符串进行编码
func Sha256Encode(inputStr string) string {
	hash := sha256.New()
	hash.Write([]byte(inputStr))
	return hex.EncodeToString(hash.Sum(nil))
}

// Sha256Verify 验证输入字符串的 SHA-256 哈希值是否与给定的哈希值匹配
func Sha256Verify(inputStr string, hashStr string) bool {
	return Sha256Encode(inputStr) == hashStr
}

// HashPassword bcrypt.DefaultCost = 14 // 设置 bcrypt 的成本
// 加密用户密码
func HashPassword(password string) (string, error) {
	// 使用 bcrypt 生成密码哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// VerifyPassword 验证用户密码
func VerifyPassword(hashedPassword, password string) error {
	// 比较哈希后的密码和原始密码是否匹配
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
