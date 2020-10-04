package main

import (
	"github.com/joho/godotenv"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"os"
	"strconv"
)

func sendMessage(telegramToken string, telegramChannel int64, message string) {
	bot, _ := tgbotapi.NewBotAPI(telegramToken)
	msg := tgbotapi.NewMessage(telegramChannel, message)
	bot.Send(msg)
}

func main() {
	godotenv.Load()
	telegramToken := os.Getenv("TELEGRAM_TOKEN")
	telegramChannel, _ := strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)
	sendMessage(telegramToken, telegramChannel, "test message")
}
