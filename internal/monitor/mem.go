// internal/monitor/mem.go
package monitor

import (
	"fmt"
	"log"
	"time"

	"github.com/shirou/gopsutil/v3/mem"
)

// monitorMemory 内存监控核心逻辑
func (m *Manager) monitorMemory() {
	defer m.wg.Done()
	ticker := time.NewTicker(m.memCfg.Interval)
	defer ticker.Stop()

	log.Println("内存监控协程已启动")
	for {
		select {
		case <-m.ctx.Done():
			log.Println("内存监控协程退出")
			return
		case <-ticker.C:
			// 获取内存信息
			memInfo, err := mem.VirtualMemory()
			if err != nil {
				log.Printf("内存监控失败: %v", err)
				continue
			}

			totalGB := float64(memInfo.Total) / 1024 / 1024 / 1024
			availableGB := float64(memInfo.Available) / 1024 / 1024 / 1024
			usedPercent := memInfo.UsedPercent

			log.Printf(
				"内存状态 | 总内存: %.2fGB | 可用内存: %.2fGB | 使用率: %.2f%% | 可用阈值: %.2fGB",
				totalGB, availableGB, usedPercent, m.memCfg.AvailableThreshold,
			)

			// 触发告警
			if availableGB < m.memCfg.AvailableThreshold {
				content := fmt.Sprintf(
					"可用内存不足！\n总内存: %.2fGB\n当前可用: %.2fGB\n内存使用率: %.2f%%\n告警阈值: %.2fGB",
					totalGB, availableGB, usedPercent, m.memCfg.AvailableThreshold,
				)
				m.sendAlerts("内存告警", content)
			}
		}
	}
}
