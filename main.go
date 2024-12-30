package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

type (
	OAIRequest struct {
		Model    string    `json:"model"`
		Stream   bool      `json:"stream"`
		Messages []Message `json:"messages"`
	}
	Message struct {
		Role    string          `json:"role"`
		Content json.RawMessage `json:"content"`
	}
	StringContent = string
	MediaContent  struct {
		Type       string `json:"type"`
		Text       string `json:"text,omitempty"`
		ImageUrl   any    `json:"image_url,omitempty"`
		InputAudio any    `json:"input_audio,omitempty"`
	}
	MessageImageUrl struct {
		Url    string `json:"url"`
		Detail string `json:"detail"`
	}
	MessageInputAudio struct {
		Data   string `json:"data"` //base64
		Format string `json:"format"`
	}
)

func main() {
	req := OAIRequest{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{Role: "user", Content: json.RawMessage(`"你好"`)},
			{Role: "assistant", Content: json.RawMessage(`"你好"`)},
		},
	}
	str, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(str))
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
