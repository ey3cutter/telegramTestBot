package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"telegramTestBot/bot"
	"telegramTestBot/db"
)

func main() {
	db.InitDb("user=postgres password=12341234 dbname=gotgtest sslmode=disable")

	var newBot = bot.NewBot("7792067052:AAFLPsJsLNyH4_75fdd9k_SNpE5Dh5vNcaA")
	log.Printf("Авторизован как %s", newBot.Bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := newBot.Bot.GetUpdatesChan(u)

	newBot.HandleUpdates(updates)
}
