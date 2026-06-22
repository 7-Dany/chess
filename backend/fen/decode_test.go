package fen

import (
	"testing"

	"github.com/7-Dany/chess/core"
	"github.com/7-Dany/chess/testutil"
)

// TestDecode verifies that Decode correctly parses every FEN field and
// rejects malformed input with a clear error.
//
// Each subtest is a self-contained scenario. The valid-parse cases build a
// FEN string, decode it, and assert the resulting TurnContext fields. The
// error cases feed malformed input and assert that Decode returns an error
// (and doesn't panic).
func TestDecode(t *testing.T) {
	parser := GetDefaultFenParser()

	// Helper: decode a FEN string and return the context + error.
	decode := func(t *testing.T, fen string) (core.TurnContext, error) {
		t.Helper()
		var ctx core.TurnContext
		err := parser.Decode(fen, &ctx)
		return ctx, err
	}

	// Helper: assert Decode succeeds.
	mustDecode := func(t *testing.T, fen string) core.TurnContext {
		t.Helper()
		ctx, err := decode(t, fen)
		if err != nil {
			t.Fatalf("Decode(%q) returned error: %v", fen, err)
		}
		return ctx
	}

	// Helper: assert Decode fails.
	mustFail := func(t *testing.T, fen string) {
		t.Helper()
		_, err := decode(t, fen)
		if err == nil {
			t.Errorf("Decode(%q) should have returned an error, got nil", fen)
		}
	}

	// =========================================================================
	// Full valid FEN strings — every field parsed correctly.
	// =========================================================================

	t.Run("the standard starting position parses with all fields correct", func(t *testing.T) {
		ctx := mustDecode(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")

		// Spot-check a few squares from each rank.
		// Rank 8 (top): black rook on A8, black king on E8.
		testutil.AssertSquareHas(t, ctx.Board, core.A8, core.ROOK, core.BLACK)
		testutil.AssertSquareHas(t, ctx.Board, core.E8, core.KING, core.BLACK)
		// Rank 7: black pawns across.
		testutil.AssertSquareHas(t, ctx.Board, core.A7, core.PAWN, core.BLACK)
		testutil.AssertSquareHas(t, ctx.Board, core.H7, core.PAWN, core.BLACK)
		// Rank 4-5: empty.
		testutil.AssertSquareEmpty(t, ctx.Board, core.E4)
		testutil.AssertSquareEmpty(t, ctx.Board, core.D5)
		// Rank 2: white pawns across.
		testutil.AssertSquareHas(t, ctx.Board, core.A2, core.PAWN, core.WHITE)
		testutil.AssertSquareHas(t, ctx.Board, core.H2, core.PAWN, core.WHITE)
		// Rank 1: white queen on D1, white king on E1.
		testutil.AssertSquareHas(t, ctx.Board, core.D1, core.QUEEN, core.WHITE)
		testutil.AssertSquareHas(t, ctx.Board, core.E1, core.KING, core.WHITE)

		// Side to move: white.
		if ctx.SideToMove != core.WHITE {
			t.Errorf("SideToMove = %v, want WHITE", ctx.SideToMove)
		}
		// Castling: all four rights.
		if !ctx.Sides[core.WHITE].CanCastleKingSide || !ctx.Sides[core.WHITE].CanCastleQueenSide {
			t.Errorf("white castling rights = %+v, want both true", ctx.Sides[core.WHITE])
		}
		if !ctx.Sides[core.BLACK].CanCastleKingSide || !ctx.Sides[core.BLACK].CanCastleQueenSide {
			t.Errorf("black castling rights = %+v, want both true", ctx.Sides[core.BLACK])
		}
		// No en passant target.
		if ctx.EnPassantTarget != core.NoPosition {
			t.Errorf("EnPassantTarget = %v, want NoPosition", ctx.EnPassantTarget)
		}
		// Clocks: 0 half-moves, move 1.
		if ctx.HalfMoveClock != 0 {
			t.Errorf("HalfMoveClock = %d, want 0", ctx.HalfMoveClock)
		}
		if ctx.FullMoveNumber != 1 {
			t.Errorf("FullMoveNumber = %d, want 1", ctx.FullMoveNumber)
		}
	})

	t.Run("the Kiwipete position parses correctly", func(t *testing.T) {
		ctx := mustDecode(t, "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1")

		// Spot-check the non-pawn pieces.
		testutil.AssertSquareHas(t, ctx.Board, core.A8, core.ROOK, core.BLACK)
		testutil.AssertSquareHas(t, ctx.Board, core.E8, core.KING, core.BLACK)
		testutil.AssertSquareHas(t, ctx.Board, core.H8, core.ROOK, core.BLACK)
		testutil.AssertSquareHas(t, ctx.Board, core.E7, core.QUEEN, core.BLACK)
		testutil.AssertSquareHas(t, ctx.Board, core.G7, core.BISHOP, core.BLACK)
		testutil.AssertSquareHas(t, ctx.Board, core.B6, core.KNIGHT, core.BLACK)
		testutil.AssertSquareHas(t, ctx.Board, core.D5, core.PAWN, core.WHITE)
		testutil.AssertSquareHas(t, ctx.Board, core.E5, core.KNIGHT, core.WHITE)
		testutil.AssertSquareHas(t, ctx.Board, core.F3, core.QUEEN, core.WHITE)
		testutil.AssertSquareHas(t, ctx.Board, core.E1, core.KING, core.WHITE)
		testutil.AssertSquareHas(t, ctx.Board, core.A1, core.ROOK, core.WHITE)
		testutil.AssertSquareHas(t, ctx.Board, core.H1, core.ROOK, core.WHITE)

		if ctx.SideToMove != core.WHITE {
			t.Errorf("SideToMove = %v, want WHITE", ctx.SideToMove)
		}
	})

	t.Run("a position with black to move parses with the correct side", func(t *testing.T) {
		ctx := mustDecode(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQkq - 0 1")
		if ctx.SideToMove != core.BLACK {
			t.Errorf("SideToMove = %v, want BLACK", ctx.SideToMove)
		}
	})

	t.Run("a position with no castling rights parses with all rights false", func(t *testing.T) {
		ctx := mustDecode(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w - - 0 1")
		if ctx.Sides[core.WHITE].CanCastleKingSide || ctx.Sides[core.WHITE].CanCastleQueenSide {
			t.Errorf("white castling rights should be false, got %+v", ctx.Sides[core.WHITE])
		}
		if ctx.Sides[core.BLACK].CanCastleKingSide || ctx.Sides[core.BLACK].CanCastleQueenSide {
			t.Errorf("black castling rights should be false, got %+v", ctx.Sides[core.BLACK])
		}
	})

	t.Run("a position with only white king-side castling parses with that single right", func(t *testing.T) {
		ctx := mustDecode(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w K - 0 1")
		if !ctx.Sides[core.WHITE].CanCastleKingSide {
			t.Errorf("white king-side right should be true")
		}
		if ctx.Sides[core.WHITE].CanCastleQueenSide {
			t.Errorf("white queen-side right should be false")
		}
		if ctx.Sides[core.BLACK].CanCastleKingSide || ctx.Sides[core.BLACK].CanCastleQueenSide {
			t.Errorf("black should have no rights, got %+v", ctx.Sides[core.BLACK])
		}
	})

	t.Run("a position with only black queen-side castling parses with that single right", func(t *testing.T) {
		ctx := mustDecode(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w q - 0 1")
		if !ctx.Sides[core.BLACK].CanCastleQueenSide {
			t.Errorf("black queen-side right should be true")
		}
		if ctx.Sides[core.BLACK].CanCastleKingSide {
			t.Errorf("black king-side right should be false")
		}
		if ctx.Sides[core.WHITE].CanCastleKingSide || ctx.Sides[core.WHITE].CanCastleQueenSide {
			t.Errorf("white should have no rights, got %+v", ctx.Sides[core.WHITE])
		}
	})

	t.Run("an en passant target on rank 3 (white just moved) parses correctly", func(t *testing.T) {
		// After 1.e4, the en passant target is e3.
		ctx := mustDecode(t, "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1")
		if ctx.EnPassantTarget != core.E3 {
			t.Errorf("EnPassantTarget = %v, want E3", ctx.EnPassantTarget)
		}
	})

	t.Run("an en passant target on rank 6 (black just moved) parses correctly", func(t *testing.T) {
		// After 1...d5, the en passant target is d6.
		ctx := mustDecode(t, "rnbqkbnr/ppp1pppp/8/3p4/4P3/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 2")
		if ctx.EnPassantTarget != core.D6 {
			t.Errorf("EnPassantTarget = %v, want D6", ctx.EnPassantTarget)
		}
	})

	t.Run("the halfmove clock and fullmove number parse as decimal numbers", func(t *testing.T) {
		ctx := mustDecode(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 47 132")
		if ctx.HalfMoveClock != 47 {
			t.Errorf("HalfMoveClock = %d, want 47", ctx.HalfMoveClock)
		}
		if ctx.FullMoveNumber != 132 {
			t.Errorf("FullMoveNumber = %d, want 132", ctx.FullMoveNumber)
		}
	})

	// =========================================================================
	// Piece placement — rank/file ordering and digit runs.
	// =========================================================================

	t.Run("digit runs place the correct number of empty squares", func(t *testing.T) {
		// Rank 8: r3k2r → rook, 3 empty, king, 2 empty, rook.
		ctx := mustDecode(t, "r3k2r/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
		testutil.AssertSquareHas(t, ctx.Board, core.A8, core.ROOK, core.BLACK)
		testutil.AssertSquareEmpty(t, ctx.Board, core.B8)
		testutil.AssertSquareEmpty(t, ctx.Board, core.C8)
		testutil.AssertSquareEmpty(t, ctx.Board, core.D8)
		testutil.AssertSquareHas(t, ctx.Board, core.E8, core.KING, core.BLACK)
		testutil.AssertSquareEmpty(t, ctx.Board, core.F8)
		testutil.AssertSquareEmpty(t, ctx.Board, core.G8)
		testutil.AssertSquareHas(t, ctx.Board, core.H8, core.ROOK, core.BLACK)
	})

	t.Run("the digit 8 fills an entire rank with empties", func(t *testing.T) {
		ctx := mustDecode(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
		for f := uint8(0); f < 8; f++ {
			testutil.AssertSquareEmpty(t, ctx.Board, core.NewPosition(core.File(f), core.RANK_4))
			testutil.AssertSquareEmpty(t, ctx.Board, core.NewPosition(core.File(f), core.RANK_5))
		}
	})

	t.Run("pieces are placed on the correct squares (FEN rank 8 = internal RANK_8)", func(t *testing.T) {
		// A single white knight on B1 (bottom-right of the FEN string).
		ctx := mustDecode(t, "8/8/8/8/8/8/8/1N6 w - - 0 1")
		testutil.AssertSquareHas(t, ctx.Board, core.B1, core.KNIGHT, core.WHITE)
		// Everything else empty.
		testutil.AssertSquareEmpty(t, ctx.Board, core.A1)
		testutil.AssertSquareEmpty(t, ctx.Board, core.C1)
		testutil.AssertSquareEmpty(t, ctx.Board, core.A8)
	})

	t.Run("pieces are placed on the correct squares (FEN rank 1 = internal RANK_8, top)", func(t *testing.T) {
		// A single black king on E8 (top of the board, first rank in FEN).
		ctx := mustDecode(t, "4k3/8/8/8/8/8/8/8 w - - 0 1")
		testutil.AssertSquareHas(t, ctx.Board, core.E8, core.KING, core.BLACK)
		testutil.AssertSquareEmpty(t, ctx.Board, core.A8)
		testutil.AssertSquareEmpty(t, ctx.Board, core.E1)
	})

	t.Run("all six piece types parse for both colors", func(t *testing.T) {
		// One of each piece type for both colors. White pieces on rank 7,
		// black pieces on rank 2, each row padded to 8 files with a digit.
		//   Rank 7 (FEN row 2): PNBRQK1p — white P N B R Q K, empty, black p
		//   Rank 2 (FEN row 7): pnbrqk1P — black p n b r q k, empty, white P
		ctx := mustDecode(t, "8/PNBRQK1p/8/8/8/8/pnbrqk1P/8 w - - 0 1")
		// White pieces on rank 7 (FEN row 2, internal RANK_7).
		testutil.AssertSquareHas(t, ctx.Board, core.A7, core.PAWN, core.WHITE)
		testutil.AssertSquareHas(t, ctx.Board, core.B7, core.KNIGHT, core.WHITE)
		testutil.AssertSquareHas(t, ctx.Board, core.C7, core.BISHOP, core.WHITE)
		testutil.AssertSquareHas(t, ctx.Board, core.D7, core.ROOK, core.WHITE)
		testutil.AssertSquareHas(t, ctx.Board, core.E7, core.QUEEN, core.WHITE)
		testutil.AssertSquareHas(t, ctx.Board, core.F7, core.KING, core.WHITE)
		testutil.AssertSquareHas(t, ctx.Board, core.H7, core.PAWN, core.BLACK) // lowercase p at file H
		// Black pieces on rank 2 (FEN row 7, internal RANK_2).
		testutil.AssertSquareHas(t, ctx.Board, core.A2, core.PAWN, core.BLACK)
		testutil.AssertSquareHas(t, ctx.Board, core.B2, core.KNIGHT, core.BLACK)
		testutil.AssertSquareHas(t, ctx.Board, core.C2, core.BISHOP, core.BLACK)
		testutil.AssertSquareHas(t, ctx.Board, core.D2, core.ROOK, core.BLACK)
		testutil.AssertSquareHas(t, ctx.Board, core.E2, core.QUEEN, core.BLACK)
		testutil.AssertSquareHas(t, ctx.Board, core.F2, core.KING, core.BLACK)
		testutil.AssertSquareHas(t, ctx.Board, core.H2, core.PAWN, core.WHITE) // uppercase P at file H
	})

	// =========================================================================
	// Error cases — piece placement field.
	// =========================================================================

	t.Run("a rank with too few files returns an error", func(t *testing.T) {
		// Rank 8 has only 7 files (rnbqkbr should be rnbqkbnr — but here it's
		// just 7 letters). Actually "rnbqkbn" is 7 — the 8th is missing.
		mustFail(t, "rnbqkbn/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	})

	t.Run("a rank with too many files returns an error", func(t *testing.T) {
		// Rank 8 has 9 files (rnbqkbnr + extra pawn).
		mustFail(t, "rnbqkbnrp/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	})

	t.Run("a digit that overflows the rank returns an error", func(t *testing.T) {
		// Rank 8: "r9" means rook + 9 empties = 10 files > 8.
		mustFail(t, "r9/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	})

	t.Run("too many ranks returns an error", func(t *testing.T) {
		// 9 ranks (extra /8 at the end of the piece field).
		mustFail(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR/8 w KQkq - 0 1")
	})

	t.Run("too few ranks returns an error", func(t *testing.T) {
		// Only 7 ranks.
		mustFail(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP w KQkq - 0 1")
	})

	t.Run("an invalid piece letter returns an error", func(t *testing.T) {
		// 'x' is not a FEN piece letter.
		mustFail(t, "rnbqkbnx/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	})

	t.Run("a piece placement field with no space terminator returns an error", func(t *testing.T) {
		// Only the piece field, nothing else.
		mustFail(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR")
	})

	// =========================================================================
	// Error cases — side to move field.
	// =========================================================================

	t.Run("a side-to-move letter other than w or b returns an error", func(t *testing.T) {
		mustFail(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR x KQkq - 0 1")
	})

	t.Run("a missing side-to-move field returns an error", func(t *testing.T) {
		// String ends after the piece placement + space.
		mustFail(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR ")
	})

	// =========================================================================
	// Error cases — castling rights field.
	// =========================================================================

	t.Run("an invalid castling-rights letter returns an error", func(t *testing.T) {
		// 'X' is not a valid castling letter.
		mustFail(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w XQkq - 0 1")
	})

	t.Run("a missing castling-rights field returns an error", func(t *testing.T) {
		// String ends after side-to-move + space.
		mustFail(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w ")
	})

	// =========================================================================
	// Error cases — en passant target field.
	// =========================================================================

	t.Run("an en passant target on an invalid rank returns an error", func(t *testing.T) {
		// e4 is not rank 3 or 6.
		mustFail(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq e4 0 1")
	})

	t.Run("an en passant target with an invalid file returns an error", func(t *testing.T) {
		// 'z' is not a valid file.
		mustFail(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq z3 0 1")
	})

	t.Run("an en passant target that is too short returns an error", func(t *testing.T) {
		// Only a file letter, no rank digit (string ends).
		mustFail(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq e")
	})

	t.Run("a missing en-passant-target field returns an error", func(t *testing.T) {
		// String ends after castling rights + space.
		mustFail(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq ")
	})

	// =========================================================================
	// Error cases — halfmove clock field.
	// =========================================================================

	t.Run("a halfmove clock that is not a number returns an error", func(t *testing.T) {
		mustFail(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - ab 1")
	})

	t.Run("a missing halfmove-clock field returns an error", func(t *testing.T) {
		// String ends after en passant + space.
		mustFail(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - ")
	})

	// =========================================================================
	// Error cases — fullmove number field.
	// =========================================================================

	t.Run("a fullmove number that is not a number returns an error", func(t *testing.T) {
		mustFail(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 ab")
	})

	t.Run("a missing fullmove-number field returns an error", func(t *testing.T) {
		// String ends after halfmove clock + space.
		mustFail(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 ")
	})

	// =========================================================================
	// Edge cases.
	// =========================================================================

	t.Run("an empty string returns an error", func(t *testing.T) {
		mustFail(t, "")
	})

	t.Run("the en passant target of '-' sets NoPosition", func(t *testing.T) {
		ctx := mustDecode(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
		if ctx.EnPassantTarget != core.NoPosition {
			t.Errorf("EnPassantTarget = %v, want NoPosition", ctx.EnPassantTarget)
		}
	})

	t.Run("castling rights can appear in any order", func(t *testing.T) {
		// kqKQ instead of KQkq — order shouldn't matter, each letter sets its right.
		ctx := mustDecode(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w kqKQ - 0 1")
		if !ctx.Sides[core.WHITE].CanCastleKingSide || !ctx.Sides[core.WHITE].CanCastleQueenSide {
			t.Errorf("white rights should both be true, got %+v", ctx.Sides[core.WHITE])
		}
		if !ctx.Sides[core.BLACK].CanCastleKingSide || !ctx.Sides[core.BLACK].CanCastleQueenSide {
			t.Errorf("black rights should both be true, got %+v", ctx.Sides[core.BLACK])
		}
	})

	t.Run("Decode can be called twice on the same ctx without leftover state", func(t *testing.T) {
		// First decode: full position with castling and en passant.
		var ctx core.TurnContext
		if err := parser.Decode("rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 5 10", &ctx); err != nil {
			t.Fatalf("first Decode failed: %v", err)
		}
		// Second decode: starting position, no en passant, white to move.
		if err := parser.Decode("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", &ctx); err != nil {
			t.Fatalf("second Decode failed: %v", err)
		}
		// The en passant target from the first decode must be gone.
		if ctx.EnPassantTarget != core.NoPosition {
			t.Errorf("EnPassantTarget = %v, want NoPosition (leftover from first decode)", ctx.EnPassantTarget)
		}
		if ctx.SideToMove != core.WHITE {
			t.Errorf("SideToMove = %v, want WHITE", ctx.SideToMove)
		}
		if ctx.HalfMoveClock != 0 {
			t.Errorf("HalfMoveClock = %d, want 0 (leftover from first decode)", ctx.HalfMoveClock)
		}
		if ctx.FullMoveNumber != 1 {
			t.Errorf("FullMoveNumber = %d, want 1", ctx.FullMoveNumber)
		}
		// The white pawn that was on E4 (from the first FEN) must be gone —
		// E4 should be empty in the starting position.
		testutil.AssertSquareEmpty(t, ctx.Board, core.E4)
		// And E2 should have a pawn again.
		testutil.AssertSquareHas(t, ctx.Board, core.E2, core.PAWN, core.WHITE)
	})

	// =========================================================================
	// King position auto-detection — the engine's castling and king-safety
	// logic depends on Sides[color].KingPosition being set correctly. FEN
	// doesn't encode this explicitly, so Decode must detect it from the board.
	// =========================================================================

	t.Run("the white king position is detected from the board", func(t *testing.T) {
		// White king on E1 (standard position).
		ctx := mustDecode(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
		if ctx.Sides[core.WHITE].KingPosition != core.E1 {
			t.Errorf("white KingPosition = %v, want E1", ctx.Sides[core.WHITE].KingPosition)
		}
	})

	t.Run("the black king position is detected from the board", func(t *testing.T) {
		// Black king on E8 (standard position).
		ctx := mustDecode(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
		if ctx.Sides[core.BLACK].KingPosition != core.E8 {
			t.Errorf("black KingPosition = %v, want E8", ctx.Sides[core.BLACK].KingPosition)
		}
	})

	t.Run("a king that has castled to G1 is detected at G1", func(t *testing.T) {
		// White king on G1, rook on F1 (after king-side castling).
		ctx := mustDecode(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/R4RK1 w kq - 0 1")
		if ctx.Sides[core.WHITE].KingPosition != core.G1 {
			t.Errorf("white KingPosition = %v, want G1", ctx.Sides[core.WHITE].KingPosition)
		}
	})

	t.Run("a king that has castled to C8 is detected at C8", func(t *testing.T) {
		// Black king on C8, rook on D8 (after queen-side castling).
		ctx := mustDecode(t, "2kr3r/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQ - 0 1")
		if ctx.Sides[core.BLACK].KingPosition != core.C8 {
			t.Errorf("black KingPosition = %v, want C8", ctx.Sides[core.BLACK].KingPosition)
		}
	})

	t.Run("a king in the middle of the board is detected at its square", func(t *testing.T) {
		// White king on E4, black king on D5 — an unusual mid-board position.
		ctx := mustDecode(t, "8/8/8/3k4/4K3/8/8/8 w - - 0 1")
		if ctx.Sides[core.WHITE].KingPosition != core.E4 {
			t.Errorf("white KingPosition = %v, want E4", ctx.Sides[core.WHITE].KingPosition)
		}
		if ctx.Sides[core.BLACK].KingPosition != core.D5 {
			t.Errorf("black KingPosition = %v, want D5", ctx.Sides[core.BLACK].KingPosition)
		}
	})

	t.Run("a missing white king leaves KingPosition at the zero value", func(t *testing.T) {
		// No white king on the board — KingPosition stays at A1 (the zero
		// value of Position). This is an illegal position, but Decode should
		// not crash; it just leaves the field unset.
		ctx := mustDecode(t, "4k3/8/8/8/8/8/8/8 w - - 0 1")
		if ctx.Sides[core.WHITE].KingPosition != core.A1 {
			t.Errorf("white KingPosition = %v, want A1 (zero value, no king found)", ctx.Sides[core.WHITE].KingPosition)
		}
	})
}
