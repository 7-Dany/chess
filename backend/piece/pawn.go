package piece

import (
	"github.com/7-Dany/chess/core"
)

// the pawn can attack diagonally, file (+1) or file (-1) -> 2 files change and one rank up
var PawnDirections = [2]int8{1, -1}

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

func (p Pawn) Attacks(attacks []core.Position, from core.Position, ctx core.BoardContext) []core.Position {
	color := ctx.Board[from].Color()
	step, _, _ := p.direction(color)

	// guard against last rank to make sure pawn will not go off board
	rank, ok := from.Rank().Add(step)
	if !ok {
		return attacks
	}

	// Pawns attack diagonally; guard file bounds for pawns on the A or H file.
	if rightFile, ok := from.File().Add(PawnDirections[0]); ok {
		attacks = append(attacks, core.NewPosition(rightFile, rank))
	}
	if leftFile, ok := from.File().Add(PawnDirections[1]); ok {
		attacks = append(attacks, core.NewPosition(leftFile, rank))
	}

	return attacks
}

func (p Pawn) PseudoLegalMoves(moves []core.Move, from core.Position, ctx core.MoveContext) []core.Move {
	// get pawn direction
	step, _, end := p.direction(ctx.SideToMove)

	// next rank with guard, pawn can't pass the last rank (raw)
	rank, ok := from.Rank().Add(step)
	if !ok {
		return nil
	}

	// rule 1: if the next rank is the end, it's a promotion
	if rank == end {
		return p.promote(moves, from, ctx)
	}

	// rule 2: handle push, at start singe or double, otherwise single
	moves = p.push(moves, from, ctx)

	// rule 3: capture on right file, or left file if occupied by enemy
	moves = p.captures(moves, from, ctx)

	return moves
}

// promote returns all promotion moves for a pawn whose next rank is the last rank,
// Each destination expands to four moves — Q, R, B, N.
func (p Pawn) promote(moves []core.Move, from core.Position, ctx core.MoveContext) []core.Move {
	pawn := core.Piece{Type: core.PAWN, Color: ctx.SideToMove}

	// rank is guaranteed valid: promote is only called when the dispatcher's
	// Add(step) succeeded and landed on the end rank.
	step, _, _ := p.direction(ctx.SideToMove)
	rank, _ := from.Rank().Add(step)

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
	if ctx.Board[forward].IsEmpty() {
		stamp(forward, false, core.Piece{})
	}

	// diagonal capture promotions — walk attack squares
	for _, direction := range PawnDirections {
		file, ok := from.File().Add(direction)
		if !ok {
			continue // pawn is on the A or H file, skip that side
		}

		to := core.NewPosition(file, rank)
		square := ctx.Board[to]
		if square.IsOccupiedBy(pawn.Color.Opponent()) {
			stamp(to, true, square.Piece())
		}
	}

	return moves
}

// push returns forward pushes only (single, and double from the start rank)
func (p Pawn) push(moves []core.Move, from core.Position, ctx core.MoveContext) []core.Move {
	step, start, _ := p.direction(ctx.SideToMove)
	pawn := core.Piece{Type: core.PAWN, Color: ctx.SideToMove}

	rank, ok := from.Rank().Add(step)
	if !ok {
		return moves
	}

	// single push: front square must be empty
	next := core.NewPosition(from.File(), rank)
	square := ctx.Board[next]
	if square.IsOccupied() {
		return moves
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
	if doubleSquare.IsOccupied() {
		return moves
	}

	moves = append(moves, core.Move{Piece: pawn, From: from, To: doublePosition, Type: core.NORMAL})

	return moves
}

// captures returns diagonal capture moves, including en passant.
func (p Pawn) captures(moves []core.Move, from core.Position, ctx core.MoveContext) []core.Move {
	pawn := core.Piece{Type: core.PAWN, Color: ctx.SideToMove}
	enemy := ctx.SideToMove.Opponent()
	step, _, _ := p.direction(ctx.SideToMove)

	// Guard: compute the forward rank once
	rank, ok := from.Rank().Add(step)
	if !ok {
		return moves
	}

	// Try right (+1) and left (-1) files
	for _, fileDelta := range [2]int8{1, -1} {
		file, ok := from.File().Add(fileDelta)
		if !ok {
			continue // pawn is on the A or H file, skip that side
		}

		attack := core.NewPosition(file, rank)

		// En passant: square is empty but still a valid capture
		if ctx.EnPassantTarget == attack {
			moves = append(moves, core.Move{
				Type:       core.EN_PASSANT,
				Piece:      pawn,
				From:       from,
				To:         attack,
				HasCapture: true,
				Captured:   core.Piece{Type: core.PAWN, Color: enemy},
			})
			continue
		}

		// Normal diagonal capture: square must be occupied by an enemy
		square := ctx.Board[attack]
		if square.IsOccupiedBy(enemy) {
			moves = append(moves, core.Move{
				Type:       core.NORMAL,
				Piece:      pawn,
				From:       from,
				To:         attack,
				HasCapture: true,
				Captured:   square.Piece(),
			})
		}
	}

	return moves
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
