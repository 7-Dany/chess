// Package history stores the move history of a game. Each entry is a
// core.Snapshot, which holds the move and the pre-move state needed to
// undo it.
//
// The default implementation, MemoryStore, is an in-memory stack. Future
// implementations (RedisStore, DatabaseStore) can satisfy the same
// HistoryStore interface for persistence.
package history

import "github.com/7-Dany/chess/core"

// HistoryStore is a stack of snapshots — the undo history of a game.
type HistoryStore interface {
	// Push appends a snapshot to the top of the stack.
	Push(entry core.Snapshot)

	// Pop removes and returns the top snapshot. Returns false if the stack
	// is empty.
	Pop() (core.Snapshot, bool)

	// Peek returns the top snapshot without removing it. Returns false if
	// the stack is empty.
	Peek() (core.Snapshot, bool)

	// Len returns the number of snapshots in the stack.
	Len() int

	// All returns a copy of all snapshots (oldest first). The returned slice
	// is safe to modify without affecting the store.
	All() []core.Snapshot
}

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
