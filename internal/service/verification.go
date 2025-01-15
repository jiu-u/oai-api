package service

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type VerificationService interface {
	SendEmailVerificationCode(ctx context.Context, to string) error
	CheckEmailVerificationCode(ctx context.Context, to string, code string) error
}

func NewVerificationService(s *Service, emailSvc EmailService) VerificationService {
	return &verificationService{
		Service:  s,
		emailSvc: emailSvc,
	}
}

type verificationService struct {
	*Service
	emailSvc EmailService
}

func (s *verificationService) SendEmailVerificationCode(ctx context.Context, to string) error {
	// 流量限制
	// 发送验证码
	code := generateVerificationCode()
	body := fmt.Sprintf(`
        <html>
        <body>
            <h3>您的验证码是：%s</h3>
            <p>验证码有效期为5分钟，请尽快使用。</p>
            <p>如果这不是您的操作，请忽略此邮件。</p>
        </body>
        </html>
    `, code)
	err := s.emailSvc.SendEmail(ctx, to, "验证码", body)
	if err != nil {
		return err
	}
	s.Cache.Set("email_"+to, code, 5*time.Minute)
	return nil
}

func (s *verificationService) CheckEmailVerificationCode(ctx context.Context, to string, code string) error {
	// 流量限制
	// 验证验证码
	vCode, exist := s.Cache.Get("email_" + to)
	if !exist || vCode == nil {
		return fmt.Errorf("验证码已过期")
	}
	if vCode.(string) != code {
		return fmt.Errorf("验证码错误")
	}
	// 清除缓存
	s.Cache.Delete("email_" + to)
	return nil
}

// generateVerificationCode 生成6位随机验证码
func generateVerificationCode() string {
	digits := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	code := make([]string, 6)
	for i := 0; i < 6; i++ {
		code[i] = digits[rand.Intn(len(digits))]
	}
	return strings.Join(code, "")
}
