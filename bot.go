package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	redis "github.com/redis/go-redis/v9"
	fasthttp "github.com/valyala/fasthttp"
)

var BotAPI *tgbotapi.BotAPI
var WebhookPath string
var MainCtx context.Context
var RedisClient *redis.Client
var MySqlClient *sql.DB

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

	MainCtx = context.Background()
	RedisClient = redis.NewClient(&redis.Options{Addr: config.RedisUrl})
	_, err = RedisClient.Ping(MainCtx).Result()
	if err != nil {
		log.Fatal(err)
	}

	MySqlClient, err = sql.Open("mysql", config.SqlUser+":"+config.SqlPass+"@/postowl")
	if err != nil {
		log.Fatal(err)
	}
	MySqlClient.SetMaxOpenConns(config.MaxUsers)
	MySqlClient.SetMaxIdleConns(config.MaxUsers / 10)
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
