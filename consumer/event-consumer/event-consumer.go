package event_consumer

import (
	"context"
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

func (c *Consumer) Start(ctx context.Context) error {
	backoff := 2 * time.Second
	failCount := 0

	for {

		// 1. Check if context is already canceled before fetching
		select {
		case <-ctx.Done():
			log.Println("consumer stopped by context") // TODO Add proper logger
			return ctx.Err()
		default:
		}

		if failCount >= criticalFailThreshold {
			log.Printf("[ERR] consumer: %s", ErrCriticalFailure)
			return e.Wrap("[ERR] consumer:", ErrCriticalFailure) // handle it in main.
		}

		gotEvents, err := c.fetcher.Fetch(ctx, c.batchSize, c.updatesTimeout)
		if err != nil {
			log.Printf("[ERR] consumer: %s", err.Error())

			// Fix: context-aware backoff delay
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
			}

			backoff *= 2
			failCount++

			if backoff > maxBackoff {
				backoff = maxBackoff
			}

			if failCount >= failThreshold {
				log.Printf("[WARN] too many transient errors, pausing %v", pauseDuration)
				// Context-aware long pause
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(pauseDuration):
				}
			}
			continue
		}

		//backoff and failcount reset in case of success
		backoff = 2 * time.Second
		failCount = 0

		if len(gotEvents) == 0 { // !!! we check this on a Fetch() level, and there we return nil, so, its not possible to get empty gotEvents
			// Context-aware short delay
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(1 * time.Second):
			}
			continue
		}
		// Pass context to handler
		if err := c.handleEvents(ctx, gotEvents); err != nil {
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

/* Original handleEvents version
func (c *Consumer) handleEvents(events []events.Event) error {

	for _, event := range events {
		log.Printf("got new event: %s", event.Text)

		if err := c.processor.Process(event); err != nil {
			log.Printf("can't handle event: %s", err.Error())

			continue
		}
	}

	return nil
}
*/

func (c *Consumer) handleEvents(ctx context.Context, evts []events.Event) error {

	// Issue: Unlimited concurrency:
	// Options: use a worker pool or a semaphore / buffered channel to limit concurrency

	/* Unbounded concurrency - numEvents == goroutines
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
	*/

	// Example of semaphore / buffered channel

	var wg sync.WaitGroup

	concurrency := 10
	sem := make(chan struct{}, concurrency)

	for _, ev := range evts {

		select {
		case <-ctx.Done():
			// Context canceled, stop processing the rest of the batch, but wait for active ones to finish
			log.Println("stopping batch processing: context canceled")
			wg.Wait()
			return ctx.Err()
		case sem <- struct{}{}:
		}

		wg.Add(1)
		sem <- struct{}{}
		go func(ev events.Event) {
			defer wg.Done()
			defer func() { <-sem }()

			log.Printf("got new event: %s", ev.Text)

			if err := c.processor.Process(ctx, ev); err != nil {
				log.Printf("can't handle event: %v", err)
			}
		}(ev)
	}

	wg.Wait()
	return nil
}
