# sys-monitor-service：系统资源监控与多渠道告警服务

一个基于 Go 开发的轻量级、高可扩展系统监控服务，支持 CPU / 内存 / 磁盘资源监控，集成钉钉 / 邮箱多渠道告警 , 自用工具。

## 核心特性

1. **模块化监控**：支持 CPU 使用率、内存可用量、磁盘使用率监控，各资源采样间隔独立配置
2. **可扩展告警**：基于 `AlertSender` 接口设计，已实现钉钉 / 邮箱告警，后续新增渠道（如短信 / 企业微信）无需改动核心代码
3. **自动注册机制**：告警渠道通过自注册模式加载，扩展时仅需新增实现类 + 注册代码
4. **配置驱动**：所有监控规则、告警开关通过 YAML 配置文件管理
5. **健壮性设计**：支持优雅退出，避免协程泄漏、告警失败重试、配置默认值填充与合法性校验
6. **多平台兼容**：基于 `gopsutil` 实现，支持 Windows/Linux/macOS 系统

## 快速开始

### 1. 环境准备

- Go 1.21+（推荐 1.21 及以上版本）
- 网络权限：监控服务需本地系统权限，告警渠道需外网访问权限（钉钉 / 邮箱 SMTP 服务器）

### 2. 项目结构（核心架构）

```plaintext
sys-monitor-service/
├── cmd/                  # 可执行程序入口
│   └── sys-monitor/
│       └── main.go       # 服务启动入口（优雅退出+监控启动）
├── configs/              # 配置相关（结构体+加载逻辑）
│   ├── alert_config/     # 告警渠道配置结构体（钉钉/邮箱）
│   ├── monitor_config/   # 监控资源配置结构体（CPU/内存/磁盘）
│   └── loader.go         # 配置加载+默认值+校验逻辑
├── internal/
│   ├── alert_services/   # 告警渠道实现（自动注册）
│   │   ├── dingtalk.go   # 钉钉告警实现
│   │   ├── email.go      # 邮箱告警实现
│   │   └── register.go   # 告警自动注册逻辑
│   ├── interfaces/       # 核心接口定义
│   │   └── alert.go      # AlertSender 告警接口
│   ├── monitor/          # 监控核心逻辑
│   │   └── manager.go    # 监控管理器（启动/停止/告警发送）
│   └── registry/         # 告警注册器（统一创建启用的告警实例）
│       └── alert.go
├── sys-monitor           # 编译的二进制可执行文件（linux）
├── sys-monitor-win.exe   # 编译的二进制可执行文件（window）
├── config.yml            # 配置文件
├── config.yml            # 配置文件
├── go.mod                # 依赖管理
└── README.md             # 项目文档
```

#### 3.编译运行

```bash
#切换平台
 $env:GOOS = "linux"
 $env:GOOS = "windows"

# 编译（压缩二进制文件）
go build -ldflags "-s -w" -trimpath -o sys-monitor ./cmd/sys-monitor/

# 运行（Windows）
./sys-monitor.exe

# 运行（Linux/macOS）
chmod +x sys-monitor
./sys-monitor
```



##  配置说明

### 1. 监控配置（monitor 节点）

| 字段名                  | 类型     | 说明                                   | 默认值                 |
| ----------------------- | -------- | -------------------------------------- | ---------------------- |
| server_name             | string   | 服务器名称（用于告警标题区分多服务器） | sys-monitor-[系统类型] |
| cpu_interval            | duration | CPU 采样间隔（支持 s/m/h，如 30s、5m） | 30s                    |
| cpu_threshold           | float64  | CPU 使用率告警阈值（0-100）            | 80.0                   |
| mem_interval            | duration | 内存采样间隔                           | 30s                    |
| mem_available_threshold | float64  | 可用内存告警阈值（单位：GB）           | 2.0                    |
| disk_interval           | duration | 磁盘采样间隔                           | 60s                    |
| disk_usage_threshold    | float64  | 磁盘使用率告警阈值（0-100）            | 85.0                   |
| monitor_disks           | []string | 需监控的磁盘分区（如 ["/", "/data"]）  | 自动识别系统磁盘       |

### 2. 告警配置（alert 节点）

#### （1）钉钉告警（dingtalk 节点）

| 字段名 | 类型   | 说明                                 |
| ------ | ------ | ------------------------------------ |
| token  | string | 钉钉机器人 Token（创建机器人时获取） |
| secret | string | 签名密钥（可选，启用加签模式需填写） |

> 提示：钉钉机器人需开启「自定义关键词」（如 “告警”），否则告警消息会被拦截

#### （2）邮箱告警（email 节点）

| 字段名    | 类型     | 说明                                                      |
| --------- | -------- | --------------------------------------------------------- |
| from      | string   | 发件人邮箱（如 alert@qq.com）                             |
| password  | string   | 邮箱密码或授权码（QQ 邮箱需用授权码）                     |
| smtp_host | string   | SMTP 服务器地址（如 [smtp.qq.com](https://smtp.qq.com/)） |
| smtp_port | int      | SMTP 端口（SSL 通常 465，非 SSL 通常 25）                 |
| to        | []string | 收件人邮箱列表（支持多个）                                |



## 部署方式

### 1. 本地测试运行

直接执行编译后的二进制文件，适合测试环境：

```bash
./sys-monitor  # Linux/macOS
./sys-monitor.exe  # Windows
```

### 2. Linux 后台运行（systemd）

#### （1）创建服务文件

```bash
sudo vim /etc/systemd/system/sys-monitor.service
```

#### （2）写入以下内容（修改路径为实际项目路径）

```ini
[Unit]
Description=System Monitor Service
After=network.target

[Service]
User=root
WorkingDirectory=/opt/system/  #运行目录
ExecStart=/opt/system/sys-monitor #运行文件
Restart=always  # 进程退出时自动重启
RestartSec=5    # 重启间隔5秒

[Install]
WantedBy=multi-user.target

```

#### （3）启动服务并设置开机自启

```bash
# 重新加载systemd配置
sudo systemctl daemon-reload

# 启动服务
sudo systemctl start sys-monitor

# 查看服务状态
sudo systemctl status sys-monitor

# 设置开机自启
sudo systemctl enable sys-monitor
```

### 3. Windows 服务部署

使用 `nssm` 工具将二进制文件注册为 Windows 服务：

```bash
# 安装nssm后执行
nssm install SysMonitorService "C:\sys-monitor-service\sys-monitor.exe"
nssm start SysMonitorService
```

