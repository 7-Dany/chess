package engine

import (
	"github.com/7-Dany/chess/core"
	"github.com/7-Dany/chess/piece"
)

func (e *DefaultEngine) IsSquareAttacked(target core.Position, color core.PieceColor, ctx core.TurnContext) bool {
	if e.GetPiece(core.KNIGHT).IsAttacking(color, target, ctx.BoardContext) {
		return true
	}
	if e.GetPiece(core.KING).IsAttacking(color, target, ctx.BoardContext) {
		return true
	}
	if e.GetPiece(core.PAWN).IsAttacking(color, target, ctx.BoardContext) {
		return true
	}

	// Inline slider scan (not Bishop/Rook/Queen.IsAttacking) to cover all three in 8 rays instead of 16.
	for _, direction := range piece.BishopDirections {
		file, rank := target.File(), target.Rank()
		for {
			nextFile, fok := file.Add(direction[0])
			nextRank, rok := rank.Add(direction[1])
			if !fok || !rok {
				break
			}

			file, rank = nextFile, nextRank

			square := ctx.Board[core.NewPosition(file, rank)]
			if square.IsOccupiedByAny(color, core.BISHOP, core.QUEEN) {
				return true
			}

			if square.IsOccupied() {
				break
			}
		}
	}

	// Orthogonal rays — rook or queen
	for _, direction := range piece.RookDirections {
		file, rank := target.File(), target.Rank()
		for {
			nextFile, fok := file.Add(direction[0])
			nextRank, rok := rank.Add(direction[1])
			if !fok || !rok {
				break
			}

			file, rank = nextFile, nextRank

			square := ctx.Board[core.NewPosition(file, rank)]
			if square.IsOccupiedByAny(color, core.ROOK, core.QUEEN) {
				return true
			}

			if square.IsOccupied() {
				break
			}
		}
	}

	return false
}
