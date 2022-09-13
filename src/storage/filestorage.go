package storage

import (
	"github.com/wispwisp/learnraft/mylogger"
)

type FileStorage struct {
	filePath string
	logger   mylogger.Logger
}

func NewFileStorage(filePath string, logger mylogger.Logger) *FileStorage {
	return &FileStorage{
		filePath: filePath,
		logger:   logger,
	}
}

func (s *FileStorage) Add(key, value string) bool {
	return false
}

func (s *FileStorage) Get(key string) (string, bool) {
	return "nil", false
}
