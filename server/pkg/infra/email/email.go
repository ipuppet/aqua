// Package email 提供基于 SMTP 协议的邮件发送功能。
package email

import (
	"gopkg.in/gomail.v2"

	"github.com/ipuppet/aqua/server/pkg/infra/models"
)

// SendEmailWithConfig 根据传入的配置构建并发送 HTML 格式的邮件。
func SendEmailWithConfig(emailCfg models.EmailConfig, to, subject, body string) error {
	if emailCfg.Password == "" || emailCfg.Host == "" {
		return nil
	}

	mailer := gomail.NewMessage()

	from := emailCfg.From
	if from == "" {
		from = emailCfg.Username
	}

	mailer.SetHeader("From", from)
	mailer.SetHeader("To", to)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", body)

	dialer := gomail.NewDialer(emailCfg.Host, emailCfg.Port, emailCfg.Username, emailCfg.Password)

	// 发送后自动关闭连接。由于此处是单次发送操作，不需要复用长连接。
	if err := dialer.DialAndSend(mailer); err != nil {
		return err
	}
	return nil
}
