package main

import (
	"flag"
	"log"
	"read-reminder-bot/clients/telegram"
)

const (
	tgBotHost = "api.telegram.org" // its better to get it the same way as Token -> via parameter.
)

func main() {

	tgClient = telegram.New(tgBotHost, mustToken())

	// fetcher = fetcher.New()

	// processor = processor.New()

	// consumer.Start(fetcher, processor)

}

func mustToken() string {
	token := flag.String(
		"token",
		"",
		"token for access to telegram bot",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
