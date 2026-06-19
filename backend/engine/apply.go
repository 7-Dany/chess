package engine

import "github.com/7-Dany/chess/core"

func (e *DefaultEngine) Apply(ctx *core.TurnContext, move core.Move) core.Snapshot {
	snapshot := ctx.Snapshot(move)

	piece := move.Piece
	if move.Type == core.PROMOTION {
		piece.Type = move.PromoteTo
	}

	// Move the piece to its destination, empty the origin.
	ctx.Board.Clear(move.From)
	ctx.Board.Place(move.To, piece)

	// Move-type specific board / rights mutations.
	switch move.Type {
	case core.NORMAL:
		e.applyNormal(ctx, move)
	case core.CASTLING:
		e.applyCastling(ctx, move)
	case core.EN_PASSANT:
		e.applyEnPassant(ctx, move)
	}

	// Capturing a rook on its home file removes the opponent's right.
	if move.HasCapture && move.Captured.Type == core.ROOK {
		clearCastlingRightByFile(ctx, move.Captured.Color, move.To.File())
	}

	// En passant target: set on double pawn push, cleared otherwise.
	if isDoublePawnPush(move) {
		ctx.EnPassantTarget = enPassantTarget(move)
	} else {
		ctx.EnPassantTarget = core.NoPosition
	}

	return snapshot
}

func (e *DefaultEngine) applyNormal(ctx *core.TurnContext, move core.Move) {
	switch move.Piece.Type {
	case core.KING:
		ctx.Sides[move.Piece.Color].KingPosition = move.To
		ctx.Sides[move.Piece.Color].CastlingRights = core.CastlingRights{}
	case core.ROOK:
		clearCastlingRightByFile(ctx, move.Piece.Color, move.From.File())
	}
}

func (e *DefaultEngine) applyCastling(ctx *core.TurnContext, move core.Move) {
	if move.To.File() > move.From.File() {
		// King-side: rook H -> F
		moveRook(ctx, move.From.Rank(), core.FILE_H, core.FILE_F)
	} else {
		// Queen-side: rook A -> D
		moveRook(ctx, move.From.Rank(), core.FILE_A, core.FILE_D)
	}
	ctx.Sides[move.Piece.Color].KingPosition = move.To
	ctx.Sides[move.Piece.Color].CastlingRights = core.CastlingRights{}
}

func (e *DefaultEngine) applyEnPassant(ctx *core.TurnContext, move core.Move) {
	capturedPawnPosition := core.NewPosition(move.To.File(), move.From.Rank())
	ctx.Board.Clear(capturedPawnPosition)
}
