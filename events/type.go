package events

import "context"

type Fetcher interface {
	Fetch(ctx context.Context, limit int, timeout int) ([]Event, error) // offset will be implemented inside
}

type Processor interface {
	Process(ctx context.Context, e Event) error
}

type Type int

const (
	Unknown Type = iota
	Message
)

type Event struct {
	Type Type
	Text string
	Meta interface{}
}
