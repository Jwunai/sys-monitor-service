// configs/alerts/email.go
package alert_config

// EmailConfig 邮箱告警专属配置
type EmailConfig struct {
	From     string   `yaml:"from"`      // 发件人邮箱
	Password string   `yaml:"password"`  // 邮箱授权码
	SmtpHost string   `yaml:"smtp_host"` // SMTP服务器地址
	SmtpPort int      `yaml:"smtp_port"` // SMTP端口（465/587）
	To       []string `yaml:"to"`        // 收件人列表
}
