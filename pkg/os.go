// pkg/os.go
package pkg

import (
	"runtime"
	"strings"
)

// GetOS 获取当前操作系统名称
func GetOS() string {
	switch runtime.GOOS {
	case "windows":
		return "Windows"
	case "linux":
		return "Linux"
	case "darwin":
		return "macOS"
	default:
		return runtime.GOOS // 其他系统返回原始值（如freebsd等）
	}
}

// IsWindows 判断是否为Windows系统
func IsWindows() bool {
	return runtime.GOOS == "windows"
}

// IsLinux 判断是否为Linux系统
func IsLinux() bool {
	return runtime.GOOS == "linux"
}

// GetDefaultDisks 自动获取系统默认监控磁盘分区
func GetDefaultDisks() []string {
	switch runtime.GOOS {
	case "windows":
		// Windows默认监控C盘（兼容中文/英文系统）
		return []string{"C:\\"}
	case "linux":
		// Linux默认监控根目录（/），可选加/data
		return []string{"/"}
	case "darwin":
		// macOS默认监控根目录
		return []string{"/"}
	default:
		// 其他系统返回空数组，由用户手动配置
		return []string{}
	}
}

// GetDiskDisplayName 格式化磁盘分区名称（告警展示用，可选）
// 比如把 "C:\\" 转成 "C盘"，"/" 转成 "根目录"
func GetDiskDisplayName(diskPath string) string {
	if IsWindows() {
		if strings.HasPrefix(diskPath, "C:\\") {
			return "C盘"
		}
		if strings.HasPrefix(diskPath, "D:\\") {
			return "D盘"
		}
		return diskPath // 其他分区返回原始值
	}
	if diskPath == "/" {
		return "根目录(/)"
	}
	return diskPath
}
