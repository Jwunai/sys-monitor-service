// configs/loader.go
package configs

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Jwunai/sys-monitor-service/configs/alert_config"  // 替换为你的实际module名
	"github.com/Jwunai/sys-monitor-service/configs/monitor_config" // 替换为你的实际module名
	"github.com/Jwunai/sys-monitor-service/pkg"             // 工具包（系统信息/脱敏等）
	"gopkg.in/yaml.v3"
)

// MonitorConfig 监控总配置（整合server_name + 各资源专属配置）
type MonitorConfig struct {
	ServerName string             `yaml:"server_name"` // 服务器名称（告警标题标识）
	CPU        monitor_config.CPUConfig  `yaml:",inline"`     // 内嵌CPU配置（匹配cpu_interval/cpu_threshold）
	Disk       monitor_config.DiskConfig `yaml:",inline"`     // 内嵌磁盘配置（匹配disk_interval等）
	Mem        monitor_config.MemConfig  `yaml:",inline"`     // 内嵌内存配置（匹配mem_interval等）
}

// AlertConfig 告警总配置（无变化，匹配alert嵌套层级）
type AlertConfig struct {
	DingTalk alert_config.DingTalkConfig `yaml:"dingtalk"` // 匹配alert.dingtalk
	Email    alert_config.EmailConfig    `yaml:"email"`    // 匹配alert.email
	//SMS      alerts.SMSConfig      `yaml:"sms"`      // 匹配alert.sms
}

// AppConfig 全局配置（匹配YAML根层级）
type AppConfig struct {
	Monitor MonitorConfig `yaml:"monitor"` // 匹配monitor根层级
	Alert   AlertConfig   `yaml:"alert"`   // 匹配alert根层级
}

// LoadConfig 加载并解析配置文件（核心逻辑不变，仅调整默认值）
func LoadConfig(configPath string) (*AppConfig, error) {
	log.Println("========== 开始加载配置文件 ==========")

	// 1. 处理配置文件路径
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, fmt.Errorf("解析配置文件路径失败：%w", err)
	}

	// 2. 读取配置文件
	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败：%w", err)
	}

	// 3. 解析YAML到结构体
	var cfg AppConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析YAML配置失败：%w", err)
	}

	// 4. 填充默认值（适配新的独立间隔）
	setDefaultConfig(&cfg)

	// 5. 校验配置合法性
	if err := validateConfig(&cfg); err != nil {
		log.Printf("【警告】配置校验不通过：%v", err)
	}

	log.Println("========== 配置文件加载完成 ==========")
	return &cfg, nil
}

// setDefaultConfig 填充配置默认值（适配独立的采样间隔）
func setDefaultConfig(cfg *AppConfig) {
	// 服务器名称默认值
	if cfg.Monitor.ServerName == "" {
		cfg.Monitor.ServerName = fmt.Sprintf("sys-monitor-%s", pkg.GetOS())
	}

	// CPU配置默认值
	if cfg.Monitor.CPU.Interval == 0 {
		cfg.Monitor.CPU.Interval = 30 * time.Second
	}
	if cfg.Monitor.CPU.Threshold == 0 {
		cfg.Monitor.CPU.Threshold = 80.0
	}

	// 内存配置默认值
	if cfg.Monitor.Mem.Interval == 0 {
		cfg.Monitor.Mem.Interval = 30 * time.Second
	}
	if cfg.Monitor.Mem.AvailableThreshold == 0 {
		cfg.Monitor.Mem.AvailableThreshold = 2.0
	}

	// 磁盘配置默认值
	if cfg.Monitor.Disk.Interval == 0 {
		cfg.Monitor.Disk.Interval = 60 * time.Second
	}
	if cfg.Monitor.Disk.UsageThreshold == 0 {
		cfg.Monitor.Disk.UsageThreshold = 85.0
	}
	if len(cfg.Monitor.Disk.MonitorDisks) == 0 {
		cfg.Monitor.Disk.MonitorDisks = pkg.GetDefaultDisks() // 自动识别系统磁盘
	}
}

// validateConfig 校验配置合法性（无变化）
func validateConfig(cfg *AppConfig) error {
	var errMsg []string

	// CPU配置校验
	if cfg.Monitor.CPU.Threshold < 0 || cfg.Monitor.CPU.Threshold > 100 {
		errMsg = append(errMsg, "CPU告警阈值必须在0-100之间")
	}
	if cfg.Monitor.CPU.Interval < 5*time.Second { // 最小间隔5秒，避免高频采样
		errMsg = append(errMsg, "CPU采样间隔不能小于5秒")
	}

	// 内存配置校验
	if cfg.Monitor.Mem.AvailableThreshold < 0 {
		errMsg = append(errMsg, "可用内存阈值不能为负数")
	}
	if cfg.Monitor.Mem.Interval < 5*time.Second {
		errMsg = append(errMsg, "内存采样间隔不能小于5秒")
	}

	// 磁盘配置校验
	if cfg.Monitor.Disk.UsageThreshold < 0 || cfg.Monitor.Disk.UsageThreshold > 100 {
		errMsg = append(errMsg, "磁盘使用率阈值必须在0-100之间")
	}
	if cfg.Monitor.Disk.Interval < 5*time.Second {
		errMsg = append(errMsg, "磁盘采样间隔不能小于5秒")
	}

	// 告警配置校验
	if cfg.Alert.DingTalk.Token != "" && cfg.Alert.DingTalk.Secret == "" {
		errMsg = append(errMsg, "配置了钉钉Token但未配置Secret")
	}
	if cfg.Alert.Email.From != "" {
		if cfg.Alert.Email.SmtpHost == "" {
			errMsg = append(errMsg, "配置了发件人邮箱但未配置SMTP服务器")
		}
		if cfg.Alert.Email.SmtpPort == 0 {
			errMsg = append(errMsg, "配置了发件人邮箱但未配置SMTP端口")
		}
	}

	if len(errMsg) > 0 {
		return fmt.Errorf("%v", errMsg)
	}
	return nil
}
