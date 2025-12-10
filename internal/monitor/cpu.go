// internal/monitor/cpu.go
package monitor

import (
	"fmt"
	"log"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
)

// monitorCPU CPU监控核心逻辑
func (m *Manager) monitorCPU() {
	defer m.wg.Done()
	ticker := time.NewTicker(m.cpuCfg.Interval)
	defer ticker.Stop()

	log.Println("CPU监控协程已启动")
	for {
		select {
		case <-m.ctx.Done():
			log.Println("CPU监控协程退出")
			return
		case <-ticker.C:
			// 获取CPU使用率
			usageList, err := cpu.Percent(0, false)
			if err != nil {
				log.Printf("CPU监控失败: %v", err)
				continue
			}
			if len(usageList) == 0 {
				log.Println("CPU监控：未获取到使用率数据")
				continue
			}
			cpuUsage := usageList[0]

			log.Printf("CPU状态 | 使用率: %.2f%% | 阈值: %.2f%%", cpuUsage, m.cpuCfg.Threshold)

			// 触发告警
			if cpuUsage > m.cpuCfg.Threshold {
				content := fmt.Sprintf(
					"CPU使用率超标！\n当前使用率: %.2f%%\n告警阈值: %.2f%%",
					cpuUsage, m.cpuCfg.Threshold,
				)
				m.sendAlerts("CPU告警", content)
			}
		}
	}
}
