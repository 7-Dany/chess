package piece

import (
	"slices"

	"github.com/7-Dany/chess/core"
)

var QueenDirections = [8][2]int8{
	// orthogonal (rook)
	{0, 1},
	{0, -1},
	{1, 0},
	{-1, 0},

	// diagonal (bishop)
	{1, 1},
	{1, -1},
	{-1, 1},
	{-1, -1},
}

type Queen struct{}

func (Queen) IsAttacking(color core.PieceColor, target core.Position, ctx core.BoardContext) bool {
	for _, direction := range QueenDirections {
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
			if square.IsOccupiedByAny(color, core.QUEEN) {
				return true
			}

			if square.Occupied {
				break
			}
		}
	}

	return false
}

func (Queen) Attacks(from core.Position, ctx core.BoardContext) []core.Position {
	attacks := make([]core.Position, 0, 27)

	for _, direction := range QueenDirections {
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
			if ctx.Board[position].Occupied {
				break
			}
		}
	}

	return slices.Clip(attacks)
}

func (q Queen) PseudoLegalMoves(from core.Position, ctx core.MoveContext) []core.Move {
	queen := core.Piece{Type: core.QUEEN, Color: ctx.SideToMove}
	moves := make([]core.Move, 0, 27)

	for _, direction := range QueenDirections {
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
				Piece: queen,
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
