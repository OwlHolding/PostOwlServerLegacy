package main

import (
	"encoding/json"
	"log"
	"os"
)

type ServerConfig struct {
	Token    string
	Url      string
	Port     string
	CertFile string
	KeyFile  string
	RedisUrl string
	SqlUser  string
	SqlPass  string
	MaxUsers int
}

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

func LoadConfigFromEnv(variable string) ServerConfig {
	value, exists := os.LookupEnv(variable)
	if !exists {
		log.Fatalf("Variable %s does not exist", variable)
	}
	return LoadConfig(value)
}
