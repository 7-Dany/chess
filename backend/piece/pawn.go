package piece

import (
	"slices"

	"github.com/7-Dany/chess/core"
)

type Pawn struct{}

func (p Pawn) IsAttacking(color core.PieceColor, target core.Position, ctx core.BoardContext) bool {
	step, _, _ := p.direction(color)

	// checkig if position is being attacked,
	// expecting the pawn to be behind the target position (-step)
	rank, ok := target.Rank().Add(-step)
	if !ok {
		return false
	}

	// Pawns attack diagonally; guard file bounds for pawns on the A or H file.
	if rightFile, ok := target.File().Add(1); ok {
		square := ctx.Board[core.NewPosition(rightFile, rank)]
		if square.IsOccupiedByAny(color, core.PAWN) {
			return true
		}
	}

	if leftFile, ok := target.File().Add(-1); ok {
		square := ctx.Board[core.NewPosition(leftFile, rank)]
		if square.IsOccupiedByAny(color, core.PAWN) {
			return true
		}
	}

	return false
}

func (p Pawn) Attacks(from core.Position, ctx core.BoardContext) []core.Position {
	color := ctx.Board[from].Piece.Color
	step, _, _ := p.direction(color)

	return p.diagonals(from, step)
}

func (p Pawn) PseudoLegalMoves(from core.Position, ctx core.MoveContext) []core.Move {
	// get pawn direction
	step, _, end := p.direction(ctx.SideToMove)

	// next rank with guard, pawn can't pass the last rank (raw)
	rank, ok := from.Rank().Add(step)
	if !ok {
		return []core.Move{}
	}

	// rule 1: if the next rank is the end, it's a promotion
	if rank == end {
		return p.promote(from, ctx)
	}

	moves := make([]core.Move, 0, 4)

	// rule 2: handle push, at start singe or double, otherwise single
	moves = append(moves, p.push(from, ctx)...)

	// rule 3: capture on right file, or left file if occupied by enemy
	moves = append(moves, p.captures(from, ctx)...)

	return slices.Clip(moves)
}

// promote returns all promotion moves for a pawn whose next rank is the last rank,
// Each destination expands to four moves — Q, R, B, N.
func (p Pawn) promote(from core.Position, ctx core.MoveContext) []core.Move {
	pawn := core.Piece{Type: core.PAWN, Color: ctx.SideToMove}

	// rank is guaranteed valid: promote is only called when the dispatcher's
	// Add(step) succeeded and landed on the end rank.
	step, _, _ := p.direction(ctx.SideToMove)
	rank, _ := from.Rank().Add(step)

	moves := make([]core.Move, 0, 12)

	// stamp appends the four promotion variants (Q, R, B, N) for a single destination.
	stamp := func(to core.Position, hasCapture bool, captured core.Piece) {
		base := core.Move{
			Type:       core.PROMOTION,
			Piece:      pawn,
			From:       from,
			To:         to,
			HasCapture: hasCapture,
			Captured:   captured,
		}

		for _, promoteTo := range [4]core.PieceType{core.QUEEN, core.ROOK, core.BISHOP, core.KNIGHT} {
			move := base
			move.PromoteTo = promoteTo
			moves = append(moves, move)
		}
	}

	// forward push promotion (only if the square is empty)
	forward := core.NewPosition(from.File(), rank)
	if !ctx.Board[forward].Occupied {
		stamp(forward, false, core.Piece{})
	}

	// diagonal capture promotions — walk attack squares
	for _, to := range p.diagonals(from, step) {
		square := ctx.Board[to]
		if square.IsOccupiedBy(pawn.Color.Opponent()) {
			stamp(to, true, square.Piece)
		}
	}

	return slices.Clip(moves)
}

// push returns forward pushes only (single, and double from the start rank)
func (p Pawn) push(from core.Position, ctx core.MoveContext) []core.Move {
	step, start, _ := p.direction(ctx.SideToMove)
	pawn := core.Piece{Type: core.PAWN, Color: ctx.SideToMove}
	moves := make([]core.Move, 0, 2)

	rank, ok := from.Rank().Add(step)
	if !ok {
		return []core.Move{}
	}

	// single push: front square must be empty
	next := core.NewPosition(from.File(), rank)
	square := ctx.Board[next]
	if square.Occupied {
		return []core.Move{}
	}

	moves = append(moves, core.Move{Piece: pawn, From: from, To: next, Type: core.NORMAL})

	// double push only from the start rank
	if from.Rank() != start {
		return moves
	}

	// Safe to skip the ok check: start rank is always 2 steps from the board edge,
	// so two single-step adds from start always land in-bounds.
	doubleRank, _ := rank.Add(step)
	doublePosition := core.NewPosition(from.File(), doubleRank)
	doubleSquare := ctx.Board[doublePosition]
	if doubleSquare.Occupied {
		return moves
	}

	moves = append(moves, core.Move{Piece: pawn, From: from, To: doublePosition, Type: core.NORMAL})

	return moves
}

// captures returns diagonal capture moves, including en passant.
func (p Pawn) captures(from core.Position, ctx core.MoveContext) []core.Move {
	pawn := core.Piece{Type: core.PAWN, Color: ctx.SideToMove}
	enemy := ctx.SideToMove.Opponent()
	moves := make([]core.Move, 0, 2)

	step, _, _ := p.direction(ctx.SideToMove)
	for _, attack := range p.diagonals(from, step) {
		// en passant square is empty, but its capture according to pawn rules
		if ctx.EnPassantTarget == attack {
			moves = append(moves, core.Move{
				Type:       core.EN_PASSANT,
				Piece:      pawn,
				To:         attack,
				From:       from,
				HasCapture: true,
				Captured:   core.Piece{Type: core.PAWN, Color: enemy},
			})

			continue
		}

		// if the square occupied by enemy, its valid capture
		square := ctx.Board[attack]
		if square.IsOccupiedBy(enemy) {
			moves = append(moves, core.Move{
				Type:       core.NORMAL,
				Piece:      pawn,
				To:         attack,
				From:       from,
				HasCapture: true,
				Captured:   square.Piece,
			})
		}
	}

	return moves
}

// diagonals returns the up-to-two squares diagonally adjacent to from,
// one rank toward step, guarding both rank and file edges.
func (p Pawn) diagonals(from core.Position, step int8) []core.Position {
	squares := make([]core.Position, 0, 2)

	// A pawn on its last rank is impossible (promotion occurs on arrival),
	// but guard defensively for correctness.
	rank, ok := from.Rank().Add(step)
	if !ok {
		return squares
	}

	// Pawns attack diagonally; guard file bounds for pawns on the A or H file.
	if rightFile, ok := from.File().Add(1); ok {
		squares = append(squares, core.NewPosition(rightFile, rank))
	}
	if leftFile, ok := from.File().Add(-1); ok {
		squares = append(squares, core.NewPosition(leftFile, rank))
	}

	return squares
}

// direction returns the movement properties for a pawn of the given color.
// White pawns move up the board, black pawns move down.
//
// Returns:
//   - step:  rank increment per move (+1 for white, -1 for black)
//   - start: the rank pawns of this color begin on, used for double-push eligibility
//   - end:   the rank a pawn must reach to promote
func (p Pawn) direction(color core.PieceColor) (step int8, start, end core.Rank) {
	if color == core.WHITE {
		return 1, core.RANK_2, core.RANK_8
	}
	return -1, core.RANK_7, core.RANK_1
}
