package errors

import (
	"fmt"
	"net/http"
)

// ErrorCode 错误码
type ErrorCode int

const (
	// ErrCodeUnknown 未知错误
	ErrCodeUnknown ErrorCode = 1000
	// ErrCodeInvalidParams 参数错误
	ErrCodeInvalidParams ErrorCode = 1001
	// ErrCodeNotFound 资源不存在
	ErrCodeNotFound ErrorCode = 1002
	// ErrCodeUnauthorized 未授权
	ErrCodeUnauthorized ErrorCode = 1003
	// ErrCodeForbidden 禁止访问
	ErrCodeForbidden ErrorCode = 1004
	// ErrCodeInternal 内部错误
	ErrCodeInternal ErrorCode = 1005
	// ErrCodeDatabase 数据库错误
	ErrCodeDatabase ErrorCode = 1006
	// ErrCodeRedis Redis错误
	ErrCodeRedis ErrorCode = 1007
)

// AppError 应用错误
type AppError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Err     error     `json:"-"`
}

// Error 实现 error 接口
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Unwrap 返回底层错误
func (e *AppError) Unwrap() error {
	return e.Err
}

// HTTPStatus 返回HTTP状态码
func (e *AppError) HTTPStatus() int {
	switch e.Code {
	case ErrCodeInvalidParams:
		return http.StatusBadRequest
	case ErrCodeNotFound:
		return http.StatusNotFound
	case ErrCodeUnauthorized:
		return http.StatusUnauthorized
	case ErrCodeForbidden:
		return http.StatusForbidden
	case ErrCodeInternal, ErrCodeDatabase, ErrCodeRedis:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// New 创建新错误
func New(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// Newf 创建格式化错误
func Newf(code ErrorCode, format string, args ...interface{}) *AppError {
	return &AppError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

// Wrap 包装错误
func Wrap(err error, code ErrorCode, message string) *AppError {
	if err == nil {
		return nil
	}
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Wrapf 包装格式化错误
func Wrapf(err error, code ErrorCode, format string, args ...interface{}) *AppError {
	if err == nil {
		return nil
	}
	return &AppError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
		Err:     err,
	}
}

// IsAppError 判断是否为应用错误
func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

// AsAppError 转换为应用错误
func AsAppError(err error) (*AppError, bool) {
	appErr, ok := err.(*AppError)
	return appErr, ok
}

// 常用错误快捷方法

// ErrInvalidParams 参数错误
func ErrInvalidParams(message string) *AppError {
	return New(ErrCodeInvalidParams, message)
}

// ErrNotFound 资源不存在
func ErrNotFound(message string) *AppError {
	return New(ErrCodeNotFound, message)
}

// ErrUnauthorized 未授权
func ErrUnauthorized(message string) *AppError {
	return New(ErrCodeUnauthorized, message)
}

// ErrForbidden 禁止访问
func ErrForbidden(message string) *AppError {
	return New(ErrCodeForbidden, message)
}

// ErrInternal 内部错误
func ErrInternal(message string) *AppError {
	return New(ErrCodeInternal, message)
}

// ErrDatabase 数据库错误
func ErrDatabase(err error) *AppError {
	return Wrap(err, ErrCodeDatabase, "数据库操作失败")
}

// ErrRedis Redis错误
func ErrRedis(err error) *AppError {
	return Wrap(err, ErrCodeRedis, "Redis操作失败")
}
