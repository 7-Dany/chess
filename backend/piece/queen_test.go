package piece

import (
	"testing"

	"github.com/7-Dany/chess/core"
	"github.com/7-Dany/chess/testutil"
)

// TestQueenIsAttacking verifies that Queen.IsAttacking correctly reports
// whether a queen of the given color attacks the target square.
//
// A queen attacks along all eight directions — four orthogonal (like a rook)
// and four diagonal (like a bishop) — up to the first blocker on each ray.
// The scan goes FROM the target outward.
func TestQueenIsAttacking(t *testing.T) {
	queen := Queen{}

	// All eight directions through E4. A queen on any of these (with nothing
	// between it and E4) attacks E4.
	t.Run("a white queen on any of the eight rays through E4 attacks E4", func(t *testing.T) {
		// One square on each ray from E4:
		//   E5  up        (same file, rank +1)
		//   E3  down      (same file, rank -1)
		//   D4  left      (file -1, same rank)
		//   F4  right     (file +1, same rank)
		//   F5  up-right  (file +1, rank +1)
		//   F3  down-right(file +1, rank -1)
		//   D5  up-left   (file -1, rank +1)
		//   D3  down-left (file -1, rank -1)
		rays := []core.Position{core.E5, core.E3, core.D4, core.F4, core.F5, core.F3, core.D5, core.D3}
		for _, from := range rays {
			var board core.Board
			board[from] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.WHITE})
			ctx := core.BoardContext{Board: &board}

			if !queen.IsAttacking(core.WHITE, core.E4, ctx) {
				t.Errorf("white queen on %v should attack E4 (clear ray)", from)
			}
		}
	})

	t.Run("a queen adjacent to the target attacks (distance 1)", func(t *testing.T) {
		var board core.Board
		board[core.E5] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		if !queen.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("queen on E5 (adjacent) should attack E4")
		}
	})

	t.Run("a queen at maximum ray distance attacks (distance 7)", func(t *testing.T) {
		// A1 to A8 is the longest possible orthogonal ray — 7 squares apart.
		var board core.Board
		board[core.A8] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		if !queen.IsAttacking(core.WHITE, core.A1, ctx) {
			t.Errorf("queen on A8 should attack A1 (full-length ray)")
		}
	})

	// A square not on any of the 8 rays through E4: a knight-L shape.
	t.Run("a queen on a non-ray square does not attack E4", func(t *testing.T) {
		// C3 is a knight-L from E4 (file -2, rank -1) — not orthogonal, not
		// diagonal.
		var board core.Board
		board[core.C3] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		if queen.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("queen on C3 should NOT attack E4 (knight-L, not a queen ray)")
		}
	})

	// A blocker between the queen and the target breaks the attack —
	// regardless of whose piece the blocker is.
	t.Run("a friendly piece between the queen and the target blocks the attack", func(t *testing.T) {
		// Queen on E8 → E6 (own pawn) → E4 (target). The pawn blocks.
		var board core.Board
		board[core.E8] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.WHITE})
		board[core.E6] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		if queen.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("friendly pawn on E6 should block the queen on E8 from attacking E4")
		}
	})

	t.Run("an enemy piece between the queen and the target blocks the attack", func(t *testing.T) {
		var board core.Board
		board[core.E8] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.WHITE})
		board[core.E6] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		ctx := core.BoardContext{Board: &board}

		if queen.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("enemy pawn on E6 should block the queen on E8 from attacking E4")
		}
	})

	t.Run("a piece behind the target does not block (target is between queen and piece)", func(t *testing.T) {
		// Queen on E8 → E4 (target) → E2 (own pawn). The pawn is past the
		// target, so it doesn't block.
		var board core.Board
		board[core.E8] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.WHITE})
		board[core.E2] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		if !queen.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("pawn on E2 is past E4 — should not block the queen on E8")
		}
	})

	t.Run("with two queens on the same ray, the closer one attacks and blocks the farther", func(t *testing.T) {
		// Queens on E8 and E6 (both white). E6 is adjacent to E5... wait,
		// target is E4. E6 → E5 → E4 (distance 2). E8 → E7 → E6 → ... → E4.
		// E6 attacks E4 (clear); E8 is blocked by E6. Scan stops at E6 → true.
		var board core.Board
		board[core.E8] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.WHITE})
		board[core.E6] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		if !queen.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("the closer queen on E6 should attack E4")
		}
	})

	t.Run("a queen of the wrong color is ignored", func(t *testing.T) {
		var board core.Board
		board[core.E8] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.BLACK})
		ctx := core.BoardContext{Board: &board}

		if queen.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("black queen should not count as a white attacker")
		}
		if !queen.IsAttacking(core.BLACK, core.E4, ctx) {
			t.Errorf("black queen on E8 should attack E4 for black")
		}
	})

	// Only a queen triggers the queen attack — a rook/bishop/etc on the
	// same ray must not falsely report a queen attack. (Important: a rook
	// on the same file and a bishop on the same diagonal both share a ray
	// with the queen, but neither should trigger queen.IsAttacking.)
	t.Run("a non-queen piece on the ray does not trigger a queen attack", func(t *testing.T) {
		nonQueens := []core.PieceType{core.ROOK, core.BISHOP, core.KNIGHT, core.KING, core.PAWN}
		for _, pt := range nonQueens {
			var board core.Board
			board[core.E8] = core.NewSquare(core.Piece{Type: pt, Color: core.WHITE})
			ctx := core.BoardContext{Board: &board}

			if queen.IsAttacking(core.WHITE, core.E4, ctx) {
				t.Errorf("%v on E8 should not trigger a queen attack", pt)
			}
		}
	})

	t.Run("an empty board reports no attack", func(t *testing.T) {
		var board core.Board
		ctx := core.BoardContext{Board: &board}

		if queen.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("empty board should not report an attack")
		}
	})

	// A corner is attacked by a queen on the opposite end of its file, its
	// rank, AND its diagonal — all three rays reach the opposite corner.
	t.Run("a corner is attacked along its file, its rank, and its diagonal", func(t *testing.T) {
		// Target A1: attacked by queen on A8 (file), H1 (rank), H8 (diagonal).
		attackers := []core.Position{core.A8, core.H1, core.H8}
		for _, from := range attackers {
			var board core.Board
			board[from] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.WHITE})
			ctx := core.BoardContext{Board: &board}

			if !queen.IsAttacking(core.WHITE, core.A1, ctx) {
				t.Errorf("queen on %v should attack A1 (full-length ray)", from)
			}
		}
	})

	t.Run("a queen sitting on the target square itself does not attack it", func(t *testing.T) {
		var board core.Board
		board[core.E4] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		if queen.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("queen on the target square should not attack itself")
		}
	})

	t.Run("among multiple queens, any matching-color queen on a clear ray attacks", func(t *testing.T) {
		// Queens on A4 (same rank), E8 (same file), H7 (diagonal). All three
		// attack E4 on a clear ray.
		var board core.Board
		board[core.A4] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.WHITE})
		board[core.E8] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.WHITE})
		board[core.H7] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		if !queen.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("at least one queen should attack E4")
		}
	})

	t.Run("multiple enemy queens with all rays blocked do not attack", func(t *testing.T) {
		// Four white queens on the ends of E4's orthogonal rays, each with a
		// black pawn blocker just before E4.
		var board core.Board
		board[core.E8] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.WHITE}) // up
		board[core.E1] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.WHITE}) // down
		board[core.A4] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.WHITE}) // left
		board[core.H4] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.WHITE}) // right
		board[core.E5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})  // blocks E8
		board[core.E3] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})  // blocks E1
		board[core.D4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})  // blocks A4
		board[core.F4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})  // blocks H4
		ctx := core.BoardContext{Board: &board}

		if queen.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("all four orthogonal rays are blocked — no queen should attack E4")
		}
	})

	t.Run("mixed-color queens: only the matching color counts", func(t *testing.T) {
		// Black queen on E8 (would attack if white), white queen on A4 (attacks).
		var board core.Board
		board[core.E8] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.BLACK})
		board[core.A4] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		if !queen.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("white queen on A4 should attack E4 even with a black queen on E8")
		}
	})

	// A blocker touching the target still blocks; a blocker touching the
	// queen also still blocks. The blocker must be a non-queen piece (a
	// queen of the matching color would itself be an attacker, not a blocker).
	t.Run("a blocker immediately adjacent to the target blocks the attack", func(t *testing.T) {
		// Queen on E8 → E5 (own pawn) → E4 (target). E5 is adjacent to E4.
		var board core.Board
		board[core.E8] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.WHITE})
		board[core.E5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		if queen.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("pawn on E5 (adjacent to target) should block the queen on E8")
		}
	})

	t.Run("a blocker immediately adjacent to the queen blocks the attack", func(t *testing.T) {
		// Queen on E8 → E7 (own pawn) → ... → E4. E7 is adjacent to E8.
		var board core.Board
		board[core.E8] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.WHITE})
		board[core.E7] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		if queen.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("pawn on E7 (adjacent to attacking queen) should block it")
		}
	})
}

// TestQueenAttacks verifies that Queen.Attacks returns every square a queen
// threatens from the given position.
//
// The queen has eight rays (4 orthogonal + 4 diagonal). Attacks includes
// squares occupied by friendly pieces and stops at the first occupied
// square on each ray (including it).
func TestQueenAttacks(t *testing.T) {
	queen := Queen{}

	// On an empty board from D4: 8 rays, 27 squares total.
	// (up=4, down=3, left=3, right=4, NE=4, SE=3, NW=3, SW=3)
	t.Run("queen on center D4 with an empty board threatens 27 squares along 8 rays", func(t *testing.T) {
		var board core.Board
		ctx := core.BoardContext{Board: &board}

		got := queen.Attacks(make([]core.Position, 0, MAX_MOVES), core.D4, ctx)

		testutil.AssertPositionsMatch(t, got, []core.Position{
			// up: D5 D6 D7 D8
			core.D5, core.D6, core.D7, core.D8,
			// down: D3 D2 D1
			core.D3, core.D2, core.D1,
			// left: C4 B4 A4
			core.C4, core.B4, core.A4,
			// right: E4 F4 G4 H4
			core.E4, core.F4, core.G4, core.H4,
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

	// A corner queen has three rays (up + right + diagonal) — 7 squares each.
	t.Run("queen on corner A1 threatens 21 squares along its three rays", func(t *testing.T) {
		var board core.Board
		ctx := core.BoardContext{Board: &board}

		got := queen.Attacks(make([]core.Position, 0, MAX_MOVES), core.A1, ctx)

		testutil.AssertPositionsMatch(t, got, []core.Position{
			// up the file: A2 A3 A4 A5 A6 A7 A8
			core.A2, core.A3, core.A4, core.A5, core.A6, core.A7, core.A8,
			// across the rank: B1 C1 D1 E1 F1 G1 H1
			core.B1, core.C1, core.D1, core.E1, core.F1, core.G1, core.H1,
			// up-right diagonal: B2 C3 D4 E5 F6 G7 H8
			core.B2, core.C3, core.D4, core.E5, core.F6, core.G7, core.H8,
		})
	})

	t.Run("queen on corner H1 threatens 21 squares along its three rays", func(t *testing.T) {
		var board core.Board
		ctx := core.BoardContext{Board: &board}

		got := queen.Attacks(make([]core.Position, 0, MAX_MOVES), core.H1, ctx)

		testutil.AssertPositionsMatch(t, got, []core.Position{
			// up: H2 H3 H4 H5 H6 H7 H8
			core.H2, core.H3, core.H4, core.H5, core.H6, core.H7, core.H8,
			// left: G1 F1 E1 D1 C1 B1 A1
			core.G1, core.F1, core.E1, core.D1, core.C1, core.B1, core.A1,
			// up-left diagonal: G2 F3 E4 D5 C6 B7 A8
			core.G2, core.F3, core.E4, core.D5, core.C6, core.B7, core.A8,
		})
	})

	t.Run("queen on corner A8 threatens 21 squares along its three rays", func(t *testing.T) {
		var board core.Board
		ctx := core.BoardContext{Board: &board}

		got := queen.Attacks(make([]core.Position, 0, MAX_MOVES), core.A8, ctx)

		testutil.AssertPositionsMatch(t, got, []core.Position{
			// down: A7 A6 A5 A4 A3 A2 A1
			core.A7, core.A6, core.A5, core.A4, core.A3, core.A2, core.A1,
			// right: B8 C8 D8 E8 F8 G8 H8
			core.B8, core.C8, core.D8, core.E8, core.F8, core.G8, core.H8,
			// down-right diagonal: B7 C6 D5 E4 F3 G2 H1
			core.B7, core.C6, core.D5, core.E4, core.F3, core.G2, core.H1,
		})
	})

	t.Run("queen on corner H8 threatens 21 squares along its three rays", func(t *testing.T) {
		var board core.Board
		ctx := core.BoardContext{Board: &board}

		got := queen.Attacks(make([]core.Position, 0, MAX_MOVES), core.H8, ctx)

		testutil.AssertPositionsMatch(t, got, []core.Position{
			// down: H7 H6 H5 H4 H3 H2 H1
			core.H7, core.H6, core.H5, core.H4, core.H3, core.H2, core.H1,
			// left: G8 F8 E8 D8 C8 B8 A8
			core.G8, core.F8, core.E8, core.D8, core.C8, core.B8, core.A8,
			// down-left diagonal: G7 F6 E5 D4 C3 B2 A1
			core.G7, core.F6, core.E5, core.D4, core.C3, core.B2, core.A1,
		})
	})

	// An edge queen has five rays (three directions fall off the board on
	// one side: for A4, the "left", "up-left", and "down-left" rays are gone).
	t.Run("queen on edge A4 threatens 21 squares along its five rays", func(t *testing.T) {
		var board core.Board
		ctx := core.BoardContext{Board: &board}

		got := queen.Attacks(make([]core.Position, 0, MAX_MOVES), core.A4, ctx)

		testutil.AssertPositionsMatch(t, got, []core.Position{
			// up: A5 A6 A7 A8
			core.A5, core.A6, core.A7, core.A8,
			// down: A3 A2 A1
			core.A3, core.A2, core.A1,
			// right: B4 C4 D4 E4 F4 G4 H4
			core.B4, core.C4, core.D4, core.E4, core.F4, core.G4, core.H4,
			// up-right: B5 C6 D7 E8
			core.B5, core.C6, core.D7, core.E8,
			// down-right: B3 C2 D1
			core.B3, core.C2, core.D1,
		})
	})

	// A blocker stops the ray but is included in the attack set.
	t.Run("a friendly blocker on the ray stops the scan but is included in the attacks", func(t *testing.T) {
		// Queen on D4, own pawn on D6 (up ray). Attacks includes D5 and D6
		// (the blocker) but nothing past D6.
		var board core.Board
		board[core.D6] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		got := queen.Attacks(make([]core.Position, 0, MAX_MOVES), core.D4, ctx)

		testutil.AssertPositionsMatch(t, got, []core.Position{
			// up stops at D6 (the blocker is included)
			core.D5, core.D6,
			// down, left, right, NE, SE, NW, SW unchanged
			core.D3, core.D2, core.D1,
			core.C4, core.B4, core.A4,
			core.E4, core.F4, core.G4, core.H4,
			core.E5, core.F6, core.G7, core.H8,
			core.E3, core.F2, core.G1,
			core.C5, core.B6, core.A7,
			core.C3, core.B2, core.A1,
		})
	})

	t.Run("an enemy blocker on the ray stops the scan but is included in the attacks", func(t *testing.T) {
		// Same as above but the blocker is an enemy pawn. Attacks doesn't
		// filter by color — it reports every threatened square.
		var board core.Board
		board[core.D6] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		ctx := core.BoardContext{Board: &board}

		got := queen.Attacks(make([]core.Position, 0, MAX_MOVES), core.D4, ctx)

		testutil.AssertPositionsMatch(t, got, []core.Position{
			core.D5, core.D6,
			core.D3, core.D2, core.D1,
			core.C4, core.B4, core.A4,
			core.E4, core.F4, core.G4, core.H4,
			core.E5, core.F6, core.G7, core.H8,
			core.E3, core.F2, core.G1,
			core.C5, core.B6, core.A7,
			core.C3, core.B2, core.A1,
		})
	})

	t.Run("a queen surrounded by pieces on all eight adjacent squares threatens only those 8 squares", func(t *testing.T) {
		// Pawns on D5, D3, C4, E4, E5, E3, C5, C3 — one step away on each ray.
		var board core.Board
		board[core.D5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.D3] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.C4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.E4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.E5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.E3] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.C5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.C3] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		got := queen.Attacks(make([]core.Position, 0, MAX_MOVES), core.D4, ctx)

		testutil.AssertPositionsMatch(t, got, []core.Position{
			core.D5, core.D3, core.C4, core.E4, core.E5, core.E3, core.C5, core.C3,
		})
	})

	t.Run("a blocker on a corner queen's ray stops it early", func(t *testing.T) {
		// Queen on A1, pawn on A3 (up ray). Attacks A2 and A3, nothing past
		// on the up ray; the right and diagonal rays are unchanged.
		var board core.Board
		board[core.A3] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		got := queen.Attacks(make([]core.Position, 0, MAX_MOVES), core.A1, ctx)

		testutil.AssertPositionsMatch(t, got, []core.Position{
			// up stops at A3: A2, A3
			core.A2, core.A3,
			// right unchanged: B1 C1 D1 E1 F1 G1 H1
			core.B1, core.C1, core.D1, core.E1, core.F1, core.G1, core.H1,
			// up-right diagonal unchanged: B2 C3 D4 E5 F6 G7 H8
			core.B2, core.C3, core.D4, core.E5, core.F6, core.G7, core.H8,
		})
	})
}

// TestQueenPseudoLegalMoves verifies that Queen.PseudoLegalMoves returns the
// correct set of moves — respecting blockers and captures, but NOT filtering
// for king safety (that's the engine's job).
//
// Key rules for queen moves:
//   - A queen slides along its eight rays until it hits a piece.
//   - A square occupied by a friendly piece is excluded (can't capture own);
//     the slide stops before that square.
//   - A square occupied by an enemy piece is included as a capture; the
//     slide stops after that square (the enemy is taken, nothing beyond).
//   - Every queen move is type NORMAL.
func TestQueenPseudoLegalMoves(t *testing.T) {
	queen := Queen{}

	// Helper: extract just the destination squares from a move list.
	destinations := func(moves []core.Move) []core.Position {
		tos := make([]core.Position, len(moves))
		for i, m := range moves {
			tos[i] = m.To
		}
		return tos
	}

	// All 27 destinations from D4 on an empty board (8 rays).
	d4Empty := []core.Position{
		core.D5, core.D6, core.D7, core.D8, // up
		core.D3, core.D2, core.D1, // down
		core.C4, core.B4, core.A4, // left
		core.E4, core.F4, core.G4, core.H4, // right
		core.E5, core.F6, core.G7, core.H8, // NE
		core.E3, core.F2, core.G1, // SE
		core.C5, core.B6, core.A7, // NW
		core.C3, core.B2, core.A1, // SW
	}

	t.Run("queen on center D4 with an empty board has 27 moves along 8 rays", func(t *testing.T) {
		var board core.Board
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := queen.PseudoLegalMoves(make([]core.Move, 0, MAX_MOVES), core.D4, ctx)

		testutil.AssertPositionsMatch(t, destinations(moves), d4Empty)
		testutil.AssertMoveCount(t, moves, 27)
	})

	t.Run("a square occupied by an enemy piece is included as a capture and stops the slide", func(t *testing.T) {
		// Enemy pawn on D6 (up ray). The move to D6 is a capture;
		// D7 and D8 beyond it are unreachable.
		var board core.Board
		board[core.D6] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := queen.PseudoLegalMoves(make([]core.Move, 0, MAX_MOVES), core.D4, ctx)

		// Replace D6 in d4Empty with just D5, D6 (up stops at D6).
		want := []core.Position{
			core.D5, core.D6, // up stops at D6 (capture included)
			core.D3, core.D2, core.D1,
			core.C4, core.B4, core.A4,
			core.E4, core.F4, core.G4, core.H4,
			core.E5, core.F6, core.G7, core.H8,
			core.E3, core.F2, core.G1,
			core.C5, core.B6, core.A7,
			core.C3, core.B2, core.A1,
		}
		testutil.AssertPositionsMatch(t, destinations(moves), want)

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
		// Three enemy pieces on three different rays from D4: a queen on D6
		// (up), a rook on D1 (down), a knight on B6 (NW diagonal). Each
		// capture move must carry the exact captured piece.
		var board core.Board
		board[core.D6] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.BLACK})
		board[core.D1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
		board[core.B6] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.BLACK})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := queen.PseudoLegalMoves(make([]core.Move, 0, MAX_MOVES), core.D4, ctx)

		wantCaptures := map[core.Position]core.Piece{
			core.D6: {Type: core.QUEEN, Color: core.BLACK},
			core.D1: {Type: core.ROOK, Color: core.BLACK},
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
		// Own pawn on D6 (up ray). D6 is NOT in the move list (can't capture
		// own); D7 and D8 beyond it are also unreachable.
		var board core.Board
		board[core.D6] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := queen.PseudoLegalMoves(make([]core.Move, 0, MAX_MOVES), core.D4, ctx)

		// Up ray stops BEFORE D6 (D6 excluded, nothing beyond).
		want := []core.Position{
			core.D5,
			core.D3, core.D2, core.D1,
			core.C4, core.B4, core.A4,
			core.E4, core.F4, core.G4, core.H4,
			core.E5, core.F6, core.G7, core.H8,
			core.E3, core.F2, core.G1,
			core.C5, core.B6, core.A7,
			core.C3, core.B2, core.A1,
		}
		testutil.AssertPositionsMatch(t, destinations(moves), want)
	})

	t.Run("a friendly piece blocks the slide; an enemy behind it is unreachable", func(t *testing.T) {
		// Own pawn on D5 (one step up), enemy pawn on D6 (two steps up).
		// The slide stops before D5; D6 is behind the blocker, unreachable.
		var board core.Board
		board[core.D5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.D6] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := queen.PseudoLegalMoves(make([]core.Move, 0, MAX_MOVES), core.D4, ctx)

		// Up ray is empty (D5 excluded, D6 unreachable). Other rays unchanged.
		want := []core.Position{
			core.D3, core.D2, core.D1,
			core.C4, core.B4, core.A4,
			core.E4, core.F4, core.G4, core.H4,
			core.E5, core.F6, core.G7, core.H8,
			core.E3, core.F2, core.G1,
			core.C5, core.B6, core.A7,
			core.C3, core.B2, core.A1,
		}
		testutil.AssertPositionsMatch(t, destinations(moves), want)
	})

	t.Run("an enemy piece blocks the slide but is capturable; nothing beyond it is reachable", func(t *testing.T) {
		// Enemy pawn on D6 (up), enemy rook on D8 (up, behind D6). D6 is
		// captured; D8 is behind the capture, unreachable.
		var board core.Board
		board[core.D6] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		board[core.D8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := queen.PseudoLegalMoves(make([]core.Move, 0, MAX_MOVES), core.D4, ctx)

		want := []core.Position{
			core.D5, core.D6, // up: D5 (quiet), D6 (capture). D8 unreachable.
			core.D3, core.D2, core.D1,
			core.C4, core.B4, core.A4,
			core.E4, core.F4, core.G4, core.H4,
			core.E5, core.F6, core.G7, core.H8,
			core.E3, core.F2, core.G1,
			core.C5, core.B6, core.A7,
			core.C3, core.B2, core.A1,
		}
		testutil.AssertPositionsMatch(t, destinations(moves), want)
	})

	t.Run("after capturing an enemy, a friendly piece behind it is unreachable", func(t *testing.T) {
		// Enemy pawn on D6 (capturable), own rook on D8 (behind D6). D6 is
		// captured; D8 is behind the capture, unreachable.
		var board core.Board
		board[core.D6] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		board[core.D8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := queen.PseudoLegalMoves(make([]core.Move, 0, MAX_MOVES), core.D4, ctx)

		want := []core.Position{
			core.D5, core.D6,
			core.D3, core.D2, core.D1,
			core.C4, core.B4, core.A4,
			core.E4, core.F4, core.G4, core.H4,
			core.E5, core.F6, core.G7, core.H8,
			core.E3, core.F2, core.G1,
			core.C5, core.B6, core.A7,
			core.C3, core.B2, core.A1,
		}
		testutil.AssertPositionsMatch(t, destinations(moves), want)
	})

	t.Run("a mix of friendly and enemy on all eight rays yields only the captures", func(t *testing.T) {
		// One step out on each of the 4 orthogonal + 4 diagonal rays, alternating
		// enemy/own. The queen captures the 4 enemy pieces; the 4 own pieces block.
		//   D5 (enemy, up), D3 (own, down)
		//   C4 (own, left), E4 (enemy, right)
		//   E5 (enemy, NE), E3 (own, SE)
		//   C5 (own, NW), C3 (enemy, SW)
		var board core.Board
		board[core.D5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		board[core.D3] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.C4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.E4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		board[core.E5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		board[core.E3] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.C5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.C3] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := queen.PseudoLegalMoves(make([]core.Move, 0, MAX_MOVES), core.D4, ctx)

		// Only the 4 enemy captures (D5, E4, E5, C3); the 4 own pieces block.
		testutil.AssertPositionsMatch(t, destinations(moves),
			[]core.Position{core.D5, core.E4, core.E5, core.C3})
	})

	t.Run("all eight rays blocked by own pieces yields no moves", func(t *testing.T) {
		var board core.Board
		board[core.D5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.D3] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.C4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.E4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.E5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.E3] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.C5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.C3] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := queen.PseudoLegalMoves(make([]core.Move, 0, MAX_MOVES), core.D4, ctx)

		testutil.AssertNoMoves(t, moves)
	})

	t.Run("a black queen treats white pieces as enemies (captures) and black as own", func(t *testing.T) {
		// Same board as the enemy-capture test, but the queen is black. The
		// white pawn on D6 is now the enemy → capturable.
		var board core.Board
		board[core.D6] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.BLACK}

		moves := queen.PseudoLegalMoves(make([]core.Move, 0, MAX_MOVES), core.D4, ctx)

		want := []core.Position{
			core.D5, core.D6,
			core.D3, core.D2, core.D1,
			core.C4, core.B4, core.A4,
			core.E4, core.F4, core.G4, core.H4,
			core.E5, core.F6, core.G7, core.H8,
			core.E3, core.F2, core.G1,
			core.C5, core.B6, core.A7,
			core.C3, core.B2, core.A1,
		}
		testutil.AssertPositionsMatch(t, destinations(moves), want)

		// The move to D6 is a capture of the white pawn.
		for _, m := range moves {
			if m.To == core.D6 {
				if !m.HasCapture || m.Captured != (core.Piece{Type: core.PAWN, Color: core.WHITE}) {
					t.Errorf("black queen capturing D6: HasCapture=%v Captured=%v, want white pawn", m.HasCapture, m.Captured)
				}
				return
			}
		}
		t.Errorf("expected a move to D6")
	})

	t.Run("a black queen treats black pieces as own (excluded)", func(t *testing.T) {
		// Black pawn on D6 is own → excluded; slide stops before it.
		var board core.Board
		board[core.D6] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.BLACK}

		moves := queen.PseudoLegalMoves(make([]core.Move, 0, MAX_MOVES), core.D4, ctx)

		want := []core.Position{
			core.D5, // up stops before D6
			core.D3, core.D2, core.D1,
			core.C4, core.B4, core.A4,
			core.E4, core.F4, core.G4, core.H4,
			core.E5, core.F6, core.G7, core.H8,
			core.E3, core.F2, core.G1,
			core.C5, core.B6, core.A7,
			core.C3, core.B2, core.A1,
		}
		testutil.AssertPositionsMatch(t, destinations(moves), want)
	})

	// Every queen move is type NORMAL (queens don't castle, en passant, or
	// promote). Each move carries the mover and its source square.
	t.Run("every generated move has type NORMAL and carries the mover and source square", func(t *testing.T) {
		var board core.Board
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := queen.PseudoLegalMoves(make([]core.Move, 0, MAX_MOVES), core.D4, ctx)

		mover := core.Piece{Type: core.QUEEN, Color: core.WHITE}
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
