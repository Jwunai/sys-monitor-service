// internal/monitor/disk_utils.go
package monitor

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/Jwunai/sys-monitor-service/pkg"
	"github.com/shirou/gopsutil/v3/disk"
)

// getAllValidDisks 获取系统所有有效磁盘分区
func (m *Manager) getAllValidDisks() []string {
	var validDisks []string

	if pkg.IsWindows() {
		disks, err := disk.Partitions(true)
		if err != nil {
			log.Printf("Windows获取分区列表失败: %v", err)
			return validDisks
		}
		for _, d := range disks {
			if d.Fstype == "NTFS" || d.Fstype == "FAT32" || d.Fstype == "exFAT" {
				if len(d.Mountpoint) == 2 && d.Mountpoint[1] == ':' && d.Mountpoint[0] != 'A' && d.Mountpoint[0] != 'B' {
					normPath := filepath.Clean(d.Mountpoint + "\\")
					validDisks = append(validDisks, normPath)
				}
			}
		}
	} else if pkg.IsLinux() {
		disks, err := disk.Partitions(true)
		if err != nil {
			log.Printf("Linux获取分区列表失败: %v", err)
			return validDisks
		}
		excludeFsTypes := map[string]bool{
			"sysfs":       true,
			"proc":        true,
			"tmpfs":       true,
			"devtmpfs":    true,
			"devpts":      true,
			"cgroup":      true,
			"overlay":     true,
			"aufs":        true,
			"squashfs":    true,
			"rpc_pipefs":  true,
			"binfmt_misc": true,
		}
		for _, d := range disks {
			if !excludeFsTypes[d.Fstype] {
				if !strings.Contains(d.Mountpoint, "/tmp") && !strings.Contains(d.Fstype, "nfs") && !strings.Contains(d.Fstype, "smb") {
					validDisks = append(validDisks, d.Mountpoint)
				}
			}
		}
	}

	// 去重
	uniqueDisks := make([]string, 0, len(validDisks))
	seen := make(map[string]bool)
	for _, d := range validDisks {
		if !seen[d] {
			seen[d] = true
			uniqueDisks = append(uniqueDisks, d)
		}
	}

	return uniqueDisks
}

// filterMonitorDisks 过滤需要监控的磁盘分区（适配配置）
func (m *Manager) filterMonitorDisks() []string {
	allValidDisks := m.getAllValidDisks()
	if len(allValidDisks) == 0 {
		log.Printf("未检测到系统有效磁盘分区，跳过磁盘监控")
		return []string{}
	}

	configDisks := m.preprocessDiskPaths()
	if len(configDisks) == 0 {
		log.Printf("monitor_disks配置为空，默认监控所有有效分区: %v", allValidDisks)
		return allValidDisks
	}

	var finalDisks []string
	var invalidConfigDisks []string
	validDiskMap := make(map[string]bool)
	for _, d := range allValidDisks {
		validDiskMap[d] = true
	}

	for _, cfgDisk := range configDisks {
		if validDiskMap[cfgDisk] {
			finalDisks = append(finalDisks, cfgDisk)
		} else {
			invalidConfigDisks = append(invalidConfigDisks, cfgDisk)
		}
	}

	if len(finalDisks) == 0 {
		log.Printf("配置的分区[%v]均不存在，降级监控所有有效分区: %v", invalidConfigDisks, allValidDisks)
		return allValidDisks
	}

	if len(invalidConfigDisks) > 0 {
		log.Printf("配置的分区[%v]不存在，仅监控有效分区: %v", invalidConfigDisks, finalDisks)
	}
	return finalDisks
}

// preprocessDiskPaths 预处理磁盘路径（格式化）
func (m *Manager) preprocessDiskPaths() []string {
	var processedPaths []string
	for _, path := range m.diskCfg.MonitorDisks {
		if path == "" {
			continue
		}
		normPath := filepath.Clean(path)

		if pkg.IsWindows() {
			if len(normPath) == 2 && normPath[1] == ':' {
				normPath += "\\"
			}
		}

		processedPaths = append(processedPaths, normPath)
	}
	return processedPaths
}
