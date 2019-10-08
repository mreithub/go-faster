package faster

import "time"

// Singleton -- global go-faster instance
var Singleton = New(true)

// TakeSnapshot -- Returns a Snapshot of the current Faster state
func TakeSnapshot() *Snapshot {
	return Singleton.TakeSnapshot()
}

// Track -- Tracks an instance of 'key' (in singleton mode)
func Track(key ...string) *Tracker {
	return Singleton.Track(key...)
}

// TrackFn -- Tracks the calling function (using ["src", "pkgName", "typeName", "fn()"] as key - omitting typeName if empty)
func TrackFn() *Tracker {
	var key = Singleton.getCaller(1)
	return Singleton.Track(key...)
}

// Reset -- resets the internal state of the singleton Faster instance
func Reset() {
	Singleton.Reset()
}

// SetTicker -- sets up periodic snapshots (in singleton mode)
func SetTicker(name string, interval time.Duration, keep int) {
	Singleton.SetTicker(name, interval, keep)
}
