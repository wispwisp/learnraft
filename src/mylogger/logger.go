package mylogger

type Level int

const (
	ERROR   Level = iota
	WARNING Level = iota
	INFO    Level = iota
	DEBUG   Level = iota
)

type Logger interface {
	Error(msg ...interface{})
	Waring(msg ...interface{})
	Info(msg ...interface{})
	Debug(msg ...interface{})
}
