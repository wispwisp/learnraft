package logger

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"
)

type FileLogger struct {
	mtx      sync.Mutex
	fileName string

	fh     *os.File
	writer *bufio.Writer
}

func getTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
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
	fmt.Fprintln(l.writer, append([]interface{}{getTime()}, msg...)...)
	l.writer.Flush()
}

func (l *FileLogger) Info(msg ...interface{}) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	fmt.Fprintln(l.writer, append([]interface{}{getTime()}, msg...)...)
	l.writer.Flush()
}

func (l *FileLogger) Error(msg ...interface{}) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	fmt.Fprintln(l.writer, append([]interface{}{getTime()}, msg...)...) // fmt.Fprintln(l.writer, msg...)
	l.writer.Flush()
}

func (l *FileLogger) Debug(msg ...interface{}) {
	// l.mtx.Lock()
	// defer l.mtx.Unlock()
	// fmt.Fprintln(l.writer, append([]interface{}{getTime()}, msg...)...)
	// l.writer.Flush()
}
