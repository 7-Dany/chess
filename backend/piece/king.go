package piece

import (
	"slices"

	"github.com/7-Dany/chess/core"
)

var KingDirections = [8][2]int8{
	{0, 1},   // Up
	{0, -1},  // Down
	{1, 0},   // Right
	{-1, 0},  // Left
	{1, 1},   // Up Right
	{1, -1},  // Down Right
	{-1, 1},  // Up Left
	{-1, -1}, // Down Left
}

type King struct{}

func (King) IsAttacking(color core.PieceColor, target core.Position, ctx core.BoardContext) bool {
	for _, direction := range KingDirections {
		nextFile, fok := target.File().Add(direction[0])
		nextRank, rok := target.Rank().Add(direction[1])

		if !fok || !rok {
			continue
		}

		position := core.NewPosition(nextFile, nextRank)
		square := ctx.Board[position]
		if square.IsOccupiedByAny(color, core.KING) {
			return true
		}
	}
	return false
}

func (King) Attacks(from core.Position, ctx core.BoardContext) []core.Position {
	attacks := make([]core.Position, 0, 8)

	for _, position := range KingDirections {
		file, fok := from.File().Add(position[0])
		rank, rok := from.Rank().Add(position[1])
		if fok && rok {
			attacks = append(attacks, core.NewPosition(file, rank))
		}
	}

	return slices.Clip(attacks)
}

func (k King) PseudoLegalMoves(from core.Position, ctx core.MoveContext) []core.Move {
	king := core.Piece{Type: core.KING, Color: ctx.SideToMove}
	moves := make([]core.Move, 0, 10)

	for _, direction := range KingDirections {
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
			Piece: king,
			From:  from,
			To:    position,
		}

		if square.IsOccupied() {
			move.HasCapture = true
			move.Captured = square.Piece()
		}

		moves = append(moves, move)
	}

	return moves
}
