package internal

import (
	"fmt"
	"log/slog"
	"url-shortener/internal/storage"
)

type CustomLogWriter struct {
	storage        storage.Storage
	collectionName string
}

func (clw *CustomLogWriter) Write(p []byte) (n int, err error) {
	clw.storage.Log(string(p), slog.LevelInfo)
	return fmt.Print(string(p))
}

func NewCustomLogWriter(storage storage.Storage, collectionName string) *CustomLogWriter {
	return &CustomLogWriter{
		storage: storage,
	}
}
