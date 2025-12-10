// internal/alert_services/register.go
package alert_services

import (
	"github.com/Jwunai/sys-monitor-service/configs"
	"github.com/Jwunai/sys-monitor-service/configs/alert_config"
	"github.com/Jwunai/sys-monitor-service/internal/interfaces"
)

// 1. 全局注册器：key=告警渠道名，value=实例化函数
var alertRegistry = make(map[string]func(cfg interface{}) interfaces.AlertSender)

// 2. 注册方法
func Register(name string, fn func(cfg interface{}) interfaces.AlertSender) {
	alertRegistry[name] = fn
}

// 3. 初始化自动注册钉钉/邮箱
func init() {
	// 注册钉钉：入参是*alert_config.DingTalkConfig
	Register("dingtalk", func(cfg interface{}) interfaces.AlertSender {
		dingCfg, ok := cfg.(*alert_config.DingTalkConfig)
		if !ok {
			return nil
		}
		return NewDingTalk(dingCfg)
	})

	// 注册邮箱：入参是*alert_config.EmailConfig
	Register("email", func(cfg interface{}) interfaces.AlertSender {
		emailCfg, ok := cfg.(*alert_config.EmailConfig)
		if !ok {
			return nil
		}
		return NewEmail(emailCfg)
	})

}

// 4. 创建所有启用告警
// 参数是*configs.AlertConfig（全局告警配置）
func GetAllEnabled(alertCfg *configs.AlertConfig) []interfaces.AlertSender {
	var senders []interfaces.AlertSender

	// 从注册器获取钉钉实例
	if creator, ok := alertRegistry["dingtalk"]; ok {
		dingTalk := creator(&alertCfg.DingTalk) // 传*alert_config.DingTalkConfig
		if dingTalk != nil && dingTalk.IsEnabled() {
			senders = append(senders, dingTalk)
		}
	}

	// 从注册器获取邮箱实例
	if creator, ok := alertRegistry["email"]; ok {
		email := creator(&alertCfg.Email) // 传*alert_config.EmailConfig
		if email != nil && email.IsEnabled() {
			senders = append(senders, email)
		}
	}

	return senders
}
