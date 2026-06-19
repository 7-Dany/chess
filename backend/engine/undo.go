package engine

import "github.com/7-Dany/chess/core"

func (e *DefaultEngine) Undo(ctx *core.TurnContext, snapshot core.Snapshot) {
	move := snapshot.Move

	// return piece to its previous position
	ctx.Board.Place(move.From, move.Piece)

	switch move.Type {
	case core.NORMAL, core.PROMOTION:
		e.restoreDestination(ctx, move)
	case core.CASTLING:
		e.restoreCastling(ctx, move)
	case core.EN_PASSANT:
		e.restoreEnPassant(ctx, move)
	}

	// restore state
	ctx.Sides = snapshot.PreviousSides
	ctx.EnPassantTarget = snapshot.PreviousEnPassantTarget
}

func (e *DefaultEngine) restoreDestination(ctx *core.TurnContext, move core.Move) {
	if move.HasCapture {
		ctx.Board.Place(move.To, move.Captured)
	} else {
		ctx.Board.Clear(move.To)
	}
}

func (e *DefaultEngine) restoreCastling(ctx *core.TurnContext, move core.Move) {
	ctx.Board.Clear(move.To)
	if move.To.File() > move.From.File() {
		// King side from (F -> H)
		moveRook(ctx, move.From.Rank(), core.FILE_F, core.FILE_H)
	} else {
		// Queen side from (D -> A)
		moveRook(ctx, move.From.Rank(), core.FILE_D, core.FILE_A)
	}
}

func (e *DefaultEngine) restoreEnPassant(ctx *core.TurnContext, move core.Move) {
	ctx.Board.Clear(move.To)
	capturedPawnPos := core.NewPosition(move.To.File(), move.From.Rank())
	ctx.Board.Place(capturedPawnPos, move.Captured)
}
