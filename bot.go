package main

import (
	"encoding/json"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	fasthttp "github.com/valyala/fasthttp"
)

var BotAPI *tgbotapi.BotAPI
var WebhookPath string

func InitBot(config ServerConfig) {
	WebhookPath = "/" + config.Token

	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		log.Fatal(err)
	}

	webhook, _ := tgbotapi.NewWebhookWithCert(config.Url+":"+config.Port+WebhookPath,
		tgbotapi.FilePath(config.CertFile))
	_, err = bot.Request(webhook)
	if err != nil {
		log.Fatal(err)
	}

	BotAPI = bot
}

func ProcessRequest(ctx *fasthttp.RequestCtx) {
	if string(ctx.Path()) != WebhookPath {
		ctx.Error("", fasthttp.StatusForbidden)
		return
	}
	var update tgbotapi.Update
	err := json.Unmarshal(ctx.PostBody(), &update)
	if err != nil {
		ctx.Error("", fasthttp.StatusBadRequest)
		return
	}

	BotAPI.Send(tgbotapi.NewMessage(
		update.Message.From.ID, update.Message.Text))
}
