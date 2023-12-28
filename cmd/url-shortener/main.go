package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"
	"url-shortener/internal/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	clientOptions := options.Client().ApplyURI(cfg.Storage_Path)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	defer client.Disconnect(context.Background())

	err = client.Ping(context.Background(), nil)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	collection  := client.Database("url-shortener").Collection("logs")
	le := newLogEntry("starting url-shortener", 2)
	collection.InsertOne(context.TODO(), le)

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
