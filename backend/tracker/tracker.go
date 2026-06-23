package tracker

type Tracker interface {
	// Record increments the count for the given hash. Call after making a
	// move (once the hash has been updated).
	Record(hash uint64)

	// Undo decrements the count for the given hash. Call before undoing a
	// move (before the hash is reverted). If the count is already 0, this
	// is a no-op.
	Undo(hash uint64)

	// Count returns how many times the given hash has been recorded. Returns
	// 0 for a hash that was never recorded (or was recorded and fully undone).
	Count(hash uint64) int
}

type PositionTracker struct {
	counter map[uint64]int
}

// NewHashPositionTracker creates an empty tracker.
func NewPositionTracker() *PositionTracker {
	return &PositionTracker{
		counter: make(map[uint64]int),
	}
}

func (t *PositionTracker) Record(hash uint64) {
	t.counter[hash]++
}

func (t *PositionTracker) Undo(hash uint64) {
	count := t.counter[hash]
	if count <= 1 {
		delete(t.counter, hash)
		return
	}
	t.counter[hash] = count - 1
}

func (t *PositionTracker) Count(hash uint64) int {
	return t.counter[hash]
}
