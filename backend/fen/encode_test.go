package fen

import (
	"strings"
	"testing"

	"github.com/7-Dany/chess/core"
	"github.com/7-Dany/chess/testutil"
)

// TestEncode verifies that Encode serializes a TurnContext into a valid FEN
// string. Each subtest builds a context (usually by decoding a known FEN),
// re-encodes it, and asserts the output matches.
func TestEncode(t *testing.T) {
	parser := GetDefaultFenParser()

	// Helper: decode a FEN, then re-encode it, return the encoded string.
	roundTrip := func(t *testing.T, fen string) string {
		t.Helper()
		var ctx core.TurnContext
		if err := parser.Decode(fen, &ctx); err != nil {
			t.Fatalf("Decode(%q) failed: %v", fen, err)
		}
		return parser.Encode(&ctx)
	}

	// =========================================================================
	// Round-trip: Decode then Encode should produce the same string.
	// This is the strongest test — if any field round-trips correctly, the
	// encoder and decoder agree on its format.
	// =========================================================================

	t.Run("the starting position round-trips to the same FEN", func(t *testing.T) {
		const fen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
		got := roundTrip(t, fen)
		if got != fen {
			t.Errorf("round-trip mismatch:\n  got:  %s\n  want: %s", got, fen)
		}
	})

	t.Run("the Kiwipete position round-trips to the same FEN", func(t *testing.T) {
		const fen = "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1"
		got := roundTrip(t, fen)
		if got != fen {
			t.Errorf("round-trip mismatch:\n  got:  %s\n  want: %s", got, fen)
		}
	})

	t.Run("a position with no castling rights round-trips", func(t *testing.T) {
		const fen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w - - 0 1"
		got := roundTrip(t, fen)
		if got != fen {
			t.Errorf("round-trip mismatch:\n  got:  %s\n  want: %s", got, fen)
		}
	})

	t.Run("a position with only some castling rights round-trips", func(t *testing.T) {
		const fen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w Kq - 0 1"
		got := roundTrip(t, fen)
		if got != fen {
			t.Errorf("round-trip mismatch:\n  got:  %s\n  want: %s", got, fen)
		}
	})

	t.Run("a position with an en passant target round-trips", func(t *testing.T) {
		const fen = "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"
		got := roundTrip(t, fen)
		if got != fen {
			t.Errorf("round-trip mismatch:\n  got:  %s\n  want: %s", got, fen)
		}
	})

	t.Run("a position with black to move round-trips", func(t *testing.T) {
		const fen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQkq - 0 1"
		got := roundTrip(t, fen)
		if got != fen {
			t.Errorf("round-trip mismatch:\n  got:  %s\n  want: %s", got, fen)
		}
	})

	t.Run("a position with multi-digit clocks round-trips", func(t *testing.T) {
		const fen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 47 132"
		got := roundTrip(t, fen)
		if got != fen {
			t.Errorf("round-trip mismatch:\n  got:  %s\n  want: %s", got, fen)
		}
	})

	t.Run("an empty board with two kings round-trips", func(t *testing.T) {
		const fen = "4k3/8/8/8/8/8/8/4K3 w - - 0 1"
		got := roundTrip(t, fen)
		if got != fen {
			t.Errorf("round-trip mismatch:\n  got:  %s\n  want: %s", got, fen)
		}
	})

	// =========================================================================
	// Individual field checks — build a context by hand and verify the
	// encoded output field by field. This catches bugs that round-trip
	// misses (e.g. if both encoder and decoder had the same off-by-one).
	// =========================================================================

	t.Run("piece placement emits rank 8 first, then 7, down to 1", func(t *testing.T) {
		// White king on E1, black king on E8, nothing else.
		var ctx core.TurnContext
		ctx.Reset()
		ctx.Board.Place(core.E1, core.Piece{Type: core.KING, Color: core.WHITE})
		ctx.Board.Place(core.E8, core.Piece{Type: core.KING, Color: core.BLACK})

		got := parser.Encode(&ctx)

		// The first field should be "4k3/8/8/8/8/8/8/4K3".
		firstField := strings.SplitN(got, " ", 2)[0]
		want := "4k3/8/8/8/8/8/8/4K3"
		if firstField != want {
			t.Errorf("piece placement = %q, want %q", firstField, want)
		}
	})

	t.Run("consecutive empty squares collapse into a single digit", func(t *testing.T) {
		// A completely empty board (illegal, but tests the digit logic).
		// Each rank should be "8".
		var ctx core.TurnContext
		ctx.Reset()
		// No pieces placed — every square empty.

		got := parser.Encode(&ctx)
		firstField := strings.SplitN(got, " ", 2)[0]
		want := "8/8/8/8/8/8/8/8"
		if firstField != want {
			t.Errorf("piece placement = %q, want %q", firstField, want)
		}
	})

	t.Run("a rank with mixed pieces and empties encodes correctly", func(t *testing.T) {
		// Rook, 3 empty, king, 2 empty, rook → "r3k2r"
		var ctx core.TurnContext
		ctx.Reset()
		ctx.Board.Place(core.A8, core.Piece{Type: core.ROOK, Color: core.BLACK})
		ctx.Board.Place(core.E8, core.Piece{Type: core.KING, Color: core.BLACK})
		ctx.Board.Place(core.H8, core.Piece{Type: core.ROOK, Color: core.BLACK})

		got := parser.Encode(&ctx)
		firstField := strings.SplitN(got, " ", 2)[0]
		// Only rank 8 has pieces; ranks 7-1 are all empty.
		want := "r3k2r/8/8/8/8/8/8/8"
		if firstField != want {
			t.Errorf("piece placement = %q, want %q", firstField, want)
		}
	})

	t.Run("all six piece types encode with the correct letter and case", func(t *testing.T) {
		// Put one of each white piece on rank 1 and each black piece on rank 8.
		var ctx core.TurnContext
		ctx.Reset()
		ctx.Board.Place(core.A1, core.Piece{Type: core.PAWN, Color: core.WHITE})
		ctx.Board.Place(core.B1, core.Piece{Type: core.KNIGHT, Color: core.WHITE})
		ctx.Board.Place(core.C1, core.Piece{Type: core.BISHOP, Color: core.WHITE})
		ctx.Board.Place(core.D1, core.Piece{Type: core.ROOK, Color: core.WHITE})
		ctx.Board.Place(core.E1, core.Piece{Type: core.QUEEN, Color: core.WHITE})
		ctx.Board.Place(core.F1, core.Piece{Type: core.KING, Color: core.WHITE})
		ctx.Board.Place(core.A8, core.Piece{Type: core.PAWN, Color: core.BLACK})
		ctx.Board.Place(core.B8, core.Piece{Type: core.KNIGHT, Color: core.BLACK})
		ctx.Board.Place(core.C8, core.Piece{Type: core.BISHOP, Color: core.BLACK})
		ctx.Board.Place(core.D8, core.Piece{Type: core.ROOK, Color: core.BLACK})
		ctx.Board.Place(core.E8, core.Piece{Type: core.QUEEN, Color: core.BLACK})
		ctx.Board.Place(core.F8, core.Piece{Type: core.KING, Color: core.BLACK})

		got := parser.Encode(&ctx)
		firstField := strings.SplitN(got, " ", 2)[0]
		// Rank 8: pnbrqk2 (black pieces + 2 empties), Rank 1: PNBRQK2 (white + 2 empties)
		want := "pnbrqk2/8/8/8/8/8/8/PNBRQK2"
		if firstField != want {
			t.Errorf("piece placement = %q, want %q", firstField, want)
		}
	})

	t.Run("side to move emits 'w' for white", func(t *testing.T) {
		var ctx core.TurnContext
		ctx.Reset()
		ctx.SideToMove = core.WHITE
		got := parser.Encode(&ctx)
		// Second field (after the piece placement + space).
		fields := strings.Split(got, " ")
		if fields[1] != "w" {
			t.Errorf("side to move = %q, want \"w\"", fields[1])
		}
	})

	t.Run("side to move emits 'b' for black", func(t *testing.T) {
		var ctx core.TurnContext
		ctx.Reset()
		ctx.SideToMove = core.BLACK
		got := parser.Encode(&ctx)
		fields := strings.Split(got, " ")
		if fields[1] != "b" {
			t.Errorf("side to move = %q, want \"b\"", fields[1])
		}
	})

	t.Run("full castling rights emit KQkq in that order", func(t *testing.T) {
		var ctx core.TurnContext
		ctx.Reset()
		ctx.Sides = testutil.DefaultSides()
		got := parser.Encode(&ctx)
		fields := strings.Split(got, " ")
		if fields[2] != "KQkq" {
			t.Errorf("castling rights = %q, want \"KQkq\"", fields[2])
		}
	})

	t.Run("no castling rights emit '-'", func(t *testing.T) {
		var ctx core.TurnContext
		ctx.Reset()
		// Sides zero-valued → no rights.
		got := parser.Encode(&ctx)
		fields := strings.Split(got, " ")
		if fields[2] != "-" {
			t.Errorf("castling rights = %q, want \"-\"", fields[2])
		}
	})

	t.Run("partial castling rights emit only the active letters", func(t *testing.T) {
		var ctx core.TurnContext
		ctx.Reset()
		ctx.Sides[core.WHITE] = testutil.Side(core.E1, true, false) // king-side only
		ctx.Sides[core.BLACK] = testutil.Side(core.E8, false, true) // queen-side only
		got := parser.Encode(&ctx)
		fields := strings.Split(got, " ")
		// White king-side (K) + black queen-side (q) → "Kq"
		if fields[2] != "Kq" {
			t.Errorf("castling rights = %q, want \"Kq\"", fields[2])
		}
	})

	t.Run("en passant target emits the square in algebraic notation", func(t *testing.T) {
		var ctx core.TurnContext
		ctx.Reset()
		ctx.EnPassantTarget = core.E3
		got := parser.Encode(&ctx)
		fields := strings.Split(got, " ")
		if fields[3] != "e3" {
			t.Errorf("en passant = %q, want \"e3\"", fields[3])
		}
	})

	t.Run("no en passant target emits '-'", func(t *testing.T) {
		var ctx core.TurnContext
		ctx.Reset()
		ctx.EnPassantTarget = core.NoPosition
		got := parser.Encode(&ctx)
		fields := strings.Split(got, " ")
		if fields[3] != "-" {
			t.Errorf("en passant = %q, want \"-\"", fields[3])
		}
	})

	t.Run("the halfmove clock emits as a decimal number", func(t *testing.T) {
		var ctx core.TurnContext
		ctx.Reset()
		ctx.HalfMoveClock = 47
		got := parser.Encode(&ctx)
		fields := strings.Split(got, " ")
		if fields[4] != "47" {
			t.Errorf("halfmove clock = %q, want \"47\"", fields[4])
		}
	})

	t.Run("a zero halfmove clock emits '0'", func(t *testing.T) {
		var ctx core.TurnContext
		ctx.Reset()
		ctx.HalfMoveClock = 0
		got := parser.Encode(&ctx)
		fields := strings.Split(got, " ")
		if fields[4] != "0" {
			t.Errorf("halfmove clock = %q, want \"0\"", fields[4])
		}
	})

	t.Run("the fullmove number emits as a decimal number", func(t *testing.T) {
		var ctx core.TurnContext
		ctx.Reset()
		ctx.FullMoveNumber = 132
		got := parser.Encode(&ctx)
		fields := strings.Split(got, " ")
		if fields[5] != "132" {
			t.Errorf("fullmove number = %q, want \"132\"", fields[5])
		}
	})

	t.Run("a fullmove number of 1 emits '1'", func(t *testing.T) {
		var ctx core.TurnContext
		ctx.Reset()
		ctx.FullMoveNumber = 1
		got := parser.Encode(&ctx)
		fields := strings.Split(got, " ")
		if fields[5] != "1" {
			t.Errorf("fullmove number = %q, want \"1\"", fields[5])
		}
	})

	// =========================================================================
	// Full output structure — verify all six fields are present and in order.
	// =========================================================================

	t.Run("the output has exactly six space-separated fields", func(t *testing.T) {
		var ctx core.TurnContext
		ctx.Reset()
		ctx.Board.Place(core.E1, core.Piece{Type: core.KING, Color: core.WHITE})
		ctx.Board.Place(core.E8, core.Piece{Type: core.KING, Color: core.BLACK})
		ctx.SideToMove = core.WHITE
		ctx.Sides = testutil.DefaultSides()
		ctx.EnPassantTarget = core.NoPosition
		ctx.HalfMoveClock = 0
		ctx.FullMoveNumber = 1

		got := parser.Encode(&ctx)
		fields := strings.Split(got, " ")
		if len(fields) != 6 {
			t.Errorf("got %d fields, want 6: %q", len(fields), got)
		}
	})

	t.Run("a complete hand-built position encodes to the expected FEN", func(t *testing.T) {
		var ctx core.TurnContext
		ctx.Reset()
		// Build the starting position by hand.
		back := []core.PieceType{core.ROOK, core.KNIGHT, core.BISHOP, core.QUEEN,
			core.KING, core.BISHOP, core.KNIGHT, core.ROOK}
		for f := range 8 {
			ctx.Board.Place(core.NewPosition(core.File(f), core.RANK_1), core.Piece{Type: back[f], Color: core.WHITE})
			ctx.Board.Place(core.NewPosition(core.File(f), core.RANK_2), core.Piece{Type: core.PAWN, Color: core.WHITE})
			ctx.Board.Place(core.NewPosition(core.File(f), core.RANK_7), core.Piece{Type: core.PAWN, Color: core.BLACK})
			ctx.Board.Place(core.NewPosition(core.File(f), core.RANK_8), core.Piece{Type: back[f], Color: core.BLACK})
		}
		ctx.SideToMove = core.WHITE
		ctx.Sides = testutil.DefaultSides()
		ctx.EnPassantTarget = core.NoPosition
		ctx.HalfMoveClock = 0
		ctx.FullMoveNumber = 1

		got := parser.Encode(&ctx)
		want := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
		if got != want {
			t.Errorf("Encode = %q, want %q", got, want)
		}
	})
}
