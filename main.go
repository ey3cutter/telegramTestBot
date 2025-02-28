package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"telegramTestBot/bot"
	"telegramTestBot/db"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)

	db.InitDb(connStr)

	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatal("Bot token not found")
	}

	var newBot = bot.NewBot(botToken)
	log.Printf("Авторизован как %s", newBot.Bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := newBot.Bot.GetUpdatesChan(u)

	newBot.HandleUpdates(updates)
}
