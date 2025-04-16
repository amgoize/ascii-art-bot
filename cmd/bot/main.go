package main

import (
	"ascii-art-server/internal/config"
	"ascii-art-server/pkg/telegram"
	"log"
)

func main() {
	cfg := config.LoadConfig()

	bot, err := telegram.NewBot(cfg.BotToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Start()
}
