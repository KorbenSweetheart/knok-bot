package main

import (
	"context"
	"flag"
	tgClient "knok-bot/clients/telegram"
	event_consumer "knok-bot/consumer/event-consumer"
	"knok-bot/events/telegram"
	"knok-bot/storage/sqlite"
	"log"
	"os/signal"
	"syscall"
	"time"
)

const (
	tgBotHost         = "api.telegram.org" // its better to get it the same way as Token -> via parameter.
	storagePath       = "file_storage"
	sqliteStoragePath = "data/sqlite/storage.db"
	batchSize         = 100
	updatesTimeout    = 60 // seconds
)

func main() {
	ctx, shutdown := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer shutdown() // Ensure resources are cleaned up

	// s := files.New(storagePath)
	s, err := sqlite.New(sqliteStoragePath)
	if err != nil {
		log.Fatalf("can't connect to storage: %v", err)
	}

	reqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := s.Init(reqCtx); err != nil {
		log.Fatalf("can't init storage: %v", err)
	}

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		s,
	)

	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize, updatesTimeout)
	if err := consumer.Start(ctx); err != nil {
		switch err {
		case event_consumer.ErrCriticalFailure:
			log.Fatal("service is stopped", event_consumer.ErrCriticalFailure) // tmp solution need to handle it gracefully
		default:
			log.Fatal("service is stopped", err)
		}
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
