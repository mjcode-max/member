package utils

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"strings"
	"time"
	"unicode"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// IsEmpty 判断字符串是否为空
func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

// IsNotEmpty 判断字符串是否不为空
func IsNotEmpty(s string) bool {
	return !IsEmpty(s)
}

// Trim 去除首尾空白字符
func Trim(s string) string {
	return strings.TrimSpace(s)
}

// ToLower 转换为小写
func ToLower(s string) string {
	return strings.ToLower(s)
}

// ToUpper 转换为大写
func ToUpper(s string) string {
	return strings.ToUpper(s)
}

// Contains 判断是否包含子串
func Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// HasPrefix 判断是否有前缀
func HasPrefix(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

// HasSuffix 判断是否有后缀
func HasSuffix(s, suffix string) bool {
	return strings.HasSuffix(s, suffix)
}

// Split 分割字符串
func Split(s, sep string) []string {
	return strings.Split(s, sep)
}

// Join 连接字符串
func Join(elems []string, sep string) string {
	return strings.Join(elems, sep)
}

// MD5 计算MD5哈希
func MD5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// RandomString 生成随机字符串
func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// RandomNumber 生成随机数字字符串
func RandomNumber(length int) string {
	const charset = "0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// IsNumeric 判断是否为数字
func IsNumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// IsAlpha 判断是否为字母
func IsAlpha(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

// IsAlphaNumeric 判断是否为字母或数字
func IsAlphaNumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// MaskPhone 掩码手机号（显示前3位和后4位）
func MaskPhone(phone string) string {
	if len(phone) < 7 {
		return phone
	}
	return phone[:3] + "****" + phone[len(phone)-4:]
}

// MaskEmail 掩码邮箱（显示前2位和@后的内容）
func MaskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email
	}
	if len(parts[0]) <= 2 {
		return email
	}
	return parts[0][:2] + "***@" + parts[1]
}
