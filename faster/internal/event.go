package internal

// event types (for internal communication):
import "time"

// EventType -- enum type for different internal go-faster events
type EventType int

const (
	// EvStop -- stop the goroutine handling this Faster instance
	EvStop EventType = iota
	// EvReset -- resets this Faster instance
	EvReset EventType = iota
	// EvSnapshot -- Takes a snapshot and sends it to snapshotChannel
	EvSnapshot EventType = iota

	// EvTrack -- increments a ref counter
	EvTrack EventType = iota
	// EvDone -- deref a ref counter (and updates the total count + time)
	EvDone EventType = iota
)

// Event -- internal events
type Event struct {
	Type EventType
	Path []string
	Took time.Duration
}
