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
