package main

import (
	"flag"
	tgClient "knok-bot/clients/telegram"
	event_consumer "knok-bot/consumer/event-consumer"
	"knok-bot/events/telegram"
	"knok-bot/storage/files"
	"log"
)

const (
	tgBotHost      = "api.telegram.org" // its better to get it the same way as Token -> via parameter.
	storagePath    = "file_storage"
	batchSize      = 100
	updatesTimeout = 60 // seconds
)

func main() {

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		files.New(storagePath),
	)

	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize, updatesTimeout)
	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}

func mustToken() string {
	token := flag.String(
		"tg-bot-token",
		"",
		"token for access to telegram bot",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
