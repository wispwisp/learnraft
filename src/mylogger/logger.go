package mylogger

type Logger interface {
	Waring(msg ...interface{})
	Info(msg ...interface{})
	Error(msg ...interface{})
	Debug(msg ...interface{})
}
