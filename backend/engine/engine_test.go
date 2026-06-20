package engine

import (
	"testing"

	"github.com/7-Dany/chess/core"
)

// =============================================================================
// Position builders
// =============================================================================

func benchStartPos() core.TurnContext {
	var board core.Board
	back := []core.PieceType{core.ROOK, core.KNIGHT, core.BISHOP, core.QUEEN,
		core.KING, core.BISHOP, core.KNIGHT, core.ROOK}

	for f := uint8(0); f < 8; f++ {
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

func benchFoolsMate() core.TurnContext {
	var board core.Board
	back := []core.PieceType{core.ROOK, core.KNIGHT, core.BISHOP, core.QUEEN,
		core.KING, core.BISHOP, core.KNIGHT, core.ROOK}

	for f := uint8(0); f < 8; f++ {
		board[core.NewPosition(core.File(f), core.RANK_1)] = core.NewSquare(core.Piece{Type: back[f], Color: core.WHITE})
		board[core.NewPosition(core.File(f), core.RANK_8)] = core.NewSquare(core.Piece{Type: back[f], Color: core.BLACK})
	}

	for _, f := range []core.File{core.FILE_A, core.FILE_B, core.FILE_C, core.FILE_D, core.FILE_E, core.FILE_H} {
		board[core.NewPosition(f, core.RANK_2)] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
	}
	board[core.F3] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
	board[core.G4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})

	for _, f := range []core.File{core.FILE_A, core.FILE_B, core.FILE_C, core.FILE_D, core.FILE_F, core.FILE_G, core.FILE_H} {
		board[core.NewPosition(f, core.RANK_7)] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
	}
	board[core.E5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})

	board[core.D8] = core.EmptySquare
	board[core.H4] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.BLACK})

	return core.TurnContext{
		MoveContext: core.MoveContext{
			BoardContext: core.BoardContext{Board: &board},
			SideToMove:   core.WHITE,
			Sides: [2]core.SideState{
				{KingPosition: core.E1},
				{KingPosition: core.E8},
			},
		},
	}
}

func benchStalemate() core.TurnContext {
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

func benchKiwipete() core.TurnContext {
	var board core.Board
	place := func(file core.File, rank core.Rank, p core.Piece) {
		board[core.NewPosition(file, rank)] = core.NewSquare(p)
	}

	place(core.FILE_A, core.RANK_8, core.Piece{Type: core.ROOK, Color: core.BLACK})
	place(core.FILE_E, core.RANK_8, core.Piece{Type: core.KING, Color: core.BLACK})
	place(core.FILE_H, core.RANK_8, core.Piece{Type: core.ROOK, Color: core.BLACK})

	place(core.FILE_A, core.RANK_7, core.Piece{Type: core.PAWN, Color: core.BLACK})
	place(core.FILE_C, core.RANK_7, core.Piece{Type: core.PAWN, Color: core.BLACK})
	place(core.FILE_D, core.RANK_7, core.Piece{Type: core.PAWN, Color: core.BLACK})
	place(core.FILE_E, core.RANK_7, core.Piece{Type: core.QUEEN, Color: core.BLACK})
	place(core.FILE_F, core.RANK_7, core.Piece{Type: core.PAWN, Color: core.BLACK})
	place(core.FILE_G, core.RANK_7, core.Piece{Type: core.BISHOP, Color: core.BLACK})

	place(core.FILE_B, core.RANK_6, core.Piece{Type: core.KNIGHT, Color: core.BLACK})
	place(core.FILE_E, core.RANK_6, core.Piece{Type: core.PAWN, Color: core.BLACK})
	place(core.FILE_F, core.RANK_6, core.Piece{Type: core.KNIGHT, Color: core.BLACK})
	place(core.FILE_G, core.RANK_6, core.Piece{Type: core.PAWN, Color: core.BLACK})

	place(core.FILE_E, core.RANK_5, core.Piece{Type: core.PAWN, Color: core.WHITE})
	place(core.FILE_D, core.RANK_5, core.Piece{Type: core.KNIGHT, Color: core.WHITE})
	place(core.FILE_B, core.RANK_5, core.Piece{Type: core.PAWN, Color: core.BLACK})

	place(core.FILE_E, core.RANK_4, core.Piece{Type: core.PAWN, Color: core.WHITE})
	place(core.FILE_H, core.RANK_4, core.Piece{Type: core.PAWN, Color: core.BLACK})

	place(core.FILE_C, core.RANK_3, core.Piece{Type: core.KNIGHT, Color: core.WHITE})
	place(core.FILE_F, core.RANK_3, core.Piece{Type: core.QUEEN, Color: core.WHITE})
	place(core.FILE_H, core.RANK_3, core.Piece{Type: core.PAWN, Color: core.BLACK})

	place(core.FILE_A, core.RANK_2, core.Piece{Type: core.PAWN, Color: core.WHITE})
	place(core.FILE_B, core.RANK_2, core.Piece{Type: core.PAWN, Color: core.WHITE})
	place(core.FILE_C, core.RANK_2, core.Piece{Type: core.PAWN, Color: core.WHITE})
	place(core.FILE_D, core.RANK_2, core.Piece{Type: core.BISHOP, Color: core.WHITE})
	place(core.FILE_E, core.RANK_2, core.Piece{Type: core.BISHOP, Color: core.WHITE})
	place(core.FILE_F, core.RANK_2, core.Piece{Type: core.PAWN, Color: core.WHITE})
	place(core.FILE_G, core.RANK_2, core.Piece{Type: core.PAWN, Color: core.WHITE})
	place(core.FILE_H, core.RANK_2, core.Piece{Type: core.PAWN, Color: core.WHITE})

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

func benchEndgame() core.TurnContext {
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

func benchEmptyBoard() core.TurnContext {
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

func benchAttackedSquare() core.TurnContext {
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
// GetPseudoLegalMoves — per piece type
// =============================================================================

func BenchmarkGetPseudoLegalMoves_Pawn(b *testing.B) {
	engine := NewDefaultEngine()
	ctx := benchStartPos()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.GetPseudoLegalMoves(core.A2, ctx)
	}
}

func BenchmarkGetPseudoLegalMoves_Knight(b *testing.B) {
	engine := NewDefaultEngine()
	ctx := benchStartPos()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.GetPseudoLegalMoves(core.B1, ctx)
	}
}

func BenchmarkGetPseudoLegalMoves_Bishop(b *testing.B) {
	engine := NewDefaultEngine()
	ctx := benchStartPos()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.GetPseudoLegalMoves(core.C1, ctx)
	}
}

func BenchmarkGetPseudoLegalMoves_Rook(b *testing.B) {
	engine := NewDefaultEngine()
	ctx := benchStartPos()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.GetPseudoLegalMoves(core.A1, ctx)
	}
}

func BenchmarkGetPseudoLegalMoves_Queen(b *testing.B) {
	engine := NewDefaultEngine()
	ctx := benchStartPos()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.GetPseudoLegalMoves(core.D1, ctx)
	}
}

func BenchmarkGetPseudoLegalMoves_King(b *testing.B) {
	engine := NewDefaultEngine()
	ctx := benchStartPos()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.GetPseudoLegalMoves(core.E1, ctx)
	}
}

func BenchmarkGetPseudoLegalMoves_KiwipeteKnight(b *testing.B) {
	engine := NewDefaultEngine()
	ctx := benchKiwipete()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.GetPseudoLegalMoves(core.D5, ctx)
	}
}

func BenchmarkGetPseudoLegalMoves_KiwipeteQueen(b *testing.B) {
	engine := NewDefaultEngine()
	ctx := benchKiwipete()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.GetPseudoLegalMoves(core.F3, ctx)
	}
}

// =============================================================================
// GetLegalMoves — per piece type
// =============================================================================

func BenchmarkGetLegalMoves_Pawn(b *testing.B) {
	engine := NewDefaultEngine()
	ctx := benchStartPos()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.GetLegalMoves(core.A2, ctx)
	}
}

func BenchmarkGetLegalMoves_Knight(b *testing.B) {
	engine := NewDefaultEngine()
	ctx := benchStartPos()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.GetLegalMoves(core.B1, ctx)
	}
}

func BenchmarkGetLegalMoves_Bishop(b *testing.B) {
	engine := NewDefaultEngine()
	ctx := benchStartPos()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.GetLegalMoves(core.C1, ctx)
	}
}

func BenchmarkGetLegalMoves_Rook(b *testing.B) {
	engine := NewDefaultEngine()
	ctx := benchStartPos()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.GetLegalMoves(core.A1, ctx)
	}
}

func BenchmarkGetLegalMoves_Queen(b *testing.B) {
	engine := NewDefaultEngine()
	ctx := benchStartPos()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.GetLegalMoves(core.D1, ctx)
	}
}

func BenchmarkGetLegalMoves_King(b *testing.B) {
	engine := NewDefaultEngine()
	ctx := benchStartPos()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.GetLegalMoves(core.E1, ctx)
	}
}

func BenchmarkGetLegalMoves_KiwipeteKnight(b *testing.B) {
	engine := NewDefaultEngine()
	ctx := benchKiwipete()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.GetLegalMoves(core.D5, ctx)
	}
}

func BenchmarkGetLegalMoves_KiwipeteQueen(b *testing.B) {
	engine := NewDefaultEngine()
	ctx := benchKiwipete()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.GetLegalMoves(core.F3, ctx)
	}
}

func BenchmarkGetLegalMoves_KiwipeteKing(b *testing.B) {
	engine := NewDefaultEngine()
	ctx := benchKiwipete()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.GetLegalMoves(core.E1, ctx)
	}
}

// =============================================================================
// HasAnyLegalMoves — multiple positions
// =============================================================================

func BenchmarkHasAnyLegalMoves_Start(b *testing.B) {
	engine := NewDefaultEngine()
	ctx := benchStartPos()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.HasAnyLegalMoves(ctx)
	}
}

func BenchmarkHasAnyLegalMoves_FoolsMate(b *testing.B) {
	engine := NewDefaultEngine()
	ctx := benchFoolsMate()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.HasAnyLegalMoves(ctx)
	}
}

func BenchmarkHasAnyLegalMoves_Stalemate(b *testing.B) {
	engine := NewDefaultEngine()
	ctx := benchStalemate()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.HasAnyLegalMoves(ctx)
	}
}

func BenchmarkHasAnyLegalMoves_Kiwipete(b *testing.B) {
	engine := NewDefaultEngine()
	ctx := benchKiwipete()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.HasAnyLegalMoves(ctx)
	}
}

func BenchmarkHasAnyLegalMoves_Endgame(b *testing.B) {
	engine := NewDefaultEngine()
	ctx := benchEndgame()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.HasAnyLegalMoves(ctx)
	}
}

// =============================================================================
// IsSquareAttacked — multiple scenarios
// =============================================================================

func BenchmarkIsSquareAttacked_Empty(b *testing.B) {
	engine := NewDefaultEngine()
	ctx := benchEmptyBoard().BoardContext
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.IsSquareAttacked(core.E4, core.BLACK, ctx)
	}
}

func BenchmarkIsSquareAttacked_Start(b *testing.B) {
	engine := NewDefaultEngine()
	ctx := benchStartPos().BoardContext
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.IsSquareAttacked(core.E1, core.BLACK, ctx)
	}
}

func BenchmarkIsSquareAttacked_Kiwipete(b *testing.B) {
	engine := NewDefaultEngine()
	ctx := benchKiwipete().BoardContext
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.IsSquareAttacked(core.E1, core.BLACK, ctx)
	}
}

func BenchmarkIsSquareAttacked_Attacked(b *testing.B) {
	engine := NewDefaultEngine()
	ctx := benchAttackedSquare().BoardContext
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.IsSquareAttacked(core.E4, core.BLACK, ctx)
	}
}

func BenchmarkIsSquareAttacked_Corner(b *testing.B) {
	engine := NewDefaultEngine()
	ctx := benchStartPos().BoardContext
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.IsSquareAttacked(core.A1, core.BLACK, ctx)
	}
}

// =============================================================================
// Apply — per move type
// =============================================================================

func BenchmarkApply_NormalPawnPush(b *testing.B) {
	engine := NewDefaultEngine()
	base := benchStartPos()
	move := core.Move{
		Type:  core.NORMAL,
		Piece: core.Piece{Type: core.PAWN, Color: core.WHITE},
		From:  core.E2,
		To:    core.E4,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := base.Copy()
		_ = engine.Apply(&ctx, move)
	}
}

func BenchmarkApply_NormalKnightMove(b *testing.B) {
	engine := NewDefaultEngine()
	base := benchStartPos()
	move := core.Move{
		Type:  core.NORMAL,
		Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE},
		From:  core.B1,
		To:    core.C3,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := base.Copy()
		_ = engine.Apply(&ctx, move)
	}
}

func BenchmarkApply_Capture(b *testing.B) {
	engine := NewDefaultEngine()
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
	move := core.Move{
		Type:       core.NORMAL,
		Piece:      core.Piece{Type: core.PAWN, Color: core.WHITE},
		From:       core.E4,
		To:         core.D5,
		HasCapture: true,
		Captured:   core.Piece{Type: core.PAWN, Color: core.BLACK},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := base.Copy()
		_ = engine.Apply(&ctx, move)
	}
}

func BenchmarkApply_EnPassant(b *testing.B) {
	engine := NewDefaultEngine()
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
	move := core.Move{
		Type:       core.EN_PASSANT,
		Piece:      core.Piece{Type: core.PAWN, Color: core.WHITE},
		From:       core.D5,
		To:         core.E6,
		HasCapture: true,
		Captured:   core.Piece{Type: core.PAWN, Color: core.BLACK},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := base.Copy()
		_ = engine.Apply(&ctx, move)
	}
}

func BenchmarkApply_CastlingKingSide(b *testing.B) {
	engine := NewDefaultEngine()
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
	move := core.Move{
		Type:  core.CASTLING,
		Piece: core.Piece{Type: core.KING, Color: core.WHITE},
		From:  core.E1,
		To:    core.G1,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := base.Copy()
		_ = engine.Apply(&ctx, move)
	}
}

func BenchmarkApply_CastlingQueenSide(b *testing.B) {
	engine := NewDefaultEngine()
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
	move := core.Move{
		Type:  core.CASTLING,
		Piece: core.Piece{Type: core.KING, Color: core.WHITE},
		From:  core.E1,
		To:    core.C1,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := base.Copy()
		_ = engine.Apply(&ctx, move)
	}
}

func BenchmarkApply_Promotion(b *testing.B) {
	engine := NewDefaultEngine()
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
	move := core.Move{
		Type:      core.PROMOTION,
		Piece:     core.Piece{Type: core.PAWN, Color: core.WHITE},
		From:      core.E7,
		To:        core.E8,
		PromoteTo: core.QUEEN,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := base.Copy()
		_ = engine.Apply(&ctx, move)
	}
}

func BenchmarkApply_DoublePawnPush(b *testing.B) {
	engine := NewDefaultEngine()
	base := benchStartPos()
	move := core.Move{
		Type:  core.NORMAL,
		Piece: core.Piece{Type: core.PAWN, Color: core.WHITE},
		From:  core.E2,
		To:    core.E4,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := base.Copy()
		_ = engine.Apply(&ctx, move)
	}
}

// =============================================================================
// Undo — per move type (apply then undo, measure undo cost)
// =============================================================================

func BenchmarkUndo_NormalPawnPush(b *testing.B) {
	engine := NewDefaultEngine()
	base := benchStartPos()
	move := core.Move{
		Type:  core.NORMAL,
		Piece: core.Piece{Type: core.PAWN, Color: core.WHITE},
		From:  core.E2,
		To:    core.E4,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := base.Copy()
		snap := engine.Apply(&ctx, move)
		engine.Undo(&ctx, snap)
	}
}

func BenchmarkUndo_NormalKnightMove(b *testing.B) {
	engine := NewDefaultEngine()
	base := benchStartPos()
	move := core.Move{
		Type:  core.NORMAL,
		Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE},
		From:  core.B1,
		To:    core.C3,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := base.Copy()
		snap := engine.Apply(&ctx, move)
		engine.Undo(&ctx, snap)
	}
}

func BenchmarkUndo_Capture(b *testing.B) {
	engine := NewDefaultEngine()
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
	move := core.Move{
		Type:       core.NORMAL,
		Piece:      core.Piece{Type: core.PAWN, Color: core.WHITE},
		From:       core.E4,
		To:         core.D5,
		HasCapture: true,
		Captured:   core.Piece{Type: core.PAWN, Color: core.BLACK},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := base.Copy()
		snap := engine.Apply(&ctx, move)
		engine.Undo(&ctx, snap)
	}
}

func BenchmarkUndo_EnPassant(b *testing.B) {
	engine := NewDefaultEngine()
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
	move := core.Move{
		Type:       core.EN_PASSANT,
		Piece:      core.Piece{Type: core.PAWN, Color: core.WHITE},
		From:       core.D5,
		To:         core.E6,
		HasCapture: true,
		Captured:   core.Piece{Type: core.PAWN, Color: core.BLACK},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := base.Copy()
		snap := engine.Apply(&ctx, move)
		engine.Undo(&ctx, snap)
	}
}

func BenchmarkUndo_CastlingKingSide(b *testing.B) {
	engine := NewDefaultEngine()
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
	move := core.Move{
		Type:  core.CASTLING,
		Piece: core.Piece{Type: core.KING, Color: core.WHITE},
		From:  core.E1,
		To:    core.G1,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := base.Copy()
		snap := engine.Apply(&ctx, move)
		engine.Undo(&ctx, snap)
	}
}

func BenchmarkUndo_Promotion(b *testing.B) {
	engine := NewDefaultEngine()
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
	move := core.Move{
		Type:      core.PROMOTION,
		Piece:     core.Piece{Type: core.PAWN, Color: core.WHITE},
		From:      core.E7,
		To:        core.E8,
		PromoteTo: core.QUEEN,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := base.Copy()
		snap := engine.Apply(&ctx, move)
		engine.Undo(&ctx, snap)
	}
}

// =============================================================================
// Apply + Undo round-trip (the real hot-path cost)
// =============================================================================

func BenchmarkApplyUndo_PawnPush(b *testing.B) {
	engine := NewDefaultEngine()
	base := benchStartPos()
	move := core.Move{
		Type:  core.NORMAL,
		Piece: core.Piece{Type: core.PAWN, Color: core.WHITE},
		From:  core.E2,
		To:    core.E4,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := base.Copy()
		snap := engine.Apply(&ctx, move)
		engine.Undo(&ctx, snap)
	}
}

func BenchmarkApplyUndo_KnightMove(b *testing.B) {
	engine := NewDefaultEngine()
	base := benchStartPos()
	move := core.Move{
		Type:  core.NORMAL,
		Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE},
		From:  core.B1,
		To:    core.C3,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := base.Copy()
		snap := engine.Apply(&ctx, move)
		engine.Undo(&ctx, snap)
	}
}

func BenchmarkApplyUndo_Castling(b *testing.B) {
	engine := NewDefaultEngine()
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
	move := core.Move{
		Type:  core.CASTLING,
		Piece: core.Piece{Type: core.KING, Color: core.WHITE},
		From:  core.E1,
		To:    core.G1,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := base.Copy()
		snap := engine.Apply(&ctx, move)
		engine.Undo(&ctx, snap)
	}
}
