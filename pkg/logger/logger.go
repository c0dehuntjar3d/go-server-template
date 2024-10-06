package logger

type Interface interface {
	Debug(message string, args ...Field)
	Info(message string, args ...Field)
	Warn(message string, args ...Field)
	Error(message string, args ...Field)
	Fatal(message string, args ...Field)
}

type Field struct {
	Key   string
	Value interface{}
}

func NewField(key string, value interface{}) Field {
	return Field{
		Key:   key,
		Value: value,
	}
}
