package encrypte

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
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
