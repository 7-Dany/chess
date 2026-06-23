package engine

import "github.com/7-Dany/chess/core"

func (e DefaultEngine) Apply(ctx *core.TurnContext, move core.Move) core.Snapshot {
	snapshot := ctx.Snapshot(move)

	switch move.Type {
	case core.NORMAL:
		e.applyNormal(ctx, move)
	case core.PROMOTION:
		e.applyPromotion(ctx, move)
	case core.CASTLING:
		e.applyCastling(ctx, move)
	case core.EN_PASSANT:
		e.applyEnPassant(ctx, move)
	}

	// En passant target: set on double pawn push, cleared otherwise.
	ctx.SetEnPassantTarget(move)

	return snapshot
}

func (e DefaultEngine) applyNormal(ctx *core.TurnContext, move core.Move) {
	// Move the piece (mover is unchanged for a NORMAL move).
	ctx.Board.Move(move.From, move.To)

	// King moves update king position and forfeit all castling rights.
	if move.Piece.Type == core.KING {
		// KingPosition is updated for the game controller's benefit; the
		// legality checker derives king position from move.To directly.
		ctx.Sides[move.Piece.Color].KingPosition = move.To
		ctx.Sides[move.Piece.Color].ClearCastlingRights()
	}

	// Rook moves and rook captures each forfeit one castling right.
	ctx.ForfeitCastlingRight(move)
}

func (e DefaultEngine) applyPromotion(ctx *core.TurnContext, move core.Move) {
	piece := move.Piece
	piece.Type = move.PromoteTo

	// Move the pawn and change the pawn with promoted piece
	ctx.Board.Move(move.From, move.To)
	ctx.Board.Place(move.To, piece)

	// Check if the pawn captured a rook to forfeit the castle rights for enemy
	ctx.ForfeitCastlingRight(move)
}

func (e DefaultEngine) applyCastling(ctx *core.TurnContext, move core.Move) {
	// Move the king.
	ctx.Board.Move(move.From, move.To)

	// Move the rook. King-side: H -> F. Queen-side: A -> D.
	rookFrom, rookTo := move.CastlingRookPositions()
	ctx.Board.Move(rookFrom, rookTo)

	// Castling forfeits all castling rights.
	ctx.Sides[move.Piece.Color].KingPosition = move.To
	ctx.Sides[move.Piece.Color].ClearCastlingRights()
}

func (e DefaultEngine) applyEnPassant(ctx *core.TurnContext, move core.Move) {
	// Move the pawn.
	ctx.Board.Move(move.From, move.To)

	// Remove the captured pawn (sits behind the destination, not on it).
	ctx.Board.Clear(move.EnPassantCapturedPosition())
}
