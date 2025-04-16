package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken string
}

func LoadConfig() Config {
	// Загружаем .env-файл
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	return Config{
		BotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
	}
}
