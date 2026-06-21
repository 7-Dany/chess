package engine

import (
	"testing"

	"github.com/7-Dany/chess/core"
	"github.com/7-Dany/chess/testutil"
)

// TestApply verifies that Apply correctly mutates the board, en passant
// target, castling rights, and king position for every move type.
//
// Each subtest is a single, named scenario. Read the test names top to
// bottom — they describe exactly what behavior each case checks.
func TestApply(t *testing.T) {
	engine := GetDefaultEngine()

	// =========================================================================
	// Normal moves — piece moves from A to B, no special rules triggered.
	// =========================================================================

	t.Run("normal knight move relocates the piece from B1 to C3", func(t *testing.T) {
		var board core.Board
		board[core.B1] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})

		ctx := testutil.NewTurn(&board, core.WHITE)

		engine.Apply(ctx, core.Move{
			Type:  core.NORMAL,
			Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE},
			From:  core.B1,
			To:    core.C3,
		})

		testutil.AssertSquareEmpty(t, &board, core.B1)
		testutil.AssertSquareHas(t, &board, core.C3, core.KNIGHT, core.WHITE)
		// A normal move leaves castling rights and en passant untouched.
		if ctx.Sides != testutil.DefaultSides() {
			t.Errorf("castling rights should be unchanged, got %+v", ctx.Sides)
		}
		if ctx.EnPassantTarget != core.NoPosition {
			t.Errorf("en passant should be cleared, got %v", ctx.EnPassantTarget)
		}
	})

	t.Run("normal king move updates king position and clears both castling rights", func(t *testing.T) {
		var board core.Board
		board[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})

		ctx := testutil.NewTurn(&board, core.WHITE)

		engine.Apply(ctx, core.Move{
			Type:  core.NORMAL,
			Piece: core.Piece{Type: core.KING, Color: core.WHITE},
			From:  core.E1,
			To:    core.F1,
		})

		testutil.AssertSquareEmpty(t, &board, core.E1)
		testutil.AssertSquareHas(t, &board, core.F1, core.KING, core.WHITE)

		// King moved: position updated to F1, both castling rights revoked.
		if ctx.Sides[core.WHITE].KingPosition != core.F1 {
			t.Errorf("white king position = %v, want F1", ctx.Sides[core.WHITE].KingPosition)
		}
		if ctx.Sides[core.WHITE].CanCastleKingSide || ctx.Sides[core.WHITE].CanCastleQueenSide {
			t.Errorf("white castling rights should be cleared after king move")
		}
	})

	t.Run("rook move from A1 clears only the queen-side castling right", func(t *testing.T) {
		var board core.Board
		board[core.A1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})

		ctx := testutil.NewTurn(&board, core.WHITE)

		engine.Apply(ctx, core.Move{
			Type:  core.NORMAL,
			Piece: core.Piece{Type: core.ROOK, Color: core.WHITE},
			From:  core.A1,
			To:    core.A3,
		})

		testutil.AssertSquareHas(t, &board, core.A3, core.ROOK, core.WHITE)
		// A-file rook moved: queen-side right lost, king-side preserved.
		if ctx.Sides[core.WHITE].CanCastleQueenSide {
			t.Errorf("queen-side right should be cleared after A1 rook move")
		}
		if !ctx.Sides[core.WHITE].CanCastleKingSide {
			t.Errorf("king-side right should be preserved")
		}
	})

	t.Run("rook move from H1 clears only the king-side castling right", func(t *testing.T) {
		var board core.Board
		board[core.H1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})

		ctx := testutil.NewTurn(&board, core.WHITE)

		engine.Apply(ctx, core.Move{
			Type:  core.NORMAL,
			Piece: core.Piece{Type: core.ROOK, Color: core.WHITE},
			From:  core.H1,
			To:    core.H3,
		})

		testutil.AssertSquareHas(t, &board, core.H3, core.ROOK, core.WHITE)
		// H-file rook moved: king-side right lost, queen-side preserved.
		if ctx.Sides[core.WHITE].CanCastleKingSide {
			t.Errorf("king-side right should be cleared after H1 rook move")
		}
		if !ctx.Sides[core.WHITE].CanCastleQueenSide {
			t.Errorf("queen-side right should be preserved")
		}
	})

	t.Run("rook move from a non-home file preserves all castling rights", func(t *testing.T) {
		var board core.Board
		board[core.C3] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})

		ctx := testutil.NewTurn(&board, core.WHITE)

		engine.Apply(ctx, core.Move{
			Type:  core.NORMAL,
			Piece: core.Piece{Type: core.ROOK, Color: core.WHITE},
			From:  core.C3,
			To:    core.C5,
		})

		testutil.AssertSquareHas(t, &board, core.C5, core.ROOK, core.WHITE)
		// Rook wasn't on A1 or H1: no castling rights affected.
		if ctx.Sides[core.WHITE] != testutil.DefaultSides()[0] {
			t.Errorf("white side state should be unchanged, got %+v", ctx.Sides[core.WHITE])
		}
	})

	t.Run("black rook move from A8 clears black's queen-side castling right", func(t *testing.T) {
		var board core.Board
		board[core.A8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})

		ctx := testutil.NewTurn(&board, core.BLACK)

		engine.Apply(ctx, core.Move{
			Type:  core.NORMAL,
			Piece: core.Piece{Type: core.ROOK, Color: core.BLACK},
			From:  core.A8,
			To:    core.A6,
		})

		testutil.AssertSquareHas(t, &board, core.A6, core.ROOK, core.BLACK)
		if ctx.Sides[core.BLACK].CanCastleQueenSide {
			t.Errorf("black queen-side right should be cleared after A8 rook move")
		}
		if !ctx.Sides[core.BLACK].CanCastleKingSide {
			t.Errorf("black king-side right should be preserved")
		}
	})

	t.Run("black rook move from H8 clears black's king-side castling right", func(t *testing.T) {
		var board core.Board
		board[core.H8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})

		ctx := testutil.NewTurn(&board, core.BLACK)

		engine.Apply(ctx, core.Move{
			Type:  core.NORMAL,
			Piece: core.Piece{Type: core.ROOK, Color: core.BLACK},
			From:  core.H8,
			To:    core.H6,
		})

		testutil.AssertSquareHas(t, &board, core.H6, core.ROOK, core.BLACK)
		if ctx.Sides[core.BLACK].CanCastleKingSide {
			t.Errorf("black king-side right should be cleared after H8 rook move")
		}
		if !ctx.Sides[core.BLACK].CanCastleQueenSide {
			t.Errorf("black queen-side right should be preserved")
		}
	})

	// =========================================================================
	// Captures — the moving piece lands on an enemy piece and removes it.
	// Capturing a rook on its home square also forfeits that castling right.
	// =========================================================================

	t.Run("capture replaces the enemy piece on the destination square", func(t *testing.T) {
		var board core.Board
		board[core.E4] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
		board[core.D5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})

		ctx := testutil.NewTurn(&board, core.WHITE)

		engine.Apply(ctx, core.Move{
			Type:       core.NORMAL,
			Piece:      core.Piece{Type: core.KNIGHT, Color: core.WHITE},
			From:       core.E4,
			To:         core.D5,
			HasCapture: true,
			Captured:   core.Piece{Type: core.PAWN, Color: core.BLACK},
		})

		testutil.AssertSquareEmpty(t, &board, core.E4)
		testutil.AssertSquareHas(t, &board, core.D5, core.KNIGHT, core.WHITE)
	})

	t.Run("capturing a rook on A8 clears the opponent's queen-side right", func(t *testing.T) {
		var board core.Board
		board[core.A6] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
		board[core.A8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})

		ctx := testutil.NewTurn(&board, core.WHITE)

		engine.Apply(ctx, core.Move{
			Type:       core.NORMAL,
			Piece:      core.Piece{Type: core.BISHOP, Color: core.WHITE},
			From:       core.A6,
			To:         core.A8,
			HasCapture: true,
			Captured:   core.Piece{Type: core.ROOK, Color: core.BLACK},
		})

		testutil.AssertSquareHas(t, &board, core.A8, core.BISHOP, core.WHITE)
		// The captured rook was on A8 (black's queen-side home): that right is lost.
		if ctx.Sides[core.BLACK].CanCastleQueenSide {
			t.Errorf("black queen-side right should be cleared after rook captured on A8")
		}
		if !ctx.Sides[core.BLACK].CanCastleKingSide {
			t.Errorf("black king-side right should be preserved")
		}
	})

	t.Run("capturing a rook on H1 clears the opponent's king-side right", func(t *testing.T) {
		var board core.Board
		board[core.H3] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.BLACK})
		board[core.H1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})

		ctx := testutil.NewTurn(&board, core.BLACK)

		engine.Apply(ctx, core.Move{
			Type:       core.NORMAL,
			Piece:      core.Piece{Type: core.BISHOP, Color: core.BLACK},
			From:       core.H3,
			To:         core.H1,
			HasCapture: true,
			Captured:   core.Piece{Type: core.ROOK, Color: core.WHITE},
		})

		testutil.AssertSquareHas(t, &board, core.H1, core.BISHOP, core.BLACK)
		// The captured rook was on H1 (white's king-side home): that right is lost.
		if ctx.Sides[core.WHITE].CanCastleKingSide {
			t.Errorf("white king-side right should be cleared after rook captured on H1")
		}
		if !ctx.Sides[core.WHITE].CanCastleQueenSide {
			t.Errorf("white queen-side right should be preserved")
		}
	})

	t.Run("capturing a non-rook piece does not affect castling rights", func(t *testing.T) {
		var board core.Board
		board[core.B5] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
		board[core.A6] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})

		ctx := testutil.NewTurn(&board, core.WHITE)

		engine.Apply(ctx, core.Move{
			Type:       core.NORMAL,
			Piece:      core.Piece{Type: core.KNIGHT, Color: core.WHITE},
			From:       core.B5,
			To:         core.A6,
			HasCapture: true,
			Captured:   core.Piece{Type: core.PAWN, Color: core.BLACK},
		})

		testutil.AssertSquareHas(t, &board, core.A6, core.KNIGHT, core.WHITE)
		// Captured piece was a pawn, not a rook: rights unchanged.
		if ctx.Sides != testutil.DefaultSides() {
			t.Errorf("castling rights should be unchanged, got %+v", ctx.Sides)
		}
	})

	// =========================================================================
	// En passant — pawn captures diagonally onto an empty square, removing
	// the enemy pawn that sits beside the destination (not on it).
	// =========================================================================

	t.Run("white en passant capture lands on E6 and removes the black pawn on E5", func(t *testing.T) {
		var board core.Board
		board[core.D5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.E5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})

		ctx := testutil.NewTurn(&board, core.WHITE, testutil.WithEnPassantTarget(core.E6))

		engine.Apply(ctx, core.Move{
			Type:       core.EN_PASSANT,
			Piece:      core.Piece{Type: core.PAWN, Color: core.WHITE},
			From:       core.D5,
			To:         core.E6,
			HasCapture: true,
			Captured:   core.Piece{Type: core.PAWN, Color: core.BLACK},
		})

		// Pawn moved to E6; both D5 (origin) and E5 (captured pawn) are empty.
		testutil.AssertSquareHas(t, &board, core.E6, core.PAWN, core.WHITE)
		testutil.AssertSquareEmpty(t, &board, core.D5)
		testutil.AssertSquareEmpty(t, &board, core.E5)
		if ctx.EnPassantTarget != core.NoPosition {
			t.Errorf("en passant target should be cleared after the capture")
		}
	})

	t.Run("black en passant capture lands on E3 and removes the white pawn on E4", func(t *testing.T) {
		var board core.Board
		board[core.E4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.D4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})

		ctx := testutil.NewTurn(&board, core.BLACK, testutil.WithEnPassantTarget(core.E3))

		engine.Apply(ctx, core.Move{
			Type:       core.EN_PASSANT,
			Piece:      core.Piece{Type: core.PAWN, Color: core.BLACK},
			From:       core.D4,
			To:         core.E3,
			HasCapture: true,
			Captured:   core.Piece{Type: core.PAWN, Color: core.WHITE},
		})

		testutil.AssertSquareHas(t, &board, core.E3, core.PAWN, core.BLACK)
		testutil.AssertSquareEmpty(t, &board, core.D4)
		testutil.AssertSquareEmpty(t, &board, core.E4)
	})

	t.Run("en passant on the A file does not clear castling rights", func(t *testing.T) {
		var board core.Board
		board[core.B5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.A5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})

		ctx := testutil.NewTurn(&board, core.WHITE, testutil.WithEnPassantTarget(core.A6))

		engine.Apply(ctx, core.Move{
			Type:       core.EN_PASSANT,
			Piece:      core.Piece{Type: core.PAWN, Color: core.WHITE},
			From:       core.B5,
			To:         core.A6,
			HasCapture: true,
			Captured:   core.Piece{Type: core.PAWN, Color: core.BLACK},
		})

		testutil.AssertSquareHas(t, &board, core.A6, core.PAWN, core.WHITE)
		// En passant never affects castling rights.
		if ctx.Sides != testutil.DefaultSides() {
			t.Errorf("castling rights should be unchanged, got %+v", ctx.Sides)
		}
	})

	// =========================================================================
	// Promotion — pawn reaching the last rank becomes a different piece.
	// A promotion can also be a capture (landing on an enemy piece).
	// =========================================================================

	t.Run("white pawn promotes to queen when reaching E8", func(t *testing.T) {
		var board core.Board
		board[core.E7] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})

		ctx := testutil.NewTurn(&board, core.WHITE)

		engine.Apply(ctx, core.Move{
			Type:      core.PROMOTION,
			Piece:     core.Piece{Type: core.PAWN, Color: core.WHITE},
			From:      core.E7,
			To:        core.E8,
			PromoteTo: core.QUEEN,
		})

		testutil.AssertSquareEmpty(t, &board, core.E7)
		testutil.AssertSquareHas(t, &board, core.E8, core.QUEEN, core.WHITE)
	})

	t.Run("black pawn promotes to knight when reaching D1", func(t *testing.T) {
		var board core.Board
		board[core.D2] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})

		ctx := testutil.NewTurn(&board, core.BLACK)

		engine.Apply(ctx, core.Move{
			Type:      core.PROMOTION,
			Piece:     core.Piece{Type: core.PAWN, Color: core.BLACK},
			From:      core.D2,
			To:        core.D1,
			PromoteTo: core.KNIGHT,
		})

		testutil.AssertSquareEmpty(t, &board, core.D2)
		testutil.AssertSquareHas(t, &board, core.D1, core.KNIGHT, core.BLACK)
	})

	t.Run("promotion with capture on a non-home file preserves castling rights", func(t *testing.T) {
		var board core.Board
		board[core.E7] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.D8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})

		ctx := testutil.NewTurn(&board, core.WHITE)

		engine.Apply(ctx, core.Move{
			Type:       core.PROMOTION,
			Piece:      core.Piece{Type: core.PAWN, Color: core.WHITE},
			From:       core.E7,
			To:         core.D8,
			PromoteTo:  core.QUEEN,
			HasCapture: true,
			Captured:   core.Piece{Type: core.ROOK, Color: core.BLACK},
		})

		// The captured rook was on D8, not A8 or H8: no castling right affected.
		testutil.AssertSquareHas(t, &board, core.D8, core.QUEEN, core.WHITE)
		if ctx.Sides[core.BLACK].CanCastleQueenSide != true || ctx.Sides[core.BLACK].CanCastleKingSide != true {
			t.Errorf("black castling rights should be unchanged, got %+v", ctx.Sides[core.BLACK])
		}
	})

	t.Run("promotion capture on H8 clears the opponent's king-side right", func(t *testing.T) {
		var board core.Board
		board[core.G7] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.H8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})

		ctx := testutil.NewTurn(&board, core.WHITE)

		engine.Apply(ctx, core.Move{
			Type:       core.PROMOTION,
			Piece:      core.Piece{Type: core.PAWN, Color: core.WHITE},
			From:       core.G7,
			To:         core.H8,
			PromoteTo:  core.QUEEN,
			HasCapture: true,
			Captured:   core.Piece{Type: core.ROOK, Color: core.BLACK},
		})

		testutil.AssertSquareHas(t, &board, core.H8, core.QUEEN, core.WHITE)
		// Captured rook was on H8 (black's king-side home): that right is lost.
		if ctx.Sides[core.BLACK].CanCastleKingSide {
			t.Errorf("black king-side right should be cleared after rook captured on H8")
		}
		if !ctx.Sides[core.BLACK].CanCastleQueenSide {
			t.Errorf("black queen-side right should be preserved")
		}
	})

	// =========================================================================
	// Castling — king moves two squares toward a rook; the rook hops over
	// to the square the king crossed. All castling rights are forfeited.
	// =========================================================================

	t.Run("white king-side castling moves king to G1 and rook from H1 to F1", func(t *testing.T) {
		var board core.Board
		board[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
		board[core.H1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})

		ctx := testutil.NewTurn(&board, core.WHITE)

		engine.Apply(ctx, core.Move{
			Type:  core.CASTLING,
			Piece: core.Piece{Type: core.KING, Color: core.WHITE},
			From:  core.E1,
			To:    core.G1,
		})

		testutil.AssertSquareHas(t, &board, core.G1, core.KING, core.WHITE)
		testutil.AssertSquareHas(t, &board, core.F1, core.ROOK, core.WHITE)
		testutil.AssertSquareEmpty(t, &board, core.E1)
		testutil.AssertSquareEmpty(t, &board, core.H1)
		// Castling forfeits all rights for the castling side.
		if ctx.Sides[core.WHITE].KingPosition != core.G1 {
			t.Errorf("white king position = %v, want G1", ctx.Sides[core.WHITE].KingPosition)
		}
		if ctx.Sides[core.WHITE].CanCastleKingSide || ctx.Sides[core.WHITE].CanCastleQueenSide {
			t.Errorf("white castling rights should be cleared after castling")
		}
	})

	t.Run("white queen-side castling moves king to C1 and rook from A1 to D1", func(t *testing.T) {
		var board core.Board
		board[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
		board[core.A1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})

		ctx := testutil.NewTurn(&board, core.WHITE)

		engine.Apply(ctx, core.Move{
			Type:  core.CASTLING,
			Piece: core.Piece{Type: core.KING, Color: core.WHITE},
			From:  core.E1,
			To:    core.C1,
		})

		testutil.AssertSquareHas(t, &board, core.C1, core.KING, core.WHITE)
		testutil.AssertSquareHas(t, &board, core.D1, core.ROOK, core.WHITE)
		testutil.AssertSquareEmpty(t, &board, core.E1)
		testutil.AssertSquareEmpty(t, &board, core.A1)
		if ctx.Sides[core.WHITE].KingPosition != core.C1 {
			t.Errorf("white king position = %v, want C1", ctx.Sides[core.WHITE].KingPosition)
		}
	})

	t.Run("black king-side castling moves king to G8 and rook from H8 to F8", func(t *testing.T) {
		var board core.Board
		board[core.E8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
		board[core.H8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})

		ctx := testutil.NewTurn(&board, core.BLACK)

		engine.Apply(ctx, core.Move{
			Type:  core.CASTLING,
			Piece: core.Piece{Type: core.KING, Color: core.BLACK},
			From:  core.E8,
			To:    core.G8,
		})

		testutil.AssertSquareHas(t, &board, core.G8, core.KING, core.BLACK)
		testutil.AssertSquareHas(t, &board, core.F8, core.ROOK, core.BLACK)
		testutil.AssertSquareEmpty(t, &board, core.E8)
		testutil.AssertSquareEmpty(t, &board, core.H8)
		if ctx.Sides[core.BLACK].KingPosition != core.G8 {
			t.Errorf("black king position = %v, want G8", ctx.Sides[core.BLACK].KingPosition)
		}
	})

	t.Run("black queen-side castling moves king to C8 and rook from A8 to D8", func(t *testing.T) {
		var board core.Board
		board[core.E8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
		board[core.A8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})

		ctx := testutil.NewTurn(&board, core.BLACK)

		engine.Apply(ctx, core.Move{
			Type:  core.CASTLING,
			Piece: core.Piece{Type: core.KING, Color: core.BLACK},
			From:  core.E8,
			To:    core.C8,
		})

		testutil.AssertSquareHas(t, &board, core.C8, core.KING, core.BLACK)
		testutil.AssertSquareHas(t, &board, core.D8, core.ROOK, core.BLACK)
		testutil.AssertSquareEmpty(t, &board, core.E8)
		testutil.AssertSquareEmpty(t, &board, core.A8)
		if ctx.Sides[core.BLACK].KingPosition != core.C8 {
			t.Errorf("black king position = %v, want C8", ctx.Sides[core.BLACK].KingPosition)
		}
	})

	// =========================================================================
	// En passant target — a double pawn push sets the EP target; any other
	// move clears it.
	// =========================================================================

	t.Run("white pawn double push from E2 to E4 sets en passant target on E3", func(t *testing.T) {
		var board core.Board
		board[core.E2] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})

		ctx := testutil.NewTurn(&board, core.WHITE)

		engine.Apply(ctx, core.Move{
			Type:  core.NORMAL,
			Piece: core.Piece{Type: core.PAWN, Color: core.WHITE},
			From:  core.E2,
			To:    core.E4,
		})

		testutil.AssertSquareHas(t, &board, core.E4, core.PAWN, core.WHITE)
		if ctx.EnPassantTarget != core.E3 {
			t.Errorf("en passant target = %v, want E3", ctx.EnPassantTarget)
		}
	})

	t.Run("black pawn double push from D7 to D5 sets en passant target on D6", func(t *testing.T) {
		var board core.Board
		board[core.D7] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})

		ctx := testutil.NewTurn(&board, core.BLACK)

		engine.Apply(ctx, core.Move{
			Type:  core.NORMAL,
			Piece: core.Piece{Type: core.PAWN, Color: core.BLACK},
			From:  core.D7,
			To:    core.D5,
		})

		testutil.AssertSquareHas(t, &board, core.D5, core.PAWN, core.BLACK)
		if ctx.EnPassantTarget != core.D6 {
			t.Errorf("en passant target = %v, want D6", ctx.EnPassantTarget)
		}
	})

	t.Run("white pawn double push from A2 to A4 sets en passant target on A3", func(t *testing.T) {
		var board core.Board
		board[core.A2] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})

		ctx := testutil.NewTurn(&board, core.WHITE)

		engine.Apply(ctx, core.Move{
			Type:  core.NORMAL,
			Piece: core.Piece{Type: core.PAWN, Color: core.WHITE},
			From:  core.A2,
			To:    core.A4,
		})

		testutil.AssertSquareHas(t, &board, core.A4, core.PAWN, core.WHITE)
		if ctx.EnPassantTarget != core.A3 {
			t.Errorf("en passant target = %v, want A3", ctx.EnPassantTarget)
		}
	})

	t.Run("single pawn push clears the previous en passant target", func(t *testing.T) {
		var board core.Board
		board[core.E2] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})

		// There was a prior double push (EP target on D3); this single push must clear it.
		ctx := testutil.NewTurn(&board, core.WHITE, testutil.WithEnPassantTarget(core.D3))

		engine.Apply(ctx, core.Move{
			Type:  core.NORMAL,
			Piece: core.Piece{Type: core.PAWN, Color: core.WHITE},
			From:  core.E2,
			To:    core.E3,
		})

		testutil.AssertSquareHas(t, &board, core.E3, core.PAWN, core.WHITE)
		if ctx.EnPassantTarget != core.NoPosition {
			t.Errorf("en passant target should be cleared, got %v", ctx.EnPassantTarget)
		}
	})

	t.Run("non-pawn move clears the previous en passant target", func(t *testing.T) {
		var board core.Board
		board[core.B1] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})

		ctx := testutil.NewTurn(&board, core.WHITE, testutil.WithEnPassantTarget(core.E3))

		engine.Apply(ctx, core.Move{
			Type:  core.NORMAL,
			Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE},
			From:  core.B1,
			To:    core.C3,
		})

		if ctx.EnPassantTarget != core.NoPosition {
			t.Errorf("en passant target should be cleared, got %v", ctx.EnPassantTarget)
		}
		// The snapshot returned by Apply must carry the pre-move EP target
		// so Undo can restore it.
		// (Verified separately in undo_test.go's round-trip tests.)
	})

	t.Run("Apply does not flip the side to move", func(t *testing.T) {
		var board core.Board
		board[core.B1] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})

		ctx := testutil.NewTurn(&board, core.WHITE, testutil.WithEnPassantTarget(core.E3))

		engine.Apply(ctx, core.Move{
			Type:  core.NORMAL,
			Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE},
			From:  core.B1,
			To:    core.C3,
		})

		// Apply mutates the board and game state, but leaves SideToMove alone.
		// The game controller (not the engine) is responsible for flipping turns.
		if ctx.SideToMove != core.WHITE {
			t.Errorf("SideToMove = %v, want WHITE (Apply should not flip side)", ctx.SideToMove)
		}
	})
}
