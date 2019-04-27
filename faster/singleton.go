package faster

import "time"

// Singleton -- global GoFaster instance
var Singleton = New()

// GetSnapshot -- Returns a Snapshot of the GoFaster (synchronously)
func GetSnapshot() *Snapshot {
	return Singleton.GetSnapshot()
}

// Ref -- References an instance of 'key' (in singleton mode)
func Ref(key ...string) *Instance {
	return Singleton.Ref(key...)
}

// Reset -- resets the internal state of the singleton GoFaster instance
func Reset() {
	Singleton.Reset()
}

func SetTicker(name string, interval time.Duration, keep int) {
	Singleton.SetTicker(name, interval, keep)
}
