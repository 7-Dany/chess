package engine

import "github.com/7-Dany/chess/core"

func (e DefaultEngine) Undo(ctx *core.TurnContext, snapshot core.Snapshot) {
	move := snapshot.Move

	switch move.Type {
	case core.NORMAL:
		e.undoNormal(ctx, move)
	case core.PROMOTION:
		e.undoPromotion(ctx, move)
	case core.CASTLING:
		e.undoCastling(ctx, move)
	case core.EN_PASSANT:
		e.undoEnPassant(ctx, move)
	}

	// Restore state captured at the start of Apply.
	ctx.Sides = snapshot.PreviousSides
	ctx.EnPassantTarget = snapshot.PreviousEnPassantTarget
	ctx.HalfMoveClock = snapshot.PreviousHalfMoveClock
	ctx.FullMoveNumber = snapshot.PreviousFullMoveNumber
}

func (e DefaultEngine) undoNormal(ctx *core.TurnContext, move core.Move) {
	// Move the piece back to its origin.
	ctx.Board.Move(move.To, move.From)

	// Restore captured piece, if any.
	if move.HasCapture {
		ctx.Board.Place(move.To, move.Captured)
	}
}

func (e DefaultEngine) undoPromotion(ctx *core.TurnContext, move core.Move) {
	// Clear the prompoted pawn.
	ctx.Board.Clear(move.To)

	// Return the pawn back to its position.
	ctx.Board.Place(move.From, move.Piece)

	// Restore captured piece, if any.
	if move.HasCapture {
		ctx.Board.Place(move.To, move.Captured)
	}
}

func (e DefaultEngine) undoCastling(ctx *core.TurnContext, move core.Move) {
	// Restore king to origin, clear destination.
	ctx.Board.Clear(move.To)
	ctx.Board.Place(move.From, move.Piece)

	// Restore rook. King-side: F -> H. Queen-side: D -> A.
	rookTo, rookFrom := move.CastlingRookPositions()
	ctx.Board.Move(rookFrom, rookTo)
}

func (e DefaultEngine) undoEnPassant(ctx *core.TurnContext, move core.Move) {
	// Restore pawn to origin, clear destination.
	ctx.Board.Clear(move.To)
	ctx.Board.Place(move.From, move.Piece)

	// Restore the captured pawn behind the destination.
	ctx.Board.Place(move.EnPassantCapturedPosition(), move.Captured)
}
