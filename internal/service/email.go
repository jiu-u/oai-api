package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/jiu-u/oai-api/internal/repository"
	"net/smtp"
)

type EmailService interface {
	SendEmail(ctx context.Context, to string, subject string, body string) error
}

func NewEmailService(repo repository.SystemRepository) EmailService {
	return &emailService{repo: repo}
}

type emailService struct {
	*Service
	repo repository.SystemRepository
}

func (s *emailService) SendEmail(ctx context.Context, to string, subject string, body string) error {
	// 获取配置
	config2, err := s.repo.GetEmailConfig(ctx)
	if err != nil {
		return errors.New("获取邮件配置失败")
	}
	// 构建邮件头
	headers := make(map[string]string)
	headers["From"] = config2.Username
	headers["To"] = to
	headers["Subject"] = subject
	headers["Content-Type"] = "text/html; charset=UTF-8"
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body
	auth := smtp.PlainAuth("", config2.Username, config2.Password, config2.Host)
	addr := fmt.Sprintf("%s:%d", config2.Host, config2.Port)
	err = smtp.SendMail(
		addr,
		auth,
		config2.Username,
		[]string{to},
		[]byte(message),
	)
	if err != nil {
		return fmt.Errorf("发送邮件失败: %v", err)
	}
	return nil
}
