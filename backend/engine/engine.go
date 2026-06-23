// Package engine is the move-generation and move-execution layer of the chess
// backend. It sits between the pure domain types in core and the
// piece-movement rules in piece, combining both to answer three questions:
//
//  1. What moves can a piece make? (pseudo-legal, then legal after king-safety
//     filtering)
//  2. Is a given square attacked? (used for check detection and castling
//     validation)
//  3. How does a move change the board — and how do we take it back?
//
// The Engine interface is the public contract; DefaultEngine is the standard
// implementation. All methods operate on core.TurnContext, which holds the
// board, side-to-move, castling rights, and en passant target in one struct.
//
// Move generation is allocation-friendly: callers pass in a []core.Move buffer
// and the methods append to it. Stack-allocated buffers (e.g. [core.MAX_MOVES]
// or [MAX_TOTAL_MOVES]) are enough to keep hot paths heap-free.
package engine

import (
	"github.com/7-Dany/chess/core"
	"github.com/7-Dany/chess/piece"
)

// Engine is the move-generation and move-execution subsystem.
// It answers queries about positions, generates legal moves, and
// applies/undoes moves on a TurnContext.
type Engine interface {
	// GetPseudoLegalMoves returns all pseudo-legal moves for the piece
	// at the given position.
	GetPseudoLegalMoves(moves []core.Move, position core.Position, ctx core.TurnContext) []core.Move

	// GetLegalMoves returns only moves that are strictly legal — pseudo-legal
	// moves that do not leave the moving side's king in check, plus castling.
	GetLegalMoves(moves []core.Move, position core.Position, ctx core.TurnContext) []core.Move

	// GetAllLegalMoves appends every legal move for the side to move into
	// moves and returns the extended slice. Iterates all 64 squares, calls
	// GetLegalMoves per friendly piece, and accumulates into one buffer.
	// Pass a buffer of at least MAX_TOTAL_MOVES (256).
	GetAllLegalMoves(moves []core.Move, ctx core.TurnContext) []core.Move

	// HasAnyLegalMoves reports whether color has at least one legal move
	// anywhere on the board. Built for checkmate/stalemate detection, which
	// need a yes/no answer to "can this side move at all".
	HasAnyLegalMoves(ctx core.TurnContext) bool

	// IsSquareAttacked reports whether any piece of the given color
	// attacks the specified position.
	IsSquareAttacked(position core.Position, attackerColor core.PieceColor, ctx core.BoardContext) bool

	// Apply mutates the TurnContext in place, applying the move and
	// returning a Snapshot for undo. Handles piece movement, captures,
	// en passant, castling, promotion, king position, en passant target,
	// castling rights.
	Apply(ctx *core.TurnContext, move core.Move) core.Snapshot

	// Undo reverses an Apply using the saved Snapshot.
	Undo(ctx *core.TurnContext, snap core.Snapshot)
}

// MAX_TOTAL_MOVES is the maximum number of legal moves an entire side can
// have in a single position. The theoretical maximum is 218 (achieved in
// some contrived positions); 256 provides headroom. Pass a buffer of at
// least this size to GetAllLegalMoves to guarantee zero heap allocations.
const MAX_TOTAL_MOVES = 256

// DefaultEngine is the standard implementation of Engine.
// It delegates piece-specific logic to a PieceProvider and adds the
// game-state-aware checks (king safety, castling) that pieces cannot
// know about on their own.
type DefaultEngine struct {
	pieces piece.Pieces
}

var defaultEngine = DefaultEngine{pieces: piece.GetDefaultPieces()}

// NewDefaultEngine creates a DefaultEngine with the standard
// PieceProvider that knows about all six piece types.
func GetDefaultEngine() DefaultEngine {
	return defaultEngine
}
