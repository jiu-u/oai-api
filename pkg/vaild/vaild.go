package vaild

import "regexp"

// IsValidEmail 验证电子邮件格式是否合法
func IsValidEmail(email string) bool {
	// 定义正则表达式，匹配标准的邮箱格式
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

// IsValidUsername 检查用户名是否合法
func IsValidUsername(username string) bool {
	// 使用正则表达式进行匹配
	var validUsernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,32}$`)

	// 检查长度并验证正则表达式
	return validUsernameRegex.MatchString(username)
}
