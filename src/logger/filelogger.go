package logger

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

type FileLogger struct {
	mtx      sync.Mutex
	fileName string

	fh     *os.File
	writer *bufio.Writer
}

func NewFileLogger(fileName string) (logger *FileLogger, err error) {
	f, err := os.Create(fileName)
	if err != nil {
		return
	}

	logger = &FileLogger{fileName: fileName}
	logger.fh = f
	logger.writer = bufio.NewWriter(f)

	return
}

func (l *FileLogger) Close() {
	l.fh.Close()
	l.writer.Flush()
}

func (l *FileLogger) Waring(msg ...interface{}) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	fmt.Fprintln(l.writer, msg...)
	l.writer.Flush()
}

func (l *FileLogger) Info(msg ...interface{}) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	fmt.Fprintln(l.writer, msg...)
	l.writer.Flush()
}

func (l *FileLogger) Error(msg ...interface{}) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	fmt.Fprintln(l.writer, msg...)
	l.writer.Flush()
}
