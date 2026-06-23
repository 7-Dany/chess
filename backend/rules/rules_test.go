package rules

import (
	"testing"

	"github.com/7-Dany/chess/core"
	"github.com/7-Dany/chess/engine"
	"github.com/7-Dany/chess/fen"
	"github.com/7-Dany/chess/tracker"
)

// TestRules verifies every game-ending condition: fifty-move rule,
// threefold repetition, insufficient material, checkmate, stalemate, and
// the fused GetGameResult that checks them all in cheapest-first order.
func TestRules(t *testing.T) {
	rules := GetDefaultRules()
	eng := engine.GetDefaultEngine()

	decode := func(t *testing.T, fenStr string) core.TurnContext {
		t.Helper()
		var ctx core.TurnContext
		if err := fen.GetDefaultFenParser().Decode(fenStr, &ctx); err != nil {
			t.Fatalf("Decode failed: %v", err)
		}
		return ctx
	}

	// =========================================================================
	// IsFiftyMoveRule
	// =========================================================================

	t.Run("IsFiftyMoveRule: halfmove clock below 100 returns false", func(t *testing.T) {
		ctx := decode(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
		if rules.IsFiftyMoveRule(ctx) {
			t.Errorf("IsFiftyMoveRule should be false when HalfMoveClock = 0")
		}
	})

	t.Run("IsFiftyMoveRule: halfmove clock at 99 returns false", func(t *testing.T) {
		ctx := decode(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 99 50")
		if rules.IsFiftyMoveRule(ctx) {
			t.Errorf("IsFiftyMoveRule should be false when HalfMoveClock = 99")
		}
	})

	t.Run("IsFiftyMoveRule: halfmove clock at 100 returns true", func(t *testing.T) {
		ctx := decode(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 100 50")
		if !rules.IsFiftyMoveRule(ctx) {
			t.Errorf("IsFiftyMoveRule should be true when HalfMoveClock = 100")
		}
	})

	t.Run("IsFiftyMoveRule: halfmove clock above 100 returns true", func(t *testing.T) {
		ctx := decode(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 150 75")
		if !rules.IsFiftyMoveRule(ctx) {
			t.Errorf("IsFiftyMoveRule should be true when HalfMoveClock = 150")
		}
	})

	// =========================================================================
	// IsThreefoldRepetition
	// =========================================================================

	t.Run("IsThreefoldRepetition: count below 3 returns false", func(t *testing.T) {
		tr := tracker.NewPositionTracker()
		tr.Record(42)
		tr.Record(42)
		if rules.IsThreefoldRepetition(tr, 42) {
			t.Errorf("should be false when count = 2")
		}
	})

	t.Run("IsThreefoldRepetition: count at 3 returns true", func(t *testing.T) {
		tr := tracker.NewPositionTracker()
		tr.Record(42)
		tr.Record(42)
		tr.Record(42)
		if !rules.IsThreefoldRepetition(tr, 42) {
			t.Errorf("should be true when count = 3")
		}
	})

	t.Run("IsThreefoldRepetition: count above 3 returns true", func(t *testing.T) {
		tr := tracker.NewPositionTracker()
		for range 5 {
			tr.Record(42)
		}
		if !rules.IsThreefoldRepetition(tr, 42) {
			t.Errorf("should be true when count = 5")
		}
	})

	t.Run("IsThreefoldRepetition: a hash that was never recorded returns false", func(t *testing.T) {
		tr := tracker.NewPositionTracker()
		tr.Record(42)
		if rules.IsThreefoldRepetition(tr, 99) {
			t.Errorf("should be false for unrecorded hash")
		}
	})

	// =========================================================================
	// IsInsufficientMaterial
	// =========================================================================

	t.Run("IsInsufficientMaterial: K vs K returns true", func(t *testing.T) {
		ctx := decode(t, "4k3/8/8/8/8/8/8/4K3 w - - 0 1")
		if !rules.IsInsufficientMaterial(ctx) {
			t.Errorf("K vs K should be insufficient")
		}
	})

	t.Run("IsInsufficientMaterial: K+N vs K returns true", func(t *testing.T) {
		ctx := decode(t, "4k3/8/8/8/8/8/8/3NK3 w - - 0 1")
		if !rules.IsInsufficientMaterial(ctx) {
			t.Errorf("K+N vs K should be insufficient")
		}
	})

	t.Run("IsInsufficientMaterial: K+B vs K returns true", func(t *testing.T) {
		ctx := decode(t, "4k3/8/8/8/8/8/8/3BK3 w - - 0 1")
		if !rules.IsInsufficientMaterial(ctx) {
			t.Errorf("K+B vs K should be insufficient")
		}
	})

	t.Run("IsInsufficientMaterial: K+B vs K+B same color returns true", func(t *testing.T) {
		// White bishop on C1 (file 2, rank 0 → 2+0=2 even → dark).
		// Black bishop on F8 (file 5, rank 7 → 5+7=12 even → dark).
		// Both dark → insufficient.
		ctx := decode(t, "5b2/8/8/8/8/8/8/2B1K2k w - - 0 1")
		if !rules.IsInsufficientMaterial(ctx) {
			t.Errorf("K+B vs K+B same color should be insufficient")
		}
	})

	t.Run("IsInsufficientMaterial: K+B vs K+B opposite colors returns false", func(t *testing.T) {
		// White bishop on C1 (file 2, rank 0 → 2+0=2 even → dark).
		// Black bishop on C8 (file 2, rank 7 → 2+7=9 odd → light).
		// Different colors → sufficient (can potentially win).
		ctx := decode(t, "2b1k3/8/8/8/8/8/8/2B1K3 w - - 0 1")
		if rules.IsInsufficientMaterial(ctx) {
			t.Errorf("K+B vs K+B opposite colors should be sufficient")
		}
	})

	t.Run("IsInsufficientMaterial: K+N+N vs K returns false (can't force mate but not a rule)", func(t *testing.T) {
		// Two knights can't FORCE checkmate, but it's not classified as
		// insufficient material (it's possible to checkmate with help).
		ctx := decode(t, "4k3/8/8/8/8/8/8/2NNK3 w - - 0 1")
		if rules.IsInsufficientMaterial(ctx) {
			t.Errorf("K+N+N vs K should be sufficient (not classified as insufficient)")
		}
	})

	t.Run("IsInsufficientMaterial: K+P vs K returns false", func(t *testing.T) {
		ctx := decode(t, "4k3/8/8/8/8/8/4P3/4K3 w - - 0 1")
		if rules.IsInsufficientMaterial(ctx) {
			t.Errorf("K+P vs K should be sufficient (pawn can promote)")
		}
	})

	t.Run("IsInsufficientMaterial: K+Q vs K returns false", func(t *testing.T) {
		ctx := decode(t, "4k3/8/8/8/8/8/8/3QK3 w - - 0 1")
		if rules.IsInsufficientMaterial(ctx) {
			t.Errorf("K+Q vs K should be sufficient")
		}
	})

	t.Run("IsInsufficientMaterial: K+R vs K returns false", func(t *testing.T) {
		ctx := decode(t, "4k3/8/8/8/8/8/8/3RK3 w - - 0 1")
		if rules.IsInsufficientMaterial(ctx) {
			t.Errorf("K+R vs K should be sufficient")
		}
	})

	t.Run("IsInsufficientMaterial: starting position returns false", func(t *testing.T) {
		ctx := decode(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
		if rules.IsInsufficientMaterial(ctx) {
			t.Errorf("starting position should be sufficient")
		}
	})

	// =========================================================================
	// IsCheckMate
	// =========================================================================

	t.Run("IsCheckMate: fool's mate position returns true", func(t *testing.T) {
		// After 1.f3 e5 2.g4 Qh4# — white is checkmated.
		ctx := decode(t, "rnb1kbnr/pppp1ppp/8/4p3/6Pq/5P2/PPPPP2P/RNBQKBNR w KQkq - 1 3")
		if !rules.IsCheckMate(ctx, eng) {
			t.Errorf("white should be checkmated")
		}
	})

	t.Run("IsCheckMate: starting position returns false", func(t *testing.T) {
		ctx := decode(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
		if rules.IsCheckMate(ctx, eng) {
			t.Errorf("starting position should not be checkmate")
		}
	})

	t.Run("IsCheckMate: a stalemated position returns false (not in check)", func(t *testing.T) {
		// Black king on A8, white queen on C7, white king on C6 — stalemate.
		ctx := decode(t, "k7/2Q5/2K5/8/8/8/8/8 b - - 0 1")
		if rules.IsCheckMate(ctx, eng) {
			t.Errorf("stalemate should not be checkmate")
		}
	})

	// =========================================================================
	// IsStaleMate
	// =========================================================================

	t.Run("IsStaleMate: a stalemated position returns true", func(t *testing.T) {
		ctx := decode(t, "k7/2Q5/2K5/8/8/8/8/8 b - - 0 1")
		if !rules.IsStaleMate(ctx, eng) {
			t.Errorf("black should be stalemated")
		}
	})

	t.Run("IsStaleMate: starting position returns false", func(t *testing.T) {
		ctx := decode(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
		if rules.IsStaleMate(ctx, eng) {
			t.Errorf("starting position should not be stalemate")
		}
	})

	t.Run("IsStaleMate: a checkmated position returns false (in check)", func(t *testing.T) {
		ctx := decode(t, "rnb1kbnr/pppp1ppp/8/4p3/6Pq/5P2/PPPPP2P/RNBQKBNR w KQkq - 1 3")
		if rules.IsStaleMate(ctx, eng) {
			t.Errorf("checkmate should not be stalemate")
		}
	})

	// =========================================================================
	// GetGameResult — the fused evaluator.
	// =========================================================================

	t.Run("GetGameResult: starting position returns InProgress", func(t *testing.T) {
		ctx := decode(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
		tr := tracker.NewPositionTracker()
		result := rules.GetGameResult(ctx, eng, tr, 0)
		if result.Status != core.InProgress {
			t.Errorf("Status = %v, want InProgress", result.Status)
		}
	})

	t.Run("GetGameResult: fifty-move rule returns Draw with FiftyMoveRule reason", func(t *testing.T) {
		ctx := decode(t, "4k3/8/8/8/8/8/8/4K3 w - - 100 50")
		tr := tracker.NewPositionTracker()
		result := rules.GetGameResult(ctx, eng, tr, 0)
		if result.Status != core.Draw {
			t.Errorf("Status = %v, want Draw", result.Status)
		}
		if result.DrawReason != core.FiftyMoveRule {
			t.Errorf("DrawReason = %v, want FiftyMoveRule", result.DrawReason)
		}
	})

	t.Run("GetGameResult: threefold repetition returns Draw with ThreefoldRepetition reason", func(t *testing.T) {
		ctx := decode(t, "4k3/8/8/8/8/8/8/4K3 w - - 0 1")
		tr := tracker.NewPositionTracker()
		tr.Record(42)
		tr.Record(42)
		tr.Record(42)
		result := rules.GetGameResult(ctx, eng, tr, 42)
		if result.Status != core.Draw {
			t.Errorf("Status = %v, want Draw", result.Status)
		}
		if result.DrawReason != core.ThreefoldRepetition {
			t.Errorf("DrawReason = %v, want ThreefoldRepetition", result.DrawReason)
		}
	})

	t.Run("GetGameResult: insufficient material returns Draw with InsufficientMaterial reason", func(t *testing.T) {
		ctx := decode(t, "4k3/8/8/8/8/8/8/4K3 w - - 0 1")
		tr := tracker.NewPositionTracker()
		result := rules.GetGameResult(ctx, eng, tr, 0)
		if result.Status != core.Draw {
			t.Errorf("Status = %v, want Draw", result.Status)
		}
		if result.DrawReason != core.InsufficientMaterial {
			t.Errorf("DrawReason = %v, want InsufficientMaterial", result.DrawReason)
		}
	})

	t.Run("GetGameResult: checkmate returns CheckMate with the winner set", func(t *testing.T) {
		// Fool's mate — white is checkmated, black wins.
		ctx := decode(t, "rnb1kbnr/pppp1ppp/8/4p3/6Pq/5P2/PPPPP2P/RNBQKBNR w KQkq - 1 3")
		tr := tracker.NewPositionTracker()
		result := rules.GetGameResult(ctx, eng, tr, 0)
		if result.Status != core.CheckMate {
			t.Errorf("Status = %v, want CheckMate", result.Status)
		}
		if !result.HasWinner {
			t.Errorf("HasWinner should be true on checkmate")
		}
		if result.Winner != core.BLACK {
			t.Errorf("Winner = %v, want BLACK", result.Winner)
		}
	})

	t.Run("GetGameResult: stalemate returns Draw with Stalemate reason", func(t *testing.T) {
		ctx := decode(t, "k7/2Q5/2K5/8/8/8/8/8 b - - 0 1")
		tr := tracker.NewPositionTracker()
		result := rules.GetGameResult(ctx, eng, tr, 0)
		if result.Status != core.Draw {
			t.Errorf("Status = %v, want Draw", result.Status)
		}
		if result.DrawReason != core.Stalemate {
			t.Errorf("DrawReason = %v, want Stalemate", result.DrawReason)
		}
	})

	// =========================================================================
	// GetGameResult: evaluation order — cheapest first.
	// FiftyMoveRule is checked before checkmate/stalemate.
	// =========================================================================

	t.Run("GetGameResult: fifty-move rule takes priority over checkmate", func(t *testing.T) {
		// A checkmate position with HalfMoveClock = 100.
		// The fifty-move rule should win (it's checked first).
		ctx := decode(t, "rnb1kbnr/pppp1ppp/8/4p3/6Pq/5P2/PPPPP2P/RNBQKBNR w KQkq - 100 50")
		tr := tracker.NewPositionTracker()
		result := rules.GetGameResult(ctx, eng, tr, 0)
		if result.DrawReason != core.FiftyMoveRule {
			t.Errorf("DrawReason = %v, want FiftyMoveRule (checked before checkmate)", result.DrawReason)
		}
	})
}
