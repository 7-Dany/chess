package engine

import (
	"testing"

	"github.com/7-Dany/chess/core"
)

// TestIsSquareAttacked verifies that IsSquareAttacked correctly reports
// whether any piece of the given color attacks the target square.
//
// The method combines five attack checks:
//   - Knight (leaper): scans the 8 L-shapes around the target.
//   - King (leaper): scans the 8 adjacent squares.
//   - Pawn (leaper, color-dependent): scans the two squares "behind" the
//     target relative to the attacker's movement direction.
//   - Bishop/Queen (diagonal sliders): scans the 4 diagonal rays.
//   - Rook/Queen (orthogonal sliders): scans the 4 orthogonal rays.
//
// A queen is covered by BOTH slider scans (diagonal + orthogonal), so it's
// detected on any of its 8 rays. Each slider ray stops at the first
// occupied square — a blocker breaks the attack regardless of color.
func TestIsSquareAttacked(t *testing.T) {
	engine := GetDefaultEngine()

	// Helper: build a BoardContext from a setup function.
	withBoard := func(setup func(*core.Board)) core.BoardContext {
		var board core.Board
		setup(&board)
		return core.BoardContext{Board: &board}
	}

	// Helper: assert the attack result.
	assertAttack := func(t *testing.T, target core.Position, color core.PieceColor, ctx core.BoardContext, want bool) {
		t.Helper()
		got := engine.IsSquareAttacked(target, color, ctx)
		if got != want {
			t.Errorf("IsSquareAttacked(%v, %v) = %v, want %v", target, color, got, want)
		}
	}

	// =========================================================================
	// Knight — leaper, scans 8 L-shapes around the target.
	// =========================================================================

	t.Run("a knight on a valid L-shape from E5 attacks E5", func(t *testing.T) {
		// D7 is an L-shape from E5 (file -1, rank +2).
		ctx := withBoard(func(b *core.Board) {
			b[core.D7] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
		})
		assertAttack(t, core.E5, core.WHITE, ctx, true)
	})

	t.Run("a knight on a non-L-shape square does not attack E5", func(t *testing.T) {
		// D6 is adjacent (not an L-shape).
		ctx := withBoard(func(b *core.Board) {
			b[core.D6] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
		})
		assertAttack(t, core.E5, core.WHITE, ctx, false)
	})

	t.Run("a knight attacks from a corner L-shape (H8 to F7)", func(t *testing.T) {
		ctx := withBoard(func(b *core.Board) {
			b[core.H8] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
		})
		assertAttack(t, core.F7, core.WHITE, ctx, true)
	})

	// =========================================================================
	// King — leaper, scans 8 adjacent squares.
	// =========================================================================

	t.Run("a king on an orthogonally adjacent square attacks E5", func(t *testing.T) {
		ctx := withBoard(func(b *core.Board) {
			b[core.E6] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
		})
		assertAttack(t, core.E5, core.WHITE, ctx, true)
	})

	t.Run("a king on a diagonally adjacent square attacks E5", func(t *testing.T) {
		ctx := withBoard(func(b *core.Board) {
			b[core.F6] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
		})
		assertAttack(t, core.E5, core.WHITE, ctx, true)
	})

	t.Run("a king two squares away does not attack E5", func(t *testing.T) {
		ctx := withBoard(func(b *core.Board) {
			b[core.E7] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
		})
		assertAttack(t, core.E5, core.WHITE, ctx, false)
	})

	// =========================================================================
	// Pawn — leaper, color-dependent direction.
	// White pawns attack the two squares diagonally ABOVE them; black pawns
	// attack the two squares diagonally BELOW them.
	// =========================================================================

	t.Run("a white pawn attacks E5 from below-left (D4)", func(t *testing.T) {
		ctx := withBoard(func(b *core.Board) {
			b[core.D4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		})
		assertAttack(t, core.E5, core.WHITE, ctx, true)
	})

	t.Run("a white pawn attacks E5 from below-right (F4)", func(t *testing.T) {
		ctx := withBoard(func(b *core.Board) {
			b[core.F4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		})
		assertAttack(t, core.E5, core.WHITE, ctx, true)
	})

	t.Run("a white pawn does not attack E5 from above (pawns attack upward only)", func(t *testing.T) {
		// D6 is above E5; a white pawn there attacks C7 and E7, not E5.
		ctx := withBoard(func(b *core.Board) {
			b[core.D6] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		})
		assertAttack(t, core.E5, core.WHITE, ctx, false)
	})

	t.Run("a white pawn does not attack E5 from the same file (pawns attack diagonally only)", func(t *testing.T) {
		ctx := withBoard(func(b *core.Board) {
			b[core.E4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		})
		assertAttack(t, core.E5, core.WHITE, ctx, false)
	})

	t.Run("a black pawn attacks E5 from above-left (D6)", func(t *testing.T) {
		ctx := withBoard(func(b *core.Board) {
			b[core.D6] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		})
		assertAttack(t, core.E5, core.BLACK, ctx, true)
	})

	t.Run("a black pawn attacks E5 from above-right (F6)", func(t *testing.T) {
		ctx := withBoard(func(b *core.Board) {
			b[core.F6] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		})
		assertAttack(t, core.E5, core.BLACK, ctx, true)
	})

	t.Run("a black pawn does not attack E5 from below (black pawns attack downward only)", func(t *testing.T) {
		ctx := withBoard(func(b *core.Board) {
			b[core.D4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		})
		assertAttack(t, core.E5, core.BLACK, ctx, false)
	})

	// =========================================================================
	// Bishop — diagonal slider, stops at first blocker.
	// =========================================================================

	t.Run("a bishop attacks along a clear long diagonal (B2 to E5)", func(t *testing.T) {
		ctx := withBoard(func(b *core.Board) {
			b[core.B2] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
		})
		assertAttack(t, core.E5, core.WHITE, ctx, true)
	})

	t.Run("a bishop attacks along a short diagonal (D4 to E5)", func(t *testing.T) {
		ctx := withBoard(func(b *core.Board) {
			b[core.D4] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
		})
		assertAttack(t, core.E5, core.WHITE, ctx, true)
	})

	t.Run("a bishop blocked by a friendly piece between it and the target does not attack", func(t *testing.T) {
		// Bishop on B2 → D4 (own knight) → E5 (target). The knight blocks.
		ctx := withBoard(func(b *core.Board) {
			b[core.B2] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
			b[core.D4] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
		})
		assertAttack(t, core.E5, core.WHITE, ctx, false)
	})

	t.Run("a bishop blocked by an enemy piece between it and the target does not attack", func(t *testing.T) {
		ctx := withBoard(func(b *core.Board) {
			b[core.B2] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
			b[core.D4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		})
		assertAttack(t, core.E5, core.WHITE, ctx, false)
	})

	t.Run("a bishop does not attack orthogonally (same file, not diagonal)", func(t *testing.T) {
		ctx := withBoard(func(b *core.Board) {
			b[core.E3] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
		})
		assertAttack(t, core.E5, core.WHITE, ctx, false)
	})

	// =========================================================================
	// Rook — orthogonal slider, stops at first blocker.
	// =========================================================================

	t.Run("a rook attacks along a clear file (E2 to E5)", func(t *testing.T) {
		ctx := withBoard(func(b *core.Board) {
			b[core.E2] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
		})
		assertAttack(t, core.E5, core.WHITE, ctx, true)
	})

	t.Run("a rook attacks along a clear rank (A5 to E5)", func(t *testing.T) {
		ctx := withBoard(func(b *core.Board) {
			b[core.A5] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
		})
		assertAttack(t, core.E5, core.WHITE, ctx, true)
	})

	t.Run("a rook blocked by a friendly piece between it and the target does not attack", func(t *testing.T) {
		// Rook on E2 → E4 (own knight) → E5 (target). The knight blocks.
		ctx := withBoard(func(b *core.Board) {
			b[core.E2] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			b[core.E4] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
		})
		assertAttack(t, core.E5, core.WHITE, ctx, false)
	})

	t.Run("a rook blocked by an enemy piece between it and the target does not attack", func(t *testing.T) {
		ctx := withBoard(func(b *core.Board) {
			b[core.E2] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			b[core.E4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		})
		assertAttack(t, core.E5, core.WHITE, ctx, false)
	})

	t.Run("a rook does not attack diagonally", func(t *testing.T) {
		ctx := withBoard(func(b *core.Board) {
			b[core.D4] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
		})
		assertAttack(t, core.E5, core.WHITE, ctx, false)
	})

	// =========================================================================
	// Queen — covered by BOTH the diagonal and orthogonal slider scans.
	// A queen on any of its 8 rays is detected.
	// =========================================================================

	t.Run("a queen on a diagonal is caught by the bishop-path scan (B2 to E5)", func(t *testing.T) {
		ctx := withBoard(func(b *core.Board) {
			b[core.B2] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.WHITE})
		})
		assertAttack(t, core.E5, core.WHITE, ctx, true)
	})

	t.Run("a queen on an orthogonal line is caught by the rook-path scan (E2 to E5)", func(t *testing.T) {
		ctx := withBoard(func(b *core.Board) {
			b[core.E2] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.WHITE})
		})
		assertAttack(t, core.E5, core.WHITE, ctx, true)
	})

	t.Run("a queen blocked on the diagonal by a friendly piece does not attack via that ray", func(t *testing.T) {
		// Queen on B2 (diagonal) blocked by own knight on D4. No other ray
		// reaches E5, so the result is false.
		ctx := withBoard(func(b *core.Board) {
			b[core.B2] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.WHITE})
			b[core.D4] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
		})
		assertAttack(t, core.E5, core.WHITE, ctx, false)
	})

	t.Run("a queen blocked orthogonally by an enemy piece does not attack via that ray", func(t *testing.T) {
		ctx := withBoard(func(b *core.Board) {
			b[core.E2] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.WHITE})
			b[core.E4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		})
		assertAttack(t, core.E5, core.WHITE, ctx, false)
	})

	t.Run("a queen blocked on one ray but clear on another still attacks", func(t *testing.T) {
		// Queen on B2 (diagonal) blocked by D4, but a rook on A5 (rank) is
		// clear → the rook provides the attack. This tests that the slider
		// scan continues to the next ray after a block.
		ctx := withBoard(func(b *core.Board) {
			b[core.B2] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.WHITE})
			b[core.D4] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE}) // blocks queen's diagonal
			b[core.A5] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})   // clear rank attack
		})
		assertAttack(t, core.E5, core.WHITE, ctx, true)
	})

	// =========================================================================
	// Color filtering — only pieces of the queried color count.
	// =========================================================================

	t.Run("a knight of the wrong color is ignored", func(t *testing.T) {
		ctx := withBoard(func(b *core.Board) {
			b[core.D7] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.BLACK})
		})
		assertAttack(t, core.E5, core.WHITE, ctx, false)
	})

	t.Run("a king of the wrong color is ignored", func(t *testing.T) {
		ctx := withBoard(func(b *core.Board) {
			b[core.E6] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
		})
		assertAttack(t, core.E5, core.WHITE, ctx, false)
	})

	t.Run("a bishop of the wrong color is ignored", func(t *testing.T) {
		ctx := withBoard(func(b *core.Board) {
			b[core.B2] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.BLACK})
		})
		assertAttack(t, core.E5, core.WHITE, ctx, false)
	})

	t.Run("a rook of the wrong color is ignored", func(t *testing.T) {
		ctx := withBoard(func(b *core.Board) {
			b[core.E2] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
		})
		assertAttack(t, core.E5, core.WHITE, ctx, false)
	})

	// =========================================================================
	// Combinations — multiple attackers.
	// =========================================================================

	t.Run("multiple attackers of the same color: any one clear attack suffices", func(t *testing.T) {
		// Knight on D7 (L-shape) + rook on E2 (file). Either alone would attack.
		ctx := withBoard(func(b *core.Board) {
			b[core.D7] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
			b[core.E2] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
		})
		assertAttack(t, core.E5, core.WHITE, ctx, true)
	})

	t.Run("multiple enemy-color pieces, none matching the queried color, do not attack", func(t *testing.T) {
		ctx := withBoard(func(b *core.Board) {
			b[core.D7] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.BLACK})
			b[core.E2] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
		})
		assertAttack(t, core.E5, core.WHITE, ctx, false)
	})

	t.Run("mixed colors: a matching-color attacker among wrong-color pieces still attacks", func(t *testing.T) {
		// Black knight on D7 (wrong color), white rook on E2 (right color).
		ctx := withBoard(func(b *core.Board) {
			b[core.D7] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.BLACK})
			b[core.E2] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
		})
		assertAttack(t, core.E5, core.WHITE, ctx, true)
	})

	t.Run("only friendly pieces on the board do not attack (for the enemy color)", func(t *testing.T) {
		ctx := withBoard(func(b *core.Board) {
			b[core.D7] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.BLACK})
			b[core.E2] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
		})
		// Asking "does WHITE attack E5?" with only black pieces → false.
		assertAttack(t, core.E5, core.WHITE, ctx, false)
		// Asking "does BLACK attack E5?" with the same board → true.
		assertAttack(t, core.E5, core.BLACK, ctx, true)
	})

	// =========================================================================
	// Slider blocking edge cases — position of the blocker relative to the
	// target and the attacker.
	// =========================================================================

	t.Run("a blocker immediately adjacent to the target blocks the slider", func(t *testing.T) {
		// Bishop on H8 → G7 → F6 (own pawn) → E5 (target). The H8-E5 diagonal
		// runs through G7 and F6; F6 is the square adjacent to E5 on that ray.
		ctx := withBoard(func(b *core.Board) {
			b[core.H8] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
			b[core.F6] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		})
		assertAttack(t, core.E5, core.WHITE, ctx, false)
	})

	t.Run("a blocker immediately adjacent to the attacker blocks the slider", func(t *testing.T) {
		// Bishop on H8 → G7 (own pawn) → ... → E5. G7 is adjacent to H8.
		ctx := withBoard(func(b *core.Board) {
			b[core.H8] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
			b[core.G7] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		})
		assertAttack(t, core.E5, core.WHITE, ctx, false)
	})

	t.Run("a piece behind the target does not block (target is between attacker and piece)", func(t *testing.T) {
		// Bishop on H8 → E5 (target) → C3 (own pawn). The pawn is past E5.
		ctx := withBoard(func(b *core.Board) {
			b[core.H8] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
			b[core.C3] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		})
		assertAttack(t, core.E5, core.WHITE, ctx, true)
	})

	// =========================================================================
	// Edge cases — empty board, target on a corner/edge.
	// =========================================================================

	t.Run("an empty board reports no attack", func(t *testing.T) {
		ctx := withBoard(func(b *core.Board) {})
		assertAttack(t, core.E5, core.WHITE, ctx, false)
	})

	t.Run("a corner target with no attackers reports no attack", func(t *testing.T) {
		ctx := withBoard(func(b *core.Board) {})
		assertAttack(t, core.A1, core.WHITE, ctx, false)
	})

	t.Run("a corner target is attacked by a rook on the opposite end of its file", func(t *testing.T) {
		ctx := withBoard(func(b *core.Board) {
			b[core.A8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
		})
		assertAttack(t, core.A1, core.WHITE, ctx, true)
	})

	t.Run("a corner target is attacked by a bishop on the opposite corner (full diagonal)", func(t *testing.T) {
		ctx := withBoard(func(b *core.Board) {
			b[core.H8] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
		})
		assertAttack(t, core.A1, core.WHITE, ctx, true)
	})

	t.Run("a corner target is attacked by a knight on its only L-shape", func(t *testing.T) {
		// A1's knight-attackers are B3 and C2.
		ctx := withBoard(func(b *core.Board) {
			b[core.B3] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
		})
		assertAttack(t, core.A1, core.WHITE, ctx, true)
	})

	t.Run("an edge target is attacked by a king on an adjacent edge square", func(t *testing.T) {
		// A4's neighbors: A3, A5, B3, B4, B5.
		ctx := withBoard(func(b *core.Board) {
			b[core.A3] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
		})
		assertAttack(t, core.A4, core.WHITE, ctx, true)
	})

	t.Run("an edge target is attacked by a pawn on the only valid diagonal", func(t *testing.T) {
		// A4 (target). White pawn must be on rank 3: B3 attacks A4.
		ctx := withBoard(func(b *core.Board) {
			b[core.B3] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		})
		assertAttack(t, core.A4, core.WHITE, ctx, true)
	})

	// =========================================================================
	// The target square's own occupancy doesn't affect the scan — IsSquareAttacked
	// only looks at squares AROUND the target, not the target itself.
	// =========================================================================

	t.Run("a piece sitting on the target square does not prevent attacks on it", func(t *testing.T) {
		// A black king on E5 (the target), and a white rook on E2 attacking it.
		// IsSquareAttacked scans outward from E5; the rook on E2 is on a clear
		// file, so it attacks E5 regardless of what sits on E5.
		ctx := withBoard(func(b *core.Board) {
			b[core.E5] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
			b[core.E2] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
		})
		assertAttack(t, core.E5, core.WHITE, ctx, true)
	})
}
