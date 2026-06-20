package engine

import (
	"github.com/7-Dany/chess/core"
)

func (e *DefaultEngine) GetLegalMoves(position core.Position, ctx core.TurnContext) []core.Move {
	moves := e.GetPseudoLegalMoves(position, ctx)
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

func (e *DefaultEngine) HasAnyLegalMoves(ctx core.TurnContext) bool {
	current := ctx.SideToMove
	enemy := current.Opponent()
	kingStart := ctx.Sides[current].KingPosition

	for i, square := range ctx.Board {
		if !square.IsOccupiedBy(current) {
			continue
		}

		pseudoMoves := e.GetPseudoLegalMoves(core.Position(i), ctx)
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
