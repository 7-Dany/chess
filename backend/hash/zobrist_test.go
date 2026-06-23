package hash

import (
	"testing"

	"github.com/7-Dany/chess/core"
	"github.com/7-Dany/chess/engine"
	"github.com/7-Dany/chess/fen"
)

// TestZobristHasher verifies that the Zobrist hasher correctly computes
// incremental hashes and that apply+undo round-trips back to the original
// hash.
//
// The core invariant: because XOR is self-inverse, calling Hash after Apply
// (to advance the hash) and then calling Hash again before Undo (to revert
// it) must return the hash to its original value. If any move type fails
// this round-trip, the incremental update is wrong.
func TestZobristHasher(t *testing.T) {
	hasher := GetDefaultHasher()
	eng := engine.GetDefaultEngine()

	decode := func(t *testing.T, fenStr string) *core.TurnContext {
		t.Helper()
		var ctx core.TurnContext
		if err := fen.GetDefaultFenParser().Decode(fenStr, &ctx); err != nil {
			t.Fatalf("Decode failed: %v", err)
		}
		return &ctx
	}

	// roundTrip asserts that apply+hash then hash+undo returns to the
	// original hash. This is the definitive correctness test for any
	// incremental Zobrist implementation.
	roundTrip := func(t *testing.T, ctx *core.TurnContext, move core.Move) {
		t.Helper()

		originalHash := hasher.InitHash(ctx)

		// Advance: apply the move, then update the hash.
		snap := eng.Apply(ctx, move)
		moveHash := core.NewMoveHash(snap, *ctx)
		advancedHash := hasher.Hash(originalHash, moveHash)

		// Revert: update the hash (XOR is self-inverse), then undo.
		revertedHash := hasher.Hash(advancedHash, moveHash)
		eng.Undo(ctx, snap)

		if revertedHash != originalHash {
			t.Errorf("hash round-trip failed: original=%d, advanced=%d, reverted=%d",
				originalHash, advancedHash, revertedHash)
		}

		// Also verify the board was actually restored.
		restoredHash := hasher.InitHash(ctx)
		if restoredHash != originalHash {
			t.Errorf("board not restored after undo: init=%d, restored=%d",
				originalHash, restoredHash)
		}
	}

	// =========================================================================
	// InitHash — the bootstrap hash from a full position.
	// =========================================================================

	t.Run("InitHash produces a non-zero hash for the starting position", func(t *testing.T) {
		ctx := decode(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
		h := hasher.InitHash(ctx)
		if h == 0 {
			t.Errorf("InitHash returned 0 — the RNG or keys are broken")
		}
	})

	t.Run("InitHash is deterministic — same position always produces the same hash", func(t *testing.T) {
		ctx1 := decode(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
		ctx2 := decode(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
		h1 := hasher.InitHash(ctx1)
		h2 := hasher.InitHash(ctx2)
		if h1 != h2 {
			t.Errorf("InitHash not deterministic: %d vs %d", h1, h2)
		}
	})

	t.Run("InitHash differs when the side to move differs", func(t *testing.T) {
		ctxWhite := decode(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
		ctxBlack := decode(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQkq - 0 1")
		hWhite := hasher.InitHash(ctxWhite)
		hBlack := hasher.InitHash(ctxBlack)
		if hWhite == hBlack {
			t.Errorf("InitHash should differ for white vs black to move, both = %d", hWhite)
		}
	})

	t.Run("InitHash differs when castling rights differ", func(t *testing.T) {
		ctxFull := decode(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
		ctxNone := decode(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w - - 0 1")
		hFull := hasher.InitHash(ctxFull)
		hNone := hasher.InitHash(ctxNone)
		if hFull == hNone {
			t.Errorf("InitHash should differ when castling rights differ, both = %d", hFull)
		}
	})

	t.Run("InitHash differs when en passant target differs", func(t *testing.T) {
		ctxEP := decode(t, "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1")
		ctxNoEP := decode(t, "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq - 0 1")
		hEP := hasher.InitHash(ctxEP)
		hNoEP := hasher.InitHash(ctxNoEP)
		if hEP == hNoEP {
			t.Errorf("InitHash should differ when EP target differs, both = %d", hEP)
		}
	})

	t.Run("InitHash differs when a piece is in a different square", func(t *testing.T) {
		ctx1 := decode(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
		ctx2 := decode(t, "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 1")
		h1 := hasher.InitHash(ctx1)
		h2 := hasher.InitHash(ctx2)
		if h1 == h2 {
			t.Errorf("InitHash should differ for different positions, both = %d", h1)
		}
	})

	// =========================================================================
	// Hash round-trip — apply + hash, then hash + undo = original.
	// This is the definitive test: if any move type fails, the incremental
	// update is asymmetric and the hash is wrong.
	// =========================================================================

	t.Run("a normal pawn push round-trips", func(t *testing.T) {
		ctx := decode(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
		move := core.Move{
			Type:  core.NORMAL,
			Piece: core.Piece{Type: core.PAWN, Color: core.WHITE},
			From:  core.E2,
			To:    core.E3,
		}
		roundTrip(t, ctx, move)
	})

	t.Run("a normal knight move round-trips", func(t *testing.T) {
		ctx := decode(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
		move := core.Move{
			Type:  core.NORMAL,
			Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE},
			From:  core.B1,
			To:    core.C3,
		}
		roundTrip(t, ctx, move)
	})

	t.Run("a capture round-trips", func(t *testing.T) {
		// White knight on C3 captures black pawn on D5 (valid L-shape: file+1, rank+2).
		ctx := decode(t, "4k3/8/8/3p4/8/2N5/8/4K3 w - - 0 1")
		move := core.Move{
			Type:       core.NORMAL,
			Piece:      core.Piece{Type: core.KNIGHT, Color: core.WHITE},
			From:       core.C3,
			To:         core.D5,
			HasCapture: true,
			Captured:   core.Piece{Type: core.PAWN, Color: core.BLACK},
		}
		roundTrip(t, ctx, move)
	})

	t.Run("a double pawn push (sets EP target) round-trips", func(t *testing.T) {
		ctx := decode(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
		move := core.Move{
			Type:  core.NORMAL,
			Piece: core.Piece{Type: core.PAWN, Color: core.WHITE},
			From:  core.E2,
			To:    core.E4,
		}
		roundTrip(t, ctx, move)
	})

	t.Run("an en passant capture round-trips", func(t *testing.T) {
		ctx := decode(t, "rnbqkbnr/ppp1pppp/8/3pP3/8/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 2")
		move := core.Move{
			Type:       core.EN_PASSANT,
			Piece:      core.Piece{Type: core.PAWN, Color: core.WHITE},
			From:       core.E5,
			To:         core.D6,
			HasCapture: true,
			Captured:   core.Piece{Type: core.PAWN, Color: core.BLACK},
		}
		roundTrip(t, ctx, move)
	})

	t.Run("a king-side castling round-trips", func(t *testing.T) {
		ctx := decode(t, "rnbqk2r/pppppppp/5bn1/8/8/5BN1/PPPPPPPP/RNBQK2R w KQkq - 4 4")
		move := core.Move{
			Type:  core.CASTLING,
			Piece: core.Piece{Type: core.KING, Color: core.WHITE},
			From:  core.E1,
			To:    core.G1,
		}
		roundTrip(t, ctx, move)
	})

	t.Run("a queen-side castling round-trips", func(t *testing.T) {
		// King on E1, rook on A1, B1/C1/D1 clear.
		ctx := decode(t, "4k3/8/8/8/8/8/8/R3K3 w Q - 0 1")
		move := core.Move{
			Type:  core.CASTLING,
			Piece: core.Piece{Type: core.KING, Color: core.WHITE},
			From:  core.E1,
			To:    core.C1,
		}
		roundTrip(t, ctx, move)
	})

	t.Run("a promotion (no capture) round-trips", func(t *testing.T) {
		// White pawn on E7 promotes to E8. Black king on A8 (not on E8).
		ctx := decode(t, "k7/4P3/8/8/8/8/8/4K3 w - - 0 1")
		move := core.Move{
			Type:      core.PROMOTION,
			Piece:     core.Piece{Type: core.PAWN, Color: core.WHITE},
			From:      core.E7,
			To:        core.E8,
			PromoteTo: core.QUEEN,
		}
		roundTrip(t, ctx, move)
	})

	t.Run("a promotion with capture round-trips", func(t *testing.T) {
		ctx := decode(t, "3rk3/4P3/8/8/8/8/8/4K3 w - - 0 1")
		move := core.Move{
			Type:       core.PROMOTION,
			Piece:      core.Piece{Type: core.PAWN, Color: core.WHITE},
			From:       core.E7,
			To:         core.D8,
			PromoteTo:  core.QUEEN,
			HasCapture: true,
			Captured:   core.Piece{Type: core.ROOK, Color: core.BLACK},
		}
		roundTrip(t, ctx, move)
	})

	t.Run("a rook move that forfeits castling rights round-trips", func(t *testing.T) {
		// Rook on H1 moves to H5 (clear path), forfeiting king-side right.
		ctx := decode(t, "4k3/8/8/8/8/8/8/4K2R w K - 0 1")
		move := core.Move{
			Type:  core.NORMAL,
			Piece: core.Piece{Type: core.ROOK, Color: core.WHITE},
			From:  core.H1,
			To:    core.H5,
		}
		roundTrip(t, ctx, move)
	})

	t.Run("a king move that forfeits both castling rights round-trips", func(t *testing.T) {
		// King on E1 moves to E2 (clear), forfeiting both rights.
		ctx := decode(t, "4k3/8/8/8/8/8/8/4K2R w K - 0 1")
		move := core.Move{
			Type:  core.NORMAL,
			Piece: core.Piece{Type: core.KING, Color: core.WHITE},
			From:  core.E1,
			To:    core.E2,
		}
		roundTrip(t, ctx, move)
	})

	// =========================================================================
	// Multi-move sequence — verify the hash tracks across a sequence of moves.
	// =========================================================================

	t.Run("a sequence of moves and undos preserves the hash", func(t *testing.T) {
		ctx := decode(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
		originalHash := hasher.InitHash(ctx)

		// Make 3 moves, tracking snapshots and hashes.
		moves := []core.Move{
			{Type: core.NORMAL, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}, From: core.E2, To: core.E4},
			{Type: core.NORMAL, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}, From: core.E7, To: core.E5},
			{Type: core.NORMAL, Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE}, From: core.G1, To: core.F3},
		}

		var snapshots []core.Snapshot
		var hashes []uint64
		currentHash := originalHash

		for _, move := range moves {
			snap := eng.Apply(ctx, move)
			snapshots = append(snapshots, snap)
			moveHash := core.NewMoveHash(snap, *ctx)
			currentHash = hasher.Hash(currentHash, moveHash)
			hashes = append(hashes, currentHash)
		}

		// Undo in reverse order.
		for i := len(moves) - 1; i >= 0; i-- {
			moveHash := core.NewMoveHash(snapshots[i], *ctx)
			currentHash = hasher.Hash(currentHash, moveHash)
			eng.Undo(ctx, snapshots[i])
		}

		if currentHash != originalHash {
			t.Errorf("hash after 3-move round-trip = %d, want %d", currentHash, originalHash)
		}

		// Verify the board is back to the start.
		restoredHash := hasher.InitHash(ctx)
		if restoredHash != originalHash {
			t.Errorf("board not restored: init=%d, restored=%d", originalHash, restoredHash)
		}
	})

	// =========================================================================
	// Transposition — same position reached by different move orders must
	// produce the same hash.
	// =========================================================================

	t.Run("transposition: same position via different move orders has the same hash", func(t *testing.T) {
		// Path 1: Nb1-c3 first, then Ng1-f3
		// Path 2: Ng1-f3 first, then Nb1-c3
		// Neither move sets an EP target, so both paths reach the exact same
		// position (same pieces, same EP, same castling, same side to move).
		ctx1 := decode(t, "4k3/8/8/8/8/8/8/1N1QK1N1 w - - 0 1")
		h1Start := hasher.InitHash(ctx1)
		snap1a := eng.Apply(ctx1, core.Move{Type: core.NORMAL, Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE}, From: core.B1, To: core.C3})
		h1AfterNc3 := hasher.Hash(h1Start, core.NewMoveHash(snap1a, *ctx1))
		snap1b := eng.Apply(ctx1, core.Move{Type: core.NORMAL, Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE}, From: core.G1, To: core.F3})
		h1Final := hasher.Hash(h1AfterNc3, core.NewMoveHash(snap1b, *ctx1))

		ctx2 := decode(t, "4k3/8/8/8/8/8/8/1N1QK1N1 w - - 0 1")
		h2Start := hasher.InitHash(ctx2)
		snap2a := eng.Apply(ctx2, core.Move{Type: core.NORMAL, Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE}, From: core.G1, To: core.F3})
		h2AfterNf3 := hasher.Hash(h2Start, core.NewMoveHash(snap2a, *ctx2))
		snap2b := eng.Apply(ctx2, core.Move{Type: core.NORMAL, Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE}, From: core.B1, To: core.C3})
		h2Final := hasher.Hash(h2AfterNf3, core.NewMoveHash(snap2b, *ctx2))

		if h1Final != h2Final {
			t.Errorf("transposition hash mismatch: Nc3-Nf3=%d, Nf3-Nc3=%d", h1Final, h2Final)
		}
	})

	// =========================================================================
	// En passant target edge case: the EP target must be correctly hashed
	// in and out.
	// =========================================================================

	t.Run("a double pawn push then its undo restores the exact original hash", func(t *testing.T) {
		ctx := decode(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
		originalHash := hasher.InitHash(ctx)

		// e2-e4 (double push, sets EP target on e3)
		move := core.Move{Type: core.NORMAL, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}, From: core.E2, To: core.E4}
		snap := eng.Apply(ctx, move)
		moveHash := core.NewMoveHash(snap, *ctx)
		afterHash := hasher.Hash(originalHash, moveHash)

		// The hash should have changed (EP target was set).
		if afterHash == originalHash {
			t.Errorf("hash unchanged after double pawn push (EP target should affect it)")
		}

		// Undo.
		revertedHash := hasher.Hash(afterHash, moveHash)
		eng.Undo(ctx, snap)

		if revertedHash != originalHash {
			t.Errorf("hash after undo = %d, want %d", revertedHash, originalHash)
		}
	})
}
