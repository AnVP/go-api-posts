package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Db   DbConfig
	Auth AuthConfig
}

type DbConfig struct {
	Host string
}

type AuthConfig struct {
	Token string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return &Config{
		Db: DbConfig{
			Host: os.Getenv("DB_HOST"),
		},
		Auth: AuthConfig{
			Token: os.Getenv("TOKEN"),
		},
	}
}
