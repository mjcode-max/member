package test

import "member-pre/pkg/logger"

// MockLogger 模拟日志记录器，用于单元测试
// 所有日志方法都是空实现，不会输出任何内容
type MockLogger struct{}

// NewMockLogger 创建新的模拟日志记录器
func NewMockLogger() *MockLogger {
	return &MockLogger{}
}

// Debug 实现 logger.Logger 接口
func (m *MockLogger) Debug(msg string, fields ...logger.Field) {}

// Info 实现 logger.Logger 接口
func (m *MockLogger) Info(msg string, fields ...logger.Field) {}

// Warn 实现 logger.Logger 接口
func (m *MockLogger) Warn(msg string, fields ...logger.Field) {}

// Error 实现 logger.Logger 接口
func (m *MockLogger) Error(msg string, fields ...logger.Field) {}

// Fatal 实现 logger.Logger 接口
func (m *MockLogger) Fatal(msg string, fields ...logger.Field) {}

// Sync 实现 logger.Logger 接口
func (m *MockLogger) Sync() error {
	return nil
}
