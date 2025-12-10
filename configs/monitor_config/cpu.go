// configs/monitor_config/cpu.go
package monitor_config

import "time"

// CPUConfig CPU监控配置
type CPUConfig struct {
	Interval  time.Duration `yaml:"cpu_interval"`  // CPU采样间隔（秒）
	Threshold float64       `yaml:"cpu_threshold"` // CPU告警阈值（%）
}
