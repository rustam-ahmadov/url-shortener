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
	fmt.Print(cfg)

	log, err := setupLogger(cfg.Env)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	storage := mongo_storage.NewStorage(cfg.Storage_Path, cfg.Storage_Name)
	storage.Log("url-shortener has been started",slog.LevelInfo)

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
