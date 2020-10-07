package main

import (
	"github.com/joho/godotenv"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"os"
	"strconv"
	"log"
	"github.com/streadway/amqp"
)
func failOnError(err error, msg string) {
  if err != nil {
    log.Fatalf("%s: %s", msg, err)
  }
}

func sendMessage(telegramToken string, telegramChannel int64, message string) {
	bot, err := tgbotapi.NewBotAPI(telegramToken)
	failOnError(err, "Failed tgbotapi.NewBotAPI initialize")
	msg := tgbotapi.NewMessage(telegramChannel, message)
	bot.Send(msg)
}

func main() {
	godotenv.Load()
	telegramToken := os.Getenv("TELEGRAM_TOKEN")
	telegramChannel, err := strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)
	failOnError(err, "Failed to parse telegramChannel ENV")
	amqpLink := os.Getenv("AMQP")

	conn, err := amqp.Dial(amqpLink)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"notify", // name
		true,     // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		amqp.Table{"x-dead-letter-exchange": "notify.dlx", "x-dead-letter-routing-key": "notify.dlx"},     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			sendMessage(telegramToken, telegramChannel, string(d.Body))
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
