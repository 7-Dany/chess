package piece

import (
	"slices"

	"github.com/7-Dany/chess/core"
)

// Bishop moves diagonally: four directions, any distance.
var BishopDirections = [4][2]int8{
	{1, 1},
	{1, -1},
	{-1, 1},
	{-1, -1},
}

type Bishop struct{}

func (Bishop) IsAttacking(color core.PieceColor, target core.Position, ctx core.BoardContext) bool {
	for _, direction := range BishopDirections {
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
			if square.IsOccupiedByAny(color, core.BISHOP) {
				return true
			}

			if square.Occupied {
				break
			}
		}
	}
	return false
}

func (Bishop) Attacks(from core.Position, ctx core.BoardContext) []core.Position {
	attacks := make([]core.Position, 0, 13)

	for _, direction := range BishopDirections {
		file, rank := from.File(), from.Rank()

		// slide
		for {
			nextFile, fok := file.Add(direction[0])
			nextRank, rok := rank.Add(direction[1])

			if !fok || !rok {
				break
			}

			file, rank = nextFile, nextRank
			position := core.NewPosition(file, rank)
			attacks = append(attacks, position)
			if ctx.Board[position].Occupied {
				break
			}
		}
	}

	return slices.Clip(attacks)
}

func (b Bishop) PseudoLegalMoves(from core.Position, ctx core.MoveContext) []core.Move {
	bishop := core.Piece{Type: core.BISHOP, Color: ctx.SideToMove}
	moves := make([]core.Move, 0, 13)

	for _, direction := range BishopDirections {
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
				Piece: bishop,
				From:  from,
				To:    position,
			}

			if square.Occupied {
				move.HasCapture = true
				move.Captured = square.Piece
				moves = append(moves, move)
				break
			}

			moves = append(moves, move)
		}
	}

	return slices.Clip(moves)
}
