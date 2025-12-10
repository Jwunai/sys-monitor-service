// internal/monitor/disk.go
package monitor

import (
	"fmt"
	"log"
	"time"

	"github.com/Jwunai/sys-monitor-service/pkg"
	"github.com/shirou/gopsutil/v3/disk"
)

// monitorDisk 磁盘监控核心逻辑
func (m *Manager) monitorDisk() {
	defer m.wg.Done()
	ticker := time.NewTicker(m.diskCfg.Interval)
	defer ticker.Stop()

	log.Println("磁盘监控协程已启动")
	finalMonitorDisks := m.filterMonitorDisks()
	if len(finalMonitorDisks) == 0 {
		log.Printf("无有效磁盘分区可监控，磁盘监控协程退出")
		return
	}

	for {
		select {
		case <-m.ctx.Done():
			log.Println("磁盘监控协程退出")
			return
		case <-ticker.C:
			log.Println("开始磁盘监控 | 系统类型:", pkg.GetOS(), "| 监控分区:", finalMonitorDisks)

			for _, path := range finalMonitorDisks {
				diskUsage, err := disk.Usage(path)
				if err != nil {
					log.Printf("分区[%s]监控失败: %v", path, err)
					continue
				}

				totalGB := float64(diskUsage.Total) / 1024 / 1024 / 1024
				usedGB := float64(diskUsage.Used) / 1024 / 1024 / 1024
				freeGB := float64(diskUsage.Free) / 1024 / 1024 / 1024
				usedPercent := diskUsage.UsedPercent

				log.Printf(
					"磁盘状态 | 分区: %s | 总空间: %.2fGB | 已用: %.2fGB | 剩余: %.2fGB | 使用率: %.2f%% | 阈值: %.2f%%",
					path, totalGB, usedGB, freeGB, usedPercent, m.diskCfg.UsageThreshold,
				)

				// 触发告警
				if usedPercent > m.diskCfg.UsageThreshold {
					content := fmt.Sprintf(
						"分区[%s]使用率超标！\n总空间: %.2fGB\n已用: %.2fGB\n剩余: %.2fGB\n当前使用率: %.2f%%\n告警阈值: %.2f%%",
						path, totalGB, usedGB, freeGB, usedPercent, m.diskCfg.UsageThreshold,
					)
					m.sendAlerts("磁盘告警", content)
				}
			}
		}
	}
}
