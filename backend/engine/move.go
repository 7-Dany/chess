package engine

import (
	"github.com/7-Dany/chess/core"
)

func (e DefaultEngine) GetLegalMoves(moves []core.Move, position core.Position, ctx core.TurnContext) []core.Move {
	moves = e.GetPseudoLegalMoves(moves, position, ctx)

	current := ctx.SideToMove
	enemy := current.Opponent()
	kingStart := ctx.Sides[current].KingPosition

	slot := 0
	for _, move := range moves {
		snapshot := e.Apply(&ctx, move)

		// After Apply, the king is at move.To if it moved, otherwise it's still where it started.
		kingPosition := kingStart
		if move.Piece.Type == core.KING {
			kingPosition = move.To
		}

		if !e.IsSquareAttacked(kingPosition, enemy, ctx.BoardContext) {
			moves[slot] = move
			slot++
		}

		e.Undo(&ctx, snapshot)
	}

	return moves[:slot]
}

func (e DefaultEngine) GetAllLegalMoves(moves []core.Move, ctx core.TurnContext) []core.Move {
	var scratch [core.MAX_MOVES]core.Move

	for i, square := range ctx.Board {
		if !square.IsOccupiedBy(ctx.SideToMove) {
			continue
		}

		pieceMoves := e.GetLegalMoves(scratch[:0], core.Position(i), ctx)
		moves = append(moves, pieceMoves...)
	}

	return moves
}

func (e DefaultEngine) HasAnyLegalMoves(ctx core.TurnContext) bool {
	// Stack allocated scratch buffer [moves], reused for every piece on the board.
	var moves [core.MAX_MOVES]core.Move

	current := ctx.SideToMove
	enemy := current.Opponent()
	kingStart := ctx.Sides[current].KingPosition

	for i, square := range ctx.Board {
		if !square.IsOccupiedBy(current) {
			continue
		}

		// if the square is occupied by enemy, we check all moves that piece can do
		// to make sure the king will be safe
		pseudoMoves := e.GetPseudoLegalMoves(moves[:0], core.Position(i), ctx)
		for _, move := range pseudoMoves {
			snapshot := e.Apply(&ctx, move)

			kingPosition := kingStart
			if move.Piece.Type == core.KING {
				kingPosition = move.To
			}

			legal := !e.IsSquareAttacked(kingPosition, enemy, ctx.BoardContext)

			e.Undo(&ctx, snapshot)

			if legal {
				return true
			}
		}
	}

	return false
}
