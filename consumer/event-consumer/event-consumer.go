package event_consumer

import (
	"knok-bot/events"
	"log"
	"time"
)

type Consumer struct {
	fetcher        events.Fetcher
	processor      events.Processor
	batchSize      int
	updatesTimeout int
}

const (
	maxBackoff = 30 * time.Second
	// maxFailCount = 10
)

func New(fetcher events.Fetcher, processor events.Processor, batchSize int, timeout int) Consumer {
	return Consumer{
		fetcher:        fetcher,
		processor:      processor,
		batchSize:      batchSize,
		updatesTimeout: timeout,
	}
}

func (c *Consumer) Start() error {
	backoff := 2 * time.Second

	for {
		// tip implement retry in fetcher with 3 tries and window of few seconds. We can even make exponential raise for window. and gave up on some limit.
		gotEvents, err := c.fetcher.Fetch(c.batchSize, c.updatesTimeout)
		if err != nil {
			log.Printf("[ERR] consumer: %s", err.Error())

			time.Sleep(backoff)
			backoff *= 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}

			continue
		}

		backoff = 2 * time.Second //backoff reset in case of success

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)

			continue
		}

		if err := c.handleEvents(gotEvents); err != nil {
			log.Print(err)

			continue
		}
	}
}

/*
Problem -> solutions
1. loss of the events -> retry, return to storage, fallback (locally in case of network issues), approval (for fetcher)
2. handle of the entire batch -> stop when error occurs, error counter (or stop when x errors occur)
3. concurrency (parallel handle)
*/

func (c *Consumer) handleEvents(events []events.Event) error {
	// for concurrency, hint: sync.WaitGroup{}
	for _, event := range events {
		log.Printf("got new event: %s", event.Text)

		if err := c.processor.Process(event); err != nil {
			log.Printf("can't handle event: %s", err.Error())

			continue
		}
	}

	return nil
}
