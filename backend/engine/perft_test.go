package engine

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/7-Dany/chess/core"
	"github.com/7-Dany/chess/fen"
)

// perft counts the number of leaf nodes in the full move tree to the given
// depth. It is the standard chess engine correctness test: generate all
// legal moves, apply each, recurse, undo. The result is compared against
// trusted reference values from the Chess Programming Wiki.
//
// Reference: https://www.chessprogramming.org/Perft_Results
//
// Note: Apply does not flip SideToMove (that's the game controller's job),
// so perft flips it manually after each Apply and flips it back before Undo.
func perft(engine DefaultEngine, ctx *core.TurnContext, depth int) uint64 {
	var buf [MAX_TOTAL_MOVES]core.Move
	moves := engine.GetAllLegalMoves(buf[:0], *ctx)

	if depth == 1 {
		return uint64(len(moves))
	}

	var nodes uint64
	for _, move := range moves {
		snap := engine.Apply(ctx, move)
		ctx.SideToMove = ctx.SideToMove.Opponent()
		nodes += perft(engine, ctx, depth-1)
		ctx.SideToMove = ctx.SideToMove.Opponent()
		engine.Undo(ctx, snap)
	}
	return nodes
}

// perftParallel is the parallel variant of perft. It spawns one goroutine
// per top-level move, each working on its own independent TurnContext copy,
// then sums the results. Sub-trees are explored sequentially by perft —
// parallelising deeper levels pays diminishing returns against goroutine
// overhead, while the top level alone gives near-linear scaling with core count.
func perftParallel(engine DefaultEngine, ctx core.TurnContext, depth int) uint64 {
	var buf [MAX_TOTAL_MOVES]core.Move
	moves := engine.GetAllLegalMoves(buf[:0], ctx)

	if depth == 1 {
		return uint64(len(moves))
	}

	var total atomic.Uint64
	var wg sync.WaitGroup

	for _, move := range moves {
		wg.Add(1)
		go func(move core.Move) {
			defer wg.Done()
			child := ctx.Copy()
			engine.Apply(&child, move)
			child.SideToMove = child.SideToMove.Opponent()
			total.Add(perft(engine, &child, depth-1))
		}(move)
	}

	wg.Wait()
	return total.Load()
}

// TestPerft validates the engine against the six standard perft positions
// from the Chess Programming Wiki. These positions stress every edge case:
// en passant (including discovered-check en passant), promotion, castling
// through/into check, pins, and double-checks.
//
// Reference: https://www.chessprogramming.org/Perft_Results
//
// Depth 1-3 runs in ~2 seconds and catches most move-generation bugs.
// Depth 4 for positions 1, 3, and 4 (the smaller node counts) is included
// below. Deeper depths (4M+ nodes) can be run manually:
//
//	go test ./engine -run TestPerft/position_2 -v -timeout=120s
func TestPerft(t *testing.T) {
	engine := GetDefaultEngine()

	decode := func(fenStr string) core.TurnContext {
		var ctx core.TurnContext
		if err := fen.GetDefaultFenParser().Decode(fenStr, &ctx); err != nil {
			t.Fatalf("Decode failed: %v", err)
		}
		return ctx
	}

	// assertPerft decodes fenStr, runs perftParallel at depth, and checks
	// the node count. Uses the parallel variant so each depth subtest
	// exploits all available cores.
	assertPerft := func(t *testing.T, fenStr string, depth int, want uint64) {
		t.Helper()
		t.Parallel()
		ctx := decode(fenStr)
		got := perftParallel(engine, ctx, depth)
		if got != want {
			t.Errorf("perft depth %d = %d, want %d", depth, got, want)
		}
	}

	// =========================================================================
	// Position 1 — the standard starting position.
	// FEN: rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1
	// =========================================================================

	t.Run("position 1 (start)", func(t *testing.T) {
		t.Parallel()
		const fen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

		t.Run("depth 1 = 20", func(t *testing.T) { assertPerft(t, fen, 1, 20) })
		t.Run("depth 2 = 400", func(t *testing.T) { assertPerft(t, fen, 2, 400) })
		t.Run("depth 3 = 8902", func(t *testing.T) { assertPerft(t, fen, 3, 8902) })
		t.Run("depth 4 = 197281", func(t *testing.T) { assertPerft(t, fen, 4, 197281) })
	})

	// =========================================================================
	// Position 2 — "Kiwipete". A dense middlegame position with many
	// sliders, knights, and castling rights. The standard stress test.
	// FEN: r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1
	// =========================================================================

	t.Run("position 2 (Kiwipete)", func(t *testing.T) {
		t.Parallel()
		const fen = "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1"

		t.Run("depth 1 = 48", func(t *testing.T) { assertPerft(t, fen, 1, 48) })
		t.Run("depth 2 = 2039", func(t *testing.T) { assertPerft(t, fen, 2, 2039) })
		t.Run("depth 3 = 97862", func(t *testing.T) { assertPerft(t, fen, 3, 97862) })
	})

	// =========================================================================
	// Position 3 — an endgame position that tests en passant discovered
	// checks. The black rook on H5 and the white king on A5 are on the same
	// rank; an en passant capture that removes both pawns from the rank would
	// expose the king to the rook.
	// FEN: 8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1
	// =========================================================================

	t.Run("position 3 (en passant edge cases)", func(t *testing.T) {
		t.Parallel()
		const fen = "8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1"

		t.Run("depth 1 = 14", func(t *testing.T) { assertPerft(t, fen, 1, 14) })
		t.Run("depth 2 = 191", func(t *testing.T) { assertPerft(t, fen, 2, 191) })
		t.Run("depth 3 = 2812", func(t *testing.T) { assertPerft(t, fen, 3, 2812) })
		t.Run("depth 4 = 43238", func(t *testing.T) { assertPerft(t, fen, 4, 43238) })
	})

	// =========================================================================
	// Position 4 — tests promotion, castling under attack, and pins.
	// FEN: r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1
	// =========================================================================

	t.Run("position 4 (promotion and pins)", func(t *testing.T) {
		t.Parallel()
		const fen = "r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1"

		t.Run("depth 1 = 6", func(t *testing.T) { assertPerft(t, fen, 1, 6) })
		t.Run("depth 2 = 264", func(t *testing.T) { assertPerft(t, fen, 2, 264) })
		t.Run("depth 3 = 9467", func(t *testing.T) { assertPerft(t, fen, 3, 9467) })
		t.Run("depth 4 = 422333", func(t *testing.T) { assertPerft(t, fen, 4, 422333) })
	})

	// =========================================================================
	// Position 5 — tests promotion-captures and castling availability.
	// FEN: rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8
	// =========================================================================

	t.Run("position 5 (promotion-captures)", func(t *testing.T) {
		t.Parallel()
		const fen = "rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8"

		t.Run("depth 1 = 44", func(t *testing.T) { assertPerft(t, fen, 1, 44) })
		t.Run("depth 2 = 1486", func(t *testing.T) { assertPerft(t, fen, 2, 1486) })
		t.Run("depth 3 = 62379", func(t *testing.T) { assertPerft(t, fen, 3, 62379) })
	})

	// =========================================================================
	// Position 6 — a complex middlegame position. Tests all move types in
	// combination.
	// FEN: r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10
	// =========================================================================

	t.Run("position 6 (complex middlegame)", func(t *testing.T) {
		t.Parallel()
		const fen = "r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10"

		t.Run("depth 1 = 46", func(t *testing.T) { assertPerft(t, fen, 1, 46) })
		t.Run("depth 2 = 2079", func(t *testing.T) { assertPerft(t, fen, 2, 2079) })
		t.Run("depth 3 = 89890", func(t *testing.T) { assertPerft(t, fen, 3, 89890) })
	})
}
