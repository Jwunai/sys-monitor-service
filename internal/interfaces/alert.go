// internal/interfaces/alert.go
package interfaces

// AlertSender 告警发送器通用接口
type AlertSender interface {
	// Name 返回告警渠道名称（如"钉钉"、"邮箱"）
	Name() string
	// SendAlert 发送告警（title：告警标题，serverName：服务器名称，content：告警内容）
	SendAlert(title, serverName, content string) error
	// IsEnabled 判断当前告警渠道是否启用（配置非空）
	IsEnabled() bool
}
