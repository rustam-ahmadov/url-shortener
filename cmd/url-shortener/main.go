package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/storage/mongo_storage"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log, err := setupLogger(cfg.Env)
	if err != nil {
		fmt.Printf("err occured while initializing logger: %s", err.Error())
		os.Exit(1)
	}

	storage, err := mongo_storage.NewStorage(cfg.Storage_Path, cfg.Storage_Name)
	if err != nil {
		fmt.Printf("err occured while initializing storage: %s", err.Error())
		os.Exit(1)
	}
	log.Info("url-shortener has been started")
	storage.Log("url-shortener has been started", slog.LevelInfo)
	log.Debug("debug messages are enabled")
}

func setupLogger(env string) (*slog.Logger, error) {
	switch env {
	case envLocal:
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})), nil
	case envDev:
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})), nil
	case envProd:
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})), nil
	}
	return nil, errors.New("logger has not been set up")
}
