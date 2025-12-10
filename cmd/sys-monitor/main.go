package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Jwunai/sys-monitor-service/configs"
	"github.com/Jwunai/sys-monitor-service/internal/monitor"
	"github.com/Jwunai/sys-monitor-service/internal/registry"
)

func main() {
	// ========== 1. 加载配置 ==========
	cfg, err := configs.LoadConfig("./config.yml")
	if err != nil {
		log.Fatalf("配置加载失败: %v", err)
	}

	// ========== 2. 创建告警实例 ==========
	alertSenders := registry.CreateAllEnabled(&cfg.Alert)
	if len(alertSenders) == 0 {
		log.Printf("未配置任何启用的告警渠道（钉钉/邮箱），告警功能将禁用")
	} else {
		// 打印启用的告警渠道
		enabledAlerts := make([]string, 0, len(alertSenders))
		for _, sender := range alertSenders {
			enabledAlerts = append(enabledAlerts, sender.Name())
		}
		log.Printf("已启用告警渠道：%v", enabledAlerts)
	}

	// ========== 3. 初始化监控管理器 ==========
	monitorMgr := monitor.NewManager(
		cfg.Monitor.ServerName,
		cfg.Monitor.CPU,
		cfg.Monitor.Mem,
		cfg.Monitor.Disk,
		alertSenders,
	)

	// ========== 4. 启动监控避免阻塞主线程 ==========
	log.Println("========== 启动监控服务 ==========")
	go func() {
		monitorMgr.Start()
	}()
	log.Println("监控服务启动成功")

	// ========== 5. 退出逻辑 ==========
	quit := make(chan os.Signal, 1)
	// 监听Ctrl+C、kill等退出信号
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-quit // 阻塞等待退出信号

	// ========== 6. 停止监控 ==========
	log.Println("接收到退出信号，正在停止监控服务...")
	monitorMgr.Stop()
	log.Println("监控服务已正常退出")
}
