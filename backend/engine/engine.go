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
	GetPseudoLegalMoves(position core.Position, ctx core.TurnContext) []core.Move

	// GetLegalMoves returns only moves that are strictly legal — pseudo-legal
	// moves that do not leave the moving side's king in check, plus castling.
	GetLegalMoves(position core.Position, ctx core.TurnContext) []core.Move

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

// DefaultEngine is the standard implementation of Engine.
// It delegates piece-specific logic to a PieceProvider and adds the
// game-state-aware checks (king safety, castling) that pieces cannot
// know about on their own.
type DefaultEngine struct {
	pieces *piece.PieceProvider
}

// NewDefaultEngine creates a DefaultEngine with the standard
// PieceProvider that knows about all six piece types.
func NewDefaultEngine() *DefaultEngine {
	return &DefaultEngine{
		pieces: piece.DefaultProvider(),
	}
}

func (e *DefaultEngine) GetPiece(pt core.PieceType) piece.Piece {
	return e.pieces.GetPiece(pt)
}
