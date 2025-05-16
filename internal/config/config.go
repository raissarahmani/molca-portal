package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Addr string `yaml:"address"`
}

type Config struct {
	Env         string `yaml:"env" env:"ENV" env-default:"dev" env-required:"true"`
	StoragePath string `yaml:"storage_path" env:"STORAGE_PATH" env-default:"./storage/storage.db" env-required:"true"`
	HTTPServer  `yaml:"http_server" env-required:"true"`
}

func MustLoad() *Config {
	var cfg Config
	var configPath string

	configPath = os.Getenv("CONFIG_PATH")

	if configPath == "" {
		flags := flag.String("config", "", "path to config file")
		flag.Parse()

		configPath = *flags

		if configPath == "" {
			log.Fatal("Config file is not set")
		}
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist: %s", configPath)
	}

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Config error: %v", err)
	}

	return &cfg
}
