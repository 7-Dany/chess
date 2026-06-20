package engine

import "github.com/7-Dany/chess/core"

func (e *DefaultEngine) Undo(ctx *core.TurnContext, snapshot core.Snapshot) {
	move := snapshot.Move

	switch move.Type {
	case core.NORMAL, core.PROMOTION:
		e.undoNormal(ctx, move)
	case core.CASTLING:
		e.undoCastling(ctx, move)
	case core.EN_PASSANT:
		e.undoEnPassant(ctx, move)
	}

	// Restore state captured at the start of Apply.
	ctx.Sides = snapshot.PreviousSides
	ctx.EnPassantTarget = snapshot.PreviousEnPassantTarget
}

func (e *DefaultEngine) undoNormal(ctx *core.TurnContext, move core.Move) {
	// Move the piece back to its origin. Promotion: Place overwrites the
	// promoted piece with the original pawn.
	ctx.Board.Move(move.To, move.From)
	if move.Type == core.PROMOTION {
		ctx.Board.Place(move.From, move.Piece)
	}

	// Restore captured piece, if any.
	if move.HasCapture {
		ctx.Board.Place(move.To, move.Captured)
	}
}

func (e *DefaultEngine) undoCastling(ctx *core.TurnContext, move core.Move) {
	// Restore king to origin, clear destination.
	ctx.Board.Clear(move.To)
	ctx.Board.Place(move.From, move.Piece)

	// Restore rook. King-side: F -> H. Queen-side: D -> A.
	rank := move.From.Rank()
	if move.To.File() > move.From.File() {
		ctx.Board.Move(core.NewPosition(core.FILE_F, rank), core.NewPosition(core.FILE_H, rank))
	} else {
		ctx.Board.Move(core.NewPosition(core.FILE_D, rank), core.NewPosition(core.FILE_A, rank))
	}
}

func (e *DefaultEngine) undoEnPassant(ctx *core.TurnContext, move core.Move) {
	// Restore pawn to origin, clear destination.
	ctx.Board.Clear(move.To)
	ctx.Board.Place(move.From, move.Piece)

	// Restore the captured pawn behind the destination.
	ctx.Board.Place(core.NewPosition(move.To.File(), move.From.Rank()), move.Captured)
}
