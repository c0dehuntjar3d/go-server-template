package logger

type Interface interface {
	Debug(message string, args ...Field)
	Info(message string, args ...Field)
	Warn(message string, args ...Field)
	Error(message string, args ...Field)
	Fatal(message string, args ...Field)
}

var (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	Gray    = "\033[37m"
	White   = "\033[97m"
)

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

func withColor(color string, message string) string {
	return color + message + Reset
}
