package mylogger

import (
	"log"
	"sync"
)

type DefaultLogger struct {
	mtx sync.Mutex
}

func (l *DefaultLogger) Waring(msg ...interface{}) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	log.Println(msg...)
}

func (l *DefaultLogger) Info(msg ...interface{}) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	log.Println(msg...)
}

func (l *DefaultLogger) Error(msg ...interface{}) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	log.Println(msg...)
}
