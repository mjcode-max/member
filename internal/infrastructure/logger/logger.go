package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

// Logger 日志接口
type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
	Sync() error
}

// ZapLogger zap日志实现
type ZapLogger struct {
	logger *zap.Logger
}

// NewZapLogger 创建zap日志实例
func NewZapLogger(level, format, output, filePath string) (*ZapLogger, error) {
	// 设置日志级别
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	// 设置编码器配置
	var encoderConfig zapcore.EncoderConfig
	if format == "json" {
		encoderConfig = zap.NewProductionEncoderConfig()
	} else {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	}
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// 设置编码器
	var encoder zapcore.Encoder
	if format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 设置输出
	var writeSyncer zapcore.WriteSyncer
	if output == "file" && filePath != "" {
		// 文件输出
		lumberjackLogger := &lumberjack.Logger{
			Filename:   filePath,
			MaxSize:    100, // MB
			MaxBackups: 7,
			MaxAge:     30, // days
			Compress:   true,
		}
		writeSyncer = zapcore.AddSync(lumberjackLogger)
	} else {
		// 标准输出
		writeSyncer = zapcore.AddSync(os.Stdout)
	}

	// 创建core
	core := zapcore.NewCore(encoder, writeSyncer, zapLevel)

	// 创建logger
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel))

	return &ZapLogger{logger: logger}, nil
}

// Debug 调试日志
func (l *ZapLogger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

// Info 信息日志
func (l *ZapLogger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

// Warn 警告日志
func (l *ZapLogger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

// Error 错误日志
func (l *ZapLogger) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

// Fatal 致命错误日志
func (l *ZapLogger) Fatal(msg string, fields ...zap.Field) {
	l.logger.Fatal(msg, fields...)
}

// Sync 同步日志
func (l *ZapLogger) Sync() error {
	return l.logger.Sync()
}
