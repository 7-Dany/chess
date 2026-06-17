package engine

import (
	"slices"

	"github.com/7-Dany/chess/core"
)

func (e *DefaultEngine) GetLegalMoves(position core.Position, ctx core.TurnContext) []core.Move {
	pseudoMoves := e.GetPseudoLegalMoves(position, ctx)

	current := ctx.SideToMove
	enemy := current.Opponent()
	moves := make([]core.Move, 0, len(pseudoMoves))

	for _, move := range pseudoMoves {
		snapshot := e.Apply(&ctx, move)
		kingPosition := ctx.Sides[current].KingPosition

		if !e.IsSquareAttacked(kingPosition, enemy, ctx) {
			moves = append(moves, move)
		}

		e.Undo(&ctx, snapshot)
	}

	return slices.Clip(moves)
}

func (e *DefaultEngine) HasAnyLegalMoves(ctx core.TurnContext) bool {
	current := ctx.SideToMove
	enemy := current.Opponent()

	for i, square := range ctx.Board {
		if !square.IsOccupiedBy(current) {
			continue
		}

		pseudoMoves := e.GetPseudoLegalMoves(core.Position(i), ctx)
		for _, move := range pseudoMoves {
			snapshot := e.Apply(&ctx, move)

			kingPosition := ctx.Sides[current].KingPosition
			legal := !e.IsSquareAttacked(kingPosition, enemy, ctx)

			e.Undo(&ctx, snapshot)

			if legal {
				return true
			}
		}
	}

	return false
}
