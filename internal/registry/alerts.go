package registry

import (
	"github.com/Jwunai/sys-monitor-service/configs"
	alertSvc "github.com/Jwunai/sys-monitor-service/internal/alert_services"
	"github.com/Jwunai/sys-monitor-service/internal/interfaces"
)

// CreateAllEnabled 自注册方法
func CreateAllEnabled(alertCfg *configs.AlertConfig) []interfaces.AlertSender {
	return alertSvc.GetAllEnabled(alertCfg)
}
