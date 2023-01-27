package main

import (
	"log"
	"net/http"
)

func main() {
	config := LoadConfig("config.json")

	InitBot(config)

	http.HandleFunc("/"+config.Token, ProcessRequest)

	log.Print("Server started")

	err := http.ListenAndServeTLS(":"+config.Port, config.CertFile, config.KeyFile, nil)
	log.Fatal(err)
}
