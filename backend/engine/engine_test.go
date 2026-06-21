package engine

import (
	"testing"

	"github.com/7-Dany/chess/core"
	"github.com/7-Dany/chess/piece"
)

// This file holds the engine benchmarks. Each measures a single engine
// operation on a representative position, and verifies the zero-allocation
// contract: every benchmark reports 0 B/op, 0 allocs/op.
//
// Positions used:
//   - startPosition: the standard chess opening (32 pieces, full castling).
//   - kiwipete: a dense middlegame position with many sliders (the standard
//     perft test position — exercises the heaviest move-generation paths).
//   - foolsMate: the position after 1.f3 e5 2.g4 Qh4# — white is checkmated
//     and has 0 legal moves. Used to stress HasAnyLegalMoves on a position
//     where every pseudo-move must be tried and rejected.
//   - stalemate: black king vs white king+queen, black to move, no legal
//     moves but not in check.
//   - endgame: king + rook vs lone king. Sparse board, fast scans.
//
// Run with:  go test ./engine -bench=. -benchmem
// Verify:    every line shows "0 B/op  0 allocs/op".

// =============================================================================
// Position builders
// =============================================================================

// startPosition returns the standard chess opening position: 32 pieces,
// both kings on home squares, full castling rights.
func startPosition() core.TurnContext {
	var board core.Board
	back := []core.PieceType{core.ROOK, core.KNIGHT, core.BISHOP, core.QUEEN,
		core.KING, core.BISHOP, core.KNIGHT, core.ROOK}

	for f := range uint8(8) {
		board[core.NewPosition(core.File(f), core.RANK_1)] = core.NewSquare(core.Piece{Type: back[f], Color: core.WHITE})
		board[core.NewPosition(core.File(f), core.RANK_2)] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
		board[core.NewPosition(core.File(f), core.RANK_7)] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
		board[core.NewPosition(core.File(f), core.RANK_8)] = core.NewSquare(core.Piece{Type: back[f], Color: core.BLACK})
	}

	return core.TurnContext{
		MoveContext: core.MoveContext{
			BoardContext: core.BoardContext{Board: &board},
			SideToMove:   core.WHITE,
			Sides: [2]core.SideState{
				{KingPosition: core.E1, CanCastleKingSide: true, CanCastleQueenSide: true},
				{KingPosition: core.E8, CanCastleKingSide: true, CanCastleQueenSide: true},
			},
		},
	}
}

// foolsMate returns the position after 1.f3 e5 2.g4 Qh4# — white is
// checkmated. White has 0 legal moves, so HasAnyLegalMoves must try every
// pseudo-move and reject it. This is the worst case for HasAnyLegalMoves.
func foolsMate() core.TurnContext {
	var board core.Board
	back := []core.PieceType{core.ROOK, core.KNIGHT, core.BISHOP, core.QUEEN,
		core.KING, core.BISHOP, core.KNIGHT, core.ROOK}

	for f := range uint8(8) {
		board[core.NewPosition(core.File(f), core.RANK_1)] = core.NewSquare(core.Piece{Type: back[f], Color: core.WHITE})
		board[core.NewPosition(core.File(f), core.RANK_8)] = core.NewSquare(core.Piece{Type: back[f], Color: core.BLACK})
	}

	// White pawns: most on rank 2, but f-pawn advanced to F3 and g-pawn to G4.
	for _, f := range []core.File{core.FILE_A, core.FILE_B, core.FILE_C, core.FILE_D, core.FILE_E, core.FILE_H} {
		board[core.NewPosition(f, core.RANK_2)] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
	}
	board[core.F3] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
	board[core.G4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})

	// Black pawns: most on rank 7, e-pawn advanced to E5.
	for _, f := range []core.File{core.FILE_A, core.FILE_B, core.FILE_C, core.FILE_D, core.FILE_F, core.FILE_G, core.FILE_H} {
		board[core.NewPosition(f, core.RANK_7)] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
	}
	board[core.E5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})

	// Black queen moved from D8 to H4 — delivering checkmate along the e1-h4 diagonal.
	board[core.D8] = core.EmptySquare
	board[core.H4] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.BLACK})

	return core.TurnContext{
		MoveContext: core.MoveContext{
			BoardContext: core.BoardContext{Board: &board},
			SideToMove:   core.WHITE, // white to move, but checkmated
			Sides: [2]core.SideState{
				{KingPosition: core.E1},
				{KingPosition: core.E8},
			},
		},
	}
}

// stalemate returns a position where black is stalemated: black king on A8,
// white queen on C7, white king on C6. Black has no legal moves but is not
// in check. HasAnyLegalMoves must try every king move and reject each.
func stalemate() core.TurnContext {
	var board core.Board
	board[core.A8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
	board[core.C7] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.WHITE})
	board[core.C6] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})

	return core.TurnContext{
		MoveContext: core.MoveContext{
			BoardContext: core.BoardContext{Board: &board},
			SideToMove:   core.BLACK,
			Sides: [2]core.SideState{
				{KingPosition: core.C6},
				{KingPosition: core.A8},
			},
		},
	}
}

// kiwipete returns the "Kiwipete" position — a standard perft test position
// with many sliders and a dense board. It exercises the heaviest
// move-generation paths (lots of sliding moves to generate and filter).
//
// FEN: r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w - -
func kiwipete() core.TurnContext {
	var board core.Board
	place := func(file core.File, rank core.Rank, p core.Piece) {
		board[core.NewPosition(file, rank)] = core.NewSquare(p)
	}

	// Black back rank: rook A8, king E8, rook H8.
	place(core.FILE_A, core.RANK_8, core.Piece{Type: core.ROOK, Color: core.BLACK})
	place(core.FILE_E, core.RANK_8, core.Piece{Type: core.KING, Color: core.BLACK})
	place(core.FILE_H, core.RANK_8, core.Piece{Type: core.ROOK, Color: core.BLACK})

	// Black rank 7: pawns + queen + bishop.
	place(core.FILE_A, core.RANK_7, core.Piece{Type: core.PAWN, Color: core.BLACK})
	place(core.FILE_C, core.RANK_7, core.Piece{Type: core.PAWN, Color: core.BLACK})
	place(core.FILE_D, core.RANK_7, core.Piece{Type: core.PAWN, Color: core.BLACK})
	place(core.FILE_E, core.RANK_7, core.Piece{Type: core.QUEEN, Color: core.BLACK})
	place(core.FILE_F, core.RANK_7, core.Piece{Type: core.PAWN, Color: core.BLACK})
	place(core.FILE_G, core.RANK_7, core.Piece{Type: core.BISHOP, Color: core.BLACK})

	// Black rank 6: two knights + two pawns.
	place(core.FILE_B, core.RANK_6, core.Piece{Type: core.KNIGHT, Color: core.BLACK})
	place(core.FILE_E, core.RANK_6, core.Piece{Type: core.PAWN, Color: core.BLACK})
	place(core.FILE_F, core.RANK_6, core.Piece{Type: core.KNIGHT, Color: core.BLACK})
	place(core.FILE_G, core.RANK_6, core.Piece{Type: core.PAWN, Color: core.BLACK})

	// Rank 5: white pawn E5, white knight D5, black pawn B5.
	place(core.FILE_E, core.RANK_5, core.Piece{Type: core.PAWN, Color: core.WHITE})
	place(core.FILE_D, core.RANK_5, core.Piece{Type: core.KNIGHT, Color: core.WHITE})
	place(core.FILE_B, core.RANK_5, core.Piece{Type: core.PAWN, Color: core.BLACK})

	// Rank 4: white pawn E4, black pawn H4.
	place(core.FILE_E, core.RANK_4, core.Piece{Type: core.PAWN, Color: core.WHITE})
	place(core.FILE_H, core.RANK_4, core.Piece{Type: core.PAWN, Color: core.BLACK})

	// Rank 3: white knight C3, white queen F3, black pawn H3.
	place(core.FILE_C, core.RANK_3, core.Piece{Type: core.KNIGHT, Color: core.WHITE})
	place(core.FILE_F, core.RANK_3, core.Piece{Type: core.QUEEN, Color: core.WHITE})
	place(core.FILE_H, core.RANK_3, core.Piece{Type: core.PAWN, Color: core.BLACK})

	// White rank 2: pawns + two bishops.
	place(core.FILE_A, core.RANK_2, core.Piece{Type: core.PAWN, Color: core.WHITE})
	place(core.FILE_B, core.RANK_2, core.Piece{Type: core.PAWN, Color: core.WHITE})
	place(core.FILE_C, core.RANK_2, core.Piece{Type: core.PAWN, Color: core.WHITE})
	place(core.FILE_D, core.RANK_2, core.Piece{Type: core.BISHOP, Color: core.WHITE})
	place(core.FILE_E, core.RANK_2, core.Piece{Type: core.BISHOP, Color: core.WHITE})
	place(core.FILE_F, core.RANK_2, core.Piece{Type: core.PAWN, Color: core.WHITE})
	place(core.FILE_G, core.RANK_2, core.Piece{Type: core.PAWN, Color: core.WHITE})
	place(core.FILE_H, core.RANK_2, core.Piece{Type: core.PAWN, Color: core.WHITE})

	// White back rank: rook A1, king E1, rook H1.
	place(core.FILE_A, core.RANK_1, core.Piece{Type: core.ROOK, Color: core.WHITE})
	place(core.FILE_E, core.RANK_1, core.Piece{Type: core.KING, Color: core.WHITE})
	place(core.FILE_H, core.RANK_1, core.Piece{Type: core.ROOK, Color: core.WHITE})

	return core.TurnContext{
		MoveContext: core.MoveContext{
			BoardContext: core.BoardContext{Board: &board},
			SideToMove:   core.WHITE,
			Sides: [2]core.SideState{
				{KingPosition: core.E1, CanCastleKingSide: true, CanCastleQueenSide: true},
				{KingPosition: core.E8, CanCastleKingSide: true, CanCastleQueenSide: true},
			},
		},
	}
}

// endgame returns a sparse position: white king + rook vs black king.
// Few pieces means fast scans — useful as a lower bound for IsSquareAttacked.
func endgame() core.TurnContext {
	var board core.Board
	board[core.H1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
	board[core.A3] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
	board[core.H8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})

	return core.TurnContext{
		MoveContext: core.MoveContext{
			BoardContext: core.BoardContext{Board: &board},
			SideToMove:   core.WHITE,
			Sides: [2]core.SideState{
				{KingPosition: core.H1},
				{KingPosition: core.H8},
			},
		},
	}
}

// emptyBoardWithKing returns a board with only a white king on E4 (both side
// states point at E4). Used to measure IsSquareAttacked on an otherwise
// empty board — the leaper checks all run, but no sliders find anything.
func emptyBoardWithKing() core.TurnContext {
	var board core.Board
	board[core.E4] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
	return core.TurnContext{
		MoveContext: core.MoveContext{
			BoardContext: core.BoardContext{Board: &board},
			SideToMove:   core.WHITE,
			Sides: [2]core.SideState{
				{KingPosition: core.E4},
				{KingPosition: core.E4},
			},
		},
	}
}

// attackedSquare returns a board where E4 (white king) is attacked by a
// black rook on E8 down the E-file. Used to measure IsSquareAttacked when
// the square IS attacked (the orthogonal slider scan finds a hit early).
func attackedSquare() core.TurnContext {
	var board core.Board
	board[core.E4] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
	board[core.E8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
	return core.TurnContext{
		MoveContext: core.MoveContext{
			BoardContext: core.BoardContext{Board: &board},
			SideToMove:   core.WHITE,
			Sides: [2]core.SideState{
				{KingPosition: core.E4},
				{KingPosition: core.E4},
			},
		},
	}
}

// =============================================================================
// GetPseudoLegalMoves — pseudo-legal move generation per piece type.
//
// Measures the cost of generating moves for a single piece, WITHOUT the
// king-safety filter. The buffer is stack-allocated (var buf [piece.MAX_MOVES]);
// passing buf[:0] each iteration means no allocation.
// =============================================================================

func BenchmarkGetPseudoLegalMoves_Pawn(b *testing.B) {
	engine := GetDefaultEngine()
	ctx := startPosition()
	var buf [piece.MAX_MOVES]core.Move
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.GetPseudoLegalMoves(buf[:0], core.A2, ctx)
	}
}

func BenchmarkGetPseudoLegalMoves_Knight(b *testing.B) {
	engine := GetDefaultEngine()
	ctx := startPosition()
	var buf [piece.MAX_MOVES]core.Move
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.GetPseudoLegalMoves(buf[:0], core.B1, ctx)
	}
}

func BenchmarkGetPseudoLegalMoves_Bishop(b *testing.B) {
	engine := GetDefaultEngine()
	ctx := startPosition()
	var buf [piece.MAX_MOVES]core.Move
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.GetPseudoLegalMoves(buf[:0], core.C1, ctx)
	}
}

func BenchmarkGetPseudoLegalMoves_Rook(b *testing.B) {
	engine := GetDefaultEngine()
	ctx := startPosition()
	var buf [piece.MAX_MOVES]core.Move
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.GetPseudoLegalMoves(buf[:0], core.A1, ctx)
	}
}

func BenchmarkGetPseudoLegalMoves_Queen(b *testing.B) {
	engine := GetDefaultEngine()
	ctx := startPosition()
	var buf [piece.MAX_MOVES]core.Move
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.GetPseudoLegalMoves(buf[:0], core.D1, ctx)
	}
}

func BenchmarkGetPseudoLegalMoves_King(b *testing.B) {
	engine := GetDefaultEngine()
	ctx := startPosition()
	var buf [piece.MAX_MOVES]core.Move
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.GetPseudoLegalMoves(buf[:0], core.E1, ctx)
	}
}

// Kiwipete positions: a knight and queen on a dense board. These exercise
// the piece logic with real blockers (the starting position above has the
// pieces hemmed in by pawns, so they generate few moves; Kiwipete gives
// them open lines).
func BenchmarkGetPseudoLegalMoves_KiwipeteKnight(b *testing.B) {
	engine := GetDefaultEngine()
	ctx := kiwipete()
	var buf [piece.MAX_MOVES]core.Move
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.GetPseudoLegalMoves(buf[:0], core.D5, ctx)
	}
}

func BenchmarkGetPseudoLegalMoves_KiwipeteQueen(b *testing.B) {
	engine := GetDefaultEngine()
	ctx := kiwipete()
	var buf [piece.MAX_MOVES]core.Move
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.GetPseudoLegalMoves(buf[:0], core.F3, ctx)
	}
}

// =============================================================================
// GetLegalMoves — legal move generation per piece type.
//
// Measures the cost of generating pseudo-legal moves AND filtering them for
// king safety. Each pseudo-move is applied, the king is checked, and the
// move is undone. This is the real per-piece cost during search.
// =============================================================================

func BenchmarkGetLegalMoves_Pawn(b *testing.B) {
	engine := GetDefaultEngine()
	ctx := startPosition()
	var buf [piece.MAX_MOVES]core.Move
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.GetLegalMoves(buf[:0], core.A2, ctx)
	}
}

func BenchmarkGetLegalMoves_Knight(b *testing.B) {
	engine := GetDefaultEngine()
	ctx := startPosition()
	var buf [piece.MAX_MOVES]core.Move
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.GetLegalMoves(buf[:0], core.B1, ctx)
	}
}

func BenchmarkGetLegalMoves_Bishop(b *testing.B) {
	engine := GetDefaultEngine()
	ctx := startPosition()
	var buf [piece.MAX_MOVES]core.Move
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.GetLegalMoves(buf[:0], core.C1, ctx)
	}
}

func BenchmarkGetLegalMoves_Rook(b *testing.B) {
	engine := GetDefaultEngine()
	ctx := startPosition()
	var buf [piece.MAX_MOVES]core.Move
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.GetLegalMoves(buf[:0], core.A1, ctx)
	}
}

func BenchmarkGetLegalMoves_Queen(b *testing.B) {
	engine := GetDefaultEngine()
	ctx := startPosition()
	var buf [piece.MAX_MOVES]core.Move
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.GetLegalMoves(buf[:0], core.D1, ctx)
	}
}

func BenchmarkGetLegalMoves_King(b *testing.B) {
	engine := GetDefaultEngine()
	ctx := startPosition()
	var buf [piece.MAX_MOVES]core.Move
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.GetLegalMoves(buf[:0], core.E1, ctx)
	}
}

// Kiwipete: a knight, queen, and king on a dense board. These have many
// pseudo-legal moves, so the king-safety filter runs many apply/undo
// cycles — the dominant cost during search.
func BenchmarkGetLegalMoves_KiwipeteKnight(b *testing.B) {
	engine := GetDefaultEngine()
	ctx := kiwipete()
	var buf [piece.MAX_MOVES]core.Move
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.GetLegalMoves(buf[:0], core.D5, ctx)
	}
}

func BenchmarkGetLegalMoves_KiwipeteQueen(b *testing.B) {
	engine := GetDefaultEngine()
	ctx := kiwipete()
	var buf [piece.MAX_MOVES]core.Move
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.GetLegalMoves(buf[:0], core.F3, ctx)
	}
}

func BenchmarkGetLegalMoves_KiwipeteKing(b *testing.B) {
	engine := GetDefaultEngine()
	ctx := kiwipete()
	var buf [piece.MAX_MOVES]core.Move
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.GetLegalMoves(buf[:0], core.E1, ctx)
	}
}

// =============================================================================
// HasAnyLegalMoves — checkmate / stalemate detection.
//
// Returns true as soon as one legal move is found; returns false only after
// trying every pseudo-move of every piece. The "false" cases (checkmate,
// stalemate) are the worst case and the most useful to benchmark.
// =============================================================================

func BenchmarkHasAnyLegalMoves_Start(b *testing.B) {
	engine := GetDefaultEngine()
	ctx := startPosition()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.HasAnyLegalMoves(ctx)
	}
}

// foolsMate: white is checkmated. Every pseudo-move must be tried and
// rejected — the worst case for HasAnyLegalMoves.
func BenchmarkHasAnyLegalMoves_FoolsMate(b *testing.B) {
	engine := GetDefaultEngine()
	ctx := foolsMate()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.HasAnyLegalMoves(ctx)
	}
}

// stalemate: black has no legal moves but isn't in check. Every king move
// must be tried and rejected (each leaves the king adjacent to the queen
// or king).
func BenchmarkHasAnyLegalMoves_Stalemate(b *testing.B) {
	engine := GetDefaultEngine()
	ctx := stalemate()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.HasAnyLegalMoves(ctx)
	}
}

func BenchmarkHasAnyLegalMoves_Kiwipete(b *testing.B) {
	engine := GetDefaultEngine()
	ctx := kiwipete()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.HasAnyLegalMoves(ctx)
	}
}

func BenchmarkHasAnyLegalMoves_Endgame(b *testing.B) {
	engine := GetDefaultEngine()
	ctx := endgame()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.HasAnyLegalMoves(ctx)
	}
}

// =============================================================================
// IsSquareAttacked — attack detection.
//
// Scans from the target outward: 3 leaper checks (knight, king, pawn) +
// 8 slider rays (4 diagonal for bishop/queen, 4 orthogonal for rook/queen).
// Never allocates. The position determines how far the slider scans run
// before hitting a blocker or the board edge.
// =============================================================================

// Empty board: every slider ray runs to the board edge (no blockers).
// The leaper checks all return false quickly.
func BenchmarkIsSquareAttacked_Empty(b *testing.B) {
	engine := GetDefaultEngine()
	ctx := emptyBoardWithKing().BoardContext
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.IsSquareAttacked(core.E4, core.BLACK, ctx)
	}
}

// Start position: asking "is E1 (white king) attacked by black?" — the
// black rooks and bishops are blocked by pawns, so all slider scans stop
// early at the first blocker.
func BenchmarkIsSquareAttacked_Start(b *testing.B) {
	engine := GetDefaultEngine()
	ctx := startPosition().BoardContext
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.IsSquareAttacked(core.E1, core.BLACK, ctx)
	}
}

func BenchmarkIsSquareAttacked_Kiwipete(b *testing.B) {
	engine := GetDefaultEngine()
	ctx := kiwipete().BoardContext
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.IsSquareAttacked(core.E1, core.BLACK, ctx)
	}
}

// Attacked square: E4 is attacked by a black rook on E8. The orthogonal
// scan finds the rook quickly (after checking E5, E6, E7) and returns true.
func BenchmarkIsSquareAttacked_Attacked(b *testing.B) {
	engine := GetDefaultEngine()
	ctx := attackedSquare().BoardContext
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.IsSquareAttacked(core.E4, core.BLACK, ctx)
	}
}

// Corner square: A1 has only 3 diagonal rays and 2 orthogonal rays (the
// others fall off the board immediately). Fewer rays to scan.
func BenchmarkIsSquareAttacked_Corner(b *testing.B) {
	engine := GetDefaultEngine()
	ctx := startPosition().BoardContext
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.IsSquareAttacked(core.A1, core.BLACK, ctx)
	}
}

// =============================================================================
// Apply — applying a move to the board.
//
// Each iteration copies the base context (so Apply has a fresh board to
// mutate) and applies the move. The Copy() is outside the timed loop's
// measurement of Apply itself — but since every iteration needs a fresh
// board, the copy is included. Copy() is allocation-free (escape analysis
// keeps it on the stack); the Apply itself is also 0-alloc.
// =============================================================================

// benchApplyMove is a helper that applies the same move to a fresh copy of
// base, b.N times. Defined to remove the repeated copy-loop boilerplate
// from each Apply benchmark.
func benchApplyMove(b *testing.B, base core.TurnContext, move core.Move) {
	b.Helper()
	engine := GetDefaultEngine()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := base.Copy()
		_ = engine.Apply(&ctx, move)
	}
}

func BenchmarkApply_NormalPawnPush(b *testing.B) {
	benchApplyMove(b, startPosition(), core.Move{
		Type:  core.NORMAL,
		Piece: core.Piece{Type: core.PAWN, Color: core.WHITE},
		From:  core.E2,
		To:    core.E4,
	})
}

func BenchmarkApply_NormalKnightMove(b *testing.B) {
	benchApplyMove(b, startPosition(), core.Move{
		Type:  core.NORMAL,
		Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE},
		From:  core.B1,
		To:    core.C3,
	})
}

func BenchmarkApply_Capture(b *testing.B) {
	// White pawn on E4 captures black pawn on D5.
	var board core.Board
	board[core.E4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
	board[core.D5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
	base := core.TurnContext{
		MoveContext: core.MoveContext{
			BoardContext: core.BoardContext{Board: &board},
			SideToMove:   core.WHITE,
			Sides: [2]core.SideState{
				{KingPosition: core.E1},
				{KingPosition: core.E8},
			},
		},
	}
	benchApplyMove(b, base, core.Move{
		Type:       core.NORMAL,
		Piece:      core.Piece{Type: core.PAWN, Color: core.WHITE},
		From:       core.E4,
		To:         core.D5,
		HasCapture: true,
		Captured:   core.Piece{Type: core.PAWN, Color: core.BLACK},
	})
}

func BenchmarkApply_EnPassant(b *testing.B) {
	// White pawn on D5 captures en passant onto E6 (EP target E6); the
	// captured black pawn sits on E5.
	var board core.Board
	board[core.D5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
	board[core.E5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
	base := core.TurnContext{
		MoveContext: core.MoveContext{
			BoardContext:    core.BoardContext{Board: &board},
			SideToMove:      core.WHITE,
			EnPassantTarget: core.E6,
			Sides: [2]core.SideState{
				{KingPosition: core.E1},
				{KingPosition: core.E8},
			},
		},
	}
	benchApplyMove(b, base, core.Move{
		Type:       core.EN_PASSANT,
		Piece:      core.Piece{Type: core.PAWN, Color: core.WHITE},
		From:       core.D5,
		To:         core.E6,
		HasCapture: true,
		Captured:   core.Piece{Type: core.PAWN, Color: core.BLACK},
	})
}

func BenchmarkApply_CastlingKingSide(b *testing.B) {
	// King on E1, rook on H1. King castles king-side to G1 (rook hops to F1).
	var board core.Board
	board[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
	board[core.H1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
	base := core.TurnContext{
		MoveContext: core.MoveContext{
			BoardContext: core.BoardContext{Board: &board},
			SideToMove:   core.WHITE,
			Sides: [2]core.SideState{
				{KingPosition: core.E1, CanCastleKingSide: true, CanCastleQueenSide: true},
				{KingPosition: core.E8},
			},
		},
	}
	benchApplyMove(b, base, core.Move{
		Type:  core.CASTLING,
		Piece: core.Piece{Type: core.KING, Color: core.WHITE},
		From:  core.E1,
		To:    core.G1,
	})
}

func BenchmarkApply_CastlingQueenSide(b *testing.B) {
	var board core.Board
	board[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
	board[core.A1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
	base := core.TurnContext{
		MoveContext: core.MoveContext{
			BoardContext: core.BoardContext{Board: &board},
			SideToMove:   core.WHITE,
			Sides: [2]core.SideState{
				{KingPosition: core.E1, CanCastleKingSide: true, CanCastleQueenSide: true},
				{KingPosition: core.E8},
			},
		},
	}
	benchApplyMove(b, base, core.Move{
		Type:  core.CASTLING,
		Piece: core.Piece{Type: core.KING, Color: core.WHITE},
		From:  core.E1,
		To:    core.C1,
	})
}

func BenchmarkApply_Promotion(b *testing.B) {
	// White pawn on E7 promotes to queen on E8.
	var board core.Board
	board[core.E7] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
	base := core.TurnContext{
		MoveContext: core.MoveContext{
			BoardContext: core.BoardContext{Board: &board},
			SideToMove:   core.WHITE,
			Sides: [2]core.SideState{
				{KingPosition: core.E1},
				{KingPosition: core.E8},
			},
		},
	}
	benchApplyMove(b, base, core.Move{
		Type:      core.PROMOTION,
		Piece:     core.Piece{Type: core.PAWN, Color: core.WHITE},
		From:      core.E7,
		To:        core.E8,
		PromoteTo: core.QUEEN,
	})
}

// Double pawn push: a NORMAL move that sets the en passant target. Same
// code path as a normal pawn push, but exercises SetEnPassantTarget's
// "this is a double push" branch.
func BenchmarkApply_DoublePawnPush(b *testing.B) {
	benchApplyMove(b, startPosition(), core.Move{
		Type:  core.NORMAL,
		Piece: core.Piece{Type: core.PAWN, Color: core.WHITE},
		From:  core.E2,
		To:    core.E4,
	})
}

// =============================================================================
// Undo — reversing an Apply.
//
// Each iteration copies the base context, applies the move, then undoes it.
// The measurement includes Apply + Copy (unavoidable setup) + Undo. To see
// the Undo cost in isolation, compare these numbers against the matching
// Apply benchmark above — the difference is roughly the Undo cost.
// =============================================================================

// benchApplyThenUndo is a helper that applies then undoes the same move,
// b.N times. Defined to remove the repeated copy-apply-undo boilerplate.
func benchApplyThenUndo(b *testing.B, base core.TurnContext, move core.Move) {
	b.Helper()
	engine := GetDefaultEngine()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := base.Copy()
		snap := engine.Apply(&ctx, move)
		engine.Undo(&ctx, snap)
	}
}

func BenchmarkUndo_NormalPawnPush(b *testing.B) {
	benchApplyThenUndo(b, startPosition(), core.Move{
		Type:  core.NORMAL,
		Piece: core.Piece{Type: core.PAWN, Color: core.WHITE},
		From:  core.E2,
		To:    core.E4,
	})
}

func BenchmarkUndo_NormalKnightMove(b *testing.B) {
	benchApplyThenUndo(b, startPosition(), core.Move{
		Type:  core.NORMAL,
		Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE},
		From:  core.B1,
		To:    core.C3,
	})
}

func BenchmarkUndo_Capture(b *testing.B) {
	var board core.Board
	board[core.E4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
	board[core.D5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
	base := core.TurnContext{
		MoveContext: core.MoveContext{
			BoardContext: core.BoardContext{Board: &board},
			SideToMove:   core.WHITE,
			Sides: [2]core.SideState{
				{KingPosition: core.E1},
				{KingPosition: core.E8},
			},
		},
	}
	benchApplyThenUndo(b, base, core.Move{
		Type:       core.NORMAL,
		Piece:      core.Piece{Type: core.PAWN, Color: core.WHITE},
		From:       core.E4,
		To:         core.D5,
		HasCapture: true,
		Captured:   core.Piece{Type: core.PAWN, Color: core.BLACK},
	})
}

func BenchmarkUndo_EnPassant(b *testing.B) {
	var board core.Board
	board[core.D5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
	board[core.E5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
	base := core.TurnContext{
		MoveContext: core.MoveContext{
			BoardContext:    core.BoardContext{Board: &board},
			SideToMove:      core.WHITE,
			EnPassantTarget: core.E6,
			Sides: [2]core.SideState{
				{KingPosition: core.E1},
				{KingPosition: core.E8},
			},
		},
	}
	benchApplyThenUndo(b, base, core.Move{
		Type:       core.EN_PASSANT,
		Piece:      core.Piece{Type: core.PAWN, Color: core.WHITE},
		From:       core.D5,
		To:         core.E6,
		HasCapture: true,
		Captured:   core.Piece{Type: core.PAWN, Color: core.BLACK},
	})
}

func BenchmarkUndo_CastlingKingSide(b *testing.B) {
	var board core.Board
	board[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
	board[core.H1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
	base := core.TurnContext{
		MoveContext: core.MoveContext{
			BoardContext: core.BoardContext{Board: &board},
			SideToMove:   core.WHITE,
			Sides: [2]core.SideState{
				{KingPosition: core.E1, CanCastleKingSide: true, CanCastleQueenSide: true},
				{KingPosition: core.E8},
			},
		},
	}
	benchApplyThenUndo(b, base, core.Move{
		Type:  core.CASTLING,
		Piece: core.Piece{Type: core.KING, Color: core.WHITE},
		From:  core.E1,
		To:    core.G1,
	})
}

func BenchmarkUndo_Promotion(b *testing.B) {
	var board core.Board
	board[core.E7] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
	base := core.TurnContext{
		MoveContext: core.MoveContext{
			BoardContext: core.BoardContext{Board: &board},
			SideToMove:   core.WHITE,
			Sides: [2]core.SideState{
				{KingPosition: core.E1},
				{KingPosition: core.E8},
			},
		},
	}
	benchApplyThenUndo(b, base, core.Move{
		Type:      core.PROMOTION,
		Piece:     core.Piece{Type: core.PAWN, Color: core.WHITE},
		From:      core.E7,
		To:        core.E8,
		PromoteTo: core.QUEEN,
	})
}

// =============================================================================
// Apply + Undo round-trip — the real search hot-path cost.
//
// During a search, every node applies a move, recurses, then undoes. The
// round-trip cost (Apply + Undo, excluding the copy) is what matters for
// nodes-per-second. These three benchmarks measure representative move
// types at search-relevant positions.
// =============================================================================

func BenchmarkApplyUndo_PawnPush(b *testing.B) {
	benchApplyThenUndo(b, startPosition(), core.Move{
		Type:  core.NORMAL,
		Piece: core.Piece{Type: core.PAWN, Color: core.WHITE},
		From:  core.E2,
		To:    core.E4,
	})
}

func BenchmarkApplyUndo_KnightMove(b *testing.B) {
	benchApplyThenUndo(b, startPosition(), core.Move{
		Type:  core.NORMAL,
		Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE},
		From:  core.B1,
		To:    core.C3,
	})
}

func BenchmarkApplyUndo_Castling(b *testing.B) {
	var board core.Board
	board[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
	board[core.H1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
	base := core.TurnContext{
		MoveContext: core.MoveContext{
			BoardContext: core.BoardContext{Board: &board},
			SideToMove:   core.WHITE,
			Sides: [2]core.SideState{
				{KingPosition: core.E1, CanCastleKingSide: true, CanCastleQueenSide: true},
				{KingPosition: core.E8},
			},
		},
	}
	benchApplyThenUndo(b, base, core.Move{
		Type:  core.CASTLING,
		Piece: core.Piece{Type: core.KING, Color: core.WHITE},
		From:  core.E1,
		To:    core.G1,
	})
}
