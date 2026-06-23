package history

import "github.com/7-Dany/chess/core"

// MemoryStore is the default in-memory HistoryStore. It uses a slice as a
// stack. The zero value is NOT ready to use — use NewMemoryStore.
type MemoryStore struct {
	entries []core.Snapshot
}

// NewMemoryStore creates an empty in-memory history store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		entries: make([]core.Snapshot, 0, 64),
	}
}

func (s *MemoryStore) Push(entry core.Snapshot) {
	s.entries = append(s.entries, entry)
}

func (s *MemoryStore) Pop() (core.Snapshot, bool) {
	if len(s.entries) == 0 {
		return core.Snapshot{}, false
	}
	last := len(s.entries) - 1
	entry := s.entries[last]
	s.entries = s.entries[:last]
	return entry, true
}

func (s *MemoryStore) Peek() (core.Snapshot, bool) {
	if len(s.entries) == 0 {
		return core.Snapshot{}, false
	}
	return s.entries[len(s.entries)-1], true
}

func (s *MemoryStore) Len() int {
	return len(s.entries)
}

func (s *MemoryStore) All() []core.Snapshot {
	out := make([]core.Snapshot, len(s.entries))
	copy(out, s.entries)
	return out
}
