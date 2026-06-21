package piece

import (
	"testing"

	"github.com/7-Dany/chess/core"
	"github.com/7-Dany/chess/testutil"
)

// TestKingIsAttacking verifies that King.IsAttacking correctly reports
// whether a king of the given color attacks the target square.
//
// A king attacks the eight squares immediately adjacent to it (orthogonal +
// diagonal). The scan checks those eight squares of the target for a king of
// the matching color — it does NOT slide.
func TestKingIsAttacking(t *testing.T) {
	king := King{}

	// All eight squares adjacent to E4. A king on any of these attacks E4.
	t.Run("a white king on any of the eight squares adjacent to E4 attacks E4", func(t *testing.T) {
		//   D3  D4  D5     (down-left, left, up-left)
		//   E3      E5     (down,         up)
		//   F3  F4  F5     (down-right, right, up-right)
		adjacent := []core.Position{
			core.D3, core.D4, core.D5,
			core.E3, core.E5,
			core.F3, core.F4, core.F5,
		}
		for _, from := range adjacent {
			var board core.Board
			board[from] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
			ctx := core.BoardContext{Board: &board}

			if !king.IsAttacking(core.WHITE, core.E4, ctx) {
				t.Errorf("white king on %v should attack E4 (adjacent)", from)
			}
		}
	})

	t.Run("a king two squares away does not attack (kings don't slide)", func(t *testing.T) {
		// E4 → E6 is distance 2 (same file). A king only reaches distance 1.
		twoAway := []core.Position{
			core.E6, // same file, 2 ranks up
			core.C4, // same rank, 2 files left
			core.G6, // diagonal, distance 2
		}
		for _, from := range twoAway {
			var board core.Board
			board[from] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
			ctx := core.BoardContext{Board: &board}

			if king.IsAttacking(core.WHITE, core.E4, ctx) {
				t.Errorf("king on %v should NOT attack E4 (kings don't slide)", from)
			}
		}
	})

	t.Run("a king of the wrong color is ignored", func(t *testing.T) {
		var board core.Board
		board[core.D4] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
		ctx := core.BoardContext{Board: &board}

		// Black king on D4 (adjacent to E4), asking "does WHITE attack E4?" → no.
		if king.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("black king should not count as a white attacker")
		}
		// Same king, asking "does BLACK attack E4?" → yes.
		if !king.IsAttacking(core.BLACK, core.E4, ctx) {
			t.Errorf("black king on D4 should attack E4 for black")
		}
	})

	// Only a king triggers the king attack — a queen/rook/etc on an
	// adjacent square must not falsely report a king attack.
	t.Run("a non-king piece on an adjacent square does not trigger a king attack", func(t *testing.T) {
		nonKings := []core.PieceType{core.QUEEN, core.ROOK, core.BISHOP, core.KNIGHT, core.PAWN}
		for _, pt := range nonKings {
			var board core.Board
			board[core.D4] = core.NewSquare(core.Piece{Type: pt, Color: core.WHITE})
			ctx := core.BoardContext{Board: &board}

			if king.IsAttacking(core.WHITE, core.E4, ctx) {
				t.Errorf("%v on D4 should not trigger a king attack", pt)
			}
		}
	})

	t.Run("an empty board reports no attack", func(t *testing.T) {
		var board core.Board
		ctx := core.BoardContext{Board: &board}

		if king.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("empty board should not report an attack")
		}
	})

	t.Run("a king sitting on the target square itself does not attack it", func(t *testing.T) {
		var board core.Board
		board[core.E4] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		if king.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("king on the target square should not attack itself")
		}
	})

	// A corner target has only 3 adjacent squares.
	t.Run("a corner target is attacked from its 3 adjacent squares only", func(t *testing.T) {
		// A1's neighbors: A2, B1, B2.
		adjacent := []core.Position{core.A2, core.B1, core.B2}
		for _, from := range adjacent {
			var board core.Board
			board[from] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
			ctx := core.BoardContext{Board: &board}

			if !king.IsAttacking(core.WHITE, core.A1, ctx) {
				t.Errorf("king on %v should attack A1 (corner neighbor)", from)
			}
		}

		// A king on the far corner H8 does not attack A1.
		var board core.Board
		board[core.H8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}
		if king.IsAttacking(core.WHITE, core.A1, ctx) {
			t.Errorf("king on H8 should NOT attack A1 (not adjacent)")
		}
	})

	// An edge target has 5 adjacent squares.
	t.Run("an edge target is attacked from its 5 adjacent squares", func(t *testing.T) {
		// A4's neighbors: A3, A5, B3, B4, B5.
		adjacent := []core.Position{core.A3, core.A5, core.B3, core.B4, core.B5}
		for _, from := range adjacent {
			var board core.Board
			board[from] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
			ctx := core.BoardContext{Board: &board}

			if !king.IsAttacking(core.WHITE, core.A4, ctx) {
				t.Errorf("king on %v should attack A4 (edge neighbor)", from)
			}
		}
	})

	t.Run("among multiple kings, any matching-color king on an adjacent square attacks", func(t *testing.T) {
		// Kings on A1 (not adjacent to E4), D4 (adjacent), H8 (not adjacent).
		var board core.Board
		board[core.A1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
		board[core.D4] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
		board[core.H8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		if !king.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("at least one king (D4) should attack E4")
		}
	})

	t.Run("mixed-color kings: only the matching color counts", func(t *testing.T) {
		// Black king on D4 (adjacent, wrong color), white king on E5 (adjacent).
		var board core.Board
		board[core.D4] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
		board[core.E5] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		if !king.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("white king on E5 should attack E4 even with a black king on D4")
		}
	})
}

// TestKingAttacks verifies that King.Attacks returns every square a king
// threatens from the given position.
//
// Unlike the sliders, King.Attacks IGNORES the board entirely — it returns
// all adjacent squares regardless of occupancy. This is correct: a king
// "attacks" every adjacent square even if a friendly piece sits there (the
// attack is what matters for check detection, not move eligibility).
func TestKingAttacks(t *testing.T) {
	king := King{}

	t.Run("king on center D4 threatens all 8 adjacent squares", func(t *testing.T) {
		var board core.Board
		ctx := core.BoardContext{Board: &board}

		got := king.Attacks(make([]core.Position, 0, core.MAX_MOVES), core.D4, ctx)

		testutil.AssertPositionsMatch(t, got, []core.Position{
			// up, down, left, right
			core.D5, core.D3, core.C4, core.E4,
			// up-right, down-right, up-left, down-left
			core.E5, core.E3, core.C5, core.C3,
		})
	})

	// A corner king has only 3 adjacent squares.
	t.Run("king on corner A1 threatens 3 squares (A2, B1, B2)", func(t *testing.T) {
		var board core.Board
		ctx := core.BoardContext{Board: &board}

		got := king.Attacks(make([]core.Position, 0, core.MAX_MOVES), core.A1, ctx)

		testutil.AssertPositionsMatch(t, got, []core.Position{core.A2, core.B1, core.B2})
	})

	t.Run("king on corner H1 threatens 3 squares (G1, G2, H2)", func(t *testing.T) {
		var board core.Board
		ctx := core.BoardContext{Board: &board}

		got := king.Attacks(make([]core.Position, 0, core.MAX_MOVES), core.H1, ctx)

		testutil.AssertPositionsMatch(t, got, []core.Position{core.G1, core.G2, core.H2})
	})

	t.Run("king on corner A8 threatens 3 squares (A7, B7, B8)", func(t *testing.T) {
		var board core.Board
		ctx := core.BoardContext{Board: &board}

		got := king.Attacks(make([]core.Position, 0, core.MAX_MOVES), core.A8, ctx)

		testutil.AssertPositionsMatch(t, got, []core.Position{core.A7, core.B7, core.B8})
	})

	t.Run("king on corner H8 threatens 3 squares (G7, G8, H7)", func(t *testing.T) {
		var board core.Board
		ctx := core.BoardContext{Board: &board}

		got := king.Attacks(make([]core.Position, 0, core.MAX_MOVES), core.H8, ctx)

		testutil.AssertPositionsMatch(t, got, []core.Position{core.G7, core.G8, core.H7})
	})

	// An edge king has 5 adjacent squares (one direction falls off the board).
	t.Run("king on edge A4 threatens 5 squares", func(t *testing.T) {
		var board core.Board
		ctx := core.BoardContext{Board: &board}

		got := king.Attacks(make([]core.Position, 0, core.MAX_MOVES), core.A4, ctx)

		testutil.AssertPositionsMatch(t, got, []core.Position{core.A3, core.A5, core.B3, core.B4, core.B5})
	})

	t.Run("king on edge D1 threatens 5 squares", func(t *testing.T) {
		var board core.Board
		ctx := core.BoardContext{Board: &board}

		got := king.Attacks(make([]core.Position, 0, core.MAX_MOVES), core.D1, ctx)

		testutil.AssertPositionsMatch(t, got, []core.Position{core.C1, core.C2, core.D2, core.E1, core.E2})
	})

	// The unique king behavior: Attacks ignores the board. A piece on an
	// adjacent square is still "attacked" (important for check detection —
	// a king defends the squares around it even if occupied by friends).
	t.Run("king attacks adjacent squares even when they are occupied (Attacks ignores the board)", func(t *testing.T) {
		// Fill every adjacent square to D4 with pieces (mix of friendly and enemy).
		var board core.Board
		board[core.C3] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.C4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.C5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.D3] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		board[core.D5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		board[core.E3] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
		board[core.E4] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
		board[core.E5] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.BLACK})
		ctx := core.BoardContext{Board: &board}

		got := king.Attacks(make([]core.Position, 0, core.MAX_MOVES), core.D4, ctx)

		// All 8 adjacent squares are returned, regardless of occupancy.
		testutil.AssertPositionsMatch(t, got, []core.Position{
			core.D5, core.D3, core.C4, core.E4,
			core.E5, core.E3, core.C5, core.C3,
		})
	})
}

// TestKingPseudoLegalMoves verifies that King.PseudoLegalMoves returns the
// correct set of moves — respecting occupancy and captures, but NOT
// filtering for king safety (that's the engine's job).
//
// Key rules for king moves:
//   - A king moves one square in any of the 8 directions (no sliding).
//   - A square occupied by a friendly piece is excluded (can't capture own).
//   - A square occupied by an enemy piece is included as a capture.
//   - Every king move is type NORMAL. Castling is NOT included here — the
//     engine adds it separately (see engine.castlingMoves).
func TestKingPseudoLegalMoves(t *testing.T) {
	king := King{}

	// Helper: extract just the destination squares from a move list.
	destinations := func(moves []core.Move) []core.Position {
		tos := make([]core.Position, len(moves))
		for i, m := range moves {
			tos[i] = m.To
		}
		return tos
	}

	// All 8 destinations from D4 on an empty board.
	d4Adjacent := []core.Position{
		core.D5, core.D3, core.C4, core.E4, // orthogonal
		core.E5, core.E3, core.C5, core.C3, // diagonal
	}

	t.Run("king on center D4 with an empty board has 8 moves", func(t *testing.T) {
		var board core.Board
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := king.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.D4, ctx)

		testutil.AssertPositionsMatch(t, destinations(moves), d4Adjacent)
		testutil.AssertMoveCount(t, moves, 8)
	})

	t.Run("a square occupied by a friendly piece is excluded from the move list", func(t *testing.T) {
		// Own pawn on D5 (up). D5 is NOT in the move list; the other 7 remain.
		var board core.Board
		board[core.D5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := king.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.D4, ctx)

		testutil.AssertPositionsMatch(t, destinations(moves), []core.Position{
			core.D3, core.C4, core.E4,
			core.E5, core.E3, core.C5, core.C3,
		})
	})

	t.Run("a square occupied by an enemy piece is included as a capture", func(t *testing.T) {
		// Enemy pawn on D5 (up). D5 IS in the move list, flagged as a capture.
		var board core.Board
		board[core.D5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := king.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.D4, ctx)

		// All 8 destinations present (D5 is a capture, the rest are quiet).
		testutil.AssertPositionsMatch(t, destinations(moves), d4Adjacent)

		// The move to D5 must be flagged as a capture of the black pawn.
		for _, m := range moves {
			if m.To == core.D5 {
				if !m.HasCapture {
					t.Errorf("move to D5 should be a capture")
				}
				if m.Captured != (core.Piece{Type: core.PAWN, Color: core.BLACK}) {
					t.Errorf("move to D5: Captured = %v, want black pawn", m.Captured)
				}
				return
			}
		}
		t.Errorf("expected a move to D5 (the capture square)")
	})

	t.Run("captures carry the exact enemy piece type and color sitting on the destination", func(t *testing.T) {
		// Three enemy pieces on three different adjacent squares: a queen on
		// D5 (up), a rook on E4 (right), a knight on C3 (down-left).
		var board core.Board
		board[core.D5] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.BLACK})
		board[core.E4] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
		board[core.C3] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.BLACK})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := king.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.D4, ctx)

		wantCaptures := map[core.Position]core.Piece{
			core.D5: {Type: core.QUEEN, Color: core.BLACK},
			core.E4: {Type: core.ROOK, Color: core.BLACK},
			core.C3: {Type: core.KNIGHT, Color: core.BLACK},
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

	t.Run("a mix of friendly and enemy on adjacent squares excludes own, includes enemy", func(t *testing.T) {
		// D5 (own, up), C4 (own, left), E4 (enemy, right), C3 (enemy, down-left).
		var board core.Board
		board[core.D5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.C4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.E4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		board[core.C3] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := king.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.D4, ctx)

		// D5 and C4 excluded (own); E4 and C3 included (enemy captures);
		// the other 4 adjacent squares are quiet moves.
		testutil.AssertPositionsMatch(t, destinations(moves), []core.Position{
			core.D3, core.E4,
			core.E5, core.E3, core.C5, core.C3,
		})
	})

	t.Run("all 8 adjacent squares blocked by own pieces yields no moves", func(t *testing.T) {
		var board core.Board
		for _, pos := range d4Adjacent {
			board[pos] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		}
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := king.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.D4, ctx)

		testutil.AssertNoMoves(t, moves)
	})

	t.Run("king on corner A1 with an empty board has 3 moves", func(t *testing.T) {
		var board core.Board
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := king.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.A1, ctx)

		testutil.AssertPositionsMatch(t, destinations(moves), []core.Position{core.A2, core.B1, core.B2})
	})

	t.Run("king on edge A4 with an empty board has 5 moves", func(t *testing.T) {
		var board core.Board
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := king.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.A4, ctx)

		testutil.AssertPositionsMatch(t, destinations(moves), []core.Position{core.A3, core.A5, core.B3, core.B4, core.B5})
	})

	t.Run("a black king treats white pieces as enemies (captures) and black as own", func(t *testing.T) {
		// White pawn on D5 (enemy for black king), black pawn on C4 (own).
		var board core.Board
		board[core.D5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.C4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.BLACK}

		moves := king.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.D4, ctx)

		// D5 included (capture of white); C4 excluded (own); 6 quiet moves.
		testutil.AssertPositionsMatch(t, destinations(moves), []core.Position{
			core.D5, core.D3, core.E4,
			core.E5, core.E3, core.C5, core.C3,
		})

		// The move to D5 is a capture of the white pawn.
		for _, m := range moves {
			if m.To == core.D5 {
				if !m.HasCapture || m.Captured != (core.Piece{Type: core.PAWN, Color: core.WHITE}) {
					t.Errorf("black king capturing D5: HasCapture=%v Captured=%v, want white pawn", m.HasCapture, m.Captured)
				}
				return
			}
		}
		t.Errorf("expected a move to D5")
	})

	// Every king move is type NORMAL. Castling is NOT generated here — the
	// engine adds it separately (a king "pseudo-legal move" is just a
	// one-step move). Each move carries the mover and its source square.
	t.Run("every generated move has type NORMAL and carries the mover and source square (no castling)", func(t *testing.T) {
		var board core.Board
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := king.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.D4, ctx)

		mover := core.Piece{Type: core.KING, Color: core.WHITE}
		for _, m := range moves {
			if m.Type != core.NORMAL {
				t.Errorf("move to %v: Type = %v, want NORMAL (castling is added by the engine, not the piece)", m.To, m.Type)
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
