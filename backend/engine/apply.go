package engine

import "github.com/7-Dany/chess/core"

func (e *DefaultEngine) Apply(ctx *core.TurnContext, move core.Move) core.Snapshot {
	snapshot := ctx.Snapshot(move)

	piece := move.Piece
	if move.Type == core.PROMOTION {
		piece.Type = move.PromoteTo
	}

	// Move the piece to its destination, empty the origin.
	ctx.Board.Move(move.From, move.To)
	if move.Type == core.PROMOTION {
		ctx.Board.Place(move.To, piece) // overwrite with promoted piece
	}

	// Move-type specific board / rights mutations.
	switch move.Type {
	case core.NORMAL:
		e.applyNormal(ctx, move)
	case core.CASTLING:
		e.applyCastling(ctx, move)
	case core.EN_PASSANT:
		e.applyEnPassant(ctx, move)
	}

	// Capturing a rook on its home rank forfeits the castling right for that file.
	// Rank guard: a rook elsewhere is not the original then that right is already gone.
	if move.HasCapture && move.Captured.Type == core.ROOK &&
		move.To.Rank() == move.Captured.Color.KingStartRank() {
		ctx.Sides[move.Captured.Color].ClearCastlingRight(move.To.File())
	}

	// En passant target: set on double pawn push, cleared otherwise.
	if move.IsDoublePawnPush() {
		ctx.EnPassantTarget = move.EnPassantTarget()
	} else {
		ctx.EnPassantTarget = core.NoPosition
	}

	return snapshot
}

func (e *DefaultEngine) applyNormal(ctx *core.TurnContext, move core.Move) {
	switch move.Piece.Type {
	case core.KING:
		ctx.Sides[move.Piece.Color].KingPosition = move.To
		ctx.Sides[move.Piece.Color].ClearCastlingRights()
	case core.ROOK:
		// A rook leaving its home rank forfeits the right for that file.
		// Rank guard: a rook elsewhere is not the original then that right is already gone.
		if move.From.Rank() == move.Piece.Color.KingStartRank() {
			ctx.Sides[move.Piece.Color].ClearCastlingRight(move.From.File())
		}
	}
}

func (e *DefaultEngine) applyCastling(ctx *core.TurnContext, move core.Move) {
	rank := move.From.Rank()
	if move.To.File() > move.From.File() {
		// King-side: rook H -> F, Move Rook from H file to F file
		ctx.Board.Move(core.NewPosition(core.FILE_H, rank), core.NewPosition(core.FILE_F, rank))
	} else {
		// Queen-side: rook A -> D, Move Rook from A file to D file
		ctx.Board.Move(core.NewPosition(core.FILE_A, rank), core.NewPosition(core.FILE_D, rank))
	}
	ctx.Sides[move.Piece.Color].KingPosition = move.To
	ctx.Sides[move.Piece.Color].ClearCastlingRights()
}

func (e *DefaultEngine) applyEnPassant(ctx *core.TurnContext, move core.Move) {
	capturedPawnPosition := core.NewPosition(move.To.File(), move.From.Rank())
	ctx.Board.Clear(capturedPawnPosition)
}
