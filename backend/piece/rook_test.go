package piece

import (
	"testing"

	"github.com/7-Dany/chess/core"
	"github.com/7-Dany/chess/testutil"
)

// TestRookIsAttacking verifies that Rook.IsAttacking correctly reports
// whether a rook of the given color attacks the target square.
//
// A rook attacks along its four orthogonal lines (up, down, left, right) —
// up to the first blocker. The scan goes FROM the target outward: "is there
// a rook of `color` on one of my four orthogonal rays, with nothing between
// us?".
func TestRookIsAttacking(t *testing.T) {
	rook := Rook{}

	// All four orthogonal directions through E4. A rook on any of these
	// (with nothing between it and E4) attacks E4.
	t.Run("a white rook on any of the four orthogonal lines through E4 attacks E4", func(t *testing.T) {
		// One square on each line from E4:
		//   E7  up    (same file, rank +3)
		//   E1  down  (same file, rank -3)
		//   A4  left  (file -4, same rank)
		//   H4  right (file +3, same rank)
		// Each shares either the file or the rank with E4.
		lines := []core.Position{core.E7, core.E1, core.A4, core.H4}
		for _, from := range lines {
			var board core.Board
			board[from] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			ctx := core.BoardContext{Board: &board}

			if !rook.IsAttacking(core.WHITE, core.E4, ctx) {
				t.Errorf("white rook on %v should attack E4 (clear line)", from)
			}
		}
	})

	t.Run("a rook adjacent to the target attacks (distance 1)", func(t *testing.T) {
		var board core.Board
		board[core.E5] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		if !rook.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("rook on E5 (adjacent) should attack E4")
		}
	})

	t.Run("a rook at maximum line distance attacks (distance 7)", func(t *testing.T) {
		// A1 to A8 is the longest possible orthogonal line — 7 squares apart.
		var board core.Board
		board[core.A8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		if !rook.IsAttacking(core.WHITE, core.A1, ctx) {
			t.Errorf("rook on A8 should attack A1 (full-length line)")
		}
	})

	// Squares NOT on any orthogonal line through E4: diagonal, or off-line.
	t.Run("a rook on a non-orthogonal square does not attack E4", func(t *testing.T) {
		nonOrthogonal := []core.Position{
			core.D5, // diagonal (file -1, rank +1)
			core.G6, // off-line (file +2, rank +2)
			core.F6, // off-line (file +1, rank +2)
		}
		for _, from := range nonOrthogonal {
			var board core.Board
			board[from] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			ctx := core.BoardContext{Board: &board}

			if rook.IsAttacking(core.WHITE, core.E4, ctx) {
				t.Errorf("rook on %v should NOT attack E4 (not on a line)", from)
			}
		}
	})

	// A blocker between the rook and the target breaks the attack —
	// regardless of whose piece the blocker is.
	t.Run("a friendly piece between the rook and the target blocks the attack", func(t *testing.T) {
		// Rook on E7 → E5 (own pawn) → E4 (target). The pawn blocks.
		var board core.Board
		board[core.E7] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
		board[core.E5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		if rook.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("friendly pawn on E5 should block the rook on E7 from attacking E4")
		}
	})

	t.Run("an enemy piece between the rook and the target blocks the attack", func(t *testing.T) {
		var board core.Board
		board[core.E7] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
		board[core.E5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		ctx := core.BoardContext{Board: &board}

		if rook.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("enemy pawn on E5 should block the rook on E7 from attacking E4")
		}
	})

	t.Run("a piece behind the target does not block (target is between rook and piece)", func(t *testing.T) {
		// Rook on E7 → E4 (target) → E2 (own pawn). The pawn is past the
		// target, so it doesn't block.
		var board core.Board
		board[core.E7] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
		board[core.E2] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		if !rook.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("pawn on E2 is past E4 — should not block the rook on E7")
		}
	})

	t.Run("with two rooks on the same line, the closer one attacks and blocks the farther", func(t *testing.T) {
		// Rooks on E7 and E5 (both white). E5 is adjacent to E4 and attacks
		// it; E7 is blocked by E5. The scan stops at E5, so it reports the
		// attack (true).
		var board core.Board
		board[core.E7] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
		board[core.E5] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		if !rook.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("the closer rook on E5 should attack E4")
		}
	})

	t.Run("a rook of the wrong color is ignored", func(t *testing.T) {
		var board core.Board
		board[core.E7] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
		ctx := core.BoardContext{Board: &board}

		// Black rook on E7, asking "does WHITE attack E4?" → no.
		if rook.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("black rook should not count as a white attacker")
		}
		// Same rook, asking "does BLACK attack E4?" → yes.
		if !rook.IsAttacking(core.BLACK, core.E4, ctx) {
			t.Errorf("black rook on E7 should attack E4 for black")
		}
	})

	// Only a rook triggers the rook attack — a queen/bishop/etc on the
	// same line must not falsely report a rook attack.
	t.Run("a non-rook piece on the line does not trigger a rook attack", func(t *testing.T) {
		nonRooks := []core.PieceType{core.QUEEN, core.BISHOP, core.KNIGHT, core.KING, core.PAWN}
		for _, pt := range nonRooks {
			var board core.Board
			board[core.E7] = core.NewSquare(core.Piece{Type: pt, Color: core.WHITE})
			ctx := core.BoardContext{Board: &board}

			if rook.IsAttacking(core.WHITE, core.E4, ctx) {
				t.Errorf("%v on E7 should not trigger a rook attack", pt)
			}
		}
	})

	t.Run("an empty board reports no attack", func(t *testing.T) {
		var board core.Board
		ctx := core.BoardContext{Board: &board}

		if rook.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("empty board should not report an attack")
		}
	})

	// A corner is attacked by a rook on the opposite end of its file and
	// the opposite end of its rank.
	t.Run("a corner is attacked along both its file and its rank", func(t *testing.T) {
		// Target A1: attacked by rook on A8 (same file) and H1 (same rank).
		pairs := []struct {
			target core.Position
			rook   core.Position
		}{
			{core.A1, core.A8}, // up the file
			{core.A1, core.H1}, // across the rank
			{core.H8, core.H1}, // down the file
			{core.H8, core.A8}, // across the rank
		}
		for _, p := range pairs {
			var board core.Board
			board[p.rook] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			ctx := core.BoardContext{Board: &board}

			if !rook.IsAttacking(core.WHITE, p.target, ctx) {
				t.Errorf("rook on %v should attack %v (full-length line)", p.rook, p.target)
			}
		}
	})

	t.Run("a rook sitting on the target square itself does not attack it", func(t *testing.T) {
		var board core.Board
		board[core.E4] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		if rook.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("rook on the target square should not attack itself")
		}
	})

	t.Run("among multiple rooks, any matching-color rook on a clear line attacks", func(t *testing.T) {
		// Rooks on A1, H1, A8. Only A1 and H1 are on E4's rank (rank 4? no —
		// A1/H1 are rank 1). A8 is on file A (not file E). So none of these
		// attack E4. Let me use rooks that share E4's file or rank:
		//   A4 (same rank, rank 4), E8 (same file, file E), H4 (same rank).
		var board core.Board
		board[core.A4] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
		board[core.E8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
		board[core.H4] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		if !rook.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("at least one rook (A4, E8, or H4) should attack E4")
		}
	})

	t.Run("multiple enemy rooks with all lines blocked do not attack", func(t *testing.T) {
		// Four white rooks on the ends of E4's lines, each with a black pawn
		// blocker just before E4.
		var board core.Board
		board[core.E8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE}) // up
		board[core.E1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE}) // down
		board[core.A4] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE}) // left
		board[core.H4] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE}) // right
		board[core.E5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK}) // blocks E8
		board[core.E3] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK}) // blocks E1
		board[core.D4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK}) // blocks A4
		board[core.F4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK}) // blocks H4
		ctx := core.BoardContext{Board: &board}

		if rook.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("all four lines are blocked — no rook should attack E4")
		}
	})

	t.Run("mixed-color rooks: only the matching color counts", func(t *testing.T) {
		// Black rook on E7 (would attack if white), white rook on A4 (attacks).
		var board core.Board
		board[core.E7] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
		board[core.A4] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		if !rook.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("white rook on A4 should attack E4 even with a black rook on E7")
		}
	})

	// A blocker touching the target still blocks; a blocker touching the
	// rook also still blocks. The blocker must be a non-rook piece (a rook
	// of the matching color would itself be an attacker, not a blocker).
	t.Run("a blocker immediately adjacent to the target blocks the attack", func(t *testing.T) {
		// Rook on E7 → E5 (own pawn) → E4 (target). E5 is adjacent to E4.
		var board core.Board
		board[core.E7] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
		board[core.E5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		if rook.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("pawn on E5 (adjacent to target) should block the rook on E7")
		}
	})

	t.Run("a blocker immediately adjacent to the rook blocks the attack", func(t *testing.T) {
		// Rook on E7 → E6 (own pawn) → ... → E4. E6 is adjacent to E7.
		var board core.Board
		board[core.E7] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
		board[core.E6] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		if rook.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("pawn on E6 (adjacent to attacking rook) should block it")
		}
	})
}

// TestRookAttacks verifies that Rook.Attacks returns every square a rook
// threatens from the given position.
//
// Unlike IsAttacking (which scans from the target), Attacks scans from the
// source: "what squares does THIS rook threaten?". The result includes
// squares occupied by friendly pieces (a piece "attacks" a square even if
// it can't move there) and stops at the first occupied square on each line
// (including it).
func TestRookAttacks(t *testing.T) {
	rook := Rook{}

	// On an empty board from D4: 4 lines, 14 squares total
	// (up=4, down=3, left=3, right=4 — the lines are uneven because D4 is
	// off-center).
	t.Run("rook on center D4 with an empty board threatens 14 squares along 4 lines", func(t *testing.T) {
		var board core.Board
		ctx := core.BoardContext{Board: &board}

		got := rook.Attacks(make([]core.Position, 0, core.MAX_MOVES), core.D4, ctx)

		testutil.AssertPositionsMatch(t, got, []core.Position{
			// up: D5 D6 D7 D8
			core.D5, core.D6, core.D7, core.D8,
			// down: D3 D2 D1
			core.D3, core.D2, core.D1,
			// left: C4 B4 A4
			core.C4, core.B4, core.A4,
			// right: E4 F4 G4 H4
			core.E4, core.F4, core.G4, core.H4,
		})
	})

	// A corner rook has two lines — 7 squares each (file + rank).
	t.Run("rook on corner A1 threatens 14 squares along its two lines", func(t *testing.T) {
		var board core.Board
		ctx := core.BoardContext{Board: &board}

		got := rook.Attacks(make([]core.Position, 0, core.MAX_MOVES), core.A1, ctx)

		testutil.AssertPositionsMatch(t, got, []core.Position{
			// up the file: A2 A3 A4 A5 A6 A7 A8
			core.A2, core.A3, core.A4, core.A5, core.A6, core.A7, core.A8,
			// across the rank: B1 C1 D1 E1 F1 G1 H1
			core.B1, core.C1, core.D1, core.E1, core.F1, core.G1, core.H1,
		})
	})

	t.Run("rook on corner H1 threatens 14 squares along its two lines", func(t *testing.T) {
		var board core.Board
		ctx := core.BoardContext{Board: &board}

		got := rook.Attacks(make([]core.Position, 0, core.MAX_MOVES), core.H1, ctx)

		testutil.AssertPositionsMatch(t, got, []core.Position{
			// up: H2 H3 H4 H5 H6 H7 H8
			core.H2, core.H3, core.H4, core.H5, core.H6, core.H7, core.H8,
			// left: G1 F1 E1 D1 C1 B1 A1
			core.G1, core.F1, core.E1, core.D1, core.C1, core.B1, core.A1,
		})
	})

	t.Run("rook on corner A8 threatens 14 squares along its two lines", func(t *testing.T) {
		var board core.Board
		ctx := core.BoardContext{Board: &board}

		got := rook.Attacks(make([]core.Position, 0, core.MAX_MOVES), core.A8, ctx)

		testutil.AssertPositionsMatch(t, got, []core.Position{
			// down: A7 A6 A5 A4 A3 A2 A1
			core.A7, core.A6, core.A5, core.A4, core.A3, core.A2, core.A1,
			// right: B8 C8 D8 E8 F8 G8 H8
			core.B8, core.C8, core.D8, core.E8, core.F8, core.G8, core.H8,
		})
	})

	t.Run("rook on corner H8 threatens 14 squares along its two lines", func(t *testing.T) {
		var board core.Board
		ctx := core.BoardContext{Board: &board}

		got := rook.Attacks(make([]core.Position, 0, core.MAX_MOVES), core.H8, ctx)

		testutil.AssertPositionsMatch(t, got, []core.Position{
			// down: H7 H6 H5 H4 H3 H2 H1
			core.H7, core.H6, core.H5, core.H4, core.H3, core.H2, core.H1,
			// left: G8 F8 E8 D8 C8 B8 A8
			core.G8, core.F8, core.E8, core.D8, core.C8, core.B8, core.A8,
		})
	})

	// An edge rook has three lines (one direction falls off the board).
	t.Run("rook on edge A4 threatens 14 squares along its three lines", func(t *testing.T) {
		var board core.Board
		ctx := core.BoardContext{Board: &board}

		got := rook.Attacks(make([]core.Position, 0, core.MAX_MOVES), core.A4, ctx)

		testutil.AssertPositionsMatch(t, got, []core.Position{
			// up: A5 A6 A7 A8
			core.A5, core.A6, core.A7, core.A8,
			// down: A3 A2 A1
			core.A3, core.A2, core.A1,
			// right: B4 C4 D4 E4 F4 G4 H4 (no left — A is the edge file)
			core.B4, core.C4, core.D4, core.E4, core.F4, core.G4, core.H4,
		})
	})

	// A blocker stops the line but is included in the attack set (the
	// rook attacks the blocker's square even if it can't move there).
	t.Run("a friendly blocker on the line stops the scan but is included in the attacks", func(t *testing.T) {
		// Rook on D4, own pawn on D6 (up line). Attacks include D5 and D6
		// (the blocker) but nothing past D6.
		var board core.Board
		board[core.D6] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		got := rook.Attacks(make([]core.Position, 0, core.MAX_MOVES), core.D4, ctx)

		testutil.AssertPositionsMatch(t, got, []core.Position{
			// up stops at D6 (the blocker is included)
			core.D5, core.D6,
			// down, left, right unchanged
			core.D3, core.D2, core.D1,
			core.C4, core.B4, core.A4,
			core.E4, core.F4, core.G4, core.H4,
		})
	})

	t.Run("an enemy blocker on the line stops the scan but is included in the attacks", func(t *testing.T) {
		// Same as above but the blocker is an enemy pawn. Attacks doesn't
		// filter by color — it reports every threatened square.
		var board core.Board
		board[core.D6] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		ctx := core.BoardContext{Board: &board}

		got := rook.Attacks(make([]core.Position, 0, core.MAX_MOVES), core.D4, ctx)

		testutil.AssertPositionsMatch(t, got, []core.Position{
			core.D5, core.D6,
			core.D3, core.D2, core.D1,
			core.C4, core.B4, core.A4,
			core.E4, core.F4, core.G4, core.H4,
		})
	})

	t.Run("a rook surrounded by pieces on all four adjacent squares threatens only those 4 squares", func(t *testing.T) {
		// Pawns on D5, D3, C4, E4 — one step away on each line.
		var board core.Board
		board[core.D5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.D3] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.C4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.E4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		got := rook.Attacks(make([]core.Position, 0, core.MAX_MOVES), core.D4, ctx)

		testutil.AssertPositionsMatch(t, got, []core.Position{core.D5, core.D3, core.C4, core.E4})
	})

	t.Run("a blocker on a corner rook's line stops it early", func(t *testing.T) {
		// Rook on A1, pawn on A3 (up line). Attacks A2 and A3, nothing past.
		var board core.Board
		board[core.A3] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		got := rook.Attacks(make([]core.Position, 0, core.MAX_MOVES), core.A1, ctx)

		// Up line: A2, A3 (blocker). Right line unchanged.
		testutil.AssertPositionsMatch(t, got, []core.Position{
			core.A2, core.A3,
			core.B1, core.C1, core.D1, core.E1, core.F1, core.G1, core.H1,
		})
	})
}

// TestRookPseudoLegalMoves verifies that Rook.PseudoLegalMoves returns the
// correct set of moves — respecting blockers and captures, but NOT
// filtering for king safety (that's the engine's job).
//
// Key rules for rook moves:
//   - A rook slides along its four orthogonal lines until it hits a piece.
//   - A square occupied by a friendly piece is excluded (can't capture own);
//     the slide stops before that square.
//   - A square occupied by an enemy piece is included as a capture; the
//     slide stops after that square (the enemy is taken, nothing beyond).
//   - Every rook move is type NORMAL.
func TestRookPseudoLegalMoves(t *testing.T) {
	rook := Rook{}

	// Helper: extract just the destination squares from a move list.
	destinations := func(moves []core.Move) []core.Position {
		tos := make([]core.Position, len(moves))
		for i, m := range moves {
			tos[i] = m.To
		}
		return tos
	}

	// All 14 destinations from D4 on an empty board (4 lines).
	d4Empty := []core.Position{
		core.D5, core.D6, core.D7, core.D8, // up
		core.D3, core.D2, core.D1, // down
		core.C4, core.B4, core.A4, // left
		core.E4, core.F4, core.G4, core.H4, // right
	}

	t.Run("rook on center D4 with an empty board has 14 moves along 4 lines", func(t *testing.T) {
		var board core.Board
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := rook.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.D4, ctx)

		testutil.AssertPositionsMatch(t, destinations(moves), d4Empty)
		testutil.AssertMoveCount(t, moves, 14)
	})

	t.Run("a square occupied by an enemy piece is included as a capture and stops the slide", func(t *testing.T) {
		// Enemy pawn on D6 (up line). The move to D6 is a capture;
		// D7 and D8 beyond it are unreachable.
		var board core.Board
		board[core.D6] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := rook.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.D4, ctx)

		testutil.AssertPositionsMatch(t, destinations(moves), []core.Position{
			// up stops at D6 (capture included)
			core.D5, core.D6,
			// down, left, right unchanged
			core.D3, core.D2, core.D1,
			core.C4, core.B4, core.A4,
			core.E4, core.F4, core.G4, core.H4,
		})

		// The move to D6 must be flagged as a capture of the black pawn.
		for _, m := range moves {
			if m.To == core.D6 {
				if !m.HasCapture {
					t.Errorf("move to D6 should be a capture")
				}
				if m.Captured != (core.Piece{Type: core.PAWN, Color: core.BLACK}) {
					t.Errorf("move to D6: Captured = %v, want black pawn", m.Captured)
				}
				return
			}
		}
		t.Errorf("expected a move to D6 (the capture square)")
	})

	t.Run("captures carry the exact enemy piece type and color sitting on the destination", func(t *testing.T) {
		// Three enemy pieces on three different lines from D4: a queen on D6
		// (up), a rook on D1 (down), a knight on B4 (left). Each capture
		// move must carry the exact captured piece, not just "an enemy".
		var board core.Board
		board[core.D6] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.BLACK})
		board[core.D1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
		board[core.B4] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.BLACK})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := rook.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.D4, ctx)

		wantCaptures := map[core.Position]core.Piece{
			core.D6: {Type: core.QUEEN, Color: core.BLACK},
			core.D1: {Type: core.ROOK, Color: core.BLACK},
			core.B4: {Type: core.KNIGHT, Color: core.BLACK},
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
		// Own pawn on D6 (up line). D6 is NOT in the move list (can't
		// capture own); D7 and D8 beyond it are also unreachable.
		var board core.Board
		board[core.D6] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := rook.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.D4, ctx)

		testutil.AssertPositionsMatch(t, destinations(moves), []core.Position{
			// up stops BEFORE D6 (D6 excluded, nothing beyond)
			core.D5,
			// down, left, right unchanged
			core.D3, core.D2, core.D1,
			core.C4, core.B4, core.A4,
			core.E4, core.F4, core.G4, core.H4,
		})
	})

	t.Run("a friendly piece blocks the slide; an enemy behind it is unreachable", func(t *testing.T) {
		// Own pawn on D5 (one step up), enemy pawn on D6 (two steps up).
		// The slide stops before D5; D6 is behind the blocker, unreachable.
		var board core.Board
		board[core.D5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.D6] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := rook.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.D4, ctx)

		// Up line is empty (D5 excluded, D6 unreachable). Other lines unchanged.
		testutil.AssertPositionsMatch(t, destinations(moves), []core.Position{
			core.D3, core.D2, core.D1,
			core.C4, core.B4, core.A4,
			core.E4, core.F4, core.G4, core.H4,
		})
	})

	t.Run("an enemy piece blocks the slide but is capturable; nothing beyond it is reachable", func(t *testing.T) {
		// Enemy pawn on D6 (up), enemy rook on D8 (up, behind D6). D6 is
		// captured; D8 is behind the capture, unreachable.
		var board core.Board
		board[core.D6] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		board[core.D8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := rook.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.D4, ctx)

		testutil.AssertPositionsMatch(t, destinations(moves), []core.Position{
			// up: D5 (quiet), D6 (capture). D8 unreachable.
			core.D5, core.D6,
			core.D3, core.D2, core.D1,
			core.C4, core.B4, core.A4,
			core.E4, core.F4, core.G4, core.H4,
		})
	})

	t.Run("after capturing an enemy, a friendly piece behind it is unreachable", func(t *testing.T) {
		// Enemy pawn on D6 (capturable), own rook on D8 (behind D6). D6 is
		// captured; D8 is behind the capture, unreachable.
		var board core.Board
		board[core.D6] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		board[core.D8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := rook.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.D4, ctx)

		testutil.AssertPositionsMatch(t, destinations(moves), []core.Position{
			core.D5, core.D6,
			core.D3, core.D2, core.D1,
			core.C4, core.B4, core.A4,
			core.E4, core.F4, core.G4, core.H4,
		})
	})

	t.Run("a mix of friendly and enemy on all four lines yields only the captures", func(t *testing.T) {
		// One step out on each line: D5 (enemy), D3 (own), C4 (enemy), E4 (own).
		// The rook can capture D5 and C4; D3 and E4 block.
		var board core.Board
		board[core.D5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		board[core.D3] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.C4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		board[core.E4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := rook.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.D4, ctx)

		// Only the two enemy captures (D5, C4); the two own pieces block.
		testutil.AssertPositionsMatch(t, destinations(moves), []core.Position{core.D5, core.C4})
	})

	t.Run("all four lines blocked by own pieces yields no moves", func(t *testing.T) {
		var board core.Board
		board[core.D5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.D3] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.C4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.E4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := rook.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.D4, ctx)

		testutil.AssertNoMoves(t, moves)
	})

	t.Run("a black rook treats white pieces as enemies (captures) and black as own", func(t *testing.T) {
		// Same board as the enemy-capture test, but the rook is black. The
		// white pawn on D6 is now the enemy → capturable.
		var board core.Board
		board[core.D6] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.BLACK}

		moves := rook.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.D4, ctx)

		testutil.AssertPositionsMatch(t, destinations(moves), []core.Position{
			core.D5, core.D6,
			core.D3, core.D2, core.D1,
			core.C4, core.B4, core.A4,
			core.E4, core.F4, core.G4, core.H4,
		})

		// The move to D6 is a capture of the white pawn.
		for _, m := range moves {
			if m.To == core.D6 {
				if !m.HasCapture || m.Captured != (core.Piece{Type: core.PAWN, Color: core.WHITE}) {
					t.Errorf("black rook capturing D6: HasCapture=%v Captured=%v, want white pawn", m.HasCapture, m.Captured)
				}
				return
			}
		}
		t.Errorf("expected a move to D6")
	})

	t.Run("a black rook treats black pieces as own (excluded)", func(t *testing.T) {
		// Black pawn on D6 is own → excluded; slide stops before it.
		var board core.Board
		board[core.D6] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.BLACK}

		moves := rook.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.D4, ctx)

		testutil.AssertPositionsMatch(t, destinations(moves), []core.Position{
			core.D5, // up stops before D6
			core.D3, core.D2, core.D1,
			core.C4, core.B4, core.A4,
			core.E4, core.F4, core.G4, core.H4,
		})
	})

	// Every rook move is type NORMAL (rooks don't castle, en passant, or
	// promote — castling is added by the engine, not the piece). Each move
	// carries the mover and its source square.
	t.Run("every generated move has type NORMAL and carries the mover and source square", func(t *testing.T) {
		var board core.Board
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := rook.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.D4, ctx)

		mover := core.Piece{Type: core.ROOK, Color: core.WHITE}
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
