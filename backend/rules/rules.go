// Package rules evaluates chess game-ending conditions and produces a
// core.GameResult after each move. It is a pure evaluator: every method is
// stateless and synchronous — no goroutines, no side effects, no knowledge of
// the Chess orchestrator.
//
// # Evaluation order
//
// GetGameResult short-circuits from cheapest to most expensive:
//
//  1. FiftyMoveRule        — O(1)  single integer comparison
//  2. ThreefoldRepetition  — O(1)  single map lookup
//  3. InsufficientMaterial — O(64) one pass over the board
//  4. HasAnyLegalMoves     — O(n)  pseudo-legal generation + king-safety filter
//     called at most once per result check;
//     if false → IsSquareAttacked on the king
//     determines checkmate vs stalemate.
//
// The individual IsCheckMate and IsStaleMate methods each call HasAnyLegalMoves
// internally. Never call both back-to-back — use GetGameResult instead, which
// fuses them into a single HasAnyLegalMoves call.
//
// # Dependency contract
//
// Each method receives only the data it actually reads:
//
//	IsFiftyMoveRule           → ctx.HalfMoveClock
//	IsThreefoldRepetition     → tracker, hash
//	IsInsufficientMaterial    → ctx.Board
//	IsCheckMate / IsStaleMate → ctx, engine
//	GetGameResult             → all of the above
//
// No method accepts a *Chess pointer. Keeping Rules decoupled from the
// orchestrator means it can be tested in isolation and reused by any
// caller that can supply a TurnContext and an Engine.
package rules

import (
	"github.com/7-Dany/chess/core"
	"github.com/7-Dany/chess/engine"
	"github.com/7-Dany/chess/tracker"
)

// Rules evaluates whether the game has ended and why.
// Implementations must be safe to call after every half-move.
type Rules interface {
	// IsFiftyMoveRule reports whether the fifty-move rule applies.
	// It returns true when ctx.HalfMoveClock has reached 100 half-moves
	// (50 full moves) without a pawn push or capture, which entitles
	// either player to claim a draw.
	IsFiftyMoveRule(ctx core.TurnContext) bool

	// IsThreefoldRepetition reports whether the current position has
	// occurred at least three times in the game. It delegates to tracker
	// with the Zobrist hash of the current position. O(1) — single map lookup.
	IsThreefoldRepetition(tracker tracker.Tracker, hash uint64) bool

	// IsInsufficientMaterial reports whether neither side has enough pieces
	// to force checkmate. The following material combinations are drawn:
	//
	//	K vs K
	//	K + N vs K
	//	K + B vs K
	//	K + B vs K + B  (only when both bishops share the same square color)
	//
	// Any pawn, rook, or queen on the board means material is sufficient.
	IsInsufficientMaterial(ctx core.TurnContext) bool

	// IsCheckMate reports whether the side to move is in checkmate:
	// their king is in check and they have no legal move to escape it.
	// Calls engine.HasAnyLegalMoves internally.
	//
	// Prefer GetGameResult when you need a full evaluation — it fuses
	// this check with IsStaleMate to avoid a redundant HasAnyLegalMoves call.
	IsCheckMate(ctx core.TurnContext, engine engine.Engine) bool

	// IsStaleMate reports whether the side to move is in stalemate:
	// they have no legal move but their king is not in check.
	// Calls engine.HasAnyLegalMoves internally.
	//
	// Prefer GetGameResult when you need a full evaluation — it fuses
	// this check with IsCheckMate to avoid a redundant HasAnyLegalMoves call.
	IsStaleMate(ctx core.TurnContext, engine engine.Engine) bool

	// GetGameResult evaluates all termination conditions in cheapest-first
	// order and returns the current game result. Call this after every move.
	//
	// Returns core.GameResult with Status == InProgress when the game
	// continues. When the game is over, Status is either CheckMate or Draw;
	// a Draw carries a DrawReason identifying which rule triggered.
	//
	// HasAnyLegalMoves is called at most once, regardless of how many
	// conditions are checked — see package-level evaluation order.
	GetGameResult(ctx core.TurnContext, engine engine.Engine, tracker tracker.Tracker, hash uint64) core.GameResult
}

type DefaultRules struct{}

var defaultRules = DefaultRules{}

func GetDefaultRules() DefaultRules {
	return defaultRules
}

func (DefaultRules) IsFiftyMoveRule(ctx core.TurnContext) bool {
	return ctx.HalfMoveClock >= 100
}

func (DefaultRules) IsThreefoldRepetition(tracker tracker.Tracker, hash uint64) bool {
	return tracker.Count(hash) >= 3
}

func (DefaultRules) IsInsufficientMaterial(ctx core.TurnContext) bool {
	wBishops, wKnights, wBishopPosition := 0, 0, core.NoPosition
	bBishops, bKnights, bBishopPosition := 0, 0, core.NoPosition

	for i, square := range ctx.Board {
		if square.IsEmpty() {
			continue
		}

		if square.IsOccupiedByAnyPiece(core.PAWN, core.QUEEN, core.ROOK) {
			return false
		}

		if square.IsOccupiedByAny(core.WHITE, core.BISHOP) {
			wBishops++
			wBishopPosition = core.Position(i)
		}
		if square.IsOccupiedByAny(core.WHITE, core.KNIGHT) {
			wKnights++
		}

		if square.IsOccupiedByAny(core.BLACK, core.BISHOP) {
			bBishops++
			bBishopPosition = core.Position(i)
		}
		if square.IsOccupiedByAny(core.BLACK, core.KNIGHT) {
			bKnights++
		}
	}

	whiteHasKingOnly := wBishops == 0 && wKnights == 0
	blackHasKingOnly := bBishops == 0 && bKnights == 0
	if whiteHasKingOnly && blackHasKingOnly {
		return true
	}

	whiteHasOneBishobBlackNone := wBishops == 1 && wKnights == 0 && bBishops == 0 && bKnights == 0
	blackHasOneBishobWhiteNone := bBishops == 1 && bKnights == 0 && wBishops == 0 && wKnights == 0
	if whiteHasOneBishobBlackNone || blackHasOneBishobWhiteNone {
		return true
	}

	whiteHasOneKnightBlackNone := wKnights == 1 && wBishops == 0 && bKnights == 0 && bBishops == 0
	blackHasOneKnightWhiteNone := bKnights == 1 && bBishops == 0 && wKnights == 0 && wBishops == 0
	if whiteHasOneKnightBlackNone || blackHasOneKnightWhiteNone {
		return true
	}

	bothHaveOneBishop := wBishops == 1 && bBishops == 1 && wKnights == 0 && bKnights == 0
	if bothHaveOneBishop {
		return wBishopPosition.SquareColor() == bBishopPosition.SquareColor()
	}

	return false
}

func (DefaultRules) IsCheckMate(ctx core.TurnContext, engine engine.Engine) bool {
	if engine.HasAnyLegalMoves(ctx) {
		return false
	}

	current := ctx.SideToMove
	enemy := current.Opponent()
	kingPosition := ctx.Sides[current].KingPosition

	return engine.IsSquareAttacked(kingPosition, enemy, ctx.BoardContext)
}

func (DefaultRules) IsStaleMate(ctx core.TurnContext, engine engine.Engine) bool {
	if engine.HasAnyLegalMoves(ctx) {
		return false
	}

	current := ctx.SideToMove
	enemy := current.Opponent()
	kingPosition := ctx.Sides[current].KingPosition

	return !engine.IsSquareAttacked(kingPosition, enemy, ctx.BoardContext)
}

func (r DefaultRules) GetGameResult(ctx core.TurnContext, engine engine.Engine, tracker tracker.Tracker, hash uint64) core.GameResult {
	if r.IsFiftyMoveRule(ctx) {
		return core.GameResult{Status: core.Draw, DrawReason: core.FiftyMoveRule}
	}

	if r.IsThreefoldRepetition(tracker, hash) {
		return core.GameResult{Status: core.Draw, DrawReason: core.ThreefoldRepetition}
	}

	if r.IsInsufficientMaterial(ctx) {
		return core.GameResult{Status: core.Draw, DrawReason: core.InsufficientMaterial}
	}

	if engine.HasAnyLegalMoves(ctx) {
		return core.GameResult{Status: core.InProgress}
	}

	current := ctx.SideToMove
	enemy := current.Opponent()
	kingPosition := ctx.Sides[current].KingPosition

	if engine.IsSquareAttacked(kingPosition, enemy, ctx.BoardContext) {
		return core.GameResult{Status: core.CheckMate, HasWinner: true, Winner: enemy}
	}

	return core.GameResult{Status: core.Draw, DrawReason: core.Stalemate}
}
