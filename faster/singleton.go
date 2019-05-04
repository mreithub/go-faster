package faster

import "time"

// Singleton -- global GoFaster instance
var Singleton = New(true)

// GetSnapshot -- Returns a Snapshot of the GoFaster (synchronously)
func GetSnapshot() *Snapshot {
	return Singleton.GetSnapshot()
}

// Track -- Tracks an instance of 'key' (in singleton mode)
func Track(key ...string) *Tracker {
	return Singleton.Track(key...)
}

// Reset -- resets the internal state of the singleton GoFaster instance
func Reset() {
	Singleton.Reset()
}

// SetTicker -- sets up periodic snapshots (in singleton mode)
func SetTicker(name string, interval time.Duration, keep int) {
	Singleton.SetTicker(name, interval, keep)
}
