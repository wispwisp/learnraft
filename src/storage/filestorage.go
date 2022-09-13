package storage

import (
	log "github.com/wispwisp/learnraft/logger"
)

type FileStorage struct {
	filePath string
	logger   *log.FileLogger
}

func NewFileStorage(filePath string, logger *log.FileLogger) *FileStorage {
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
