package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/jiu-u/oai-api/pkg/sid"
	"strings"
)

func main() {
	sid := sid.NewSid()
	fmt.Println(sid.GenUint64())

}

func GenerateOpenAIKey() string {
	// 生成 32 字节的随机数据
	randomBytes := make([]byte, 50)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err) // 处理随机数生成错误
	}

	// 将随机数据编码为 Base64
	key := base64.StdEncoding.EncodeToString(randomBytes)

	// 移除 Base64 中的特殊字符（如 '+' 和 '/'），并截取前 48 个字符
	key = strings.ReplaceAll(key, "+", "")
	key = strings.ReplaceAll(key, "/", "")
	key = key[:48]

	// 添加 "sk-" 前缀
	return "sk-" + key
}
