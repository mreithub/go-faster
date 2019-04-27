package internal

// event types (for internal communication):
import "time"

// EventType -- enum type for different internal GoFaster events
type EventType int

const (
	// EvStop -- stop the goroutine handling this GoFaster instance
	EvStop EventType = iota
	// EvReset -- resets this GoFaster instance
	EvReset EventType = iota
	// EvSnapshot -- Takes a snapshot and sends it to snapshotChannel
	EvSnapshot EventType = iota
	// EvRef increments a ref counter
	EvRef EventType = iota
	// EvDeref a ref counter (and updates the total count + time)
	EvDeref EventType = iota
)

// Event -- internal events
type Event struct {
	Type EventType
	Path []string
	Took time.Duration
}
