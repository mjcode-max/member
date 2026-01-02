package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"member-pre/pkg/logger"
	"os"
)

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

// 实现pkg/logger.Logger接口

// Debug 调试日志
func (l *ZapLogger) Debug(msg string, fields ...logger.Field) {
	l.logger.Debug(msg, convertFields(fields)...)
}

// Info 信息日志
func (l *ZapLogger) Info(msg string, fields ...logger.Field) {
	l.logger.Info(msg, convertFields(fields)...)
}

// Warn 警告日志
func (l *ZapLogger) Warn(msg string, fields ...logger.Field) {
	l.logger.Warn(msg, convertFields(fields)...)
}

// Error 错误日志
func (l *ZapLogger) Error(msg string, fields ...logger.Field) {
	l.logger.Error(msg, convertFields(fields)...)
}

// Fatal 致命错误日志
func (l *ZapLogger) Fatal(msg string, fields ...logger.Field) {
	l.logger.Fatal(msg, convertFields(fields)...)
}

// convertFields 将pkg/logger.Field转换为zap.Field
func convertFields(fields []logger.Field) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for _, f := range fields {
		zapFields = append(zapFields, convertField(f))
	}
	return zapFields
}

// convertField 将单个pkg/logger.Field转换为zap.Field
func convertField(f logger.Field) zap.Field {
	key := f.Key()
	value := f.Value()

	// 根据值的类型选择合适的zap.Field
	switch v := value.(type) {
	case string:
		return zap.String(key, v)
	case int:
		return zap.Int(key, v)
	case int64:
		return zap.Int64(key, v)
	case uint:
		return zap.Uint(key, v)
	case uint64:
		return zap.Uint64(key, v)
	case float64:
		return zap.Float64(key, v)
	case bool:
		return zap.Bool(key, v)
	case error:
		return zap.Error(v)
	default:
		return zap.Any(key, v)
	}
}

// Sync 同步日志
func (l *ZapLogger) Sync() error {
	return l.logger.Sync()
}
