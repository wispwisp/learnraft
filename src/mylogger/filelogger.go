package mylogger

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
)

type FileLogger struct {
	mtx      sync.Mutex
	fileName string
	level    Level
	fh       *os.File
	writer   *bufio.Writer
}

func getTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func NewFileLogger(level Level, fileName string) (logger *FileLogger, err error) {
	f, err := os.Create(fileName)
	if err != nil {
		return
	}

	if level < ERROR || level > DEBUG {
		err = errors.New("wrong level")
		return
	}

	logger = &FileLogger{
		fileName: fileName,
		level:    level,
		fh:       f,
		writer:   bufio.NewWriter(f),
	}

	return
}

func (l *FileLogger) Close() {
	l.fh.Close()
	l.writer.Flush()
}

func (l *FileLogger) Error(msg ...interface{}) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	if l.level == ERROR {
		// fmt.Fprintln(l.writer, msg...)
		fmt.Fprintln(l.writer, append([]interface{}{getTime(), "[ERROR]"}, msg...)...)
		l.writer.Flush()
	}
}

func (l *FileLogger) Waring(msg ...interface{}) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	if l.level >= WARNING {
		fmt.Fprintln(l.writer, append([]interface{}{getTime(), "[WARNING]"}, msg...)...)
		l.writer.Flush()
	}
}

func (l *FileLogger) Info(msg ...interface{}) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	if l.level >= INFO {
		fmt.Fprintln(l.writer, append([]interface{}{getTime(), "[INFO]"}, msg...)...)
		l.writer.Flush()
	}
}

func (l *FileLogger) Debug(msg ...interface{}) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	if l.level == DEBUG {
		fmt.Fprintln(l.writer, append([]interface{}{getTime(), "[DEBUG]"}, msg...)...)
		l.writer.Flush()
	}
}
