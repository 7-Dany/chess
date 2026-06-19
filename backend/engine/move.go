package engine

import (
	"github.com/7-Dany/chess/core"
)

func (e *DefaultEngine) GetLegalMoves(position core.Position, ctx core.TurnContext) []core.Move {
	moves := e.GetPseudoLegalMoves(position, ctx)
	current := ctx.SideToMove
	enemy := current.Opponent()

	square := ctx.Board[position]
	isKing := square.Type() == core.KING
	staticKing := ctx.Sides[current].KingPosition

	slot := 0
	for _, move := range moves {
		snapshot := e.Apply(&ctx, move)

		// after Apply, KingPosition is only updated when the king itself moved.
		// for all other pieces, staticKing stays valid throughout the loop.
		var kingPosition core.Position
		if isKing {
			kingPosition = ctx.Sides[current].KingPosition
		} else {
			kingPosition = staticKing
		}

		// king is safe, write legal moves back into the same slice to avoid a second allocation
		if !e.IsSquareAttacked(kingPosition, enemy, ctx) {
			// update move in place, ignoring invalid moves.
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
