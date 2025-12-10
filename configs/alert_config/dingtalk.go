// configs/alerts_config/dingtalk.go
package alert_config

// DingTalkConfig 钉钉告警专属配置
type DingTalkConfig struct {
	Token  string `yaml:"token"`  // 钉钉机器人token
	Secret string `yaml:"secret"` // 钉钉机器人secret
}
