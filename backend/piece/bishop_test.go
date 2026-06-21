package piece

import (
	"testing"

	"github.com/7-Dany/chess/core"
	"github.com/7-Dany/chess/testutil"
)

// TestBishopIsAttacking verifies that Bishop.IsAttacking correctly reports
// whether a bishop of the given color attacks the target square.
//
// A bishop attacks along its four diagonals — up to the first blocker. The
// scan goes FROM the target outward: "is there a bishop of `color` on one of
// my four diagonal rays, with nothing between us?".
func TestBishopIsAttacking(t *testing.T) {
	bishop := Bishop{}

	// All four diagonals through E4. A bishop on any of these (with nothing
	// between it and E4) attacks E4.
	t.Run("a white bishop on any of the four diagonals through E4 attacks E4", func(t *testing.T) {
		// One square on each diagonal direction from E4:
		//   H7  up-right   (+file, +rank)
		//   A8  up-left    (-file, +rank)
		//   H1  down-right (+file, -rank)
		//   B1  down-left  (-file, -rank)
		// Each has equal |file-diff| and |rank-diff| to E4.
		diagonals := []core.Position{core.H7, core.A8, core.H1, core.B1}
		for _, from := range diagonals {
			var board core.Board
			board[from] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
			ctx := core.BoardContext{Board: &board}

			if !bishop.IsAttacking(core.WHITE, core.E4, ctx) {
				t.Errorf("white bishop on %v should attack E4 (clear diagonal)", from)
			}
		}
	})

	t.Run("a bishop adjacent to the target attacks (distance 1)", func(t *testing.T) {
		var board core.Board
		board[core.D3] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		if !bishop.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("bishop on D3 (adjacent diagonal) should attack E4")
		}
	})

	t.Run("a bishop at maximum diagonal distance attacks (distance 7)", func(t *testing.T) {
		// H8 to A1 is the longest possible diagonal — 7 squares apart.
		var board core.Board
		board[core.H8] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		if !bishop.IsAttacking(core.WHITE, core.A1, ctx) {
			t.Errorf("bishop on H8 should attack A1 (full-length diagonal)")
		}
	})

	// Squares NOT on any diagonal through E4: same file, same rank, or
	// off-diagonal (file-diff != rank-diff).
	t.Run("a bishop on a non-diagonal square does not attack E4", func(t *testing.T) {
		nonDiagonal := []core.Position{
			core.E7, // same file (file-diff 0, rank-diff 3)
			core.H4, // same rank (file-diff 3, rank-diff 0)
			core.F6, // off-diagonal (file-diff 1, rank-diff 2)
		}
		for _, from := range nonDiagonal {
			var board core.Board
			board[from] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
			ctx := core.BoardContext{Board: &board}

			if bishop.IsAttacking(core.WHITE, core.E4, ctx) {
				t.Errorf("bishop on %v should NOT attack E4 (not on a diagonal)", from)
			}
		}
	})

	// A blocker between the bishop and the target breaks the attack —
	// regardless of whose piece the blocker is.
	t.Run("a friendly piece between the bishop and the target blocks the attack", func(t *testing.T) {
		// Bishop on H7 → F5 (own pawn) → E4 (target). The pawn blocks.
		var board core.Board
		board[core.H7] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
		board[core.F5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		if bishop.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("friendly pawn on F5 should block the bishop on H7 from attacking E4")
		}
	})

	t.Run("an enemy piece between the bishop and the target blocks the attack", func(t *testing.T) {
		var board core.Board
		board[core.H7] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
		board[core.F5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		ctx := core.BoardContext{Board: &board}

		if bishop.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("enemy pawn on F5 should block the bishop on H7 from attacking E4")
		}
	})

	t.Run("a piece behind the target does not block (target is between bishop and piece)", func(t *testing.T) {
		// Bishop on H7 → E4 (target) → C2 (own pawn). The pawn is past the
		// target, so it doesn't block.
		var board core.Board
		board[core.H7] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
		board[core.C2] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		if !bishop.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("pawn on C2 is past E4 — should not block the bishop on H7")
		}
	})

	t.Run("with two bishops on the same diagonal, the closer one attacks and blocks the farther", func(t *testing.T) {
		// Bishops on H7 and F5 (both white). F5 is adjacent to E4 and attacks
		// it; H7 is blocked by F5. The scan stops at F5, so it reports the
		// attack (true).
		var board core.Board
		board[core.H7] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
		board[core.F5] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		if !bishop.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("the closer bishop on F5 should attack E4")
		}
	})

	t.Run("a bishop of the wrong color is ignored", func(t *testing.T) {
		var board core.Board
		board[core.H7] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.BLACK})
		ctx := core.BoardContext{Board: &board}

		// Black bishop on H7, asking "does WHITE attack E4?" → no.
		if bishop.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("black bishop should not count as a white attacker")
		}
		// Same bishop, asking "does BLACK attack E4?" → yes.
		if !bishop.IsAttacking(core.BLACK, core.E4, ctx) {
			t.Errorf("black bishop on H7 should attack E4 for black")
		}
	})

	// Only a bishop triggers the bishop attack — a queen/rook/etc on the
	// same diagonal must not falsely report a bishop attack.
	t.Run("a non-bishop piece on the diagonal does not trigger a bishop attack", func(t *testing.T) {
		nonBishops := []core.PieceType{core.QUEEN, core.ROOK, core.PAWN, core.KING, core.KNIGHT}
		for _, pt := range nonBishops {
			var board core.Board
			board[core.H7] = core.NewSquare(core.Piece{Type: pt, Color: core.WHITE})
			ctx := core.BoardContext{Board: &board}

			if bishop.IsAttacking(core.WHITE, core.E4, ctx) {
				t.Errorf("%v on H7 should not trigger a bishop attack", pt)
			}
		}
	})

	t.Run("an empty board reports no attack", func(t *testing.T) {
		var board core.Board
		ctx := core.BoardContext{Board: &board}

		if bishop.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("empty board should not report an attack")
		}
	})

	// The four corner-to-corner diagonals: each corner attacks the opposite.
	t.Run("corner-to-corner diagonals: each corner is attacked by the opposite corner", func(t *testing.T) {
		corners := []struct {
			target core.Position
			bishop core.Position
		}{
			{core.A1, core.H8},
			{core.H8, core.A1},
			{core.A8, core.H1},
			{core.H1, core.A8},
		}
		for _, c := range corners {
			var board core.Board
			board[c.bishop] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
			ctx := core.BoardContext{Board: &board}

			if !bishop.IsAttacking(core.WHITE, c.target, ctx) {
				t.Errorf("bishop on %v should attack %v (full-length diagonal)", c.bishop, c.target)
			}
		}
	})

	t.Run("a bishop sitting on the target square itself does not attack it", func(t *testing.T) {
		var board core.Board
		board[core.E4] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		if bishop.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("bishop on the target square should not attack itself")
		}
	})

	t.Run("among multiple bishops, any matching-color bishop on a clear diagonal attacks", func(t *testing.T) {
		// Bishops on A1, H1, A8. Only H1 is on a diagonal to E4 (down-right
		// diagonal: H1→G2→F3→E4).
		var board core.Board
		board[core.A1] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
		board[core.H1] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
		board[core.A8] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		if !bishop.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("bishop on H1 should attack E4 via the down-right diagonal")
		}
	})

	t.Run("multiple enemy bishops with all diagonals blocked do not attack", func(t *testing.T) {
		// Four white bishops on the corners of the diagonals through E4, but
		// each diagonal has a black pawn blocker just before E4.
		var board core.Board
		board[core.H7] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
		board[core.A7] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
		board[core.H1] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
		board[core.B1] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
		board[core.F5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK}) // blocks H7
		board[core.D5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK}) // blocks A7
		board[core.F3] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK}) // blocks H1
		board[core.D3] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK}) // blocks B1
		ctx := core.BoardContext{Board: &board}

		if bishop.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("all four diagonals are blocked — no bishop should attack E4")
		}
	})

	t.Run("mixed-color bishops: only the matching color counts", func(t *testing.T) {
		// Black bishop on H7 (blocked from E4's perspective by nothing — it
		// would attack if it were white). White bishop on A8 attacks E4.
		var board core.Board
		board[core.H7] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.BLACK})
		board[core.A8] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		if !bishop.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("white bishop on A8 should attack E4 even with a black bishop on H7")
		}
	})

	// A blocker touching the target still blocks; a blocker touching the
	// bishop also still blocks.
	t.Run("a blocker immediately adjacent to the target blocks the attack", func(t *testing.T) {
		// Bishop on H7 → F5 (rook) → E4 (target). F5 is adjacent to E4.
		var board core.Board
		board[core.H7] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
		board[core.F5] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		if bishop.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("rook on F5 (adjacent to target) should block the bishop")
		}
	})

	t.Run("a blocker immediately adjacent to the bishop blocks the attack", func(t *testing.T) {
		// Bishop on H7 → G6 (rook) → ... → E4. G6 is adjacent to H7.
		var board core.Board
		board[core.H7] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
		board[core.G6] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		if bishop.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("rook on G6 (adjacent to bishop) should block the bishop")
		}
	})
}

// TestBishopAttacks verifies that Bishop.Attacks returns every square a
// bishop threatens from the given position.
//
// Unlike IsAttacking (which scans from the target), Attacks scans from the
// source: "what squares does THIS bishop threaten?". The result includes
// squares occupied by friendly pieces (a piece "attacks" a square even if
// it can't move there) and stops at the first occupied square on each
// diagonal (including it).
func TestBishopAttacks(t *testing.T) {
	bishop := Bishop{}

	// On an empty board from D4: 4 diagonals, 13 squares total (NE=4, SE=3,
	// NW=3, SW=3 — the diagonals are uneven because D4 is off-center).
	t.Run("bishop on center D4 with an empty board threatens 13 squares along 4 diagonals", func(t *testing.T) {
		var board core.Board
		ctx := core.BoardContext{Board: &board}

		got := bishop.Attacks(make([]core.Position, 0, core.MAX_MOVES), core.D4, ctx)

		testutil.AssertPositionsMatch(t, got, []core.Position{
			// NE: E5 F6 G7 H8
			core.E5, core.F6, core.G7, core.H8,
			// SE: E3 F2 G1
			core.E3, core.F2, core.G1,
			// NW: C5 B6 A7
			core.C5, core.B6, core.A7,
			// SW: C3 B2 A1
			core.C3, core.B2, core.A1,
		})
	})

	// A corner bishop has only one diagonal — 7 squares to the opposite corner.
	t.Run("bishop on corner A1 threatens 7 squares along its single diagonal", func(t *testing.T) {
		var board core.Board
		ctx := core.BoardContext{Board: &board}

		got := bishop.Attacks(make([]core.Position, 0, core.MAX_MOVES), core.A1, ctx)

		testutil.AssertPositionsMatch(t, got, []core.Position{
			core.B2, core.C3, core.D4, core.E5, core.F6, core.G7, core.H8,
		})
	})

	t.Run("bishop on corner H1 threatens 7 squares along its single diagonal", func(t *testing.T) {
		var board core.Board
		ctx := core.BoardContext{Board: &board}

		got := bishop.Attacks(make([]core.Position, 0, core.MAX_MOVES), core.H1, ctx)

		testutil.AssertPositionsMatch(t, got, []core.Position{
			core.G2, core.F3, core.E4, core.D5, core.C6, core.B7, core.A8,
		})
	})

	t.Run("bishop on corner A8 threatens 7 squares along its single diagonal", func(t *testing.T) {
		var board core.Board
		ctx := core.BoardContext{Board: &board}

		got := bishop.Attacks(make([]core.Position, 0, core.MAX_MOVES), core.A8, ctx)

		testutil.AssertPositionsMatch(t, got, []core.Position{
			core.B7, core.C6, core.D5, core.E4, core.F3, core.G2, core.H1,
		})
	})

	t.Run("bishop on corner H8 threatens 7 squares along its single diagonal", func(t *testing.T) {
		var board core.Board
		ctx := core.BoardContext{Board: &board}

		got := bishop.Attacks(make([]core.Position, 0, core.MAX_MOVES), core.H8, ctx)

		testutil.AssertPositionsMatch(t, got, []core.Position{
			core.G7, core.F6, core.E5, core.D4, core.C3, core.B2, core.A1,
		})
	})

	// An edge bishop has two diagonals.
	t.Run("bishop on edge A4 threatens 7 squares along its two diagonals", func(t *testing.T) {
		var board core.Board
		ctx := core.BoardContext{Board: &board}

		got := bishop.Attacks(make([]core.Position, 0, core.MAX_MOVES), core.A4, ctx)

		testutil.AssertPositionsMatch(t, got, []core.Position{
			// NE: B5 C6 D7 E8
			core.B5, core.C6, core.D7, core.E8,
			// SE: B3 C2 D1
			core.B3, core.C2, core.D1,
		})
	})

	// A blocker stops the diagonal but is included in the attack set (the
	// bishop attacks the blocker's square even if it can't move there).
	t.Run("a blocker on the diagonal stops the scan but is included in the attacks", func(t *testing.T) {
		// Bishop on D4, pawn on F6 (NE diagonal). Attacks include E5 and F6
		// (the blocker) but nothing past F6.
		var board core.Board
		board[core.F6] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		got := bishop.Attacks(make([]core.Position, 0, core.MAX_MOVES), core.D4, ctx)

		testutil.AssertPositionsMatch(t, got, []core.Position{
			// NE stops at F6 (the blocker is included)
			core.E5, core.F6,
			// SE, NW, SW unchanged
			core.E3, core.F2, core.G1,
			core.C5, core.B6, core.A7,
			core.C3, core.B2, core.A1,
		})
	})

	t.Run("a bishop surrounded by pieces on all four adjacent diagonals threatens only those 4 squares", func(t *testing.T) {
		// Pawns on E5, E3, C5, C3 — one step away on each diagonal. The
		// bishop attacks each of those 4 squares (the blockers) and nothing
		// beyond.
		var board core.Board
		board[core.E5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.E3] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.C5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.C3] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		got := bishop.Attacks(make([]core.Position, 0, core.MAX_MOVES), core.D4, ctx)

		testutil.AssertPositionsMatch(t, got, []core.Position{core.E5, core.E3, core.C5, core.C3})
	})

	t.Run("a blocker on a corner bishop's diagonal stops it early", func(t *testing.T) {
		// Bishop on A1, pawn on C3. Attacks B2 and C3 (the blocker), nothing past.
		var board core.Board
		board[core.C3] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		got := bishop.Attacks(make([]core.Position, 0, core.MAX_MOVES), core.A1, ctx)

		testutil.AssertPositionsMatch(t, got, []core.Position{core.B2, core.C3})
	})

	t.Run("an enemy piece on the diagonal is included in the attacks ( Attacks doesn't filter by color)", func(t *testing.T) {
		// Bishop on D4, enemy pawn on F6. Attacks includes F6 (the enemy
		// square) and stops there — same as a friendly blocker. The Attacks
		// method reports every threatened square regardless of occupancy.
		var board core.Board
		board[core.F6] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		ctx := core.BoardContext{Board: &board}

		got := bishop.Attacks(make([]core.Position, 0, core.MAX_MOVES), core.D4, ctx)

		testutil.AssertPositionsMatch(t, got, []core.Position{
			core.E5, core.F6, // NE stops at the enemy
			core.E3, core.F2, core.G1,
			core.C5, core.B6, core.A7,
			core.C3, core.B2, core.A1,
		})
	})
}

// TestBishopPseudoLegalMoves verifies that Bishop.PseudoLegalMoves returns
// the correct set of moves — respecting blockers and captures, but NOT
// filtering for king safety (that's the engine's job).
//
// Key rules for bishop moves:
//   - A bishop slides along its four diagonals until it hits a piece.
//   - A square occupied by a friendly piece is excluded (can't capture own);
//     the slide stops before that square.
//   - A square occupied by an enemy piece is included as a capture; the
//     slide stops after that square (the enemy is taken, nothing beyond).
//   - Every bishop move is type NORMAL.
func TestBishopPseudoLegalMoves(t *testing.T) {
	bishop := Bishop{}

	// Helper: extract just the destination squares from a move list.
	destinations := func(moves []core.Move) []core.Position {
		tos := make([]core.Position, len(moves))
		for i, m := range moves {
			tos[i] = m.To
		}
		return tos
	}

	// All 13 destinations from D4 on an empty board (4 diagonals).
	d4Empty := []core.Position{
		core.E5, core.F6, core.G7, core.H8, // NE
		core.E3, core.F2, core.G1, // SE
		core.C5, core.B6, core.A7, // NW
		core.C3, core.B2, core.A1, // SW
	}

	t.Run("bishop on center D4 with an empty board has 13 moves along 4 diagonals", func(t *testing.T) {
		var board core.Board
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := bishop.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.D4, ctx)

		testutil.AssertPositionsMatch(t, destinations(moves), d4Empty)
		testutil.AssertMoveCount(t, moves, 13)
	})

	t.Run("a square occupied by an enemy piece is included as a capture and stops the slide", func(t *testing.T) {
		// Enemy pawn on F6 (NE diagonal). The move to F6 is a capture;
		// G7 and H8 beyond it are unreachable.
		var board core.Board
		board[core.F6] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := bishop.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.D4, ctx)

		testutil.AssertPositionsMatch(t, destinations(moves), []core.Position{
			// NE stops at F6 (capture included)
			core.E5, core.F6,
			// SE, NW, SW unchanged
			core.E3, core.F2, core.G1,
			core.C5, core.B6, core.A7,
			core.C3, core.B2, core.A1,
		})

		// The move to F6 must be flagged as a capture of the black pawn.
		for _, m := range moves {
			if m.To == core.F6 {
				if !m.HasCapture {
					t.Errorf("move to F6 should be a capture")
				}
				if m.Captured != (core.Piece{Type: core.PAWN, Color: core.BLACK}) {
					t.Errorf("move to F6: Captured = %v, want black pawn", m.Captured)
				}
				return
			}
		}
		t.Errorf("expected a move to F6 (the capture square)")
	})

	t.Run("captures carry the exact enemy piece type and color sitting on the destination", func(t *testing.T) {
		// Three enemy pieces on three different diagonals from D4: a queen
		// on F6 (NE), a rook on F3 (SE), a knight on B6 (NW). Each capture
		// move must carry the exact captured piece, not just "an enemy".
		var board core.Board
		board[core.F6] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.BLACK})
		board[core.F3] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
		board[core.B6] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.BLACK})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := bishop.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.D4, ctx)

		wantCaptures := map[core.Position]core.Piece{
			core.F6: {Type: core.QUEEN, Color: core.BLACK},
			core.F3: {Type: core.ROOK, Color: core.BLACK},
			core.B6: {Type: core.KNIGHT, Color: core.BLACK},
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

	t.Run("a square occupied by a friendly piece is excluded and stops the slide", func(t *testing.T) {
		// Own pawn on F6 (NE diagonal). F6 is NOT in the move list (can't
		// capture own); G7 and H8 beyond it are also unreachable.
		var board core.Board
		board[core.F6] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := bishop.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.D4, ctx)

		testutil.AssertPositionsMatch(t, destinations(moves), []core.Position{
			// NE stops BEFORE F6 (F6 excluded, nothing beyond)
			core.E5,
			// SE, NW, SW unchanged
			core.E3, core.F2, core.G1,
			core.C5, core.B6, core.A7,
			core.C3, core.B2, core.A1,
		})
	})

	t.Run("a friendly piece blocks the slide; an enemy behind it is unreachable", func(t *testing.T) {
		// Own pawn on E5 (one step NE), enemy pawn on F6 (two steps NE).
		// The slide stops before E5; F6 is behind the blocker, unreachable.
		var board core.Board
		board[core.E5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.F6] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := bishop.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.D4, ctx)

		// NE diagonal is empty (E5 excluded, F6 unreachable). Other diagonals unchanged.
		testutil.AssertPositionsMatch(t, destinations(moves), []core.Position{
			core.E3, core.F2, core.G1,
			core.C5, core.B6, core.A7,
			core.C3, core.B2, core.A1,
		})
	})

	t.Run("an enemy piece blocks the slide but is capturable; nothing beyond it is reachable", func(t *testing.T) {
		// Enemy pawn on F6 (NE), enemy rook on H8 (NE, behind F6). F6 is
		// captured; H8 is behind the capture, unreachable.
		var board core.Board
		board[core.F6] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		board[core.H8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := bishop.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.D4, ctx)

		testutil.AssertPositionsMatch(t, destinations(moves), []core.Position{
			// NE: E5 (quiet), F6 (capture). H8 unreachable.
			core.E5, core.F6,
			core.E3, core.F2, core.G1,
			core.C5, core.B6, core.A7,
			core.C3, core.B2, core.A1,
		})
	})

	t.Run("after capturing an enemy, a friendly piece behind it is unreachable", func(t *testing.T) {
		// Enemy pawn on F6 (capturable), own rook on H8 (behind F6). F6 is
		// captured; H8 is behind the capture, unreachable.
		var board core.Board
		board[core.F6] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		board[core.H8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := bishop.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.D4, ctx)

		testutil.AssertPositionsMatch(t, destinations(moves), []core.Position{
			core.E5, core.F6,
			core.E3, core.F2, core.G1,
			core.C5, core.B6, core.A7,
			core.C3, core.B2, core.A1,
		})
	})

	t.Run("a mix of friendly and enemy on all four diagonals yields only the captures", func(t *testing.T) {
		// One step out on each diagonal: E5 (enemy), E3 (own), C5 (enemy),
		// C3 (own). The bishop can capture E5 and C5; E3 and C3 block.
		var board core.Board
		board[core.E5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		board[core.E3] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.C5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		board[core.C3] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := bishop.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.D4, ctx)

		// Only the two enemy captures (E5, C5); the two own pieces block.
		testutil.AssertPositionsMatch(t, destinations(moves), []core.Position{core.E5, core.C5})
	})

	t.Run("all four diagonals blocked by own pieces yields no moves", func(t *testing.T) {
		var board core.Board
		board[core.E5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.E3] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.C5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.C3] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := bishop.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.D4, ctx)

		testutil.AssertNoMoves(t, moves)
	})

	t.Run("a black bishop treats white pieces as enemies (captures) and black as own", func(t *testing.T) {
		// Same board as the enemy-capture test, but the bishop is black. The
		// white pawn on F6 is now the enemy → capturable.
		var board core.Board
		board[core.F6] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.BLACK}

		moves := bishop.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.D4, ctx)

		testutil.AssertPositionsMatch(t, destinations(moves), []core.Position{
			core.E5, core.F6,
			core.E3, core.F2, core.G1,
			core.C5, core.B6, core.A7,
			core.C3, core.B2, core.A1,
		})

		// The move to F6 is a capture of the white pawn.
		for _, m := range moves {
			if m.To == core.F6 {
				if !m.HasCapture || m.Captured != (core.Piece{Type: core.PAWN, Color: core.WHITE}) {
					t.Errorf("black bishop capturing F6: HasCapture=%v Captured=%v, want white pawn", m.HasCapture, m.Captured)
				}
				return
			}
		}
		t.Errorf("expected a move to F6")
	})

	t.Run("a black bishop treats black pieces as own (excluded)", func(t *testing.T) {
		// Black pawn on F6 is own → excluded; slide stops before it.
		var board core.Board
		board[core.F6] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.BLACK}

		moves := bishop.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.D4, ctx)

		testutil.AssertPositionsMatch(t, destinations(moves), []core.Position{
			core.E5, // NE stops before F6
			core.E3, core.F2, core.G1,
			core.C5, core.B6, core.A7,
			core.C3, core.B2, core.A1,
		})
	})

	// Every bishop move is type NORMAL (bishops don't castle, en passant,
	// or promote). Each move carries the mover and its source square.
	t.Run("every generated move has type NORMAL and carries the mover and source square", func(t *testing.T) {
		var board core.Board
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := bishop.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.D4, ctx)

		mover := core.Piece{Type: core.BISHOP, Color: core.WHITE}
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
