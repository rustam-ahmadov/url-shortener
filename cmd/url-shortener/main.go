package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/http-server/hadlers/url/redirect"
	"url-shortener/internal/http-server/hadlers/url/save"
	"url-shortener/internal/http-server/middleware/logger"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage/mongo_storage"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
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
		log.Error("failed to init storage", sl.Err(err))
		storage.Log("failded to init storage", slog.LevelError)
		os.Exit(1)
	}
	log.Info("url-shortener started...")
	storage.Log("url-shortener started...", slog.LevelInfo)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(logger.New(storage, log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/url", save.New(log, storage))
	router.Get("/*", redirect.New(log, storage))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
		storage.Log("failed to start server", slog.LevelError)
	}

	log.Error("server stopped")
	storage.Log("server stopped", slog.LevelError)
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
