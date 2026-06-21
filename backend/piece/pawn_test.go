package piece

import (
	"testing"

	"github.com/7-Dany/chess/core"
	"github.com/7-Dany/chess/testutil"
)

// TestPawnIsAttacking verifies that Pawn.IsAttacking correctly reports
// whether a pawn of the given color attacks the target square.
//
// Pawns are the only piece whose attack direction depends on color: a white
// pawn attacks the two squares diagonally AHEAD of it (up-left, up-right);
// a black pawn attacks the two squares diagonally ahead of it (down-left,
// down-right). Pawns do NOT attack the square directly ahead — only
// diagonally.
//
// The scan goes FROM the target outward: "is there a pawn of `color` on one
// of the two squares behind me (relative to the attacker's movement
// direction)?".
func TestPawnIsAttacking(t *testing.T) {
	pawn := Pawn{}

	// White pawns attack upward. Target E4 is attacked by a white pawn on D3
	// (down-left) or F3 (down-right).
	t.Run("a white pawn attacks E4 from down-left (D3)", func(t *testing.T) {
		var board core.Board
		board[core.D3] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		if !pawn.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("white pawn on D3 should attack E4")
		}
	})

	t.Run("a white pawn attacks E4 from down-right (F3)", func(t *testing.T) {
		var board core.Board
		board[core.F3] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		if !pawn.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("white pawn on F3 should attack E4")
		}
	})

	// Black pawns attack downward. Target E4 is attacked by a black pawn on
	// D5 (up-left from black's perspective) or F5 (up-right).
	t.Run("a black pawn attacks E4 from up-left (D5)", func(t *testing.T) {
		var board core.Board
		board[core.D5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		ctx := core.BoardContext{Board: &board}

		if !pawn.IsAttacking(core.BLACK, core.E4, ctx) {
			t.Errorf("black pawn on D5 should attack E4")
		}
	})

	t.Run("a black pawn attacks E4 from up-right (F5)", func(t *testing.T) {
		var board core.Board
		board[core.F5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		ctx := core.BoardContext{Board: &board}

		if !pawn.IsAttacking(core.BLACK, core.E4, ctx) {
			t.Errorf("black pawn on F5 should attack E4")
		}
	})

	// A pawn on the same file (directly ahead or behind) does NOT attack —
	// pawns only attack diagonally.
	t.Run("a pawn on the same file does not attack (pawns attack diagonally only)", func(t *testing.T) {
		// White pawn on E3 (directly behind E4 from white's perspective).
		var board core.Board
		board[core.E3] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		if pawn.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("white pawn on E3 should NOT attack E4 (same file, not diagonal)")
		}
	})

	t.Run("a pawn of the wrong color is ignored", func(t *testing.T) {
		// Black pawn on D5 attacks E4 (for black). Asking "does WHITE attack
		// E4?" → no (the pawn is black). Asking "does BLACK attack E4?" → yes.
		var board core.Board
		board[core.D5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		ctx := core.BoardContext{Board: &board}

		if pawn.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("black pawn should not count as a white attacker")
		}
		if !pawn.IsAttacking(core.BLACK, core.E4, ctx) {
			t.Errorf("black pawn on D5 should attack E4 for black")
		}
	})

	t.Run("a non-pawn piece on the attack diagonal does not trigger a pawn attack", func(t *testing.T) {
		nonPawns := []core.PieceType{core.QUEEN, core.ROOK, core.BISHOP, core.KNIGHT, core.KING}
		for _, pt := range nonPawns {
			var board core.Board
			board[core.D3] = core.NewSquare(core.Piece{Type: pt, Color: core.WHITE})
			ctx := core.BoardContext{Board: &board}

			if pawn.IsAttacking(core.WHITE, core.E4, ctx) {
				t.Errorf("%v on D3 should not trigger a pawn attack", pt)
			}
		}
	})

	t.Run("an empty board reports no attack", func(t *testing.T) {
		var board core.Board
		ctx := core.BoardContext{Board: &board}

		if pawn.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("empty board should not report an attack")
		}
	})

	t.Run("a pawn sitting on the target square itself does not attack it", func(t *testing.T) {
		var board core.Board
		board[core.E4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		if pawn.IsAttacking(core.WHITE, core.E4, ctx) {
			t.Errorf("pawn on the target square should not attack itself")
		}
	})

	// A pawn on the A file can only attack to the right (no left diagonal).
	t.Run("a pawn on the A file attacks only to the right", func(t *testing.T) {
		// White pawn on B3 attacks A4 (down-left from B3's perspective... wait).
		// Target A4: a white pawn must be on the rank below (rank 3) and one
		// file to the right (B). So B3 attacks A4.
		var board core.Board
		board[core.B3] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		if !pawn.IsAttacking(core.WHITE, core.A4, ctx) {
			t.Errorf("white pawn on B3 should attack A4 (the only diagonal from A file)")
		}
	})

	// A pawn on the H file can only attack to the left (no right diagonal).
	t.Run("a pawn on the H file attacks only to the left", func(t *testing.T) {
		// Target H4: a white pawn must be on rank 3, one file left (G). So G3.
		var board core.Board
		board[core.G3] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		if !pawn.IsAttacking(core.WHITE, core.H4, ctx) {
			t.Errorf("white pawn on G3 should attack H4 (the only diagonal from H file)")
		}
	})

	// A target on rank 1 cannot be attacked by a white pawn (white pawns
	// would have to be on rank 0, which doesn't exist — pawns promote on
	// arrival, so no white pawn sits on rank 1).
	t.Run("a target on rank 1 is not attacked by any white pawn (no rank below)", func(t *testing.T) {
		var board core.Board
		// Put white pawns on rank 2 (the closest possible). They attack rank 3, not rank 1.
		board[core.D2] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		ctx := core.BoardContext{Board: &board}

		if pawn.IsAttacking(core.WHITE, core.D1, ctx) {
			t.Errorf("no white pawn can attack rank 1 (would need to be on rank 0)")
		}
	})
}

// TestPawnAttacks verifies that Pawn.Attacks returns every square a pawn
// threatens from the given position.
//
// Pawns attack the two squares diagonally ahead (color-dependent). A pawn on
// the A or H file has only one attack square. Like King.Attacks, this method
// ignores occupancy — a pawn "attacks" a square even if a friendly piece
// sits there.
func TestPawnAttacks(t *testing.T) {
	pawn := Pawn{}

	// Helper: place a pawn of the given color on the board and return the context.
	pawnCtx := func(pos core.Position, color core.PieceColor) (core.Board, core.BoardContext) {
		var board core.Board
		board[pos] = core.NewSquare(core.Piece{Type: core.PAWN, Color: color})
		return board, core.BoardContext{Board: &board}
	}

	t.Run("a white pawn on E4 threatens D5 and F5 (the two squares diagonally ahead)", func(t *testing.T) {
		_, ctx := pawnCtx(core.E4, core.WHITE)

		got := pawn.Attacks(make([]core.Position, 0, core.MAX_MOVES), core.E4, ctx)

		testutil.AssertPositionsMatch(t, got, []core.Position{core.D5, core.F5})
	})

	t.Run("a black pawn on E4 threatens D3 and F3 (the two squares diagonally ahead)", func(t *testing.T) {
		_, ctx := pawnCtx(core.E4, core.BLACK)

		got := pawn.Attacks(make([]core.Position, 0, core.MAX_MOVES), core.E4, ctx)

		testutil.AssertPositionsMatch(t, got, []core.Position{core.D3, core.F3})
	})

	// A pawn on the A file has no left diagonal — only one attack square.
	t.Run("a white pawn on A4 threatens only B5 (no left diagonal)", func(t *testing.T) {
		_, ctx := pawnCtx(core.A4, core.WHITE)

		got := pawn.Attacks(make([]core.Position, 0, core.MAX_MOVES), core.A4, ctx)

		testutil.AssertPositionsMatch(t, got, []core.Position{core.B5})
	})

	t.Run("a white pawn on H4 threatens only G5 (no right diagonal)", func(t *testing.T) {
		_, ctx := pawnCtx(core.H4, core.WHITE)

		got := pawn.Attacks(make([]core.Position, 0, core.MAX_MOVES), core.H4, ctx)

		testutil.AssertPositionsMatch(t, got, []core.Position{core.G5})
	})

	t.Run("a black pawn on A4 threatens only B3 (no left diagonal)", func(t *testing.T) {
		_, ctx := pawnCtx(core.A4, core.BLACK)

		got := pawn.Attacks(make([]core.Position, 0, core.MAX_MOVES), core.A4, ctx)

		testutil.AssertPositionsMatch(t, got, []core.Position{core.B3})
	})

	t.Run("a black pawn on H4 threatens only G3 (no right diagonal)", func(t *testing.T) {
		_, ctx := pawnCtx(core.H4, core.BLACK)

		got := pawn.Attacks(make([]core.Position, 0, core.MAX_MOVES), core.H4, ctx)

		testutil.AssertPositionsMatch(t, got, []core.Position{core.G3})
	})

	// Attacks ignores occupancy — a pawn attacks a square even if a friendly
	// piece sits there (matters for check detection).
	t.Run("a pawn attacks adjacent squares even when occupied by friendly pieces", func(t *testing.T) {
		var board core.Board
		board[core.E4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.D5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE}) // friendly on left diagonal
		board[core.F5] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE}) // friendly on right diagonal
		ctx := core.BoardContext{Board: &board}

		got := pawn.Attacks(make([]core.Position, 0, core.MAX_MOVES), core.E4, ctx)

		// Both diagonals returned regardless of occupancy.
		testutil.AssertPositionsMatch(t, got, []core.Position{core.D5, core.F5})
	})
}

// TestPawnPseudoLegalMoves verifies that Pawn.PseudoLegalMoves returns the
// correct set of moves — respecting blockers, captures, en passant, and
// promotion, but NOT filtering for king safety (that's the engine's job).
//
// Pawn moves are the most complex of any piece:
//   - Single push: one square forward (if empty).
//   - Double push: two squares forward, only from the start rank, and only
//     if BOTH squares are empty.
//   - Diagonal capture: one square diagonally forward, if an enemy piece is
//     there.
//   - En passant: a diagonal capture onto an empty square that is the
//     en-passant target.
//   - Promotion: when the pawn reaches the last rank, the push/capture
//     becomes a PROMOTION move (4 variants: Q, R, B, N).
func TestPawnPseudoLegalMoves(t *testing.T) {
	pawn := Pawn{}

	// Helper: extract just the destination squares from a move list.
	destinations := func(moves []core.Move) []core.Position {
		tos := make([]core.Position, len(moves))
		for i, m := range moves {
			tos[i] = m.To
		}
		return tos
	}

	// =========================================================================
	// Single and double push
	// =========================================================================

	t.Run("a white pawn on its start rank with both squares empty can single or double push", func(t *testing.T) {
		var board core.Board
		board[core.E2] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := pawn.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.E2, ctx)

		testutil.AssertPositionsMatch(t, destinations(moves), []core.Position{core.E3, core.E4})
		testutil.AssertMoveCount(t, moves, 2)
	})

	t.Run("a black pawn on its start rank with both squares empty can single or double push", func(t *testing.T) {
		var board core.Board
		board[core.E7] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.BLACK}

		moves := pawn.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.E7, ctx)

		testutil.AssertPositionsMatch(t, destinations(moves), []core.Position{core.E6, core.E5})
		testutil.AssertMoveCount(t, moves, 2)
	})

	t.Run("a pawn not on its start rank can only single push", func(t *testing.T) {
		// White pawn on E3 (not rank 2) → only E4, no double push.
		var board core.Board
		board[core.E3] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := pawn.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.E3, ctx)

		testutil.AssertPositionsMatch(t, destinations(moves), []core.Position{core.E4})
		testutil.AssertMoveCount(t, moves, 1)
	})

	t.Run("a pawn whose front square is occupied cannot push at all (even from start rank)", func(t *testing.T) {
		// White pawn on E2, enemy pawn on E3 blocks both single and double push.
		var board core.Board
		board[core.E2] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.E3] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := pawn.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.E2, ctx)

		// No push moves (E3 blocked). No captures either (E3 is same file, not diagonal).
		testutil.AssertNoMoves(t, moves)
	})

	t.Run("a pawn on its start rank whose double-push square is occupied can only single push", func(t *testing.T) {
		// White pawn on E2, E3 empty, E4 occupied → single push to E3 only.
		var board core.Board
		board[core.E2] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.E4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := pawn.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.E2, ctx)

		testutil.AssertPositionsMatch(t, destinations(moves), []core.Position{core.E3})
		testutil.AssertMoveCount(t, moves, 1)
	})

	// =========================================================================
	// Diagonal captures
	// =========================================================================

	t.Run("a white pawn captures an enemy on its right diagonal", func(t *testing.T) {
		// White pawn on E4, enemy on F5 → capture F5 + push E5.
		var board core.Board
		board[core.E4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.F5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := pawn.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.E4, ctx)

		testutil.AssertPositionsMatch(t, destinations(moves), []core.Position{core.E5, core.F5})

		// The capture move must be flagged correctly.
		for _, m := range moves {
			if m.To == core.F5 {
				if !m.HasCapture || m.Captured != (core.Piece{Type: core.PAWN, Color: core.BLACK}) {
					t.Errorf("capture of F5: HasCapture=%v Captured=%v, want black pawn", m.HasCapture, m.Captured)
				}
				return
			}
		}
		t.Errorf("expected a capture move to F5")
	})

	t.Run("a white pawn captures an enemy on its left diagonal", func(t *testing.T) {
		var board core.Board
		board[core.E4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.D5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := pawn.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.E4, ctx)

		testutil.AssertPositionsMatch(t, destinations(moves), []core.Position{core.E5, core.D5})
	})

	t.Run("a white pawn captures enemies on both diagonals", func(t *testing.T) {
		var board core.Board
		board[core.E4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.D5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		board[core.F5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := pawn.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.E4, ctx)

		// Push E5 + captures D5, F5.
		testutil.AssertPositionsMatch(t, destinations(moves), []core.Position{core.E5, core.D5, core.F5})
		testutil.AssertMoveCount(t, moves, 3)
	})

	t.Run("a friendly piece on the diagonal is not capturable", func(t *testing.T) {
		// White pawn on E4, own pawn on F5 → no capture, just push E5.
		var board core.Board
		board[core.E4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.F5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := pawn.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.E4, ctx)

		testutil.AssertPositionsMatch(t, destinations(moves), []core.Position{core.E5})
	})

	t.Run("captures carry the exact enemy piece type (not just pawns)", func(t *testing.T) {
		// White pawn on E4, enemy queen on D5, enemy rook on F5.
		var board core.Board
		board[core.E4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.D5] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.BLACK})
		board[core.F5] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := pawn.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.E4, ctx)

		wantCaptures := map[core.Position]core.Piece{
			core.D5: {Type: core.QUEEN, Color: core.BLACK},
			core.F5: {Type: core.ROOK, Color: core.BLACK},
		}
		for _, m := range moves {
			if want, ok := wantCaptures[m.To]; ok {
				if !m.HasCapture || m.Captured != want {
					t.Errorf("capture of %v: HasCapture=%v Captured=%v, want %v", m.To, m.HasCapture, m.Captured, want)
				}
			}
		}
	})

	// =========================================================================
	// En passant
	// =========================================================================

	t.Run("a white pawn captures en passant on its right diagonal", func(t *testing.T) {
		// White pawn on E5, enemy just double-pushed to F5 (EP target F6).
		// The capture lands on F6 (empty), removing the pawn on F5.
		var board core.Board
		board[core.E5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.F5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		ctx := core.MoveContext{
			BoardContext:    core.BoardContext{Board: &board},
			SideToMove:      core.WHITE,
			EnPassantTarget: core.F6,
		}

		moves := pawn.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.E5, ctx)

		// Push E6 + en passant capture F6.
		testutil.AssertPositionsMatch(t, destinations(moves), []core.Position{core.E6, core.F6})

		// The en passant move must be type EN_PASSANT, flagged as a capture
		// of a pawn.
		for _, m := range moves {
			if m.To == core.F6 {
				if m.Type != core.EN_PASSANT {
					t.Errorf("en passant move: Type = %v, want EN_PASSANT", m.Type)
				}
				if !m.HasCapture || m.Captured != (core.Piece{Type: core.PAWN, Color: core.BLACK}) {
					t.Errorf("en passant: HasCapture=%v Captured=%v, want black pawn", m.HasCapture, m.Captured)
				}
				return
			}
		}
		t.Errorf("expected an en passant move to F6")
	})

	t.Run("a white pawn captures en passant on its left diagonal", func(t *testing.T) {
		var board core.Board
		board[core.E5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.D5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		ctx := core.MoveContext{
			BoardContext:    core.BoardContext{Board: &board},
			SideToMove:      core.WHITE,
			EnPassantTarget: core.D6,
		}

		moves := pawn.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.E5, ctx)

		testutil.AssertPositionsMatch(t, destinations(moves), []core.Position{core.E6, core.D6})
	})

	t.Run("a black pawn captures en passant", func(t *testing.T) {
		// Black pawn on E4, enemy just double-pushed to D4 (EP target D3).
		var board core.Board
		board[core.E4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		board[core.D4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		ctx := core.MoveContext{
			BoardContext:    core.BoardContext{Board: &board},
			SideToMove:      core.BLACK,
			EnPassantTarget: core.D3,
		}

		moves := pawn.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.E4, ctx)

		// Push E3 + en passant capture D3.
		testutil.AssertPositionsMatch(t, destinations(moves), []core.Position{core.E3, core.D3})

		for _, m := range moves {
			if m.To == core.D3 {
				if m.Type != core.EN_PASSANT {
					t.Errorf("en passant move: Type = %v, want EN_PASSANT", m.Type)
				}
				if m.Captured != (core.Piece{Type: core.PAWN, Color: core.WHITE}) {
					t.Errorf("en passant: Captured = %v, want white pawn", m.Captured)
				}
				return
			}
		}
		t.Errorf("expected an en passant move to D3")
	})

	t.Run("an en passant target not on a diagonal is ignored", func(t *testing.T) {
		// EP target is on E6 (directly ahead), not on a diagonal — no en passant.
		var board core.Board
		board[core.E5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		ctx := core.MoveContext{
			BoardContext:    core.BoardContext{Board: &board},
			SideToMove:      core.WHITE,
			EnPassantTarget: core.E6, // directly ahead, not diagonal
		}

		moves := pawn.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.E5, ctx)

		// Only the push to E6 (which happens to be the EP target, but it's
		// a push not a capture).
		testutil.AssertPositionsMatch(t, destinations(moves), []core.Position{core.E6})
		// No en passant move.
		for _, m := range moves {
			if m.Type == core.EN_PASSANT {
				t.Errorf("should not generate an en passant move when target is not on a diagonal")
			}
		}
	})

	// =========================================================================
	// Promotion
	// =========================================================================

	t.Run("a white pawn reaching the last rank by forward push promotes to Q, R, B, or N", func(t *testing.T) {
		// White pawn on E7, E8 empty → 4 promotion moves.
		var board core.Board
		board[core.E7] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := pawn.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.E7, ctx)

		// 4 moves, all to E8, all type PROMOTION, with PromoteTo = Q/R/B/N.
		testutil.AssertMoveCount(t, moves, 4)
		promoteTypes := map[core.PieceType]bool{}
		for _, m := range moves {
			if m.To != core.E8 {
				t.Errorf("promotion move should go to E8, got %v", m.To)
			}
			if m.Type != core.PROMOTION {
				t.Errorf("move to E8: Type = %v, want PROMOTION", m.Type)
			}
			if m.HasCapture {
				t.Errorf("forward promotion should not be a capture")
			}
			promoteTypes[m.PromoteTo] = true
		}
		// All four promotion types must be present.
		for _, pt := range []core.PieceType{core.QUEEN, core.ROOK, core.BISHOP, core.KNIGHT} {
			if !promoteTypes[pt] {
				t.Errorf("missing promotion to %v", pt)
			}
		}
	})

	t.Run("a black pawn reaching the last rank by forward push promotes", func(t *testing.T) {
		// Black pawn on E2, E1 empty → 4 promotion moves.
		var board core.Board
		board[core.E2] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.BLACK}

		moves := pawn.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.E2, ctx)

		testutil.AssertMoveCount(t, moves, 4)
		for _, m := range moves {
			if m.To != core.E1 {
				t.Errorf("promotion move should go to E1, got %v", m.To)
			}
			if m.Type != core.PROMOTION {
				t.Errorf("move to E1: Type = %v, want PROMOTION", m.Type)
			}
		}
	})

	t.Run("a promotion with a diagonal capture produces 4 capture-promotion moves", func(t *testing.T) {
		// White pawn on E7, enemy rook on D8 → 4 promotion-capture moves to D8.
		var board core.Board
		board[core.E7] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.D8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := pawn.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.E7, ctx)

		// 4 forward promotions (E8) + 4 capture promotions (D8) = 8 total.
		testutil.AssertMoveCount(t, moves, 8)

		// The 4 capture moves to D8 must carry the captured rook.
		captureCount := 0
		for _, m := range moves {
			if m.To == core.D8 {
				captureCount++
				if !m.HasCapture || m.Captured != (core.Piece{Type: core.ROOK, Color: core.BLACK}) {
					t.Errorf("promotion capture of D8: HasCapture=%v Captured=%v, want black rook", m.HasCapture, m.Captured)
				}
				if m.Type != core.PROMOTION {
					t.Errorf("promotion capture: Type = %v, want PROMOTION", m.Type)
				}
			}
		}
		if captureCount != 4 {
			t.Errorf("expected 4 capture-promotion moves to D8, got %d", captureCount)
		}
	})

	t.Run("a promotion with captures on both diagonals plus forward produces 12 moves", func(t *testing.T) {
		// White pawn on E7, E8 empty, enemy on D8 and F8.
		var board core.Board
		board[core.E7] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.D8] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.BLACK})
		board[core.F8] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.BLACK})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := pawn.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.E7, ctx)

		// 4 forward (E8) + 4 capture D8 + 4 capture F8 = 12.
		testutil.AssertMoveCount(t, moves, 12)

		// Verify the captures carry the right pieces.
		for _, m := range moves {
			if m.To == core.D8 && m.HasCapture && m.Captured.Type != core.QUEEN {
				t.Errorf("capture of D8: Captured type = %v, want QUEEN", m.Captured.Type)
			}
			if m.To == core.F8 && m.HasCapture && m.Captured.Type != core.KNIGHT {
				t.Errorf("capture of F8: Captured type = %v, want KNIGHT", m.Captured.Type)
			}
		}
	})

	t.Run("a pawn whose promotion square is blocked can still capture diagonally", func(t *testing.T) {
		// White pawn on E7, E8 occupied by enemy → no forward promotion,
		// but diagonal captures to D8/F8 (if enemies there) still available.
		var board core.Board
		board[core.E7] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.E8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})  // blocks forward
		board[core.D8] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.BLACK}) // capturable
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := pawn.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.E7, ctx)

		// 4 capture-promotions to D8 only (forward blocked, F8 empty).
		testutil.AssertMoveCount(t, moves, 4)
		for _, m := range moves {
			if m.To != core.D8 {
				t.Errorf("expected all moves to go to D8, got %v", m.To)
			}
		}
	})

	t.Run("a pawn with no promotion moves (forward blocked, no diagonal captures) yields nothing", func(t *testing.T) {
		// White pawn on E7, E8 occupied by own piece, D8 and F8 empty.
		var board core.Board
		board[core.E7] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.E8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE}) // own, blocks forward
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := pawn.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.E7, ctx)

		testutil.AssertNoMoves(t, moves)
	})

	// =========================================================================
	// Edge cases: A-file and H-file pawns
	// =========================================================================

	t.Run("a pawn on the A file has no left diagonal (only right capture)", func(t *testing.T) {
		// White pawn on A4, enemy on B5 → capture B5 + push A5.
		var board core.Board
		board[core.A4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.B5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := pawn.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.A4, ctx)

		testutil.AssertPositionsMatch(t, destinations(moves), []core.Position{core.A5, core.B5})
	})

	t.Run("a pawn on the H file has no right diagonal (only left capture)", func(t *testing.T) {
		// White pawn on H4, enemy on G5 → capture G5 + push H5.
		var board core.Board
		board[core.H4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.G5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := pawn.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.H4, ctx)

		testutil.AssertPositionsMatch(t, destinations(moves), []core.Position{core.H5, core.G5})
	})

	// =========================================================================
	// Black pawn perspective (color flip of all the above)
	// =========================================================================

	t.Run("a black pawn captures downward (white pieces are enemies)", func(t *testing.T) {
		// Black pawn on E5, white pawn on D4 → capture D4 + push E4.
		var board core.Board
		board[core.E5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		board[core.D4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.BLACK}

		moves := pawn.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.E5, ctx)

		testutil.AssertPositionsMatch(t, destinations(moves), []core.Position{core.E4, core.D4})

		for _, m := range moves {
			if m.To == core.D4 {
				if m.Captured != (core.Piece{Type: core.PAWN, Color: core.WHITE}) {
					t.Errorf("black pawn capture: Captured = %v, want white pawn", m.Captured)
				}
				return
			}
		}
		t.Errorf("expected a capture move to D4")
	})

	t.Run("a black pawn treats black pieces as own (not capturable)", func(t *testing.T) {
		// Black pawn on E5, own pawn on D4 → no capture, just push E4.
		var board core.Board
		board[core.E5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		board[core.D4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.BLACK}

		moves := pawn.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.E5, ctx)

		testutil.AssertPositionsMatch(t, destinations(moves), []core.Position{core.E4})
	})

	// =========================================================================
	// Move metadata
	// =========================================================================

	t.Run("every non-promotion move has type NORMAL and carries the mover and source", func(t *testing.T) {
		// White pawn on E2 with open board: push E3, push E4. Both NORMAL.
		var board core.Board
		board[core.E2] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		ctx := core.MoveContext{BoardContext: core.BoardContext{Board: &board}, SideToMove: core.WHITE}

		moves := pawn.PseudoLegalMoves(make([]core.Move, 0, core.MAX_MOVES), core.E2, ctx)

		mover := core.Piece{Type: core.PAWN, Color: core.WHITE}
		for _, m := range moves {
			if m.Type != core.NORMAL {
				t.Errorf("move to %v: Type = %v, want NORMAL", m.To, m.Type)
			}
			if m.From != core.E2 {
				t.Errorf("move to %v: From = %v, want E2", m.To, m.From)
			}
			if m.Piece != mover {
				t.Errorf("move to %v: Piece = %v, want %v", m.To, m.Piece, mover)
			}
		}
	})
}
