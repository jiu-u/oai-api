package main

import (
	"fmt"
	"math/rand"
	"net/smtp"
	"strings"
	"time"
)

// EmailConfig 邮件配置
type EmailConfig struct {
	Host     string // SMTP服务器地址
	Port     int    // SMTP服务器端口
	Username string // 发件人邮箱
	Password string // 发件人密码或授权码
}

// SendVerificationCode 发送验证码邮件
func SendVerificationCode(config EmailConfig, toEmail string) (string, error) {
	// 生成6位随机验证码
	code := generateVerificationCode()

	// 构建邮件内容
	subject := "验证码"
	body := fmt.Sprintf(`
        <html>
        <body>
            <h3>您的验证码是：%s</h3>
            <p>验证码有效期为5分钟，请尽快使用。</p>
            <p>如果这不是您的操作，请忽略此邮件。</p>
        </body>
        </html>
    `, code)

	// 构建邮件头
	headers := make(map[string]string)
	headers["From"] = config.Username
	headers["To"] = toEmail
	headers["Subject"] = subject
	headers["Content-Type"] = "text/html; charset=UTF-8"

	// 组装邮件内容
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// 连接SMTP服务器并发送邮件
	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)

	err := smtp.SendMail(
		addr,
		auth,
		config.Username,
		[]string{toEmail},
		[]byte(message),
	)

	if err != nil {
		return "", fmt.Errorf("发送邮件失败: %v", err)
	}

	return code, nil
}

// generateVerificationCode 生成6位随机验证码
func generateVerificationCode() string {
	rand.Seed(time.Now().UnixNano())
	digits := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	code := make([]string, 6)
	for i := 0; i < 6; i++ {
		code[i] = digits[rand.Intn(len(digits))]
	}
	return strings.Join(code, "")
}

func main() {
	// 邮件配置示例
	config := EmailConfig{
		Host:     "mail14", // 以163邮箱为例
		Port:     587,
		Username: "example@mail.com",
		Password: "examplePassword", // 使用邮箱的授权码
	}

	// 发送验证码
	toEmail := "jiu999xyz666@2925.com"
	code, err := SendVerificationCode(config, toEmail)
	if err != nil {
		fmt.Printf("发送失败：%v\n", err)
		return
	}
	fmt.Printf("验证码 %s 已发送到 %s\n", code, toEmail)
}
