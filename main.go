package main

import (
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kelseyhightower/envconfig"
)

type BotConfig struct {
	Debug      bool
	Token      string
	WebhookUrl string
}

var serviceKeyboard = tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{
	tgbotapi.NewInlineKeyboardButtonData("Satisfactory", "restart-satisfactory"),
	tgbotapi.NewInlineKeyboardButtonData("Teamspeak", "restart-teamspeak"),
})

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

		if update.Message != nil {
			// Process a command
			command := update.Message.Command()
			switch command {
			case "restart":
				reply := tgbotapi.NewMessage(update.Message.Chat.ID, "Which service?")
				reply.ReplyMarkup = serviceKeyboard
				_, err := bot.Request(reply)
				if err != nil {
					log.Printf("reset message not set: %v", err.Error())
				}
			case "close":
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "")
				msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
				if _, err := bot.Send(msg); err != nil {
					panic(err)
				}
			}
		}

		// Process a callback
		if update.CallbackQuery != nil {
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := bot.Request(callback); err != nil {
				panic(err)
			}

			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Okay restarting...")
			if _, err := bot.Send(msg); err != nil {
				panic(err)
			}
		}

	}
}