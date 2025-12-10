// internal/alert_services/email.go
package alert_services

import (
	"bytes"
	"fmt"
	"net/smtp"
	"strings"
	"time"

	"github.com/Jwunai/sys-monitor-service/configs/alert_config" // 替换为实际module名
	"github.com/Jwunai/sys-monitor-service/internal/interfaces"
)

// Email 邮箱告警实现
type Email struct {
	cfg *alert_config.EmailConfig // 邮箱配置
}

// NewEmail 创建邮箱告警实例
func NewEmail(cfg *alert_config.EmailConfig) interfaces.AlertSender {
	return &Email{cfg: cfg}
}

// Name 返回告警渠道名称
func (e *Email) Name() string {
	return "邮箱"
}

// IsEnabled 判断是否启用（发件人、SMTP服务器、端口非空）
func (e *Email) IsEnabled() bool {
	return e.cfg != nil && e.cfg.From != "" && e.cfg.SmtpHost != "" && e.cfg.SmtpPort != 0 && len(e.cfg.To) > 0
}

// SendAlert 发送邮箱告警
func (e *Email) SendAlert(title, serverName, content string) error {
	if !e.IsEnabled() {
		return fmt.Errorf("邮箱告警未启用配置缺失")
	}

	// 1. 构造邮件内容
	subject := fmt.Sprintf("【%s】%s", serverName, title)
	emailContent := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>%s</title>
</head>
<body>
    <div style="font-family: Arial, sans-serif; font-size: 14px; line-height: 1.6;">
        <p>%s</p>
        <p>告警时间：%s</p>
    </div>
</body>
</html>
`, subject, content, time.Now().Format("2006-01-02 15:04:05"))

	// 2. 构造SMTP认证信息
	auth := smtp.PlainAuth(
		"",
		e.cfg.From,
		e.cfg.Password,
		e.cfg.SmtpHost,
	)

	// 3. 构造邮件头
	var msg bytes.Buffer
	msg.WriteString(fmt.Sprintf("From: %s\r\n", e.cfg.From))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(e.cfg.To, ",")))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/html; charset=UTF-8\r\n\r\n")
	msg.WriteString(emailContent)

	// 4. 发送邮件
	smtpAddr := fmt.Sprintf("%s:%d", e.cfg.SmtpHost, e.cfg.SmtpPort)
	err := smtp.SendMail(
		smtpAddr,
		auth,
		e.cfg.From,
		e.cfg.To,
		msg.Bytes(),
	)
	if err != nil {
		return fmt.Errorf("发送邮箱告警失败: %w", err)
	}

	return nil
}

var _ interfaces.AlertSender = (*Email)(nil)
