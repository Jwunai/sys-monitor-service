package monitor

import (
	"context"
	"log"
	"sync"

	"github.com/Jwunai/sys-monitor-service/configs/monitor_config"
	"github.com/Jwunai/sys-monitor-service/internal/interfaces"
	"github.com/Jwunai/sys-monitor-service/pkg"
)

// Manager结构体定义
type Manager struct {
	ctx          context.Context           // 退出上下文
	cancel       context.CancelFunc        // 取消函数
	serverName   string                    // 服务器名称（告警标识）
	cpuCfg       monitor_config.CPUConfig  // CPU专属配置（匹配monitor_config包）
	memCfg       monitor_config.MemConfig  // 内存专属配置
	diskCfg      monitor_config.DiskConfig // 磁盘专属配置
	alertSenders []interfaces.AlertSender  // 所有启用的告警实例
	wg           sync.WaitGroup            // 协程等待组
}

// 匹配monitor_config包，且字段名大写
func NewManager(
	serverName string,
	cpuCfg monitor_config.CPUConfig,
	memCfg monitor_config.MemConfig,
	diskCfg monitor_config.DiskConfig,
	alertSenders []interfaces.AlertSender,
) *Manager {
	ctx, cancel := context.WithCancel(context.Background())
	return &Manager{
		ctx:          ctx,
		cancel:       cancel,
		serverName:   serverName,
		cpuCfg:       cpuCfg,
		memCfg:       memCfg,
		diskCfg:      diskCfg,
		alertSenders: alertSenders,
	}
}

// Start 启动所有监控协程
func (m *Manager) Start() {
	log.Printf(
		"监控服务启动 | 服务器名称: %s | 系统类型: %s | CPU间隔: %v | 内存间隔: %v | 磁盘间隔: %v",
		m.serverName, pkg.GetOS(),
		m.cpuCfg.Interval, m.memCfg.Interval, m.diskCfg.Interval,
	)

	m.wg.Add(1)
	go m.monitorCPU()

	m.wg.Add(1)
	go m.monitorMemory()

	m.wg.Add(1)
	go m.monitorDisk()
}

// Stop 停止所有监控协程
func (m *Manager) Stop() {
	m.cancel()
	m.wg.Wait()
	log.Println("监控服务已完全停止")
}

// sendAlerts 通用异步告警方法
func (m *Manager) sendAlerts(title, content string) {
	if len(m.alertSenders) == 0 {
		log.Println("⚠️  无启用的告警渠道，跳过告警发送")
		return
	}

	for _, sender := range m.alertSenders {
		go func(s interfaces.AlertSender) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("❌ 告警[%s]协程panic: %v", s.Name(), r)
				}
			}()

			err := s.SendAlert(title, m.serverName, content)
			if err != nil {
				log.Printf("❌ 告警[%s]发送失败: %v", s.Name(), err)
			} else {
				log.Printf("✅ 告警[%s]发送成功", s.Name())
			}
		}(sender)
	}
}
