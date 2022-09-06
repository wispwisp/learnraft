package logger

import (
	"log"
	"sync"
)

type Logger struct {
	mtx sync.Mutex
}

var logger Logger

func (l *Logger) Waring(msg ...interface{}) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	log.Println(msg...)
}

func Waring(msg ...interface{}) {
	logger.Waring(msg...)
}

// ***

func (l *Logger) Info(msg ...interface{}) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	log.Println(msg...)
}

func Info(msg ...interface{}) {
	logger.Info(msg...)
}

// ***

func (l *Logger) Error(msg ...interface{}) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	log.Println(msg...)
}

func Error(msg ...interface{}) {
	logger.Error(msg...)
}
