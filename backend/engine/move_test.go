package engine

import (
	"testing"

	"github.com/7-Dany/chess/core"
	"github.com/7-Dany/chess/fen"
	"github.com/7-Dany/chess/testutil"
)

// TestGetLegalMoves verifies that GetLegalMoves filters pseudo-legal moves
// for king safety — only moves that do not leave the moving side's king
// attacked are returned.
//
// The king-safety filter works by applying each pseudo-legal move, checking
// if the king is attacked, then undoing. If the moving piece is the king,
// the king's position is taken from move.To (not the original square).
//
// This test focuses on the filtering logic — pin scenarios, check escape,
// castling legality, en passant edge cases, and promotion filtering. The
// raw move generation is tested in the piece package.
func TestGetLegalMoves(t *testing.T) {
	engine := GetDefaultEngine()

	// Helper: run GetLegalMoves on the given position and return the moves.
	legalMoves := func(board *core.Board, side core.PieceColor, pos core.Position, opts ...testutil.TurnOption) []core.Move {
		ctx := testutil.NewTurn(board, side, opts...)
		return engine.GetLegalMoves(make([]core.Move, 0, core.MAX_MOVES), pos, *ctx)
	}

	// =========================================================================
	// Pins — a piece pinned to its king can only move along the pin line.
	// =========================================================================

	t.Run("a pinned rook can only move along the pin line (forward/back, not sideways)", func(t *testing.T) {
		// White king on E1, white rook on E2, black rook on E8 (pinning down the E-file).
		// The rook can move along the E-file (E3, capturing E8) but not sideways (D2, F2).
		var board core.Board
		board[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
		board[core.E2] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
		board[core.E8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})

		moves := legalMoves(&board, core.WHITE, core.E2)

		testutil.AssertMovePresent(t, moves, core.E2, core.E3) // forward along pin line
		testutil.AssertMovePresent(t, moves, core.E2, core.E8) // capture the pinner
		testutil.AssertMoveAbsent(t, moves, core.E2, core.D2)  // sideways — illegal (exposes king)
		testutil.AssertMoveAbsent(t, moves, core.E2, core.F2)  // sideways — illegal (exposes king)
		testutil.AssertMoveCount(t, moves, 6)                  // E3..E7 (5) + E8 capture (1) = 6
	})

	t.Run("a pinned bishop has no legal moves (can't move along a rook pin line)", func(t *testing.T) {
		// White king on E1, white bishop on E2, black rook on E8 (pinning on the E-file).
		// A bishop can only move diagonally — no diagonal move stays on the E-file, so
		// every move exposes the king.
		var board core.Board
		board[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
		board[core.E2] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
		board[core.E8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
		board[core.H8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})

		moves := legalMoves(&board, core.WHITE, core.E2, testutil.WithSides(
			testutil.Side(core.E1, false, false),
			testutil.Side(core.H8, true, true),
		))

		testutil.AssertNoMoves(t, moves)
	})

	// =========================================================================
	// King in check — must escape. Can move away, capture the checker, or
	// block (if not a king move).
	// =========================================================================

	t.Run("a king in check on the E-file must escape sideways (cannot stay on the file)", func(t *testing.T) {
		// White king on E1, black rook on E8 gives check down the E-file.
		var board core.Board
		board[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
		board[core.E8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})

		moves := legalMoves(&board, core.WHITE, core.E1)

		// Can move sideways (D1, F1) or diagonally off the file (D2, F2).
		testutil.AssertMovePresent(t, moves, core.E1, core.D1)
		testutil.AssertMovePresent(t, moves, core.E1, core.F1)
		testutil.AssertMovePresent(t, moves, core.E1, core.D2)
		testutil.AssertMovePresent(t, moves, core.E1, core.F2)
		// Cannot move up the file (still in check from the rook).
		testutil.AssertMoveAbsent(t, moves, core.E1, core.E2)
		// Cannot castle out of check.
		testutil.AssertMoveAbsent(t, moves, core.E1, core.G1) // no castling
		testutil.AssertMoveCount(t, moves, 4)
	})

	t.Run("a king in check on rank 1 cannot stay on the rank", func(t *testing.T) {
		// White king on E1, black rook on A1 gives check along rank 1.
		var board core.Board
		board[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
		board[core.A1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})

		moves := legalMoves(&board, core.WHITE, core.E1)

		// Can only escape to rank 2 (D2, E2, F2). D1 and F1 are still on the attacked rank.
		testutil.AssertMovePresent(t, moves, core.E1, core.D2)
		testutil.AssertMovePresent(t, moves, core.E1, core.E2)
		testutil.AssertMovePresent(t, moves, core.E1, core.F2)
		testutil.AssertMoveAbsent(t, moves, core.E1, core.D1) // still on rank 1
		testutil.AssertMoveAbsent(t, moves, core.E1, core.F1) // still on rank 1
		testutil.AssertMoveCount(t, moves, 3)
	})

	t.Run("a king can capture an undefended checker", func(t *testing.T) {
		// White king on E4, black rook on E5 (giving check). The rook is
		// undefended, so the king can capture it.
		var board core.Board
		board[core.E4] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
		board[core.E5] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
		board[core.E8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})

		moves := legalMoves(&board, core.WHITE, core.E4, testutil.WithSides(
			testutil.Side(core.E4, false, false),
			testutil.Side(core.E8, true, true),
		))

		testutil.AssertMovePresent(t, moves, core.E4, core.E5) // capture the checker
		// Can also escape to safe squares (D3, D4, F3, F4).
		testutil.AssertMovePresent(t, moves, core.E4, core.D3)
		testutil.AssertMovePresent(t, moves, core.E4, core.D4)
		testutil.AssertMovePresent(t, moves, core.E4, core.F3)
		testutil.AssertMovePresent(t, moves, core.E4, core.F4)
		// Squares still in check (D5, E3, F5 are attacked by... actually these
		// are attacked by the enemy king or still on the check line).
		testutil.AssertMoveAbsent(t, moves, core.E4, core.D5) // adjacent to enemy king on E8? no — E8 is far. D5 is attacked by the rook? No, rook is on E5. Wait — let me think.
		// Actually D5 is NOT attacked by the rook on E5 (rook attacks rank 5 and file E). D5 IS on rank 5 → attacked!
		testutil.AssertMoveAbsent(t, moves, core.E4, core.E3) // still on the E-file (rook attacks it)
		testutil.AssertMoveAbsent(t, moves, core.E4, core.F5) // on rank 5 → attacked by rook
		testutil.AssertMoveCount(t, moves, 5)
	})

	t.Run("a king cannot capture a defended checker", func(t *testing.T) {
		// White king on E4, black rook on E5 (checker), black pawn on D6
		// defends E5. Capturing the rook would move the king to E5, which is
		// defended by the pawn → illegal.
		var board core.Board
		board[core.E4] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
		board[core.E5] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
		board[core.D6] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		board[core.E8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})

		moves := legalMoves(&board, core.WHITE, core.E4, testutil.WithSides(
			testutil.Side(core.E4, false, false),
			testutil.Side(core.E8, true, true),
		))

		testutil.AssertMoveAbsent(t, moves, core.E4, core.E5) // defended → cannot capture
		// The four safe escapes remain.
		testutil.AssertMovePresent(t, moves, core.E4, core.D3)
		testutil.AssertMovePresent(t, moves, core.E4, core.D4)
		testutil.AssertMovePresent(t, moves, core.E4, core.F3)
		testutil.AssertMovePresent(t, moves, core.E4, core.F4)
		testutil.AssertMoveCount(t, moves, 4)
	})

	t.Run("a king cannot capture a defended adjacent enemy piece", func(t *testing.T) {
		// White king on E1, black knight on D2 (adjacent), black pawn on C3
		// defends D2. The king cannot capture the knight because the pawn
		// would then attack the king.
		var board core.Board
		board[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
		board[core.D2] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.BLACK})
		board[core.C3] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		board[core.E8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})

		moves := legalMoves(&board, core.WHITE, core.E1, testutil.WithSides(
			testutil.Side(core.E1, false, false),
			testutil.Side(core.E8, true, true),
		))

		testutil.AssertMoveAbsent(t, moves, core.E1, core.D2) // defended by pawn
		testutil.AssertMoveAbsent(t, moves, core.E1, core.F1) // adjacent to... actually F1 is attacked by the knight on D2 (knight L: D2→F1)? D2 is file 3 rank 1, F1 is file 5 rank 0. Diff (2,1) → yes, knight L. So F1 is attacked.
		testutil.AssertMovePresent(t, moves, core.E1, core.D1)
		testutil.AssertMovePresent(t, moves, core.E1, core.E2)
		testutil.AssertMovePresent(t, moves, core.E1, core.F2)
		testutil.AssertMoveCount(t, moves, 3)
	})

	t.Run("a king cannot move adjacent to the enemy king", func(t *testing.T) {
		// White king on E4, black king on E6. The three squares between them
		// on rank 5 (D5, E5, F5) are adjacent to the black king → illegal.
		var board core.Board
		board[core.E4] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
		board[core.E6] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})

		moves := legalMoves(&board, core.WHITE, core.E4, testutil.WithSides(
			testutil.Side(core.E4, false, false),
			testutil.Side(core.E6, true, true),
		))

		testutil.AssertMovePresent(t, moves, core.E4, core.D3)
		testutil.AssertMovePresent(t, moves, core.E4, core.D4)
		testutil.AssertMovePresent(t, moves, core.E4, core.E3)
		testutil.AssertMovePresent(t, moves, core.E4, core.F3)
		testutil.AssertMovePresent(t, moves, core.E4, core.F4)
		testutil.AssertMoveAbsent(t, moves, core.E4, core.D5) // adjacent to enemy king
		testutil.AssertMoveAbsent(t, moves, core.E4, core.E5) // adjacent to enemy king
		testutil.AssertMoveAbsent(t, moves, core.E4, core.F5) // adjacent to enemy king
		testutil.AssertMoveCount(t, moves, 5)
	})

	// =========================================================================
	// Resolving check by blocking or capturing (non-king moves).
	// =========================================================================

	t.Run("a knight can block check by interposing on the check line", func(t *testing.T) {
		// White king on E1, black rook on E8 (check down E-file), white knight
		// on D6. The knight can interpose on E4 (blocking) or capture the rook
		// on E8. Other knight moves would leave the king in check.
		var board core.Board
		board[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
		board[core.D6] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
		board[core.E8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
		board[core.H8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})

		moves := legalMoves(&board, core.WHITE, core.D6, testutil.WithSides(
			testutil.Side(core.E1, false, false),
			testutil.Side(core.H8, true, true),
		))

		testutil.AssertMovePresent(t, moves, core.D6, core.E4) // block on E4
		testutil.AssertMovePresent(t, moves, core.D6, core.E8) // capture the rook
		testutil.AssertMoveAbsent(t, moves, core.D6, core.B5)  // doesn't resolve check
		testutil.AssertMoveAbsent(t, moves, core.D6, core.F7)  // doesn't resolve check
		testutil.AssertMoveAbsent(t, moves, core.D6, core.C8)  // doesn't resolve check
		testutil.AssertMoveCount(t, moves, 2)
	})

	t.Run("a knight can capture the checker or interpose to resolve check", func(t *testing.T) {
		// White king on E1, black rook on E8 gives check down the E-file.
		// White knight on D6 can capture the rook (E8) or interpose on E4.
		var board core.Board
		board[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
		board[core.D6] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
		board[core.E8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
		board[core.H8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})

		moves := legalMoves(&board, core.WHITE, core.D6, testutil.WithSides(
			testutil.Side(core.E1, false, false),
			testutil.Side(core.H8, true, true),
		))

		testutil.AssertMovePresent(t, moves, core.D6, core.E8) // capture the rook
		testutil.AssertMovePresent(t, moves, core.D6, core.E4) // interpose on E4
		testutil.AssertMoveAbsent(t, moves, core.D6, core.B5)
		testutil.AssertMoveAbsent(t, moves, core.D6, core.F7)
		testutil.AssertMoveAbsent(t, moves, core.D6, core.C8)
		testutil.AssertMoveCount(t, moves, 2)
	})

	// =========================================================================
	// En passant edge case — capturing en passant can expose the king on the
	// rank (the captured pawn and the capturing pawn both leave the rank).
	// =========================================================================

	t.Run("an en passant capture that exposes the king on the rank is illegal", func(t *testing.T) {
		// White king on H5, white pawn on F5, black pawn on E5 (just double-
		// pushed, EP target E6), black rook on A5. The en passant capture
		// (F5→E6) removes BOTH the white pawn from F5 AND the black pawn from
		// E5, leaving the H5 king exposed to the A5 rook along rank 5.
		var board core.Board
		board[core.H5] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
		board[core.F5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.E5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		board[core.A5] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
		board[core.E8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})

		moves := legalMoves(&board, core.WHITE, core.F5, testutil.WithSides(
			testutil.Side(core.H5, false, false),
			testutil.Side(core.E8, true, true),
		), testutil.WithEnPassantTarget(core.E6))

		// The en passant capture (F5→E6) is ILLEGAL — it exposes the king.
		testutil.AssertMoveAbsent(t, moves, core.F5, core.E6)
		// The normal push (F5→F6) is still legal.
		testutil.AssertMovePresent(t, moves, core.F5, core.F6)
		testutil.AssertMoveCount(t, moves, 1)
	})

	t.Run("a promotion push that leaves the king in check is filtered out", func(t *testing.T) {
		// White king on E1, white pawn on D7 (ready to promote), black rook
		// on E8 (giving check down the E-file). The pawn cannot push to D8
		// (doesn't resolve check). But it CAN capture the rook on E8
		// (promotion-capture), which resolves the check.
		var board core.Board
		board[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
		board[core.D7] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.E8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
		board[core.H8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})

		moves := legalMoves(&board, core.WHITE, core.D7, testutil.WithSides(
			testutil.Side(core.E1, false, false),
			testutil.Side(core.H8, true, true),
		))

		// The push to D8 is illegal (doesn't block or capture the checker).
		testutil.AssertMoveAbsent(t, moves, core.D7, core.D8)
		// The capture-promotion on E8 is legal (captures the rook, resolving
		// check). It produces 4 moves (Q, R, B, N) all to E8.
		testutil.AssertMovePresent(t, moves, core.D7, core.E8)
		testutil.AssertMoveCount(t, moves, 4) // 4 promotion-capture variants
	})

	// =========================================================================
	// Castling integration — GetLegalMoves includes castling when eligible,
	// and filters it out when the conditions aren't met.
	// =========================================================================

	t.Run("castling is available when all conditions are met", func(t *testing.T) {
		// King on E1, rooks on A1 and H1, all paths clear, no enemy attackers.
		var board core.Board
		board[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
		board[core.A1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
		board[core.H1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
		board[core.E8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})

		moves := legalMoves(&board, core.WHITE, core.E1)

		// Both castling moves present, plus normal king moves (D1, F1, D2, E2, F2).
		testutil.AssertMovePresent(t, moves, core.E1, core.G1) // king-side
		testutil.AssertMovePresent(t, moves, core.E1, core.C1) // queen-side
		testutil.AssertMovePresent(t, moves, core.E1, core.D1)
		testutil.AssertMovePresent(t, moves, core.E1, core.F1)
		testutil.AssertMoveCount(t, moves, 7)
	})

	t.Run("castling is removed when the king is in check", func(t *testing.T) {
		// King on E1, rooks on A1/H1, but black rook on E8 gives check.
		var board core.Board
		board[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
		board[core.A1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
		board[core.H1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
		board[core.E8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})

		moves := legalMoves(&board, core.WHITE, core.E1)

		testutil.AssertMoveAbsent(t, moves, core.E1, core.G1) // can't castle out of check
		testutil.AssertMoveAbsent(t, moves, core.E1, core.C1)
		testutil.AssertMoveCount(t, moves, 4) // D1, F1, D2, F2 (escape sideways)
	})

	t.Run("king-side castling is removed when F1 is occupied", func(t *testing.T) {
		var board core.Board
		board[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
		board[core.A1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
		board[core.H1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
		board[core.F1] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
		board[core.E8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})

		moves := legalMoves(&board, core.WHITE, core.E1)

		testutil.AssertMoveAbsent(t, moves, core.E1, core.G1)  // F1 blocked
		testutil.AssertMovePresent(t, moves, core.E1, core.C1) // queen-side still available
		testutil.AssertMoveAbsent(t, moves, core.E1, core.F1)  // own bishop there
		testutil.AssertMoveCount(t, moves, 5)
	})

	t.Run("king-side castling is removed when F1 is attacked", func(t *testing.T) {
		// Black rook on F8 attacks F1 (down the file).
		var board core.Board
		board[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
		board[core.A1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
		board[core.H1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
		board[core.F8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
		board[core.E8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})

		moves := legalMoves(&board, core.WHITE, core.E1)

		testutil.AssertMoveAbsent(t, moves, core.E1, core.G1)  // F1 attacked
		testutil.AssertMovePresent(t, moves, core.E1, core.C1) // queen-side still available
		// F1 is attacked → moving there is illegal.
		testutil.AssertMoveAbsent(t, moves, core.E1, core.F1)
		// F2 is also attacked by the rook on F8 (same file).
		testutil.AssertMoveAbsent(t, moves, core.E1, core.F2)
		// D1 is safe (not on the F-file, not on the E-file from E8).
		testutil.AssertMovePresent(t, moves, core.E1, core.D1)
		testutil.AssertMoveCount(t, moves, 4)
	})

	t.Run("no castling when rights have been lost", func(t *testing.T) {
		var board core.Board
		board[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
		board[core.A1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
		board[core.H1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
		board[core.E8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})

		moves := legalMoves(&board, core.WHITE, core.E1, testutil.WithSides(
			testutil.Side(core.E1, false, false), // no castling rights
			testutil.FullBlack(),
		))

		testutil.AssertMoveAbsent(t, moves, core.E1, core.G1)
		testutil.AssertMoveAbsent(t, moves, core.E1, core.C1)
		// Normal king moves still available (5: D1, D2, E2, F2, F1).
		testutil.AssertMoveCount(t, moves, 5)
	})

	t.Run("a non-king piece never generates castling moves", func(t *testing.T) {
		// Rook on A1 with full castling rights — the engine only adds
		// castling when the piece is the king (at KingPosition).
		var board core.Board
		board[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
		board[core.A1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
		board[core.E8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})

		moves := legalMoves(&board, core.WHITE, core.A1)

		// No castling moves from a rook.
		for _, m := range moves {
			if m.Type == core.CASTLING {
				t.Errorf("rook should not generate castling moves, got %v", m)
			}
		}
	})
}

// TestGetAllLegalMoves verifies the accumulation logic of GetAllLegalMoves —
// that it collects every friendly piece's legal moves, excludes enemy moves,
// and handles edge cases (check, checkmate, stalemate, pins, castling).
func TestGetAllLegalMoves(t *testing.T) {
	engine := GetDefaultEngine()

	decode := func(t *testing.T, fenStr string) core.TurnContext {
		t.Helper()
		var ctx core.TurnContext
		if err := fen.GetDefaultFenParser().Decode(fenStr, &ctx); err != nil {
			t.Fatalf("Decode failed: %v", err)
		}
		return ctx
	}

	allLegalMoves := func(ctx core.TurnContext) []core.Move {
		var buf [MAX_TOTAL_MOVES]core.Move
		return engine.GetAllLegalMoves(buf[:0], ctx)
	}

	// =========================================================================
	// Accumulation — only the side-to-move's pieces generate moves.
	// =========================================================================

	t.Run("only the side-to-move's pieces generate moves", func(t *testing.T) {
		ctx := decode(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
		moves := allLegalMoves(ctx)
		for _, m := range moves {
			if m.Piece.Color != core.WHITE {
				t.Errorf("found a move for a black piece (%v) when white is to move", m.Piece)
			}
		}
	})

	t.Run("when black is to move, only black pieces generate moves", func(t *testing.T) {
		ctx := decode(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQkq - 0 1")
		moves := allLegalMoves(ctx)
		for _, m := range moves {
			if m.Piece.Color != core.BLACK {
				t.Errorf("found a move for a white piece (%v) when black is to move", m.Piece)
			}
		}
	})

	t.Run("the total equals the sum of per-piece GetLegalMoves counts", func(t *testing.T) {
		ctx := decode(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
		allMoves := allLegalMoves(ctx)

		var perPieceTotal int
		var scratch [core.MAX_MOVES]core.Move
		for i, square := range ctx.Board {
			if !square.IsOccupiedBy(ctx.SideToMove) {
				continue
			}
			pieceMoves := engine.GetLegalMoves(scratch[:0], core.Position(i), ctx)
			perPieceTotal += len(pieceMoves)
		}
		if len(allMoves) != perPieceTotal {
			t.Errorf("GetAllLegalMoves = %d, sum of per-piece = %d", len(allMoves), perPieceTotal)
		}
	})

	t.Run("no move appears more than once", func(t *testing.T) {
		ctx := decode(t, "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1")
		moves := allLegalMoves(ctx)
		seen := make(map[core.Move]bool, len(moves))
		for _, m := range moves {
			if seen[m] {
				t.Errorf("duplicate move: %v", m)
			}
			seen[m] = true
		}
	})

	// =========================================================================
	// King safety — the filter is applied to every piece's moves.
	// =========================================================================

	t.Run("a king in check only has moves that resolve the check", func(t *testing.T) {
		ctx := decode(t, "4r3/8/8/8/8/8/8/4K3 w - - 0 1")
		moves := allLegalMoves(ctx)
		testutil.AssertMoveCount(t, moves, 4)
		for _, m := range moves {
			if m.Piece.Type != core.KING {
				t.Errorf("expected only king moves, got %v", m.Piece)
			}
			if m.To.File() == core.FILE_E {
				t.Errorf("king move to %v stays on the E-file (still in check)", m.To)
			}
		}
	})

	t.Run("a checkmated side has 0 legal moves", func(t *testing.T) {
		ctx := decode(t, "7k/5Q2/5K2/8/8/8/8/8 b - - 0 1")
		moves := allLegalMoves(ctx)
		testutil.AssertNoMoves(t, moves)
	})

	t.Run("a stalemated side has 0 legal moves", func(t *testing.T) {
		ctx := decode(t, "k7/2Q5/2K5/8/8/8/8/8 b - - 0 1")
		moves := allLegalMoves(ctx)
		testutil.AssertNoMoves(t, moves)
	})

	t.Run("a pinned piece's illegal moves are excluded", func(t *testing.T) {
		ctx := decode(t, "4r3/8/8/8/8/8/4R3/4K3 w - - 0 1")
		moves := allLegalMoves(ctx)
		testutil.AssertMoveCount(t, moves, 10)
		for _, m := range moves {
			if m.Piece.Type == core.ROOK && m.From == core.E2 {
				if m.To == core.D2 || m.To == core.F2 {
					t.Errorf("pinned rook should not move sideways to %v", m.To)
				}
			}
		}
	})

	// =========================================================================
	// Castling and en passant are included.
	// =========================================================================

	t.Run("castling moves are included for the king", func(t *testing.T) {
		ctx := decode(t, "4k3/8/8/8/8/8/8/R3K2R w KQ - 0 1")
		moves := allLegalMoves(ctx)
		testutil.AssertMovePresent(t, moves, core.E1, core.G1)
		testutil.AssertMovePresent(t, moves, core.E1, core.C1)
	})

	t.Run("an en passant capture is included when legal", func(t *testing.T) {
		ctx := decode(t, "rnbqkbnr/ppp1pppp/8/8/3pP3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 2")
		moves := allLegalMoves(ctx)
		found := false
		for _, m := range moves {
			if m.From == core.D4 && m.To == core.E3 && m.Type == core.EN_PASSANT {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("en passant move D4→E3 not found in %d moves", len(moves))
		}
	})

	// =========================================================================
	// Edge cases.
	// =========================================================================

	t.Run("a side with no pieces returns 0 moves", func(t *testing.T) {
		ctx := decode(t, "4k3/8/8/8/8/8/8/8 w - - 0 1")
		moves := allLegalMoves(ctx)
		testutil.AssertNoMoves(t, moves)
	})

	t.Run("a side with only a king returns the king's legal moves", func(t *testing.T) {
		ctx := decode(t, "4k3/8/8/8/8/8/8/4K3 w - - 0 1")
		moves := allLegalMoves(ctx)
		testutil.AssertMoveCount(t, moves, 5)
		for _, m := range moves {
			if m.Piece.Type != core.KING {
				t.Errorf("expected only king moves, got %v", m.Piece)
			}
		}
	})
}

// TestHasAnyLegalMoves verifies the yes/no "can this side move at all?" check.
// It's used for checkmate/stalemate detection: if the side to move has no
// legal moves, the game is over (checkmate if in check, stalemate if not).
func TestHasAnyLegalMoves(t *testing.T) {
	engine := GetDefaultEngine()

	// Helper: build a context and check HasAnyLegalMoves.
	hasMoves := func(board *core.Board, side core.PieceColor, opts ...testutil.TurnOption) bool {
		ctx := testutil.NewTurn(board, side, opts...)
		return engine.HasAnyLegalMoves(*ctx)
	}

	t.Run("a side with at least one legal move returns true", func(t *testing.T) {
		var board core.Board
		board[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
		board[core.B1] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
		board[core.E8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})

		if !hasMoves(&board, core.WHITE) {
			t.Errorf("HasAnyLegalMoves = false, want true (knight has moves)")
		}
	})

	t.Run("a checkmated side returns false", func(t *testing.T) {
		// Black king on H8, white queen on G7 (covers all escape squares except
		// H7 which is covered by the white king on F6), white king on F6.
		// This is a checkmate position.
		var board core.Board
		board[core.H8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
		board[core.G7] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.WHITE})
		board[core.F6] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})

		if hasMoves(&board, core.BLACK, testutil.WithSides(
			testutil.Side(core.F6, false, false),
			testutil.Side(core.H8, false, false),
		)) {
			t.Errorf("HasAnyLegalMoves = true, want false (checkmate)")
		}
	})

	t.Run("a stalemated side returns false (no legal moves, but not in check)", func(t *testing.T) {
		// Black king on A1, white queen on B3, white king on C2.
		// The queen covers A2, B1, B2, A3 — but NOT A1 (so it's not check).
		// The king on C2 covers B1, B2. So every black king move is attacked,
		// but the king itself is not in check → stalemate.
		var board core.Board
		board[core.A1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
		board[core.B3] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.WHITE})
		board[core.C2] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})

		if hasMoves(&board, core.BLACK, testutil.WithSides(
			testutil.Side(core.C2, false, false),
			testutil.Side(core.A1, false, false),
		)) {
			t.Errorf("HasAnyLegalMoves = true, want false (stalemate)")
		}
	})

	t.Run("only the side-to-move's pieces are checked (white stalemated, black has moves)", func(t *testing.T) {
		// Same board as the stalemate test, but colors swapped: white king on
		// A1 (stalemated), black queen on B3, black king on C2.
		var board core.Board
		board[core.A1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
		board[core.B3] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.BLACK})
		board[core.C2] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})

		// White to move → stalemated → false.
		if hasMoves(&board, core.WHITE, testutil.WithSides(
			testutil.Side(core.A1, false, false),
			testutil.Side(core.C2, false, false),
		)) {
			t.Errorf("HasAnyLegalMoves(white) = true, want false (stalemated)")
		}
	})

	t.Run("the same board with black to move returns true (black has moves)", func(t *testing.T) {
		// Same board as above, but black to move. Black's queen and king
		// have plenty of moves.
		var board core.Board
		board[core.A1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
		board[core.B3] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.BLACK})
		board[core.C2] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})

		if !hasMoves(&board, core.BLACK, testutil.WithSides(
			testutil.Side(core.A1, false, false),
			testutil.Side(core.C2, false, false),
		)) {
			t.Errorf("HasAnyLegalMoves(black) = false, want true (black has moves)")
		}
	})

	t.Run("a side with no pieces on the board returns false", func(t *testing.T) {
		// Only a black king on the board; white to move but has no pieces.
		var board core.Board
		board[core.E8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})

		if hasMoves(&board, core.WHITE, testutil.WithSides(
			testutil.Side(core.NoPosition, false, false),
			testutil.Side(core.E8, false, false),
		)) {
			t.Errorf("HasAnyLegalMoves = true, want false (no pieces for white)")
		}
	})

	t.Run("if the first piece is blocked, the scan continues to the next piece", func(t *testing.T) {
		// White pawn on A2 blocked by own pawn on A3 (can't move). But the
		// knight on B1 has moves → true.
		var board core.Board
		board[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
		board[core.A2] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.A3] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.B1] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
		board[core.E8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})

		if !hasMoves(&board, core.WHITE) {
			t.Errorf("HasAnyLegalMoves = false, want true (knight on B1 has moves)")
		}
	})

	t.Run("a pinned piece with no legal moves doesn't prevent other pieces from moving", func(t *testing.T) {
		// White king on E1, white bishop on E2 (pinned, no moves), white
		// knight on B1 (has moves). Black rook on E8 pins the bishop.
		var board core.Board
		board[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
		board[core.E2] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
		board[core.B1] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
		board[core.E8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
		board[core.H8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})

		if !hasMoves(&board, core.WHITE, testutil.WithSides(
			testutil.Side(core.E1, false, false),
			testutil.Side(core.H8, true, true),
		)) {
			t.Errorf("HasAnyLegalMoves = false, want true (knight on B1 has moves despite pinned bishop)")
		}
	})
}

// TestIsLegalMove verifies that IsLegalMove correctly identifies whether a
// specific move is legal in the given position. It delegates to GetLegalMoves
// internally, so these tests focus on the API contract rather than exhaustive
// move-generation coverage (which is tested in TestGetLegalMoves).
func TestIsLegalMove(t *testing.T) {
	engine := GetDefaultEngine()

	isLegal := func(board *core.Board, side core.PieceColor, move core.Move, opts ...testutil.TurnOption) bool {
		ctx := testutil.NewTurn(board, side, opts...)
		return engine.IsLegalMove(move, *ctx)
	}

	// =========================================================================
	// Basic legality.
	// =========================================================================

	t.Run("a normal legal move returns true", func(t *testing.T) {
		var board core.Board
		board[core.E2] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
		board[core.E8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})

		move := core.Move{Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}, From: core.E2, To: core.E4, Type: core.NORMAL}
		if !isLegal(&board, core.WHITE, move) {
			t.Errorf("IsLegalMove = false, want true for a legal pawn push")
		}
	})

	t.Run("a move that exposes the king returns false", func(t *testing.T) {
		// White king on E1, white rook on E2 pinned by black rook on E8.
		// Moving the rook sideways exposes the king.
		var board core.Board
		board[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
		board[core.E2] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
		board[core.E8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
		board[core.H8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})

		move := core.Move{Piece: core.Piece{Type: core.ROOK, Color: core.WHITE}, From: core.E2, To: core.D2, Type: core.NORMAL}
		if isLegal(&board, core.WHITE, move, testutil.WithSides(
			testutil.Side(core.E1, false, false),
			testutil.Side(core.H8, true, true),
		)) {
			t.Errorf("IsLegalMove = true, want false (move exposes king on E-file)")
		}
	})

	t.Run("a move from an empty square returns false", func(t *testing.T) {
		var board core.Board
		board[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
		board[core.E8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})

		move := core.Move{Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}, From: core.E4, To: core.E5, Type: core.NORMAL}
		if isLegal(&board, core.WHITE, move) {
			t.Errorf("IsLegalMove = true, want false (no piece on E4)")
		}
	})

	t.Run("a move for the opponent's piece returns false", func(t *testing.T) {
		var board core.Board
		board[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
		board[core.E7] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		board[core.E8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})

		move := core.Move{Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}, From: core.E7, To: core.E5, Type: core.NORMAL}
		if isLegal(&board, core.WHITE, move) {
			t.Errorf("IsLegalMove = true, want false (piece belongs to opponent)")
		}
	})

	// =========================================================================
	// Special move types.
	// =========================================================================

	t.Run("a legal castling move returns true", func(t *testing.T) {
		var board core.Board
		board[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
		board[core.H1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
		board[core.E8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})

		move := core.Move{Piece: core.Piece{Type: core.KING, Color: core.WHITE}, From: core.E1, To: core.G1, Type: core.CASTLING}
		if !isLegal(&board, core.WHITE, move) {
			t.Errorf("IsLegalMove = false, want true for king-side castling")
		}
	})

	t.Run("a legal en passant move returns true", func(t *testing.T) {
		var board core.Board
		board[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
		board[core.E5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.D5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		board[core.E8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})

		move := core.Move{
			Piece:      core.Piece{Type: core.PAWN, Color: core.WHITE},
			From:       core.E5,
			To:         core.D6,
			Type:       core.EN_PASSANT,
			Captured:   core.Piece{Type: core.PAWN, Color: core.BLACK},
			HasCapture: true,
		}
		if !isLegal(&board, core.WHITE, move, testutil.WithSides(
			testutil.Side(core.E1, false, false),
			testutil.Side(core.E8, false, false),
		), testutil.WithEnPassantTarget(core.D6)) {
			t.Errorf("IsLegalMove = false, want true for a legal en passant")
		}
	})

	t.Run("an en passant that exposes the king on the rank returns false", func(t *testing.T) {
		// White king H5, white pawn F5, black pawn E5 (EP target E6), black rook A5.
		// Capturing en passant removes both pawns from rank 5, exposing the king.
		var board core.Board
		board[core.H5] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
		board[core.F5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.E5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		board[core.A5] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
		board[core.E8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})

		move := core.Move{
			Piece:      core.Piece{Type: core.PAWN, Color: core.WHITE},
			From:       core.F5,
			To:         core.E6,
			Type:       core.EN_PASSANT,
			Captured:   core.Piece{Type: core.PAWN, Color: core.BLACK},
			HasCapture: true,
		}
		if isLegal(&board, core.WHITE, move, testutil.WithSides(
			testutil.Side(core.H5, false, false),
			testutil.Side(core.E8, true, true),
		), testutil.WithEnPassantTarget(core.E6)) {
			t.Errorf("IsLegalMove = true, want false (en passant exposes king on rank 5)")
		}
	})

	t.Run("a legal promotion move returns true", func(t *testing.T) {
		var board core.Board
		board[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
		board[core.D7] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.H8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})

		move := core.Move{
			Piece:     core.Piece{Type: core.PAWN, Color: core.WHITE},
			From:      core.D7,
			To:        core.D8,
			Type:      core.PROMOTION,
			PromoteTo: core.QUEEN,
		}
		if !isLegal(&board, core.WHITE, move, testutil.WithSides(
			testutil.Side(core.E1, false, false),
			testutil.Side(core.H8, false, false),
		)) {
			t.Errorf("IsLegalMove = false, want true for a legal promotion")
		}
	})
}
