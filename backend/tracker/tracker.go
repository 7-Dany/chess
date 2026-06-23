// Package tracker counts how many times each board position has occurred
// during a game, identified by its Zobrist hash. It is used to detect
// threefold repetition: when the same position appears three times, the game
// can be claimed as a draw.
//
// The Tracker interface is the public contract. PositionTracker is the default
// in-memory implementation backed by a map. Record is called after each move
// (once the hash is updated) and Undo is called before each move is reversed
// (before the hash is reverted).
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
