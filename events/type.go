package events

type Fetcher interface {
	Fetch(limit int) ([]Event, error) // offset will be implemented inside
}

type Processor interface {
	Process(e Event) error
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
