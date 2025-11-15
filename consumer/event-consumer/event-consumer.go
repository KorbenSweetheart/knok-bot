package event_consumer

import (
	"errors"
	"knok-bot/events"
	"knok-bot/lib/e"
	"log"
	"sync"
	"time"
)

type Consumer struct {
	fetcher        events.Fetcher
	processor      events.Processor
	batchSize      int
	updatesTimeout int
}

const (
	maxBackoff            = 30 * time.Second
	failThreshold         = 10
	pauseDuration         = 15 * time.Minute
	criticalFailThreshold = 20
)

var ErrCriticalFailure = errors.New("too many consecutive failures")

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
	failCount := 0

	for {

		// refactor with context and select

		if failCount >= criticalFailThreshold {
			log.Printf("[ERR] consumer: %s", ErrCriticalFailure)
			return e.Wrap("[ERR] consumer:", ErrCriticalFailure) // handle it in main.
		}

		gotEvents, err := c.fetcher.Fetch(c.batchSize, c.updatesTimeout)
		if err != nil {
			log.Printf("[ERR] consumer: %s", err.Error())

			time.Sleep(backoff)
			backoff *= 2
			failCount++

			if backoff > maxBackoff {
				backoff = maxBackoff
			}

			if failCount >= failThreshold {
				log.Printf("[WARN] too many transient errors, pausing %v", pauseDuration)
				time.Sleep(pauseDuration)
			}

			continue
		}

		//backoff and failcount reset in case of success
		backoff = 2 * time.Second
		failCount = 0

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

func (c *Consumer) handleEvents(evts []events.Event) error {
	// for concurrency, hint: sync.WaitGroup{}
	// Issue: Unlimited concurrency:
	// use a worker pool
	// or a semaphore / buffered channel to limit concurrency

	var wg sync.WaitGroup

	for _, event := range evts {

		log.Printf("got new event: %s", event.Text)

		wg.Add(1)
		go func(event events.Event) {
			defer wg.Done()

			if err := c.processor.Process(event); err != nil {
				log.Printf("can't handle event: %s", err.Error())
			}
		}(event)
	}

	wg.Wait()

	return nil
}
