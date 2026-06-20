package engine

import "github.com/7-Dany/chess/core"

func (e *DefaultEngine) Apply(ctx *core.TurnContext, move core.Move) core.Snapshot {
	snapshot := ctx.Snapshot(move)

	switch move.Type {
	case core.NORMAL, core.PROMOTION:
		e.applyNormal(ctx, move)
	case core.CASTLING:
		e.applyCastling(ctx, move)
	case core.EN_PASSANT:
		e.applyEnPassant(ctx, move)
	}

	return snapshot
}

func (e *DefaultEngine) applyNormal(ctx *core.TurnContext, move core.Move) {
	// Move the piece. Promotion overwrites the destination with the promoted type.
	piece := move.Piece
	if move.Type == core.PROMOTION {
		piece.Type = move.PromoteTo
	}
	ctx.Board.Move(move.From, move.To)
	if move.Type == core.PROMOTION {
		ctx.Board.Place(move.To, piece)
	}

	// King moves update king position and forfeit all castling rights.
	if move.Piece.Type == core.KING {
		// KingPosition is updated for the game controller's benefit; the
		// legality checker derives king position from move.To directly.
		ctx.Sides[move.Piece.Color].KingPosition = move.To
		ctx.Sides[move.Piece.Color].ClearCastlingRights()
	}

	// Rook moves and rook captures each forfeit one castling right.
	ctx.ForfeitCastlingRight(move)

	// En passant target: set on double pawn push, cleared otherwise.
	ctx.SetEnPassantTarget(move)
}

func (e *DefaultEngine) applyCastling(ctx *core.TurnContext, move core.Move) {
	// Move the king.
	ctx.Board.Move(move.From, move.To)

	// Move the rook. King-side: H -> F. Queen-side: A -> D.
	rank := move.From.Rank()
	if move.To.File() > move.From.File() {
		ctx.Board.Move(core.NewPosition(core.FILE_H, rank), core.NewPosition(core.FILE_F, rank))
	} else {
		ctx.Board.Move(core.NewPosition(core.FILE_A, rank), core.NewPosition(core.FILE_D, rank))
	}

	// Castling forfeits all castling rights.
	ctx.Sides[move.Piece.Color].KingPosition = move.To
	ctx.Sides[move.Piece.Color].ClearCastlingRights()

	// En passant target: cleared (not a double pawn push).
	ctx.SetEnPassantTarget(move)
}

func (e *DefaultEngine) applyEnPassant(ctx *core.TurnContext, move core.Move) {
	// Move the pawn.
	ctx.Board.Move(move.From, move.To)

	// Remove the captured pawn (sits behind the destination, not on it).
	ctx.Board.Clear(core.NewPosition(move.To.File(), move.From.Rank()))

	// En passant target: cleared (not a double pawn push).
	ctx.SetEnPassantTarget(move)
}
