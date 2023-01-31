all:
	@go build -o bin/postowlserver main.go bot.go config.go messages.go redis.go database.go