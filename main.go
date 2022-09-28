package main

import (
	"log"
	"net/http"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kelseyhightower/envconfig"
)

type BotConfig struct {
	Debug      bool
	Token      string
	WebhookUrl string
}

func main() {

	var conf BotConfig
	err := envconfig.Process("", &conf)
	if err != nil {
		log.Fatal(err.Error())
	}

	bot, err := tgbotapi.NewBotAPI(conf.Token)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = conf.Debug

	log.Printf("Authorized on account %s", bot.Self.UserName)

	wh, _ := tgbotapi.NewWebhook(conf.WebhookUrl + bot.Token)

	_, err = bot.Request(wh)
	if err != nil {
		log.Fatal(err)
	}
	updates := bot.ListenForWebhook("/" + bot.Token)
	go http.ListenAndServe("0.0.0.0:8080", nil)

	for update := range updates {
		log.Printf("Got update: %+v\n", update)
		receivedMessage := update.Message
		if receivedMessage != nil {
			log.Printf("Got message: %si from user %s\n", receivedMessage.Text, receivedMessage.From.UserName)
		}
	}
}