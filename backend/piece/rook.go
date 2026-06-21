package piece

import (
	"github.com/7-Dany/chess/core"
)

// Rook moves orthogonally: four directions, any distance.
var RookDirections = [4][2]int8{
	{0, 1},
	{0, -1},
	{1, 0},
	{-1, 0},
}

type Rook struct{}

func (Rook) IsAttacking(color core.PieceColor, target core.Position, ctx core.BoardContext) bool {
	for _, direction := range RookDirections {
		file, rank := target.File(), target.Rank()

		// slide
		for {
			nextFile, fok := file.Add(direction[0])
			nextRank, rok := rank.Add(direction[1])

			if !fok || !rok {
				break
			}

			file, rank = nextFile, nextRank
			position := core.NewPosition(file, rank)
			square := ctx.Board[position]
			if square.IsOccupiedByAny(color, core.ROOK) {
				return true
			}

			if square.IsOccupied() {
				break
			}
		}
	}

	return false
}

func (Rook) Attacks(attacks []core.Position, from core.Position, ctx core.BoardContext) []core.Position {
	for _, direction := range RookDirections {
		file, rank := from.File(), from.Rank()

		for {
			nextFile, fok := file.Add(direction[0])
			nextRank, rok := rank.Add(direction[1])

			if !fok || !rok {
				break
			}

			file, rank = nextFile, nextRank
			position := core.NewPosition(file, rank)
			attacks = append(attacks, position)
			if ctx.Board[position].IsOccupied() {
				break
			}
		}
	}

	return attacks
}

func (r Rook) PseudoLegalMoves(moves []core.Move, from core.Position, ctx core.MoveContext) []core.Move {
	rook := core.Piece{Type: core.ROOK, Color: ctx.SideToMove}

	for _, direction := range RookDirections {
		file, rank := from.File(), from.Rank()

		for {
			nextFile, fok := file.Add(direction[0])
			nextRank, rok := rank.Add(direction[1])

			if !fok || !rok {
				break
			}

			file, rank = nextFile, nextRank
			position := core.NewPosition(file, rank)
			square := ctx.Board[position]

			if square.IsOccupiedBy(ctx.SideToMove) {
				break
			}

			move := core.Move{
				Type:  core.NORMAL,
				Piece: rook,
				From:  from,
				To:    position,
			}
			if square.IsOccupied() {
				move.HasCapture = true
				move.Captured = square.Piece()
				moves = append(moves, move)
				break
			}
			moves = append(moves, move)
		}
	}

	return moves
}
