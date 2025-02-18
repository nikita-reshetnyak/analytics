package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string     `yaml:"env" env-default:"local"`
	StoragePath string     `yaml:"storage_path" env-required:"true"`
	GRPC        GrpcConfig `yaml:"grpc"`
}
type GrpcConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
	config_path := os.Getenv("CONFIG_PATH")
	if config_path == "" {
		log.Fatal("CONFIG_PATH env is not set")
	}
	if _, err := os.Stat(config_path); err != nil {
		log.Fatalf("error opening config file: %s", err)
	}
	var cfg Config
	err := cleanenv.ReadConfig(config_path, &cfg)
	if err != nil {
		log.Fatalf("error reading config file: %s", err)
	}
	return &cfg
}
