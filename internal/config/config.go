package config

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env          string `yml: env-required: "true"`
	Storage_Path string `yml: env-required: "true"`
	Storage_Name string `yml: env-required: "true"`
	HttpServer   `yaml:"http_server"`
}

type HttpServer struct {
	Address     string        `yml: "address" env-default: "localhost:8080"`
	Timeout     time.Duration `yml: "timeout" env-default:"4s`
	IdleTimeout time.Duration `yml: "idle_timeout" env-default: "60s"`
}

func MustLoad() *Config {
	configPath, _ := os.Getwd()
	configPath = filepath.Join(configPath, "../../config/local.yml") //has to be dynamic from env var
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	_, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("connot read config: %s", err)
	}
	return &cfg
}
