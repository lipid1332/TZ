package config

import (
	"log"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type DBConfig struct {
	Host     string `env:"POSTGRES_HOST"`
	Port     int    `env:"POSTGRES_PORT"`
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	DBname   string `env:"POSTGRES_DB"`
	AppPort  string `env:"APP_PORT"`
}

func New(configPath string) *DBConfig {
	cfg := DBConfig{}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("confing file does not exist: %s", configPath)
	}

	if err := godotenv.Load(configPath); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
