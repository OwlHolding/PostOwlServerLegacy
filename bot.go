package main

import (
	"encoding/json"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	fasthttp "github.com/valyala/fasthttp"
)

type ServerConfig struct {
	Token    string
	Url      string
	Port     string
	CertFile string
	KeyFile  string
}

var BotAPI *tgbotapi.BotAPI
var WebhookPath string

func LoadConfig(path string) ServerConfig {
	byte_config, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	var config ServerConfig
	err = json.Unmarshal(byte_config, &config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}

func CreateBot(config ServerConfig) *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		log.Fatal(err)
	}

	webhook, _ := tgbotapi.NewWebhookWithCert(config.Url+":"+config.Port+"/"+bot.Token,
		tgbotapi.FilePath(config.CertFile))
	_, err = bot.Request(webhook)
	if err != nil {
		log.Fatal(err)
	}

	return bot
}

func InitBot(config ServerConfig) {
	BotAPI = CreateBot(config)
	WebhookPath = "/" + config.Token
}

func ProcessRequest(ctx *fasthttp.RequestCtx) {
	if string(ctx.Path()) != WebhookPath {
		ctx.Error("", fasthttp.StatusForbidden)
	} else {
		var update tgbotapi.Update
		err := json.Unmarshal(ctx.PostBody(), &update)
		if err != nil {
			ctx.Error("", fasthttp.StatusBadRequest)
			return
		}
		go BotAPI.Send(tgbotapi.NewMessage(update.Message.From.ID, update.Message.Text))
	}
}
