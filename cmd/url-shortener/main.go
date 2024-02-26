package main

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"url-shortener/internal"
	"url-shortener/internal/config"
	"url-shortener/internal/http-server/hadlers/url/redirect"
	"url-shortener/internal/http-server/hadlers/url/save"
	"url-shortener/internal/http-server/middleware/logger"
	"url-shortener/internal/http-server/middleware/verify"
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

	storage, err := mongo_storage.NewStorage(cfg.Storage_Path, cfg.Storage_Name)
	if err != nil {
		fmt.Println("failed to setup storage")
		os.Exit(1)
	}

	clw := internal.NewCustomLogWriter(storage, "logs")
	log, err := setupLogger(clw, cfg.Env)
	if err != nil {
		fmt.Printf("err occured while initializing logger: %s", err.Error())
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Group(func(r chi.Router) {
		router.Use(verify.JwtMiddleware)

		router.Use(middleware.RequestID)
		router.Use(middleware.Logger)
		router.Use(logger.New(log))
		router.Use(middleware.Recoverer)
		router.Use(middleware.URLFormat)

		router.Post("/url", save.New(log, storage))
	})
	router.Group(func(r chi.Router) {
		router.Get("/*", redirect.New(log, storage))
	})

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
}

func setupLogger(writer io.Writer, env string) (*slog.Logger, error) {
	switch env {
	case envLocal:
		return slog.New(slog.NewJSONHandler(writer, &slog.HandlerOptions{Level: slog.LevelDebug})), nil
	case envDev:
		return slog.New(slog.NewJSONHandler(writer, &slog.HandlerOptions{Level: slog.LevelDebug})), nil
	case envProd:
		return slog.New(slog.NewJSONHandler(writer, &slog.HandlerOptions{Level: slog.LevelInfo})), nil
	}
	return nil, errors.New("logger has not been set up")
}
