package engine

import (
	"github.com/7-Dany/chess/core"
	"github.com/7-Dany/chess/piece"
)

func (e *DefaultEngine) IsSquareAttacked(target core.Position, color core.PieceColor, ctx core.BoardContext) bool {
	if e.pieces.Knight().IsAttacking(color, target, ctx) {
		return true
	}
	if e.pieces.King().IsAttacking(color, target, ctx) {
		return true
	}
	if e.pieces.Pawn().IsAttacking(color, target, ctx) {
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
