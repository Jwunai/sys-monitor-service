// pkg/desensitize.go
package pkg

import (
	"strings"
)

// Desensitize 通用脱敏函数
// 规则：保留前6位 + 后4位，中间用*填充（长度不足10位时，保留前3后2，中间*）
// 示例：
// - "ae8042c730b88ce55acced5dde21e9bd5fde419c7c9b671fff48eaf94c67ac2e" → "ae8042**************************7ac2e"
// - "1234567890" → "123****890"
// - "12345" → "12***45"
// - 空字符串/长度<2 → 返回原字符串
func Desensitize(s string) string {
	if s == "" {
		return ""
	}
	length := len(s)
	// 长度不足2，无需脱敏
	if length < 2 {
		return s
	}
	// 长度10及以上：前6后4
	if length >= 10 {
		prefix := s[:6]
		suffix := s[length-4:]
		star := strings.Repeat("*", length-10)
		return prefix + star + suffix
	}
	// 长度2-9：前3后2（不足3则前1后1）
	prefixLen := 3
	if length < 5 {
		prefixLen = 1
	}
	prefix := s[:prefixLen]
	suffix := s[length-2:]
	star := strings.Repeat("*", length-prefixLen-2)
	return prefix + star + suffix
}

// DesensitizeEmail 邮箱脱敏
// 示例："test123@qq.com" → "te****@qq.com"
func DesensitizeEmail(email string) string {
	if email == "" {
		return ""
	}
	atIdx := strings.Index(email, "@")
	if atIdx <= 2 {
		// 前缀不足2位，保留原前缀
		return email
	}
	prefix := email[:2]
	domain := email[atIdx:]
	star := strings.Repeat("*", atIdx-2)
	return prefix + star + domain
}

// DesensitizePhone 手机号脱敏
// 示例："13812345678" → "138****5678"
func DesensitizePhone(phone string) string {
	if len(phone) != 11 {
		// 非11位手机号，用通用脱敏
		return Desensitize(phone)
	}
	return phone[:3] + "****" + phone[7:]
}

// DesensitizeSMTP 邮箱SMTP密码脱敏
func DesensitizeSMTP(pwd string) string {
	return Desensitize(pwd)
}

// DesensitizeSMSKey 短信密钥脱敏
func DesensitizeSMSKey(key string) string {
	return Desensitize(key)
}
