package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DatabaseClient *sql.DB

func InitDatabase(config ServerConfig) {
	MySqlClient, err := sql.Open("mysql", config.SqlUser+":"+config.SqlPass+"@/postowl")
	if err != nil {
		log.Fatal(err)
	}
	MySqlClient.SetMaxOpenConns(config.MaxSqlConns)
	MySqlClient.SetMaxIdleConns(config.MaxSqlIdleConns / 10)
}
