package datautils

import "strings"

// MaskEmail 隐私化处理邮箱
func MaskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email // 如果不是有效的邮箱格式，直接返回原字符串
	}

	username := parts[0]
	domain := parts[1]

	// 将用户名部分替换为 *
	maskedUsername := string(username[0]) + strings.Repeat("*", len(username)-1)

	return maskedUsername + "@" + domain
}

// ProcessEmails 批量处理邮箱数据
func ProcessEmails(emails []string) []string {
	var maskedEmails []string
	for _, email := range emails {
		maskedEmails = append(maskedEmails, MaskEmail(email))
	}
	return maskedEmails
}
