package game

import (
	"testing"

	"github.com/7-Dany/chess/core"
)

// TestChess verifies the top-level orchestrator: New, MakeMove, UndoMove,
// GameResult, LegalMoves, Hash, and the full 1v1 game lifecycle.
func TestChess(t *testing.T) {
	// =========================================================================
	// New — construction and bootstrap.
	// =========================================================================

	t.Run("New creates a game at the starting position", func(t *testing.T) {
		g, err := New()
		if err != nil {
			t.Fatalf("New() error: %v", err)
		}
		ctx := g.TurnContext()
		if ctx.SideToMove != core.WHITE {
			t.Errorf("SideToMove = %v, want WHITE", ctx.SideToMove)
		}
		// Spot-check a few squares.
		if square := ctx.Board[core.E1]; !square.IsOccupied() || square.Type() != core.KING || square.Color() != core.WHITE {
			t.Errorf("E1 should have white king, got %v", square)
		}
		if square := ctx.Board[core.E8]; !square.IsOccupied() || square.Type() != core.KING || square.Color() != core.BLACK {
			t.Errorf("E8 should have black king, got %v", square)
		}
	})

	t.Run("New with a custom FEN starts at that position", func(t *testing.T) {
		g, err := New(WithFEN("4k3/8/8/8/8/8/8/4K3 w - - 0 1"))
		if err != nil {
			t.Fatalf("New() error: %v", err)
		}
		ctx := g.TurnContext()
		// Only two kings on the board.
		if square := ctx.Board[core.E1]; !square.IsOccupied() || square.Type() != core.KING {
			t.Errorf("E1 should have a king")
		}
		if square := ctx.Board[core.A1]; square.IsOccupied() {
			t.Errorf("A1 should be empty")
		}
	})

	t.Run("New with an invalid FEN returns an error", func(t *testing.T) {
		_, err := New(WithFEN("invalid fen string"))
		if err == nil {
			t.Errorf("New() should return an error for an invalid FEN")
		}
	})

	// =========================================================================
	// Hash — incremental Zobrist, maintained across moves.
	// =========================================================================

	t.Run("the initial hash is non-zero", func(t *testing.T) {
		g, _ := New()
		if g.Hash() == 0 {
			t.Errorf("initial hash should be non-zero")
		}
	})

	t.Run("the hash changes after a move", func(t *testing.T) {
		g, _ := New()
		h1 := g.Hash()

		if err := g.MakeMove(core.Move{
			Type:  core.NORMAL,
			Piece: core.Piece{Type: core.PAWN, Color: core.WHITE},
			From:  core.E2,
			To:    core.E4,
		}); err != nil {
			t.Fatalf("MakeMove failed: %v", err)
		}

		h2 := g.Hash()
		if h1 == h2 {
			t.Errorf("hash should change after a move (both = %d)", h1)
		}
	})

	t.Run("the hash is restored after undo", func(t *testing.T) {
		g, _ := New()
		h1 := g.Hash()

		g.MakeMove(core.Move{
			Type:  core.NORMAL,
			Piece: core.Piece{Type: core.PAWN, Color: core.WHITE},
			From:  core.E2,
			To:    core.E4,
		})
		g.UndoMove()

		if g.Hash() != h1 {
			t.Errorf("hash after undo = %d, want %d", g.Hash(), h1)
		}
	})

	// =========================================================================
	// MakeMove — validation and application.
	// =========================================================================

	t.Run("a legal move is applied successfully", func(t *testing.T) {
		g, _ := New()
		err := g.MakeMove(core.Move{
			Type:  core.NORMAL,
			Piece: core.Piece{Type: core.PAWN, Color: core.WHITE},
			From:  core.E2,
			To:    core.E4,
		})
		if err != nil {
			t.Errorf("MakeMove returned error for a legal move: %v", err)
		}
		ctx := g.TurnContext()
		if ctx.Board[core.E2].IsOccupied() {
			t.Errorf("E2 should be empty after e4")
		}
		if !ctx.Board[core.E4].IsOccupied() || ctx.Board[core.E4].Type() != core.PAWN {
			t.Errorf("E4 should have a pawn after e4")
		}
	})

	t.Run("an illegal move returns ErrIllegalMove", func(t *testing.T) {
		g, _ := New()
		// White pawn on E2 cannot jump to E5 in one move.
		err := g.MakeMove(core.Move{
			Type:  core.NORMAL,
			Piece: core.Piece{Type: core.PAWN, Color: core.WHITE},
			From:  core.E2,
			To:    core.E5,
		})
		if err != ErrIllegalMove {
			t.Errorf("MakeMove returned %v, want ErrIllegalMove", err)
		}
		// Board should be unchanged.
		ctx := g.TurnContext()
		if !ctx.Board[core.E2].IsOccupied() {
			t.Errorf("E2 should still have a pawn (move was rejected)")
		}
	})

	t.Run("a move that leaves the king in check is rejected", func(t *testing.T) {
		g, _ := New(WithFEN("4k3/8/8/8/8/8/4r3/4K3 w - - 0 1"))
		// White king on E1, black rook on E2 gives check.
		// King cannot move to D2 or F2 — those are attacked by the rook's
		// adjacent squares? No, the rook on E2 attacks the entire E-file
		// and rank 2. So D2 and F2 are on rank 2 → attacked.
		// King can only move to D1 or F1 (off the E-file and off rank 2).
		// Moving to E2 (capturing the rook) would be legal IF the rook
		// is undefended. It is undefended, so Kxe2 is legal.
		// Let's try a clearly illegal move: K to D2.
		err := g.MakeMove(core.Move{
			Type:  core.NORMAL,
			Piece: core.Piece{Type: core.KING, Color: core.WHITE},
			From:  core.E1,
			To:    core.D2,
		})
		if err != ErrIllegalMove {
			t.Errorf("moving the king into check should return ErrIllegalMove, got %v", err)
		}
	})

	// =========================================================================
	// UndoMove — reverting moves.
	// =========================================================================

	t.Run("undo reverts the last move", func(t *testing.T) {
		g, _ := New()
		g.MakeMove(core.Move{
			Type:  core.NORMAL,
			Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE},
			From:  core.B1,
			To:    core.C3,
		})
		err := g.UndoMove()
		if err != nil {
			t.Fatalf("UndoMove error: %v", err)
		}
		ctx := g.TurnContext()
		if !ctx.Board[core.B1].IsOccupied() || ctx.Board[core.B1].Type() != core.KNIGHT {
			t.Errorf("B1 should have a knight again after undo")
		}
		if ctx.Board[core.C3].IsOccupied() {
			t.Errorf("C3 should be empty after undo")
		}
	})

	t.Run("undo on an empty history returns ErrNothingToUndo", func(t *testing.T) {
		g, _ := New()
		err := g.UndoMove()
		if err != ErrNothingToUndo {
			t.Errorf("UndoMove on empty history = %v, want ErrNothingToUndo", err)
		}
	})

	t.Run("a full move-undo round-trip restores the exact position", func(t *testing.T) {
		g, _ := New()
		ctx1 := g.TurnContext()
		h1 := g.Hash()

		g.MakeMove(core.Move{
			Type:  core.NORMAL,
			Piece: core.Piece{Type: core.PAWN, Color: core.WHITE},
			From:  core.E2,
			To:    core.E4,
		})
		g.UndoMove()

		ctx2 := g.TurnContext()
		h2 := g.Hash()

		// Compare every square.
		for i := range 64 {
			pos := core.Position(i)
			if ctx1.Board[pos] != ctx2.Board[pos] {
				t.Errorf("board[%v] = %v, want %v (after round-trip)", pos, ctx2.Board[pos], ctx1.Board[pos])
			}
		}
		if h1 != h2 {
			t.Errorf("hash = %d, want %d (after round-trip)", h2, h1)
		}
		if ctx1.SideToMove != ctx2.SideToMove {
			t.Errorf("SideToMove = %v, want %v", ctx2.SideToMove, ctx1.SideToMove)
		}
	})

	// =========================================================================
	// GameResult — checkmate, stalemate, in-progress.
	// =========================================================================

	t.Run("the starting position is InProgress", func(t *testing.T) {
		g, _ := New()
		result := g.GameResult()
		if result.Status != core.InProgress {
			t.Errorf("Status = %v, want InProgress", result.Status)
		}
	})

	t.Run("a checkmate position returns CheckMate with the winner", func(t *testing.T) {
		// Fool's mate: white is checkmated.
		g, _ := New(WithFEN("rnb1kbnr/pppp1ppp/8/4p3/6Pq/5P2/PPPPP2P/RNBQKBNR w KQkq - 1 3"))
		result := g.GameResult()
		if result.Status != core.CheckMate {
			t.Errorf("Status = %v, want CheckMate", result.Status)
		}
		if result.Winner != core.BLACK {
			t.Errorf("Winner = %v, want BLACK", result.Winner)
		}
	})

	t.Run("a stalemate position returns Draw with Stalemate reason", func(t *testing.T) {
		g, _ := New(WithFEN("k7/2Q5/2K5/8/8/8/8/8 b - - 0 1"))
		result := g.GameResult()
		if result.Status != core.Draw {
			t.Errorf("Status = %v, want Draw", result.Status)
		}
		if result.DrawReason != core.Stalemate {
			t.Errorf("DrawReason = %v, want Stalemate", result.DrawReason)
		}
	})

	// =========================================================================
	// 1v1 game — a short game with moves, undos, and result detection.
	// =========================================================================

	t.Run("a short 1v1 game: e4 e5 Nf3 Nc6 Bc4 — all moves apply and undo", func(t *testing.T) {
		g, _ := New()
		moves := []core.Move{
			{Type: core.NORMAL, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}, From: core.E2, To: core.E4},
			{Type: core.NORMAL, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}, From: core.E7, To: core.E5},
			{Type: core.NORMAL, Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE}, From: core.G1, To: core.F3},
			{Type: core.NORMAL, Piece: core.Piece{Type: core.KNIGHT, Color: core.BLACK}, From: core.B8, To: core.C6},
			{Type: core.NORMAL, Piece: core.Piece{Type: core.BISHOP, Color: core.WHITE}, From: core.F1, To: core.C4},
		}

		// Play all moves.
		for _, m := range moves {
			if err := g.MakeMove(m); err != nil {
				t.Fatalf("MakeMove(%v→%v) failed: %v", m.From, m.To, err)
			}
		}

		// After 5 moves, it's black's turn (move 3).
		ctx := g.TurnContext()
		if ctx.SideToMove != core.BLACK {
			t.Errorf("SideToMove = %v, want BLACK (after 5 half-moves)", ctx.SideToMove)
		}

		// Undo all 5 moves.
		for range 5 {
			if err := g.UndoMove(); err != nil {
				t.Fatalf("UndoMove failed: %v", err)
			}
		}

		// Should be back to the start.
		ctx = g.TurnContext()
		if ctx.SideToMove != core.WHITE {
			t.Errorf("SideToMove = %v, want WHITE (after full undo)", ctx.SideToMove)
		}
		if !ctx.Board[core.E2].IsOccupied() || ctx.Board[core.E2].Type() != core.PAWN {
			t.Errorf("E2 should have a pawn again (after full undo)")
		}
	})

	t.Run("a scholar's mate game ends in checkmate", func(t *testing.T) {
		// Scholar's mate: 1.e4 e5 2.Bc4 Nc6 3.Qh5 Nf6?? 4.Qxf7#
		g, _ := New()
		moves := []core.Move{
			{Type: core.NORMAL, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}, From: core.E2, To: core.E4},
			{Type: core.NORMAL, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}, From: core.E7, To: core.E5},
			{Type: core.NORMAL, Piece: core.Piece{Type: core.BISHOP, Color: core.WHITE}, From: core.F1, To: core.C4},
			{Type: core.NORMAL, Piece: core.Piece{Type: core.KNIGHT, Color: core.BLACK}, From: core.B8, To: core.C6},
			{Type: core.NORMAL, Piece: core.Piece{Type: core.QUEEN, Color: core.WHITE}, From: core.D1, To: core.H5},
			{Type: core.NORMAL, Piece: core.Piece{Type: core.KNIGHT, Color: core.BLACK}, From: core.G8, To: core.F6},
			{Type: core.NORMAL, Piece: core.Piece{Type: core.QUEEN, Color: core.WHITE}, From: core.H5, To: core.F7, HasCapture: true, Captured: core.Piece{Type: core.PAWN, Color: core.BLACK}},
		}

		for i, m := range moves {
			if err := g.MakeMove(m); err != nil {
				t.Fatalf("move %d (%v→%v) failed: %v", i, m.From, m.To, err)
			}
		}

		result := g.GameResult()
		if result.Status != core.CheckMate {
			t.Errorf("Status = %v, want CheckMate (scholar's mate)", result.Status)
		}
		if result.Winner != core.WHITE {
			t.Errorf("Winner = %v, want WHITE", result.Winner)
		}
	})

	// =========================================================================
	// LegalMoves — querying available moves.
	// =========================================================================

	t.Run("LegalMoves for a knight on B1 in the starting position returns 2 moves", func(t *testing.T) {
		g, _ := New()
		moves := g.LegalMoves(core.B1)
		if len(moves) != 2 {
			t.Errorf("LegalMoves(B1) = %d moves, want 2", len(moves))
		}
	})

	t.Run("LegalMoves for an empty square returns no moves", func(t *testing.T) {
		g, _ := New()
		moves := g.LegalMoves(core.E4)
		if len(moves) != 0 {
			t.Errorf("LegalMoves(E4) = %d moves, want 0 (empty square)", len(moves))
		}
	})
}
