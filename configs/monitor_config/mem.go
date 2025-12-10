// configs/monitor_config/mem.go
package monitor_config

import "time"

// MemConfig 内存监控专属配置
type MemConfig struct {
	Interval           time.Duration `yaml:"mem_interval"`            // 内存采样间隔（秒）
	AvailableThreshold float64       `yaml:"mem_available_threshold"` // 可用内存告警阈值（GB）
}
