package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	redis "github.com/redis/go-redis/v9"
	fasthttp "github.com/valyala/fasthttp"
)

var BotAPI *tgbotapi.BotAPI
var WebhookPath string
var MainCtx context.Context
var RedisClient *redis.Client

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
	MainCtx = context.Background()
	RedisClient = redis.NewClient(&redis.Options{Addr: config.RedisUrl})
	_, err := RedisClient.Ping(MainCtx).Result()
	if err != nil {
		log.Fatal(err)
	}
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
		// go BotAPI.Send(tgbotapi.NewMessage(update.Message.From.ID, update.Message.Text))
		go Router(update.Message.From.ID, update.Message.Text)
	}
}

func Router(chatID int64, message string) {
	chat := fmt.Sprint(chatID)
	value, err := RedisClient.Get(MainCtx, chat).Result()
	if err != nil {
		if err == redis.Nil {
			value = "Hello"
		} else {
			log.Fatal(err)
		}
	}
	BotAPI.Send(tgbotapi.NewMessage(chatID, value))
	err = RedisClient.Set(MainCtx, chat, message, 0).Err()
	if err != nil {
		log.Fatal(err)
	}
}
