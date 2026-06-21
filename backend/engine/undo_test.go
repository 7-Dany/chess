package engine

import (
	"testing"

	"github.com/7-Dany/chess/core"
	"github.com/7-Dany/chess/testutil"
)

// TestUndo exercises Undo in isolation. The board is set up as it would be
// AFTER Apply ran, the snapshot is built manually with the pre-move state,
// and Undo is called directly without ever calling Apply. This isolates Undo
// bugs from Apply bugs.
func TestUndo(t *testing.T) {
	engine := GetDefaultEngine()
	defaultSides := testutil.DefaultSides()

	// Helper: run Undo on a post-move context and assert the result.
	// The snapshot carries the pre-move state that Undo restores.
	runUndo := func(t *testing.T, postMoveBoard *core.Board, postMoveSides [2]core.SideState, postMoveEP core.Position, snap core.Snapshot, wantBoard *core.Board, wantSides [2]core.SideState, wantEP core.Position) {
		t.Helper()
		ctx := &core.TurnContext{
			MoveContext: core.MoveContext{
				BoardContext:    core.BoardContext{Board: postMoveBoard},
				SideToMove:      core.BLACK, // Undo doesn't touch SideToMove
				Sides:           postMoveSides,
				EnPassantTarget: postMoveEP,
			},
		}

		engine.Undo(ctx, snap)

		// Board: compare every square.
		for i := range 64 {
			pos := core.Position(i)
			if ctx.Board[pos] != wantBoard[pos] {
				t.Errorf("board[%v] = %v, want %v", pos, ctx.Board[pos], wantBoard[pos])
			}
		}
		// Sides and EP restored from snapshot.
		if ctx.Sides != wantSides {
			t.Errorf("Sides = %+v, want %+v", ctx.Sides, wantSides)
		}
		if ctx.EnPassantTarget != wantEP {
			t.Errorf("EnPassantTarget = %v, want %v", ctx.EnPassantTarget, wantEP)
		}
	}

	// =========================================================================
	// Normal moves
	// =========================================================================

	t.Run("undoing a normal knight move returns it to its origin", func(t *testing.T) {
		// Post-move: knight on C3 (moved from B1).
		postMove := core.Board{}
		postMove[core.C3] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})

		snap := core.Snapshot{
			Move:                    core.Move{Type: core.NORMAL, Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE}, From: core.B1, To: core.C3},
			PreviousSides:           defaultSides,
			PreviousEnPassantTarget: core.NoPosition,
		}

		wantBoard := core.Board{}
		wantBoard[core.B1] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})

		runUndo(t, &postMove, defaultSides, core.NoPosition, snap, &wantBoard, defaultSides, core.NoPosition)
	})

	t.Run("undoing a king move restores the king position and castling rights", func(t *testing.T) {
		// Post-move: king on F1 (moved from E1), no castling rights (forfeited by king move).
		postMove := core.Board{}
		postMove[core.F1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
		postMoveSides := [2]core.SideState{{KingPosition: core.F1}, defaultSides[1]}

		snap := core.Snapshot{
			Move:                    core.Move{Type: core.NORMAL, Piece: core.Piece{Type: core.KING, Color: core.WHITE}, From: core.E1, To: core.F1},
			PreviousSides:           defaultSides,
			PreviousEnPassantTarget: core.NoPosition,
		}

		wantBoard := core.Board{}
		wantBoard[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})

		runUndo(t, &postMove, postMoveSides, core.NoPosition, snap, &wantBoard, defaultSides, core.NoPosition)
	})

	t.Run("undoing a rook move from A1 restores the queen-side castling right", func(t *testing.T) {
		// Post-move: rook on A3 (moved from A1), queen-side right lost.
		postMove := core.Board{}
		postMove[core.A3] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
		postMoveSides := [2]core.SideState{{KingPosition: core.E1, CanCastleKingSide: true}, defaultSides[1]}

		snap := core.Snapshot{
			Move:                    core.Move{Type: core.NORMAL, Piece: core.Piece{Type: core.ROOK, Color: core.WHITE}, From: core.A1, To: core.A3},
			PreviousSides:           defaultSides,
			PreviousEnPassantTarget: core.NoPosition,
		}

		wantBoard := core.Board{}
		wantBoard[core.A1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})

		runUndo(t, &postMove, postMoveSides, core.NoPosition, snap, &wantBoard, defaultSides, core.NoPosition)
	})

	t.Run("undoing a black rook move from H8 restores the king-side castling right", func(t *testing.T) {
		postMove := core.Board{}
		postMove[core.H6] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
		postMoveSides := [2]core.SideState{defaultSides[0], {KingPosition: core.E8, CanCastleQueenSide: true}}

		snap := core.Snapshot{
			Move:                    core.Move{Type: core.NORMAL, Piece: core.Piece{Type: core.ROOK, Color: core.BLACK}, From: core.H8, To: core.H6},
			PreviousSides:           defaultSides,
			PreviousEnPassantTarget: core.NoPosition,
		}

		wantBoard := core.Board{}
		wantBoard[core.H8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})

		runUndo(t, &postMove, postMoveSides, core.NoPosition, snap, &wantBoard, defaultSides, core.NoPosition)
	})

	// =========================================================================
	// Captures
	// =========================================================================

	t.Run("undoing a capture restores both the mover and the captured piece", func(t *testing.T) {
		// Post-move: white knight on D5 (captured black pawn that was there).
		postMove := core.Board{}
		postMove[core.D5] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})

		snap := core.Snapshot{
			Move: core.Move{
				Type: core.NORMAL, Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE},
				From: core.E4, To: core.D5,
				HasCapture: true, Captured: core.Piece{Type: core.PAWN, Color: core.BLACK},
			},
			PreviousSides:           defaultSides,
			PreviousEnPassantTarget: core.NoPosition,
		}

		wantBoard := core.Board{}
		wantBoard[core.E4] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
		wantBoard[core.D5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})

		runUndo(t, &postMove, defaultSides, core.NoPosition, snap, &wantBoard, defaultSides, core.NoPosition)
	})

	t.Run("undoing a capture of a rook on A8 restores the rook and the queen-side right", func(t *testing.T) {
		// Post-move: white bishop on A8 (captured black rook), black lost queen-side right.
		postMove := core.Board{}
		postMove[core.A8] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
		postMoveSides := [2]core.SideState{defaultSides[0], {KingPosition: core.E8, CanCastleKingSide: true}}

		snap := core.Snapshot{
			Move: core.Move{
				Type: core.NORMAL, Piece: core.Piece{Type: core.BISHOP, Color: core.WHITE},
				From: core.A6, To: core.A8,
				HasCapture: true, Captured: core.Piece{Type: core.ROOK, Color: core.BLACK},
			},
			PreviousSides:           defaultSides,
			PreviousEnPassantTarget: core.NoPosition,
		}

		wantBoard := core.Board{}
		wantBoard[core.A6] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
		wantBoard[core.A8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})

		runUndo(t, &postMove, postMoveSides, core.NoPosition, snap, &wantBoard, defaultSides, core.NoPosition)
	})

	t.Run("undoing a capture of a non-rook piece does not affect castling rights", func(t *testing.T) {
		// Post-move: white knight on A6 (captured black pawn).
		postMove := core.Board{}
		postMove[core.A6] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})

		snap := core.Snapshot{
			Move: core.Move{
				Type: core.NORMAL, Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE},
				From: core.B5, To: core.A6,
				HasCapture: true, Captured: core.Piece{Type: core.PAWN, Color: core.BLACK},
			},
			PreviousSides:           defaultSides,
			PreviousEnPassantTarget: core.NoPosition,
		}

		wantBoard := core.Board{}
		wantBoard[core.B5] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
		wantBoard[core.A6] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})

		runUndo(t, &postMove, defaultSides, core.NoPosition, snap, &wantBoard, defaultSides, core.NoPosition)
	})

	// =========================================================================
	// En passant
	// =========================================================================

	t.Run("undoing a white en passant capture restores both pawns", func(t *testing.T) {
		// Post-move: white pawn on E6 (captured en passant), black pawn on E5 gone.
		postMove := core.Board{}
		postMove[core.E6] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})

		snap := core.Snapshot{
			Move: core.Move{
				Type: core.EN_PASSANT, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE},
				From: core.D5, To: core.E6,
				HasCapture: true, Captured: core.Piece{Type: core.PAWN, Color: core.BLACK},
			},
			PreviousSides:           defaultSides,
			PreviousEnPassantTarget: core.E6,
		}

		wantBoard := core.Board{}
		wantBoard[core.D5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		wantBoard[core.E5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})

		runUndo(t, &postMove, defaultSides, core.NoPosition, snap, &wantBoard, defaultSides, core.E6)
	})

	t.Run("undoing a black en passant capture restores both pawns", func(t *testing.T) {
		// Post-move: black pawn on E3 (captured en passant), white pawn on E4 gone.
		postMove := core.Board{}
		postMove[core.E3] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})

		snap := core.Snapshot{
			Move: core.Move{
				Type: core.EN_PASSANT, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK},
				From: core.D4, To: core.E3,
				HasCapture: true, Captured: core.Piece{Type: core.PAWN, Color: core.WHITE},
			},
			PreviousSides:           defaultSides,
			PreviousEnPassantTarget: core.E3,
		}

		wantBoard := core.Board{}
		wantBoard[core.D4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		wantBoard[core.E4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})

		runUndo(t, &postMove, defaultSides, core.NoPosition, snap, &wantBoard, defaultSides, core.E3)
	})

	t.Run("undoing en passant on the A file restores the pawn but does not affect castling rights", func(t *testing.T) {
		// Post-move: white pawn on A6 (captured en passant), black pawn on A5 gone.
		postMove := core.Board{}
		postMove[core.A6] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})

		snap := core.Snapshot{
			Move: core.Move{
				Type: core.EN_PASSANT, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE},
				From: core.B5, To: core.A6,
				HasCapture: true, Captured: core.Piece{Type: core.PAWN, Color: core.BLACK},
			},
			PreviousSides:           defaultSides,
			PreviousEnPassantTarget: core.A6,
		}

		wantBoard := core.Board{}
		wantBoard[core.B5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		wantBoard[core.A5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})

		runUndo(t, &postMove, defaultSides, core.NoPosition, snap, &wantBoard, defaultSides, core.A6)
	})

	// =========================================================================
	// Promotion
	// =========================================================================

	t.Run("undoing a promotion to queen restores the original pawn", func(t *testing.T) {
		// Post-move: white queen on E8 (promoted from pawn on E7).
		postMove := core.Board{}
		postMove[core.E8] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.WHITE})

		snap := core.Snapshot{
			Move: core.Move{
				Type: core.PROMOTION, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE},
				From: core.E7, To: core.E8, PromoteTo: core.QUEEN,
			},
			PreviousSides:           defaultSides,
			PreviousEnPassantTarget: core.NoPosition,
		}

		wantBoard := core.Board{}
		wantBoard[core.E7] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})

		runUndo(t, &postMove, defaultSides, core.NoPosition, snap, &wantBoard, defaultSides, core.NoPosition)
	})

	t.Run("undoing a promotion to knight restores the original pawn", func(t *testing.T) {
		postMove := core.Board{}
		postMove[core.D1] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.BLACK})

		snap := core.Snapshot{
			Move: core.Move{
				Type: core.PROMOTION, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK},
				From: core.D2, To: core.D1, PromoteTo: core.KNIGHT,
			},
			PreviousSides:           defaultSides,
			PreviousEnPassantTarget: core.NoPosition,
		}

		wantBoard := core.Board{}
		wantBoard[core.D2] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})

		runUndo(t, &postMove, defaultSides, core.NoPosition, snap, &wantBoard, defaultSides, core.NoPosition)
	})

	t.Run("undoing a promotion with capture restores the pawn and the captured piece", func(t *testing.T) {
		// Post-move: white queen on D8 (promoted, captured black rook on D8).
		postMove := core.Board{}
		postMove[core.D8] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.WHITE})

		snap := core.Snapshot{
			Move: core.Move{
				Type: core.PROMOTION, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE},
				From: core.E7, To: core.D8, PromoteTo: core.QUEEN,
				HasCapture: true, Captured: core.Piece{Type: core.ROOK, Color: core.BLACK},
			},
			PreviousSides:           defaultSides,
			PreviousEnPassantTarget: core.NoPosition,
		}

		wantBoard := core.Board{}
		wantBoard[core.E7] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		wantBoard[core.D8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})

		runUndo(t, &postMove, defaultSides, core.NoPosition, snap, &wantBoard, defaultSides, core.NoPosition)
	})

	// =========================================================================
	// Castling
	// =========================================================================

	t.Run("undoing white king-side castling restores king to E1 and rook to H1", func(t *testing.T) {
		// Post-move: king on G1, rook on F1 (castled).
		postMove := core.Board{}
		postMove[core.G1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
		postMove[core.F1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
		postMoveSides := [2]core.SideState{{KingPosition: core.G1}, defaultSides[1]}

		snap := core.Snapshot{
			Move: core.Move{
				Type: core.CASTLING, Piece: core.Piece{Type: core.KING, Color: core.WHITE},
				From: core.E1, To: core.G1,
			},
			PreviousSides:           defaultSides,
			PreviousEnPassantTarget: core.NoPosition,
		}

		wantBoard := core.Board{}
		wantBoard[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
		wantBoard[core.H1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})

		runUndo(t, &postMove, postMoveSides, core.NoPosition, snap, &wantBoard, defaultSides, core.NoPosition)
	})

	t.Run("undoing white queen-side castling restores king to E1 and rook to A1", func(t *testing.T) {
		postMove := core.Board{}
		postMove[core.C1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
		postMove[core.D1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
		postMoveSides := [2]core.SideState{{KingPosition: core.C1}, defaultSides[1]}

		snap := core.Snapshot{
			Move: core.Move{
				Type: core.CASTLING, Piece: core.Piece{Type: core.KING, Color: core.WHITE},
				From: core.E1, To: core.C1,
			},
			PreviousSides:           defaultSides,
			PreviousEnPassantTarget: core.NoPosition,
		}

		wantBoard := core.Board{}
		wantBoard[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
		wantBoard[core.A1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})

		runUndo(t, &postMove, postMoveSides, core.NoPosition, snap, &wantBoard, defaultSides, core.NoPosition)
	})

	t.Run("undoing black king-side castling restores king to E8 and rook to H8", func(t *testing.T) {
		postMove := core.Board{}
		postMove[core.G8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
		postMove[core.F8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
		postMoveSides := [2]core.SideState{defaultSides[0], {KingPosition: core.G8}}

		snap := core.Snapshot{
			Move: core.Move{
				Type: core.CASTLING, Piece: core.Piece{Type: core.KING, Color: core.BLACK},
				From: core.E8, To: core.G8,
			},
			PreviousSides:           defaultSides,
			PreviousEnPassantTarget: core.NoPosition,
		}

		wantBoard := core.Board{}
		wantBoard[core.E8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
		wantBoard[core.H8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})

		runUndo(t, &postMove, postMoveSides, core.NoPosition, snap, &wantBoard, defaultSides, core.NoPosition)
	})

	t.Run("undoing black queen-side castling restores king to E8 and rook to A8", func(t *testing.T) {
		postMove := core.Board{}
		postMove[core.C8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
		postMove[core.D8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
		postMoveSides := [2]core.SideState{defaultSides[0], {KingPosition: core.C8}}

		snap := core.Snapshot{
			Move: core.Move{
				Type: core.CASTLING, Piece: core.Piece{Type: core.KING, Color: core.BLACK},
				From: core.E8, To: core.C8,
			},
			PreviousSides:           defaultSides,
			PreviousEnPassantTarget: core.NoPosition,
		}

		wantBoard := core.Board{}
		wantBoard[core.E8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
		wantBoard[core.A8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})

		runUndo(t, &postMove, postMoveSides, core.NoPosition, snap, &wantBoard, defaultSides, core.NoPosition)
	})

	// =========================================================================
	// En passant target restoration
	// =========================================================================

	t.Run("undoing a double pawn push restores the no-EP state (EP was set by the push)", func(t *testing.T) {
		// Post-move: white pawn on E4 (double-pushed from E2), EP target E3 set.
		postMove := core.Board{}
		postMove[core.E4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})

		snap := core.Snapshot{
			Move: core.Move{
				Type: core.NORMAL, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE},
				From: core.E2, To: core.E4,
			},
			PreviousSides:           defaultSides,
			PreviousEnPassantTarget: core.NoPosition, // before the double push, there was no EP target
		}

		wantBoard := core.Board{}
		wantBoard[core.E2] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})

		runUndo(t, &postMove, defaultSides, core.E3, snap, &wantBoard, defaultSides, core.NoPosition)
	})

	t.Run("undoing a non-pawn move restores the previous EP target", func(t *testing.T) {
		// Post-move: knight on C3 (moved from B1), EP was cleared by the move.
		// But the snapshot carries the PREVIOUS EP target (E3), which Undo restores.
		postMove := core.Board{}
		postMove[core.C3] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})

		snap := core.Snapshot{
			Move: core.Move{
				Type: core.NORMAL, Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE},
				From: core.B1, To: core.C3,
			},
			PreviousSides:           defaultSides,
			PreviousEnPassantTarget: core.E3, // there was an EP target before this move
		}

		wantBoard := core.Board{}
		wantBoard[core.B1] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})

		runUndo(t, &postMove, defaultSides, core.NoPosition, snap, &wantBoard, defaultSides, core.E3)
	})
}

// TestApplyThenUndo verifies that Apply followed by Undo returns the context
// to its original state. This catches asymmetry bugs between Apply and Undo
// that the isolated TestUndo would miss (e.g. if Apply writes X but Undo
// reads Y from the snapshot).
//
// Each case sets up a board + move, snapshots the full original state,
// applies the move, undoes it, and asserts every field matches the original.
func TestApplyThenUndo(t *testing.T) {
	engine := GetDefaultEngine()
	defaultSides := testutil.DefaultSides()

	// Helper: snapshot the original context, apply+undo, assert nothing changed.
	roundTrip := func(t *testing.T, board *core.Board, side core.PieceColor, sides [2]core.SideState, ep core.Position, move core.Move) {
		t.Helper()

		// Build the context.
		ctx := &core.TurnContext{
			MoveContext: core.MoveContext{
				BoardContext:    core.BoardContext{Board: board},
				SideToMove:      side,
				Sides:           sides,
				EnPassantTarget: ep,
			},
		}

		// Snapshot the original state for comparison.
		originalBoard := *board
		originalSideToMove := ctx.SideToMove
		originalSides := ctx.Sides
		originalEP := ctx.EnPassantTarget
		originalHalfMoveClock := ctx.HalfMoveClock
		originalFullMoveNumber := ctx.FullMoveNumber

		// Apply then Undo.
		snap := engine.Apply(ctx, move)
		engine.Undo(ctx, snap)

		// Compare every board square.
		for i := range 64 {
			pos := core.Position(i)
			if ctx.Board[pos] != originalBoard[pos] {
				t.Errorf("board[%v] = %v, want %v (after round trip)", pos, ctx.Board[pos], originalBoard[pos])
			}
		}
		// Compare every context field.
		if ctx.SideToMove != originalSideToMove {
			t.Errorf("SideToMove = %v, want %v", ctx.SideToMove, originalSideToMove)
		}
		if ctx.Sides != originalSides {
			t.Errorf("Sides = %+v, want %+v", ctx.Sides, originalSides)
		}
		if ctx.EnPassantTarget != originalEP {
			t.Errorf("EnPassantTarget = %v, want %v", ctx.EnPassantTarget, originalEP)
		}
		if ctx.HalfMoveClock != originalHalfMoveClock {
			t.Errorf("HalfMoveClock = %v, want %v", ctx.HalfMoveClock, originalHalfMoveClock)
		}
		if ctx.FullMoveNumber != originalFullMoveNumber {
			t.Errorf("FullMoveNumber = %v, want %v", ctx.FullMoveNumber, originalFullMoveNumber)
		}
	}

	// placement pairs a position with a piece for building test boards.
	type placement struct {
		pos   core.Position
		piece core.Piece
	}

	// Helper: build a board with the given pieces.
	board := func(pieces ...placement) *core.Board {
		var b core.Board
		for _, p := range pieces {
			b[p.pos] = core.NewSquare(p.piece)
		}
		return &b
	}

	// Convenience piece constructors.
	wk := core.Piece{Type: core.KING, Color: core.WHITE}
	bk := core.Piece{Type: core.KING, Color: core.BLACK}
	wp := core.Piece{Type: core.PAWN, Color: core.WHITE}
	bp := core.Piece{Type: core.PAWN, Color: core.BLACK}
	wn := core.Piece{Type: core.KNIGHT, Color: core.WHITE}
	wr := core.Piece{Type: core.ROOK, Color: core.WHITE}
	br := core.Piece{Type: core.ROOK, Color: core.BLACK}
	wb := core.Piece{Type: core.BISHOP, Color: core.WHITE}
	bb := core.Piece{Type: core.BISHOP, Color: core.BLACK}

	// =========================================================================
	// Normal moves
	// =========================================================================

	t.Run("normal knight move round trip", func(t *testing.T) {
		b := board(
			placement{core.B1, wn},
			placement{core.E1, wk},
			placement{core.E8, bk},
		)
		roundTrip(t, b, core.WHITE, defaultSides, core.NoPosition,
			core.Move{Type: core.NORMAL, Piece: wn, From: core.B1, To: core.C3})
	})

	t.Run("normal king move round trip", func(t *testing.T) {
		b := board(
			placement{core.E1, wk},
			placement{core.E8, bk},
		)
		roundTrip(t, b, core.WHITE, defaultSides, core.NoPosition,
			core.Move{Type: core.NORMAL, Piece: wk, From: core.E1, To: core.F1})
	})

	t.Run("rook move from A1 round trip (forfeits queen-side right)", func(t *testing.T) {
		b := board(
			placement{core.A1, wr},
			placement{core.E1, wk},
			placement{core.E8, bk},
		)
		roundTrip(t, b, core.WHITE, defaultSides, core.NoPosition,
			core.Move{Type: core.NORMAL, Piece: wr, From: core.A1, To: core.A3})
	})

	t.Run("rook move from a non-home file round trip (rights preserved)", func(t *testing.T) {
		b := board(
			placement{core.C3, wr},
			placement{core.E1, wk},
			placement{core.E8, bk},
		)
		roundTrip(t, b, core.WHITE, defaultSides, core.NoPosition,
			core.Move{Type: core.NORMAL, Piece: wr, From: core.C3, To: core.C5})
	})

	t.Run("black rook move from H8 round trip (forfeits king-side right)", func(t *testing.T) {
		b := board(
			placement{core.H8, br},
			placement{core.E1, wk},
			placement{core.E8, bk},
		)
		roundTrip(t, b, core.BLACK, defaultSides, core.NoPosition,
			core.Move{Type: core.NORMAL, Piece: br, From: core.H8, To: core.H6})
	})

	t.Run("black rook move from A8 round trip (forfeits queen-side right)", func(t *testing.T) {
		b := board(
			placement{core.A8, br},
			placement{core.E1, wk},
			placement{core.E8, bk},
		)
		roundTrip(t, b, core.BLACK, defaultSides, core.NoPosition,
			core.Move{Type: core.NORMAL, Piece: br, From: core.A8, To: core.A6})
	})

	// =========================================================================
	// Captures
	// =========================================================================

	t.Run("capture round trip", func(t *testing.T) {
		b := board(
			placement{core.E4, wn},
			placement{core.D5, bp},
			placement{core.E1, wk},
			placement{core.E8, bk},
		)
		roundTrip(t, b, core.WHITE, defaultSides, core.NoPosition,
			core.Move{Type: core.NORMAL, Piece: wn, From: core.E4, To: core.D5, HasCapture: true, Captured: bp})
	})

	t.Run("capture of a rook on A8 round trip (affects queen-side right)", func(t *testing.T) {
		b := board(
			placement{core.A6, wb},
			placement{core.A8, br},
			placement{core.E1, wk},
			placement{core.E8, bk},
		)
		roundTrip(t, b, core.WHITE, defaultSides, core.NoPosition,
			core.Move{Type: core.NORMAL, Piece: wb, From: core.A6, To: core.A8, HasCapture: true, Captured: br})
	})

	t.Run("capture of a rook on H1 round trip (affects king-side right)", func(t *testing.T) {
		b := board(
			placement{core.H3, bb},
			placement{core.H1, wr},
			placement{core.E1, wk},
			placement{core.E8, bk},
		)
		roundTrip(t, b, core.BLACK, defaultSides, core.NoPosition,
			core.Move{Type: core.NORMAL, Piece: bb, From: core.H3, To: core.H1, HasCapture: true, Captured: wr})
	})

	t.Run("capture of a non-rook on the A file round trip (rights preserved)", func(t *testing.T) {
		b := board(
			placement{core.B5, wn},
			placement{core.A6, bp},
			placement{core.E1, wk},
			placement{core.E8, bk},
		)
		roundTrip(t, b, core.WHITE, defaultSides, core.NoPosition,
			core.Move{Type: core.NORMAL, Piece: wn, From: core.B5, To: core.A6, HasCapture: true, Captured: bp})
	})

	// =========================================================================
	// En passant
	// =========================================================================

	t.Run("white en passant round trip", func(t *testing.T) {
		b := board(
			placement{core.D5, wp},
			placement{core.E5, bp},
			placement{core.E1, wk},
			placement{core.E8, bk},
		)
		roundTrip(t, b, core.WHITE, defaultSides, core.E6,
			core.Move{Type: core.EN_PASSANT, Piece: wp, From: core.D5, To: core.E6, HasCapture: true, Captured: bp})
	})

	t.Run("black en passant round trip", func(t *testing.T) {
		b := board(
			placement{core.E4, wp},
			placement{core.D4, bp},
			placement{core.E1, wk},
			placement{core.E8, bk},
		)
		roundTrip(t, b, core.BLACK, defaultSides, core.E3,
			core.Move{Type: core.EN_PASSANT, Piece: bp, From: core.D4, To: core.E3, HasCapture: true, Captured: wp})
	})

	t.Run("en passant on the A file round trip", func(t *testing.T) {
		b := board(
			placement{core.B5, wp},
			placement{core.A5, bp},
			placement{core.E1, wk},
			placement{core.E8, bk},
		)
		roundTrip(t, b, core.WHITE, defaultSides, core.A6,
			core.Move{Type: core.EN_PASSANT, Piece: wp, From: core.B5, To: core.A6, HasCapture: true, Captured: bp})
	})

	// =========================================================================
	// Promotion
	// =========================================================================

	t.Run("promotion to queen round trip", func(t *testing.T) {
		b := board(
			placement{core.E7, wp},
			placement{core.E1, wk},
			placement{core.H8, bk}, // black king off the promotion square
		)
		roundTrip(t, b, core.WHITE, [2]core.SideState{defaultSides[0], testutil.Side(core.H8, true, true)}, core.NoPosition,
			core.Move{Type: core.PROMOTION, Piece: wp, From: core.E7, To: core.E8, PromoteTo: core.QUEEN})
	})

	t.Run("promotion to knight round trip", func(t *testing.T) {
		b := board(
			placement{core.D2, bp},
			placement{core.A1, wk}, // white king off the promotion square
			placement{core.E8, bk},
		)
		roundTrip(t, b, core.BLACK, [2]core.SideState{testutil.Side(core.A1, true, true), defaultSides[1]}, core.NoPosition,
			core.Move{Type: core.PROMOTION, Piece: bp, From: core.D2, To: core.D1, PromoteTo: core.KNIGHT})
	})

	t.Run("promotion with capture round trip", func(t *testing.T) {
		b := board(
			placement{core.E7, wp},
			placement{core.D8, br},
			placement{core.E1, wk},
			placement{core.H8, bk}, // black king off D8
		)
		roundTrip(t, b, core.WHITE, [2]core.SideState{defaultSides[0], testutil.Side(core.H8, true, true)}, core.NoPosition,
			core.Move{Type: core.PROMOTION, Piece: wp, From: core.E7, To: core.D8, PromoteTo: core.QUEEN, HasCapture: true, Captured: br})
	})

	// =========================================================================
	// Castling
	// =========================================================================

	t.Run("white king-side castling round trip", func(t *testing.T) {
		b := board(
			placement{core.E1, wk},
			placement{core.H1, wr},
			placement{core.E8, bk},
		)
		roundTrip(t, b, core.WHITE, defaultSides, core.NoPosition,
			core.Move{Type: core.CASTLING, Piece: wk, From: core.E1, To: core.G1})
	})

	t.Run("white queen-side castling round trip", func(t *testing.T) {
		b := board(
			placement{core.E1, wk},
			placement{core.A1, wr},
			placement{core.E8, bk},
		)
		roundTrip(t, b, core.WHITE, defaultSides, core.NoPosition,
			core.Move{Type: core.CASTLING, Piece: wk, From: core.E1, To: core.C1})
	})

	t.Run("black king-side castling round trip", func(t *testing.T) {
		b := board(
			placement{core.E8, bk},
			placement{core.H8, br},
			placement{core.E1, wk},
		)
		roundTrip(t, b, core.BLACK, defaultSides, core.NoPosition,
			core.Move{Type: core.CASTLING, Piece: bk, From: core.E8, To: core.G8})
	})

	t.Run("black queen-side castling round trip", func(t *testing.T) {
		b := board(
			placement{core.E8, bk},
			placement{core.A8, br},
			placement{core.E1, wk},
		)
		roundTrip(t, b, core.BLACK, defaultSides, core.NoPosition,
			core.Move{Type: core.CASTLING, Piece: bk, From: core.E8, To: core.C8})
	})

	// =========================================================================
	// Double pawn push (sets en passant target)
	// =========================================================================

	t.Run("white double pawn push round trip (sets EP target)", func(t *testing.T) {
		b := board(
			placement{core.E2, wp},
			placement{core.E1, wk},
			placement{core.E8, bk},
		)
		roundTrip(t, b, core.WHITE, defaultSides, core.NoPosition,
			core.Move{Type: core.NORMAL, Piece: wp, From: core.E2, To: core.E4})
	})

	t.Run("black double pawn push round trip (sets EP target)", func(t *testing.T) {
		b := board(
			placement{core.D7, bp},
			placement{core.E1, wk},
			placement{core.E8, bk},
		)
		roundTrip(t, b, core.BLACK, defaultSides, core.NoPosition,
			core.Move{Type: core.NORMAL, Piece: bp, From: core.D7, To: core.D5})
	})

	t.Run("white double pawn push from the A file round trip", func(t *testing.T) {
		b := board(
			placement{core.A2, wp},
			placement{core.E1, wk},
			placement{core.E8, bk},
		)
		roundTrip(t, b, core.WHITE, defaultSides, core.NoPosition,
			core.Move{Type: core.NORMAL, Piece: wp, From: core.A2, To: core.A4})
	})

	// =========================================================================
	// Non-pawn move with an existing EP target (clears it)
	// =========================================================================

	t.Run("non-pawn move with an existing EP target round trip (clears EP)", func(t *testing.T) {
		b := board(
			placement{core.B1, wn},
			placement{core.E1, wk},
			placement{core.E8, bk},
		)
		roundTrip(t, b, core.WHITE, defaultSides, core.E3,
			core.Move{Type: core.NORMAL, Piece: wn, From: core.B1, To: core.C3})
	})
}
