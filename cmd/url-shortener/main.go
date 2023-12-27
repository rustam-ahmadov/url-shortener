package main

import (
	"fmt"
	"log/slog"
	"url-shortener/internal/config"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	fmt.Print(cfg)

}

func setupLogger(env string) {
	var log *slog.Logger

	switch env {
	case envLocal:

	}
}
