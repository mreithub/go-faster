package faster

// Singleton -- global GoFaster instance
var Singleton = New()

// GetInstance -- Returns a scoped instance (matching the given scope path)
func GetInstance(path ...string) *Faster {
	return Singleton.GetChild(path...)
}

// GetSnapshot -- Returns a Snapshot of the GoFaster (synchronously)
func GetSnapshot() Snapshot {
	return Singleton.GetSnapshot()
}

// Ref -- References an instance of 'key' (in singleton mode)
func Ref(key string) *Instance {
	return Singleton.Ref(key)
}

// Reset -- resets the internal state of the singleton GoFaster instance
func Reset() {
	Singleton.Reset()
}
