// internal/alert_services/dingtalk.go
package alert_services

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/Jwunai/sys-monitor-service/configs/alert_config"
	"github.com/Jwunai/sys-monitor-service/internal/interfaces"
)

// DingTalk 钉钉告警实现
type DingTalk struct {
	cfg *alert_config.DingTalkConfig // 钉钉配置
}

// NewDingTalk 创建钉钉告警实例
func NewDingTalk(cfg *alert_config.DingTalkConfig) interfaces.AlertSender {
	return &DingTalk{cfg: cfg}
}

// Name 返回告警渠道名称
func (d *DingTalk) Name() string {
	return "钉钉"
}

// IsEnabled 判断是否启用（token和secret非空）
func (d *DingTalk) IsEnabled() bool {
	return d.cfg != nil && d.cfg.Token != "" && d.cfg.Secret != ""
}

// SendAlert 发送钉钉告警（支持签名验证）
func (d *DingTalk) SendAlert(title, serverName, content string) error {
	if !d.IsEnabled() {
		return fmt.Errorf("钉钉告警未启用配置缺失")
	}

	// 1. 构造钉钉webhook地址
	webhookURL := fmt.Sprintf("https://oapi.dingtalk.com/robot/send?access_token=%s", d.cfg.Token)

	// 2. 计算签名（如果配置了secret）
	if d.cfg.Secret != "" {
		timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
		stringToSign := fmt.Sprintf("%s\n%s", timestamp, d.cfg.Secret)

		// HMAC-SHA256加密
		h := hmac.New(sha256.New, []byte(d.cfg.Secret))
		h.Write([]byte(stringToSign))
		sign := url.QueryEscape(base64.StdEncoding.EncodeToString(h.Sum(nil)))

		// 拼接签名到URL
		webhookURL = fmt.Sprintf("%s&timestamp=%s&sign=%s", webhookURL, timestamp, sign)
	}

	// 3. 构造告警消息体
	msg := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": fmt.Sprintf("[%s] %s\n%s", serverName, title, content),
		},
	}
	jsonData, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("构造钉钉消息失败: %w", err)
	}

	// 4. 发送HTTP请求
	resp, err := http.Post(webhookURL, "application/json", io.NopCloser(bytes.NewBuffer(jsonData)))
	if err != nil {
		return fmt.Errorf("发送钉钉告警请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 5. 解析响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取钉钉响应失败: %w", err)
	}

	// 钉钉响应格式：{"errcode":0,"errmsg":"ok"}
	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return fmt.Errorf("解析钉钉响应失败: %w, 响应内容: %s", err, string(respBody))
	}

	if errcode, ok := result["errcode"].(float64); ok && errcode != 0 {
		return fmt.Errorf("钉钉告警发送失败: %s (errcode: %.0f)", result["errmsg"], errcode)
	}

	return nil
}

var _ interfaces.AlertSender = (*DingTalk)(nil)
