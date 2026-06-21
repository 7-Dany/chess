// Package testutil holds small helpers shared by the engine and piece test suites.
package testutil

import (
	"testing"

	"github.com/7-Dany/chess/core"
)

// DefaultSides returns both kings on home squares with full castling rights.
func DefaultSides() [2]core.SideState {
	return [2]core.SideState{FullWhite(), FullBlack()}
}

// FullWhite returns white's default side state: king on E1, full castling rights.
func FullWhite() core.SideState {
	return core.SideState{KingPosition: core.E1, CanCastleKingSide: true, CanCastleQueenSide: true}
}

// FullBlack returns black's default side state: king on E8, full castling rights.
func FullBlack() core.SideState {
	return core.SideState{KingPosition: core.E8, CanCastleKingSide: true, CanCastleQueenSide: true}
}

// Side builds a SideState with the given king position and castling rights.
func Side(kingPos core.Position, kingSide, queenSide bool) core.SideState {
	return core.SideState{KingPosition: kingPos, CanCastleKingSide: kingSide, CanCastleQueenSide: queenSide}
}

// TurnOption modifies a TurnContext during construction (see NewTurn).
type TurnOption func(*core.TurnContext)

// WithEnPassantTarget sets the en passant target square before a move is applied.
func WithEnPassantTarget(pos core.Position) TurnOption {
	return func(ctx *core.TurnContext) { ctx.EnPassantTarget = pos }
}

// WithSides replaces both side states entirely.
func WithSides(white, black core.SideState) TurnOption {
	return func(ctx *core.TurnContext) {
		ctx.Sides[0] = white
		ctx.Sides[1] = black
	}
}

// NewTurn builds a TurnContext for testing. King positions are auto-detected
// from the board; castling rights default to full. Pass options to override.
func NewTurn(board *core.Board, side core.PieceColor, options ...TurnOption) *core.TurnContext {
	ctx := &core.TurnContext{
		MoveContext: core.MoveContext{
			BoardContext: core.BoardContext{Board: board},
			SideToMove:   side,
			Sides:        DefaultSides(),
		},
	}
	for i := range 64 {
		sq := board[core.Position(i)]
		if sq.IsOccupied() && sq.Type() == core.KING {
			ctx.Sides[sq.Color()].KingPosition = core.Position(i)
		}
	}
	for _, option := range options {
		option(ctx)
	}
	return ctx
}

// AssertSquareHas checks that pos holds a piece of the given type and color.
func AssertSquareHas(t *testing.T, b *core.Board, pos core.Position, pt core.PieceType, color core.PieceColor) {
	t.Helper()
	got := b[pos]
	if !got.IsOccupied() || got.Type() != pt || got.Color() != color {
		t.Errorf("board[%v] = %v, want %v %v", pos, got, color, pt)
	}
}

// AssertSquareEmpty checks that pos holds no piece.
func AssertSquareEmpty(t *testing.T, b *core.Board, pos core.Position) {
	t.Helper()
	if b[pos].IsOccupied() {
		t.Errorf("board[%v] should be empty, got %v", pos, b[pos])
	}
}

// HasMove reports whether a move from→to exists in moves.
func HasMove(moves []core.Move, from, to core.Position) bool {
	for _, m := range moves {
		if m.From == from && m.To == to {
			return true
		}
	}
	return false
}

// AssertMovePresent checks that a move from→to exists in moves.
func AssertMovePresent(t *testing.T, moves []core.Move, from, to core.Position) {
	t.Helper()
	if !HasMove(moves, from, to) {
		t.Errorf("expected move %v→%v in %d moves", from, to, len(moves))
	}
}

// AssertMoveAbsent checks that no move from→to exists in moves.
func AssertMoveAbsent(t *testing.T, moves []core.Move, from, to core.Position) {
	t.Helper()
	for _, m := range moves {
		if m.From == from && m.To == to {
			t.Errorf("move %v→%v should not be present", from, to)
		}
	}
}

// AssertMoveCount checks that moves has exactly n entries.
func AssertMoveCount(t *testing.T, moves []core.Move, n int) {
	t.Helper()
	if len(moves) != n {
		t.Errorf("got %d moves, want %d", len(moves), n)
	}
}

// AssertNoMoves checks that moves is empty.
func AssertNoMoves(t *testing.T, moves []core.Move) {
	t.Helper()
	if len(moves) > 0 {
		t.Errorf("expected no moves, got %d: %v", len(moves), moves)
	}
}

// AssertPositionsMatch checks that got contains exactly the positions in want
// (order doesn't matter, but count must match).
func AssertPositionsMatch(t *testing.T, got, want []core.Position) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf("got %d positions, want %d\n  got:  %v\n  want: %v", len(got), len(want), got, want)
		return
	}
	wantSet := make(map[core.Position]struct{}, len(want))
	for _, p := range want {
		wantSet[p] = struct{}{}
	}
	for _, p := range got {
		if _, ok := wantSet[p]; !ok {
			t.Errorf("unexpected position %v (not in expected set %v)", p, want)
		}
	}
}
