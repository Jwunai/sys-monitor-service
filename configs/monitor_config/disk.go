// configs/monitor_config/disk.go
package monitor_config

import "time"

// DiskConfig 磁盘监控配置
type DiskConfig struct {
	Interval       time.Duration `yaml:"disk_interval"`        // 磁盘采样间隔（秒）
	UsageThreshold float64       `yaml:"disk_usage_threshold"` // 磁盘使用率阈值（%）
	MonitorDisks   []string      `yaml:"monitor_disks"`        // 监控磁盘分区（空数组自动监控所有）
}
