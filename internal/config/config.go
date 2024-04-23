package config

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" env-required:"true"`
	StorageName string `yaml:"storage_name" env-required:"true" env-default:"url-shortener"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HttpServer  `yaml:"http_server"`
}

type HttpServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	User        string        `yaml:"user" env-required:"true" `
	Password    string        `yaml:"password" env-required:"true" env:"HTTP_SERVER_PASSWORD"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
}

func MustLoad() *Config {
	configPath, _ := os.Getwd()
	configPath = filepath.Join(configPath, "/config/local.yaml") //has to be dynamic from env var
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
