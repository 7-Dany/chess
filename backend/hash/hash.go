// Package hash provides Zobrist hashing for chess positions. A Zobrist hash
// is a 64-bit integer that acts as a fingerprint for a board position: each
// piece-on-square, castling right, en passant file, and side-to-move
// contributes a random constant, combined with XOR.
//
// Because XOR is its own inverse, the same operation both applies and reverts
// a change — so Hash can be called identically after Apply or before Undo.
// InitHash computes the full hash from scratch (e.g. after a FEN decode);
// Hash then updates it incrementally for every subsequent move.
//
// The Hasher interface is the public contract; Zobrist is the default
// implementation backed by a fixed pseudo-random table seeded from a
// deterministic PCG source, so hashes are stable across runs.
package hash

import "github.com/7-Dany/chess/core"

// Hasher computes incremental Zobrist hashes for position identity.
// A single Hash call handles both apply and undo — because XOR is
// self-inverse, the same operations that update the hash forward
// also revert it backward.
type Hasher interface {
	// InitHash computes the full hash from scratch for the given context.
	// Call this once to bootstrap (e.g. after FEN decode), then use Hash
	// for every subsequent Apply/Undo.
	InitHash(ctx *core.TurnContext) uint64
	// Hash updates current by XOR-ing out the facts that became false
	// and XOR-ing in the facts that became true, as described by snap.
	// Call after Apply to advance the hash, or before Undo to revert it —
	// the result is identical either way.
	Hash(current uint64, move core.MoveHash) uint64
}
