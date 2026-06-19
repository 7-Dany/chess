package engine

import "github.com/7-Dany/chess/core"

func isDoublePawnPush(move core.Move) bool {
	if move.Piece.Type != core.PAWN {
		return false
	}
	rankDiff := int(move.To.Rank()) - int(move.From.Rank())
	if move.Piece.Color == core.WHITE {
		return rankDiff == 2
	}
	return rankDiff == -2
}

func enPassantTarget(move core.Move) core.Position {
	step := int8(1)
	if move.Piece.Color == core.BLACK {
		step = -1
	}

	rank, _ := move.From.Rank().Add(step)
	return core.NewPosition(move.To.File(), rank)
}

// moveRook moves a rook along a rank: clears from, places it at to.
// Shared by apply (castling) and undo (un-castling).
func moveRook(ctx *core.TurnContext, rank core.Rank, from, to core.File) {
	fromPos := core.NewPosition(from, rank)
	toPos := core.NewPosition(to, rank)

	ctx.Board[toPos] = ctx.Board[fromPos]
	ctx.Board.Clear(fromPos)
}

