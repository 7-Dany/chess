package tracker

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
