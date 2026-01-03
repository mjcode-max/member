package validator

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

var (
	// 手机号正则
	phoneRegex = regexp.MustCompile(`^1[3-9]\d{9}$`)
	// 邮箱正则
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	// 身份证号正则（18位）
	idCardRegex = regexp.MustCompile(`^[1-9]\d{5}(18|19|20)\d{2}(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01])\d{3}[\dXx]$`)
)

// Validator 验证器
type Validator struct {
	errors []string
}

// New 创建验证器
func New() *Validator {
	return &Validator{
		errors: make([]string, 0),
	}
}

// Validate 执行验证
func (v *Validator) Validate() error {
	if len(v.errors) == 0 {
		return nil
	}
	return errors.New(strings.Join(v.errors, "; "))
}

// Errors 获取所有错误
func (v *Validator) Errors() []string {
	return v.errors
}

// HasErrors 是否有错误
func (v *Validator) HasErrors() bool {
	return len(v.errors) > 0
}

// AddError 添加错误
func (v *Validator) AddError(field, message string) {
	v.errors = append(v.errors, fmt.Sprintf("%s: %s", field, message))
}

// Required 必填验证
func (v *Validator) Required(field, value string) *Validator {
	if strings.TrimSpace(value) == "" {
		v.AddError(field, "不能为空")
	}
	return v
}

// MinLength 最小长度验证
func (v *Validator) MinLength(field, value string, min int) *Validator {
	if len(value) < min {
		v.AddError(field, fmt.Sprintf("长度不能少于%d个字符", min))
	}
	return v
}

// MaxLength 最大长度验证
func (v *Validator) MaxLength(field, value string, max int) *Validator {
	if len(value) > max {
		v.AddError(field, fmt.Sprintf("长度不能超过%d个字符", max))
	}
	return v
}

// Length 长度范围验证
func (v *Validator) Length(field, value string, min, max int) *Validator {
	if len(value) < min || len(value) > max {
		v.AddError(field, fmt.Sprintf("长度必须在%d-%d个字符之间", min, max))
	}
	return v
}

// Email 邮箱验证
func (v *Validator) Email(field, value string) *Validator {
	if value != "" && !emailRegex.MatchString(value) {
		v.AddError(field, "邮箱格式不正确")
	}
	return v
}

// Phone 手机号验证
func (v *Validator) Phone(field, value string) *Validator {
	if value != "" && !phoneRegex.MatchString(value) {
		v.AddError(field, "手机号格式不正确")
	}
	return v
}

// IDCard 身份证号验证
func (v *Validator) IDCard(field, value string) *Validator {
	if value != "" && !idCardRegex.MatchString(value) {
		v.AddError(field, "身份证号格式不正确")
	}
	return v
}

// Numeric 数字验证
func (v *Validator) Numeric(field, value string) *Validator {
	if value != "" {
		for _, r := range value {
			if !unicode.IsDigit(r) {
				v.AddError(field, "必须是数字")
				break
			}
		}
	}
	return v
}

// Min 最小值验证（数字）
func (v *Validator) Min(field string, value, min int) *Validator {
	if value < min {
		v.AddError(field, fmt.Sprintf("不能小于%d", min))
	}
	return v
}

// Max 最大值验证（数字）
func (v *Validator) Max(field string, value, max int) *Validator {
	if value > max {
		v.AddError(field, fmt.Sprintf("不能大于%d", max))
	}
	return v
}

// Range 范围验证（数字）
func (v *Validator) Range(field string, value, min, max int) *Validator {
	if value < min || value > max {
		v.AddError(field, fmt.Sprintf("必须在%d-%d之间", min, max))
	}
	return v
}

// In 枚举值验证
func (v *Validator) In(field, value string, options []string) *Validator {
	if value != "" {
		found := false
		for _, opt := range options {
			if value == opt {
				found = true
				break
			}
		}
		if !found {
			v.AddError(field, fmt.Sprintf("必须是以下值之一: %s", strings.Join(options, ", ")))
		}
	}
	return v
}

// Equal 相等验证
func (v *Validator) Equal(field, value1, value2 string) *Validator {
	if value1 != value2 {
		v.AddError(field, "两次输入不一致")
	}
	return v
}

// Regex 正则验证
func (v *Validator) Regex(field, value, pattern string) *Validator {
	if value != "" {
		matched, err := regexp.MatchString(pattern, value)
		if err != nil || !matched {
			v.AddError(field, "格式不正确")
		}
	}
	return v
}

// 快捷验证函数

// ValidatePhone 验证手机号
func ValidatePhone(phone string) bool {
	return phoneRegex.MatchString(phone)
}

// ValidateEmail 验证邮箱
func ValidateEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// ValidateIDCard 验证身份证号
func ValidateIDCard(idCard string) bool {
	return idCardRegex.MatchString(idCard)
}

// ValidateRequired 验证必填
func ValidateRequired(value string) bool {
	return strings.TrimSpace(value) != ""
}

// ValidateLength 验证长度
func ValidateLength(value string, min, max int) bool {
	length := len(value)
	return length >= min && length <= max
}
