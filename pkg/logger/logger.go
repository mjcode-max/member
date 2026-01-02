package logger

// Logger 日志接口，定义在pkg层，供domain层使用
// domain层只依赖此接口，不依赖具体实现
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	Sync() error
}

// Field 日志字段接口，用于传递键值对
type Field interface {
	Key() string
	Value() interface{}
}

// NewField 创建日志字段
func NewField(key string, value interface{}) Field {
	return &field{
		key:   key,
		value: value,
	}
}

type field struct {
	key   string
	value interface{}
}

func (f *field) Key() string {
	return f.key
}

func (f *field) Value() interface{} {
	return f.value
}
