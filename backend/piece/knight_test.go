package piece

import (
	"testing"

	"github.com/7-Dany/chess/core"
	"github.com/7-Dany/chess/testutil"
)

// TestKnightIsAttacking verifies that Knight.IsAttacking correctly reports
// whether a knight of the given color attacks the target square.
//
// IsAttacking scans FROM the target outward: it asks "is there a knight of
// `color` on any of the 8 L-shape squares around target?". This is the
// inverse of Attacks — useful for check detection where you know the king's
// square but not the attacker's.
func TestKnightIsAttacking(t *testing.T) {
	knight := Knight{}

	// The 8 L-shapes around E4: C3, C5, D2, D6, F2, F6, G3, G5.
	// A white knight on any of these attacks E4.
	t.Run("a white knight on any of the 8 L-shape squares around E4 attacks E4", func(t *testing.T) {
		lShapes := []core.Position{core.C3, core.C5, core.D2, core.D6, core.F2, core.F6, core.G3, core.G5}
		for _, from := range lShapes {
			var board core.Board
			board[from] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
			ctx := core.BoardContext{Board: &board}

			if !knight.IsAttacking(core.WHITE, core.E4, ctx) {
				t.Errorf("white knight on %v should attack E4", from)
			}
		}
	})

	// Squares that are NOT L-shapes: adjacent, same file/rank, diagonal.
	t.Run("a knight on a non-L-shape square does not attack E4", func(t *testing.T) {
		nonAttackers := []core.Position{
			core.D3, // adjacent
			core.E6, // same file
			core.H4, // same rank
			core.G6, // diagonal
			core.H7, // long diagonal
		}
		for _, from := range nonAttackers {
			var board core.Board
			board[from] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
			ctx := core.BoardContext{Board: &board}

			if knight.IsAttacking(core.WHITE, core.E4, ctx) {
				t.Errorf("knight on %v should NOT attack E4 (not an L-shape)", from)
			}
		}
	})

	t.Run("a knight of the wrong color is ignored", func(t *testing.T) {
		var board core.Board
		board[core.D6] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.BLACK})
		ctx := core.BoardContext{Board: &board}

		// Black knight on D6, but we ask "does WHITE attack E4?" → no.
		if knight.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("black knight should not count as a white attacker")
		}
		// Same knight, asking "does BLACK attack E4?" → yes.
		if !knight.IsAttacking(core.BLACK, core.E4, ctx) {
			t.Errorf("black knight on D6 should attack E4 for black")
		}
	})

	// Only a knight triggers the knight attack — a queen/rook/etc on the
	// same L-shape square must not falsely report a knight attack.
	t.Run("a non-knight piece on an L-shape square does not trigger a knight attack", func(t *testing.T) {
		nonKnights := []core.PieceType{core.QUEEN, core.ROOK, core.BISHOP, core.KING, core.PAWN}
		for _, pt := range nonKnights {
			var board core.Board
			board[core.D6] = core.NewSquare(core.Piece{Type: pt, Color: core.WHITE})
			ctx := core.BoardContext{Board: &board}

			if knight.IsAttacking(core.WHITE, core.E4, ctx) {
				t.Errorf("%v on D6 should not trigger a knight attack", pt)
			}
		}
	})

	t.Run("an empty board reports no attack", func(t *testing.T) {
		var board core.Board
		ctx := core.BoardContext{Board: &board}

		if knight.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("empty board should not report an attack")
		}
	})

	t.Run("a knight sitting on the target square itself does not attack it", func(t *testing.T) {
		var board core.Board
		board[core.E4] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		if knight.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("knight on the target square should not attack itself")
		}
	})

	// Corner targets have only 2 L-shape attackers; far corners have none.
	t.Run("corner A1 is attacked from B3 and C2 only", func(t *testing.T) {
		attackers := []core.Position{core.B3, core.C2}
		for _, from := range attackers {
			var board core.Board
			board[from] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
			ctx := core.BoardContext{Board: &board}

			if !knight.IsAttacking(core.WHITE, core.A1, ctx) {
				t.Errorf("knight on %v should attack A1", from)
			}
		}

		// A knight on the far corner H8 does not attack A1.
		var board core.Board
		board[core.H8] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}
		if knight.IsAttacking(core.WHITE, core.A1, ctx) {
			t.Errorf("knight on H8 should NOT attack A1")
		}
	})

	t.Run("corner H8 is attacked from F7 and G6", func(t *testing.T) {
		for _, from := range []core.Position{core.F7, core.G6} {
			var board core.Board
			board[from] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
			ctx := core.BoardContext{Board: &board}

			if !knight.IsAttacking(core.WHITE, core.H8, ctx) {
				t.Errorf("knight on %v should attack H8", from)
			}
		}
	})

	t.Run("corner A8 is attacked from B6", func(t *testing.T) {
		var board core.Board
		board[core.B6] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		if !knight.IsAttacking(core.WHITE, core.A8, ctx) {
			t.Errorf("knight on B6 should attack A8")
		}
	})

	t.Run("corner H1 is attacked from F2", func(t *testing.T) {
		var board core.Board
		board[core.F2] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		if !knight.IsAttacking(core.WHITE, core.H1, ctx) {
			t.Errorf("knight on F2 should attack H1")
		}
	})

	// Edge targets have 4 L-shape attackers.
	t.Run("edge target A4 is attacked from its 4 L-shape squares", func(t *testing.T) {
		for _, from := range []core.Position{core.B6, core.C5, core.B2, core.C3} {
			var board core.Board
			board[from] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
			ctx := core.BoardContext{Board: &board}

			if !knight.IsAttacking(core.WHITE, core.A4, ctx) {
				t.Errorf("knight on %v should attack A4", from)
			}
		}
	})

	t.Run("among multiple knights, any matching-color knight on an L-shape attacks", func(t *testing.T) {
		// D6 (white) attacks E4; A1 and H8 (white) don't.
		var board core.Board
		board[core.A1] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
		board[core.D6] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
		board[core.H8] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		if !knight.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("at least one white knight (D6) should attack E4")
		}
	})

	t.Run("multiple enemy-color knights do not count as attackers for us", func(t *testing.T) {
		var board core.Board
		board[core.D6] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.BLACK})
		board[core.F6] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.BLACK})
		ctx := core.BoardContext{Board: &board}

		if knight.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("black knights should not count as white attackers")
		}
	})

	t.Run("mixed-color knights: only the matching color counts", func(t *testing.T) {
		// D6 (black) doesn't count for white; F6 (white) does.
		var board core.Board
		board[core.D6] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.BLACK})
		board[core.F6] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		if !knight.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("white knight on F6 should attack E4 even with a black knight on D6")
		}
	})
}

// TestKnightAttacks verifies that Knight.Attacks returns every square a
// knight threatens from the given position.
//
// Unlike IsAttacking (which scans from the target), Attacks scans from the
// source: "what squares does THIS knight threaten?". The result includes
// squares occupied by friendly pieces (a piece "attacks" a square even if
// it can't move there) — that distinction matters for check detection.
func TestKnightAttacks(t *testing.T) {
	knight := Knight{}

	// On an empty board, Attacks returns the same squares as PseudoLegalMoves
	// (no blockers). The count depends on position: 8 from center, 4 from
	// edge, 2 from corner.

	t.Run("knight on center D4 threatens all 8 L-shape squares", func(t *testing.T) {
		var board core.Board
		ctx := core.BoardContext{Board: &board}

		got := knight.Attacks(make([]core.Position, 0, MAX_MOVES), core.D4, ctx)

		testutil.AssertPositionsMatch(t, got, []core.Position{
			core.E6, core.E2, core.C6, core.C2,
			core.F5, core.F3, core.B5, core.B3,
		})
	})

	t.Run("knight on corner A1 threatens 2 squares (B3 and C2)", func(t *testing.T) {
		var board core.Board
		ctx := core.BoardContext{Board: &board}

		got := knight.Attacks(make([]core.Position, 0, MAX_MOVES), core.A1, ctx)

		testutil.AssertPositionsMatch(t, got, []core.Position{core.B3, core.C2})
	})

	t.Run("knight on corner H1 threatens 2 squares (G3 and F2)", func(t *testing.T) {
		var board core.Board
		ctx := core.BoardContext{Board: &board}

		got := knight.Attacks(make([]core.Position, 0, MAX_MOVES), core.H1, ctx)

		testutil.AssertPositionsMatch(t, got, []core.Position{core.G3, core.F2})
	})

	t.Run("knight on corner A8 threatens 2 squares (B6 and C7)", func(t *testing.T) {
		var board core.Board
		ctx := core.BoardContext{Board: &board}

		got := knight.Attacks(make([]core.Position, 0, MAX_MOVES), core.A8, ctx)

		testutil.AssertPositionsMatch(t, got, []core.Position{core.B6, core.C7})
	})

	t.Run("knight on corner H8 threatens 2 squares (G6 and F7)", func(t *testing.T) {
		var board core.Board
		ctx := core.BoardContext{Board: &board}

		got := knight.Attacks(make([]core.Position, 0, MAX_MOVES), core.H8, ctx)

		testutil.AssertPositionsMatch(t, got, []core.Position{core.G6, core.F7})
	})

	// Edge squares have 4 L-shapes (some fall off the board).
	t.Run("knight on edge A4 threatens 4 squares", func(t *testing.T) {
		var board core.Board
		ctx := core.BoardContext{Board: &board}

		got := knight.Attacks(make([]core.Position, 0, MAX_MOVES), core.A4, ctx)

		testutil.AssertPositionsMatch(t, got, []core.Position{core.B6, core.B2, core.C5, core.C3})
	})

	t.Run("knight on edge H4 threatens 4 squares", func(t *testing.T) {
		var board core.Board
		ctx := core.BoardContext{Board: &board}

		got := knight.Attacks(make([]core.Position, 0, MAX_MOVES), core.H4, ctx)

		testutil.AssertPositionsMatch(t, got, []core.Position{core.G6, core.G2, core.F5, core.F3})
	})

	t.Run("knight on edge D1 threatens 4 squares", func(t *testing.T) {
		var board core.Board
		ctx := core.BoardContext{Board: &board}

		got := knight.Attacks(make([]core.Position, 0, MAX_MOVES), core.D1, ctx)

		testutil.AssertPositionsMatch(t, got, []core.Position{core.E3, core.C3, core.F2, core.B2})
	})

	t.Run("knight on edge D8 threatens 4 squares", func(t *testing.T) {
		var board core.Board
		ctx := core.BoardContext{Board: &board}

		got := knight.Attacks(make([]core.Position, 0, MAX_MOVES), core.D8, ctx)

		testutil.AssertPositionsMatch(t, got, []core.Position{core.E6, core.C6, core.F7, core.B7})
	})

	t.Run("knight on near-corner B2 threatens 4 squares", func(t *testing.T) {
		var board core.Board
		ctx := core.BoardContext{Board: &board}

		got := knight.Attacks(make([]core.Position, 0, MAX_MOVES), core.B2, ctx)

		testutil.AssertPositionsMatch(t, got, []core.Position{core.D3, core.D1, core.A4, core.C4})
	})
}

// TestKnightPseudoLegalMoves verifies that Knight.PseudoLegalMoves returns
// the correct set of moves — respecting blockers and captures, but NOT
// filtering for king safety (that's the engine's job).
//
// Key rules for knight moves:
//   - A knight jumps to its 8 L-shape squares (no sliding, no blocking).
//   - A square occupied by a friendly piece is excluded (can't capture own).
//   - A square occupied by an enemy piece is included as a capture.
//   - Every knight move is type NORMAL.
func TestKnightPseudoLegalMoves(t *testing.T) {
	knight := Knight{}

	// Helper: extract just the destination squares from a move list.
	destinations := func(moves []core.Move) []core.Position {
		tos := make([]core.Position, len(moves))
		for i, m := range moves {
			tos[i] = m.To
		}
		return tos
	}

	// All 8 L-shape destinations from D4 (center, empty board → 8 moves).
	d4Attacks := []core.Position{core.E6, core.E2, core.C6, core.C2, core.F5, core.F3, core.B5, core.B3}

	t.Run("knight on center D4 with an empty board has 8 moves", func(t *testing.T) {
		var board core.Board
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := knight.PseudoLegalMoves(make([]core.Move, 0, MAX_MOVES), core.D4, ctx)

		testutil.AssertPositionsMatch(t, destinations(moves), d4Attacks)
		testutil.AssertMoveCount(t, moves, 8)
	})

	t.Run("knight on corner A1 with an empty board has 2 moves", func(t *testing.T) {
		var board core.Board
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := knight.PseudoLegalMoves(make([]core.Move, 0, MAX_MOVES), core.A1, ctx)

		testutil.AssertPositionsMatch(t, destinations(moves), []core.Position{core.B3, core.C2})
	})

	t.Run("knight on corner H8 with an empty board has 2 moves", func(t *testing.T) {
		var board core.Board
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := knight.PseudoLegalMoves(make([]core.Move, 0, MAX_MOVES), core.H8, ctx)

		testutil.AssertPositionsMatch(t, destinations(moves), []core.Position{core.G6, core.F7})
	})

	t.Run("knight on edge A4 with an empty board has 4 moves", func(t *testing.T) {
		var board core.Board
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := knight.PseudoLegalMoves(make([]core.Move, 0, MAX_MOVES), core.A4, ctx)

		testutil.AssertPositionsMatch(t, destinations(moves), []core.Position{core.B6, core.B2, core.C5, core.C3})
	})

	t.Run("a square occupied by a friendly piece is excluded from the move list", func(t *testing.T) {
		var board core.Board
		board[core.E6] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE}) // own piece on E6
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := knight.PseudoLegalMoves(make([]core.Move, 0, MAX_MOVES), core.D4, ctx)

		// E6 is gone; the other 7 L-shapes remain.
		testutil.AssertPositionsMatch(t, destinations(moves),
			[]core.Position{core.E2, core.C6, core.C2, core.F5, core.F3, core.B5, core.B3})
	})

	t.Run("a square occupied by an enemy piece is included as a capture", func(t *testing.T) {
		var board core.Board
		board[core.E6] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.BLACK}) // enemy on E6
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := knight.PseudoLegalMoves(make([]core.Move, 0, MAX_MOVES), core.D4, ctx)

		// All 8 L-shapes present (E6 is a capture, not a block).
		testutil.AssertPositionsMatch(t, destinations(moves), d4Attacks)

		// The move to E6 must be flagged as a capture carrying the enemy knight.
		for _, m := range moves {
			if m.To == core.E6 {
				if !m.HasCapture {
					t.Errorf("move to E6 should be a capture")
				}
				if m.Captured != (core.Piece{Type: core.KNIGHT, Color: core.BLACK}) {
					t.Errorf("move to E6: Captured = %v, want black knight", m.Captured)
				}
				return
			}
		}
		t.Errorf("expected a move to E6 (the capture square)")
	})

	t.Run("captures carry the exact enemy piece sitting on the destination", func(t *testing.T) {
		var board core.Board
		board[core.E6] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.BLACK})
		board[core.F5] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
		board[core.C2] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.BLACK})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := knight.PseudoLegalMoves(make([]core.Move, 0, MAX_MOVES), core.D4, ctx)

		// All 8 L-shapes present (3 are captures, 5 are quiet).
		testutil.AssertPositionsMatch(t, destinations(moves), d4Attacks)

		// Verify each capture carries the right piece.
		wantCaptures := map[core.Position]core.Piece{
			core.E6: {Type: core.QUEEN, Color: core.BLACK},
			core.F5: {Type: core.ROOK, Color: core.BLACK},
			core.C2: {Type: core.BISHOP, Color: core.BLACK},
		}
		for _, m := range moves {
			if want, ok := wantCaptures[m.To]; ok {
				if !m.HasCapture {
					t.Errorf("move to %v should be a capture", m.To)
				} else if m.Captured != want {
					t.Errorf("move to %v: Captured = %v, want %v", m.To, m.Captured, want)
				}
			} else {
				if m.HasCapture {
					t.Errorf("move to %v should NOT be a capture", m.To)
				}
			}
		}
	})

	t.Run("a mix of friendly and enemy blockers excludes own, includes enemy", func(t *testing.T) {
		var board core.Board
		board[core.E6] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE}) // own → excluded
		board[core.C6] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE}) // own → excluded
		board[core.F5] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.BLACK}) // enemy → capture
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := knight.PseudoLegalMoves(make([]core.Move, 0, MAX_MOVES), core.D4, ctx)

		// E6 and C6 excluded (own); F5 included (enemy capture); rest quiet.
		testutil.AssertPositionsMatch(t, destinations(moves),
			[]core.Position{core.E2, core.C2, core.F5, core.F3, core.B5, core.B3})
	})

	t.Run("all 8 squares blocked by friendly pieces yields no moves", func(t *testing.T) {
		var board core.Board
		for _, pos := range d4Attacks {
			board[pos] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
		}
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := knight.PseudoLegalMoves(make([]core.Move, 0, MAX_MOVES), core.D4, ctx)

		testutil.AssertNoMoves(t, moves)
	})

	t.Run("a black knight treats white pieces as enemies (captures) and black as own", func(t *testing.T) {
		var board core.Board
		board[core.E6] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE}) // enemy → capture
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.BLACK}

		moves := knight.PseudoLegalMoves(make([]core.Move, 0, MAX_MOVES), core.D4, ctx)

		// E6 included (capture of white); all 8 present.
		testutil.AssertPositionsMatch(t, destinations(moves), d4Attacks)

		// The move to E6 is a capture of the white knight.
		for _, m := range moves {
			if m.To == core.E6 {
				if !m.HasCapture || m.Captured != (core.Piece{Type: core.KNIGHT, Color: core.WHITE}) {
					t.Errorf("black knight capturing E6: HasCapture=%v Captured=%v, want white knight", m.HasCapture, m.Captured)
				}
				return
			}
		}
		t.Errorf("expected a move to E6")
	})

	t.Run("a black knight treats black pieces as own (excluded)", func(t *testing.T) {
		var board core.Board
		board[core.E6] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.BLACK}) // own → excluded
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.BLACK}

		moves := knight.PseudoLegalMoves(make([]core.Move, 0, MAX_MOVES), core.D4, ctx)

		// E6 excluded; the other 7 L-shapes remain.
		testutil.AssertPositionsMatch(t, destinations(moves),
			[]core.Position{core.E2, core.C6, core.C2, core.F5, core.F3, core.B5, core.B3})
	})

	// Every knight move is type NORMAL (knights don't castle, en passant, or promote).
	t.Run("every generated move has type NORMAL and carries the mover and source square", func(t *testing.T) {
		var board core.Board
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := knight.PseudoLegalMoves(make([]core.Move, 0, MAX_MOVES), core.D4, ctx)

		mover := core.Piece{Type: core.KNIGHT, Color: core.WHITE}
		for _, m := range moves {
			if m.Type != core.NORMAL {
				t.Errorf("move to %v: Type = %v, want NORMAL", m.To, m.Type)
			}
			if m.From != core.D4 {
				t.Errorf("move to %v: From = %v, want D4", m.To, m.From)
			}
			if m.Piece != mover {
				t.Errorf("move to %v: Piece = %v, want %v", m.To, m.Piece, mover)
			}
		}
	})
}
