package piece

import (
	"slices"

	"github.com/7-Dany/chess/core"
)

// Knight L-shapes: all combinations of ±1 and ±2 offsets
var KnightDirections = [8][2]int8{
	{1, 2},
	{1, -2},
	{-1, 2},
	{-1, -2},
	{2, 1},
	{2, -1},
	{-2, 1},
	{-2, -1},
}

type Knight struct{}

func (Knight) IsAttacking(color core.PieceColor, target core.Position, ctx core.BoardContext) bool {
	for _, direction := range KnightDirections {
		nextFile, fok := target.File().Add(direction[0])
		nextRank, rok := target.Rank().Add(direction[1])

		if !fok || !rok {
			continue
		}

		position := core.NewPosition(nextFile, nextRank)
		square := ctx.Board[position]
		if square.IsOccupiedByAny(color, core.KNIGHT) {
			return true
		}
	}

	return false
}

func (Knight) Attacks(from core.Position, _ core.BoardContext) []core.Position {
	attacks := make([]core.Position, 0, 8)

	for _, direction := range KnightDirections {
		file, fok := from.File().Add(direction[0])
		rank, rok := from.Rank().Add(direction[1])
		if rok && fok {
			attacks = append(attacks, core.NewPosition(file, rank))
		}
	}

	return attacks
}

func (k Knight) PseudoLegalMoves(from core.Position, ctx core.MoveContext) []core.Move {
	knight := core.Piece{Type: core.KNIGHT, Color: ctx.SideToMove}
	moves := make([]core.Move, 0, 8)

	for _, direction := range KnightDirections {
		file, fok := from.File().Add(direction[0])
		rank, rok := from.Rank().Add(direction[1])
		if !fok || !rok {
			continue
		}

		position := core.NewPosition(file, rank)
		square := ctx.Board[position]

		// if the square is occupied by our side, ignore
		if square.IsOccupiedBy(ctx.SideToMove) {
			continue
		}

		move := core.Move{
			Type:  core.NORMAL,
			Piece: knight,
			From:  from,
			To:    position,
		}

		if square.IsOccupied() {
			move.HasCapture = true
			move.Captured = square.Piece()
		}

		moves = append(moves, move)
	}

	return slices.Clip(moves)
}
