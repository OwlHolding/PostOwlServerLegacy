package main

import (
	"log"

	"github.com/valyala/fasthttp"
)

func main() {
	config := LoadConfig("config.json")

	InitBot(config)

	log.Print("Server started")

	err := fasthttp.ListenAndServeTLS(":"+config.Port, config.CertFile, config.KeyFile,
		ProcessRequest)
	log.Fatal(err)
}
