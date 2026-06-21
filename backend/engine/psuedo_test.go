package engine

import (
	"testing"

	"github.com/7-Dany/chess/core"
	"github.com/7-Dany/chess/piece"
	"github.com/7-Dany/chess/testutil"
)

// TestGetPseudoLegalMoves verifies the engine's dispatch layer — the thin
// wrapper that checks whose turn it is, delegates to the piece, and adds
// castling for the king.
//
// The per-piece move logic itself is tested exhaustively in the piece
// package; here we only verify the dispatch: empty square, enemy piece, and
// own piece. Castling is tested separately in TestCastlingMoves.
func TestGetPseudoLegalMoves(t *testing.T) {
	engine := GetDefaultEngine()

	t.Run("an empty square returns no moves", func(t *testing.T) {
		// E4 is empty; querying it should yield nothing.
		var board core.Board
		board[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
		ctx := testutil.NewTurn(&board, core.WHITE)

		moves := engine.GetPseudoLegalMoves(make([]core.Move, 0, piece.MAX_MOVES), core.E4, *ctx)

		testutil.AssertNoMoves(t, moves)
	})

	t.Run("an enemy piece returns no moves (can't move the opponent's piece)", func(t *testing.T) {
		// Black king on E8, white to move → querying E8 returns nothing.
		var board core.Board
		board[core.E8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
		ctx := testutil.NewTurn(&board, core.WHITE)

		moves := engine.GetPseudoLegalMoves(make([]core.Move, 0, piece.MAX_MOVES), core.E8, *ctx)

		testutil.AssertNoMoves(t, moves)
	})

	t.Run("a friendly piece delegates to its piece implementation and returns moves", func(t *testing.T) {
		// White knight on B1 → should return the knight's moves (e.g. A3, C3).
		var board core.Board
		board[core.B1] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
		ctx := testutil.NewTurn(&board, core.WHITE)

		moves := engine.GetPseudoLegalMoves(make([]core.Move, 0, piece.MAX_MOVES), core.B1, *ctx)

		// A knight on B1 can reach A3 and C3 (among others). We just check
		// those two are present — the full knight move set is tested in the
		// piece package.
		testutil.AssertMovePresent(t, moves, core.B1, core.A3)
		testutil.AssertMovePresent(t, moves, core.B1, core.C3)
	})
}

// TestCastlingMoves verifies that GetPseudoLegalMoves adds castling moves
// for the king when the conditions are met.
//
// Castling rules (checked by castlingMoves / canCastleKingSide /
// canCastleQueenSide):
//   - The king must be on its home file (E). If not, bail (corrupt state).
//   - The king must not be in check.
//   - King-side: CanCastleKingSide right, F and G empty, F and G not attacked.
//   - Queen-side: CanCastleQueenSide right, B, C, D empty, C and D not
//     attacked (B may be attacked — the king doesn't pass through it).
func TestCastlingMoves(t *testing.T) {
	engine := GetDefaultEngine()

	// Helper: build a context with a king on E1 or E8 and the given side states.
	// The king is placed on the board so GetPseudoLegalMoves finds it.
	kingCtx := func(kingPos core.Position, side core.PieceColor, sides [2]core.SideState, setup func(*core.Board)) core.TurnContext {
		var board core.Board
		king := core.Piece{Type: core.KING, Color: side}
		board[kingPos] = core.NewSquare(king)
		setup(&board)
		ctx := core.TurnContext{
			MoveContext: core.MoveContext{
				BoardContext: core.BoardContext{Board: &board},
				SideToMove:   side,
				Sides:        sides,
			},
		}
		return ctx
	}

	// Standard starting side states: both kings on home squares, full rights.
	defaultSides := testutil.DefaultSides()

	t.Run("a king not on the E file returns no castling moves (corrupt state guard)", func(t *testing.T) {
		// King on D1 (not E1) but rights still set — the guard bails.
		ctx := kingCtx(core.D1, core.WHITE, defaultSides, func(b *core.Board) {})

		moves := engine.GetPseudoLegalMoves(make([]core.Move, 0, piece.MAX_MOVES), core.D1, ctx)

		// No castling moves (only the king's one-step moves, which we don't check here).
		for _, m := range moves {
			if m.Type == core.CASTLING {
				t.Errorf("king not on E file should not generate castling moves, got %v", m)
			}
		}
	})

	t.Run("a king in check returns no castling moves", func(t *testing.T) {
		// King on E1, enemy rook on E8 attacks the king (same file, clear).
		ctx := kingCtx(core.E1, core.WHITE, defaultSides, func(b *core.Board) {
			b[core.E8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
		})

		moves := engine.GetPseudoLegalMoves(make([]core.Move, 0, piece.MAX_MOVES), core.E1, ctx)

		for _, m := range moves {
			if m.Type == core.CASTLING {
				t.Errorf("king in check should not generate castling moves, got %v", m)
			}
		}
	})

	t.Run("with both rights and clear paths, the king gets two castling moves", func(t *testing.T) {
		// King on E1, F1/G1/B1/C1/D1 all empty, no enemy attackers.
		ctx := kingCtx(core.E1, core.WHITE, defaultSides, func(b *core.Board) {})

		moves := engine.GetPseudoLegalMoves(make([]core.Move, 0, piece.MAX_MOVES), core.E1, ctx)

		// King-side: E1 → G1. Queen-side: E1 → C1.
		testutil.AssertMovePresent(t, moves, core.E1, core.G1)
		testutil.AssertMovePresent(t, moves, core.E1, core.C1)
	})

	t.Run("with only the king-side right, the king gets one castling move (to G1)", func(t *testing.T) {
		kingSideOnly := [2]core.SideState{
			{KingPosition: core.E1, CanCastleKingSide: true},
			testutil.FullBlack(),
		}
		ctx := kingCtx(core.E1, core.WHITE, kingSideOnly, func(b *core.Board) {})

		moves := engine.GetPseudoLegalMoves(make([]core.Move, 0, piece.MAX_MOVES), core.E1, ctx)

		testutil.AssertMovePresent(t, moves, core.E1, core.G1)
		testutil.AssertMoveAbsent(t, moves, core.E1, core.C1)
	})

	t.Run("with only the queen-side right, the king gets one castling move (to C1)", func(t *testing.T) {
		queenSideOnly := [2]core.SideState{
			{KingPosition: core.E1, CanCastleQueenSide: true},
			testutil.FullBlack(),
		}
		ctx := kingCtx(core.E1, core.WHITE, queenSideOnly, func(b *core.Board) {})

		moves := engine.GetPseudoLegalMoves(make([]core.Move, 0, piece.MAX_MOVES), core.E1, ctx)

		testutil.AssertMovePresent(t, moves, core.E1, core.C1)
		testutil.AssertMoveAbsent(t, moves, core.E1, core.G1)
	})

	t.Run("black king on E8 with both rights and clear paths gets two castling moves", func(t *testing.T) {
		ctx := kingCtx(core.E8, core.BLACK, defaultSides, func(b *core.Board) {})

		moves := engine.GetPseudoLegalMoves(make([]core.Move, 0, piece.MAX_MOVES), core.E8, ctx)

		// King-side: E8 → G8. Queen-side: E8 → C8.
		testutil.AssertMovePresent(t, moves, core.E8, core.G8)
		testutil.AssertMovePresent(t, moves, core.E8, core.C8)
	})

	t.Run("both sides blocked by pieces returns no castling moves", func(t *testing.T) {
		// F1, G1, B1, C1, D1 all occupied by own pieces.
		ctx := kingCtx(core.E1, core.WHITE, defaultSides, func(b *core.Board) {
			b[core.F1] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
			b[core.G1] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
			b[core.B1] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
			b[core.C1] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
			b[core.D1] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.WHITE})
		})

		moves := engine.GetPseudoLegalMoves(make([]core.Move, 0, piece.MAX_MOVES), core.E1, ctx)

		for _, m := range moves {
			if m.Type == core.CASTLING {
				t.Errorf("blocked paths should not generate castling moves, got %v", m)
			}
		}
	})

	t.Run("queen-side castling is removed when B1 is occupied (path not clear)", func(t *testing.T) {
		// B1 occupied → queen-side blocked. King-side (F1, G1 empty) still works.
		ctx := kingCtx(core.E1, core.WHITE, defaultSides, func(b *core.Board) {
			b[core.B1] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
		})

		moves := engine.GetPseudoLegalMoves(make([]core.Move, 0, piece.MAX_MOVES), core.E1, ctx)

		testutil.AssertMovePresent(t, moves, core.E1, core.G1) // king-side still available
		testutil.AssertMoveAbsent(t, moves, core.E1, core.C1)  // queen-side blocked by B1
	})

	t.Run("queen-side castling is removed when an enemy piece sits on C1 (path blocked)", func(t *testing.T) {
		// Enemy knight on C1 blocks the queen-side path (and attacks D1 too).
		ctx := kingCtx(core.E1, core.WHITE, defaultSides, func(b *core.Board) {
			b[core.C1] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.BLACK})
		})

		moves := engine.GetPseudoLegalMoves(make([]core.Move, 0, piece.MAX_MOVES), core.E1, ctx)

		testutil.AssertMoveAbsent(t, moves, core.E1, core.C1)
	})
}

// TestCanCastleKingSide verifies the king-side castling eligibility check in
// isolation. The king-side castling requires:
//   - The CanCastleKingSide right.
//   - F1 (path) and G1 (destination) both empty.
//   - F1 and G1 both not attacked by the enemy.
func TestCanCastleKingSide(t *testing.T) {
	engine := GetDefaultEngine()

	// Helper: build a context and check king-side eligibility.
	checkKingSide := func(side core.PieceColor, rank core.Rank, sides [2]core.SideState, setup func(*core.Board)) bool {
		var board core.Board
		ctx := core.TurnContext{
			MoveContext: core.MoveContext{
				BoardContext: core.BoardContext{Board: &board},
				SideToMove:   side,
				Sides:        sides,
			},
		}
		setup(&board)
		return engine.canCastleKingSide(rank, ctx)
	}

	defaultSides := testutil.DefaultSides()
	noKingSide := [2]core.SideState{
		{KingPosition: core.E1, CanCastleQueenSide: true},
		testutil.FullBlack(),
	}

	t.Run("no king-side right returns false", func(t *testing.T) {
		got := checkKingSide(core.WHITE, core.RANK_1, noKingSide, func(b *core.Board) {})
		if got {
			t.Errorf("canCastleKingSide = true, want false (no right)")
		}
	})

	t.Run("F1 occupied by a friendly piece returns false", func(t *testing.T) {
		got := checkKingSide(core.WHITE, core.RANK_1, defaultSides, func(b *core.Board) {
			b[core.F1] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
		})
		if got {
			t.Errorf("canCastleKingSide = true, want false (F1 occupied by own piece)")
		}
	})

	t.Run("G1 occupied by a friendly piece returns false", func(t *testing.T) {
		got := checkKingSide(core.WHITE, core.RANK_1, defaultSides, func(b *core.Board) {
			b[core.G1] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
		})
		if got {
			t.Errorf("canCastleKingSide = true, want false (G1 occupied by own piece)")
		}
	})

	t.Run("F1 occupied by an enemy piece returns false (path blocked)", func(t *testing.T) {
		got := checkKingSide(core.WHITE, core.RANK_1, defaultSides, func(b *core.Board) {
			b[core.F1] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.BLACK})
		})
		if got {
			t.Errorf("canCastleKingSide = true, want false (F1 occupied by enemy)")
		}
	})

	t.Run("F1 attacked by an enemy rook returns false", func(t *testing.T) {
		// Rook on F8 attacks F1 down the file (clear path).
		got := checkKingSide(core.WHITE, core.RANK_1, defaultSides, func(b *core.Board) {
			b[core.F8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
		})
		if got {
			t.Errorf("canCastleKingSide = true, want false (F1 attacked)")
		}
	})

	t.Run("G1 attacked by an enemy rook returns false", func(t *testing.T) {
		got := checkKingSide(core.WHITE, core.RANK_1, defaultSides, func(b *core.Board) {
			b[core.G8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
		})
		if got {
			t.Errorf("canCastleKingSide = true, want false (G1 attacked)")
		}
	})

	t.Run("F1 attacked by an enemy bishop returns false", func(t *testing.T) {
		// Bishop on A6 attacks F1 via the a6-f1 diagonal (a6→b5→c4→d3→e2→f1).
		got := checkKingSide(core.WHITE, core.RANK_1, defaultSides, func(b *core.Board) {
			b[core.A6] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.BLACK})
		})
		if got {
			t.Errorf("canCastleKingSide = true, want false (F1 attacked by bishop)")
		}
	})

	t.Run("all clear on white rank 1 returns true", func(t *testing.T) {
		got := checkKingSide(core.WHITE, core.RANK_1, defaultSides, func(b *core.Board) {})
		if !got {
			t.Errorf("canCastleKingSide = false, want true (all clear)")
		}
	})

	t.Run("F8 occupied by a friendly piece returns false (black rank 8)", func(t *testing.T) {
		got := checkKingSide(core.BLACK, core.RANK_8, defaultSides, func(b *core.Board) {
			b[core.F8] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.BLACK})
		})
		if got {
			t.Errorf("canCastleKingSide = true, want false (F8 occupied by own piece)")
		}
	})

	t.Run("F8 attacked by a white rook returns false (black rank 8)", func(t *testing.T) {
		got := checkKingSide(core.BLACK, core.RANK_8, defaultSides, func(b *core.Board) {
			b[core.F1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
		})
		if got {
			t.Errorf("canCastleKingSide = true, want false (F8 attacked by white rook)")
		}
	})

	t.Run("all clear on black rank 8 returns true", func(t *testing.T) {
		got := checkKingSide(core.BLACK, core.RANK_8, defaultSides, func(b *core.Board) {})
		if !got {
			t.Errorf("canCastleKingSide = false, want true (all clear black rank 8)")
		}
	})
}

// TestCanCastleQueenSide verifies the queen-side castling eligibility check.
// The queen-side castling requires:
//   - The CanCastleQueenSide right.
//   - B1, C1 (destination), and D1 (path) all empty. (B1 is between the rook
//     and the king's path — the king doesn't pass through it, but the rook
//     does, so it must be empty.)
//   - C1 and D1 both not attacked by the enemy. (B1 may be attacked — the
//     king never goes there.)
func TestCanCastleQueenSide(t *testing.T) {
	engine := GetDefaultEngine()

	// Helper: build a context and check queen-side eligibility.
	checkQueenSide := func(side core.PieceColor, rank core.Rank, sides [2]core.SideState, setup func(*core.Board)) bool {
		var board core.Board
		ctx := core.TurnContext{
			MoveContext: core.MoveContext{
				BoardContext: core.BoardContext{Board: &board},
				SideToMove:   side,
				Sides:        sides,
			},
		}
		setup(&board)
		return engine.canCastleQueenSide(rank, ctx)
	}

	defaultSides := testutil.DefaultSides()
	noQueenSide := [2]core.SideState{
		{KingPosition: core.E1, CanCastleKingSide: true},
		testutil.FullBlack(),
	}

	t.Run("no queen-side right returns false", func(t *testing.T) {
		got := checkQueenSide(core.WHITE, core.RANK_1, noQueenSide, func(b *core.Board) {})
		if got {
			t.Errorf("canCastleQueenSide = true, want false (no right)")
		}
	})

	t.Run("B1 occupied by a friendly piece returns false (rook can't pass)", func(t *testing.T) {
		got := checkQueenSide(core.WHITE, core.RANK_1, defaultSides, func(b *core.Board) {
			b[core.B1] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
		})
		if got {
			t.Errorf("canCastleQueenSide = true, want false (B1 occupied)")
		}
	})

	t.Run("C1 occupied by a friendly piece returns false (destination blocked)", func(t *testing.T) {
		got := checkQueenSide(core.WHITE, core.RANK_1, defaultSides, func(b *core.Board) {
			b[core.C1] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
		})
		if got {
			t.Errorf("canCastleQueenSide = true, want false (C1 occupied)")
		}
	})

	t.Run("D1 occupied by a friendly piece returns false (path blocked)", func(t *testing.T) {
		got := checkQueenSide(core.WHITE, core.RANK_1, defaultSides, func(b *core.Board) {
			b[core.D1] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.WHITE})
		})
		if got {
			t.Errorf("canCastleQueenSide = true, want false (D1 occupied)")
		}
	})

	t.Run("C1 occupied by an enemy piece returns false (path blocked)", func(t *testing.T) {
		got := checkQueenSide(core.WHITE, core.RANK_1, defaultSides, func(b *core.Board) {
			b[core.C1] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.BLACK})
		})
		if got {
			t.Errorf("canCastleQueenSide = true, want false (C1 occupied by enemy)")
		}
	})

	t.Run("D1 attacked by an enemy rook returns false (king passes through D1)", func(t *testing.T) {
		// Rook on D8 attacks D1 down the file.
		got := checkQueenSide(core.WHITE, core.RANK_1, defaultSides, func(b *core.Board) {
			b[core.D8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
		})
		if got {
			t.Errorf("canCastleQueenSide = true, want false (D1 attacked)")
		}
	})

	t.Run("C1 attacked by an enemy rook returns false (king lands on C1)", func(t *testing.T) {
		got := checkQueenSide(core.WHITE, core.RANK_1, defaultSides, func(b *core.Board) {
			b[core.C8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
		})
		if got {
			t.Errorf("canCastleQueenSide = true, want false (C1 attacked)")
		}
	})

	t.Run("D1 attacked by an enemy bishop returns false", func(t *testing.T) {
		// Bishop on A4 attacks D1 via the a4-d1 diagonal (a4→b3→c2→d1).
		got := checkQueenSide(core.WHITE, core.RANK_1, defaultSides, func(b *core.Board) {
			b[core.A4] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.BLACK})
		})
		if got {
			t.Errorf("canCastleQueenSide = true, want false (D1 attacked by bishop)")
		}
	})

	t.Run("B1 attacked but C1 and D1 safe returns true (king never passes through B1)", func(t *testing.T) {
		// Rook on B8 attacks B1, but the king's path (D1, C1) is safe.
		got := checkQueenSide(core.WHITE, core.RANK_1, defaultSides, func(b *core.Board) {
			b[core.B8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
		})
		if !got {
			t.Errorf("canCastleQueenSide = false, want true (B1 attacked is OK; king doesn't pass through it)")
		}
	})

	t.Run("all clear on white rank 1 returns true", func(t *testing.T) {
		got := checkQueenSide(core.WHITE, core.RANK_1, defaultSides, func(b *core.Board) {})
		if !got {
			t.Errorf("canCastleQueenSide = false, want true (all clear)")
		}
	})

	t.Run("B8 occupied by a friendly piece returns false (black rank 8)", func(t *testing.T) {
		got := checkQueenSide(core.BLACK, core.RANK_8, defaultSides, func(b *core.Board) {
			b[core.B8] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.BLACK})
		})
		if got {
			t.Errorf("canCastleQueenSide = true, want false (B8 occupied by own piece)")
		}
	})

	t.Run("D8 attacked by a white rook returns false (black rank 8)", func(t *testing.T) {
		got := checkQueenSide(core.BLACK, core.RANK_8, defaultSides, func(b *core.Board) {
			b[core.D1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
		})
		if got {
			t.Errorf("canCastleQueenSide = true, want false (D8 attacked by white rook)")
		}
	})

	t.Run("B8 attacked but C8 and D8 safe returns true (black queen-side castling allowed)", func(t *testing.T) {
		// White rook on B1 attacks B8, but the black king's path (D8, C8) is safe.
		got := checkQueenSide(core.BLACK, core.RANK_8, defaultSides, func(b *core.Board) {
			b[core.B1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
		})
		if !got {
			t.Errorf("canCastleQueenSide = false, want true (B8 attacked is OK for black)")
		}
	})

	t.Run("all clear on black rank 8 returns true", func(t *testing.T) {
		got := checkQueenSide(core.BLACK, core.RANK_8, defaultSides, func(b *core.Board) {})
		if !got {
			t.Errorf("canCastleQueenSide = false, want true (all clear black rank 8)")
		}
	})
}
