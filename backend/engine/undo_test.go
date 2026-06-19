package engine

import (
	"testing"

	"github.com/7-Dany/chess/core"
)

// TestUndo exercises Undo in isolation. The board is set up as it would
// be AFTER Apply ran, the snapshot is built manually with the pre-move
// state, and Undo is called directly without ever calling Apply.
// This isolates Undo bugs from Apply bugs.
func TestUndo(t *testing.T) {
	defaultSides := [2]core.SideState{
		{KingPosition: core.E1, CastlingRights: core.CastlingRights{KingSide: true, QueenSide: true}},
		{KingPosition: core.E8, CastlingRights: core.CastlingRights{KingSide: true, QueenSide: true}},
	}

	tests := []struct {
		name        string
		setupBoard  func(*core.Board) // post-move board state
		sideToMove  core.PieceColor   // post-move SideToMove (Undo doesn't touch this)
		inputEP     core.Position     // post-move EP target (Undo doesn't touch this)
		sides       [2]core.SideState // post-move sides
		snapshot    core.Snapshot     // carries pre-move state
		expectAt    map[core.Position]core.Square
		expectEmpty []core.Position
		expectEP    core.Position
		expectSides [2]core.SideState
	}{
		// ==================== Normal moves ====================
		{
			name: "normal knight move undone",
			setupBoard: func(b *core.Board) {
				b[core.C3] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
			},
			sideToMove: core.BLACK, // post-move side (Undo doesn't touch it)
			inputEP:    core.NoPosition,
			sides:      defaultSides,
			snapshot: core.Snapshot{
				Move: core.Move{
					Type: core.NORMAL, Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE},
					From: core.B1, To: core.C3,
				},
				PreviousSides:           defaultSides,
				PreviousEnPassantTarget: core.NoPosition,
			},
			expectAt: map[core.Position]core.Square{
				core.B1: core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE}),
			},
			expectEmpty: []core.Position{core.C3},
			expectEP:    core.NoPosition,
			expectSides: defaultSides,
		},
		{
			name: "normal king move undone restores rights",
			setupBoard: func(b *core.Board) {
				b[core.F1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
			},
			sideToMove: core.BLACK,
			inputEP:    core.NoPosition,
			sides: [2]core.SideState{
				{KingPosition: core.F1, CastlingRights: core.CastlingRights{}},
				defaultSides[1],
			},
			snapshot: core.Snapshot{
				Move: core.Move{
					Type: core.NORMAL, Piece: core.Piece{Type: core.KING, Color: core.WHITE},
					From: core.E1, To: core.F1,
				},
				PreviousSides:           defaultSides,
				PreviousEnPassantTarget: core.NoPosition,
			},
			expectAt: map[core.Position]core.Square{
				core.E1: core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE}),
			},
			expectEmpty: []core.Position{core.F1},
			expectEP:    core.NoPosition,
			expectSides: defaultSides,
		},
		{
			name: "rook move from A1 undone restores queen-side right",
			setupBoard: func(b *core.Board) {
				b[core.A3] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			sideToMove: core.BLACK,
			inputEP:    core.NoPosition,
			sides: [2]core.SideState{
				{KingPosition: core.E1, CastlingRights: core.CastlingRights{KingSide: true}},
				defaultSides[1],
			},
			snapshot: core.Snapshot{
				Move: core.Move{
					Type: core.NORMAL, Piece: core.Piece{Type: core.ROOK, Color: core.WHITE},
					From: core.A1, To: core.A3,
				},
				PreviousSides:           defaultSides,
				PreviousEnPassantTarget: core.NoPosition,
			},
			expectAt: map[core.Position]core.Square{
				core.A1: core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE}),
			},
			expectEmpty: []core.Position{core.A3},
			expectEP:    core.NoPosition,
			expectSides: defaultSides,
		},
		{
			name: "black rook move from H8 undone restores king-side right",
			setupBoard: func(b *core.Board) {
				b[core.H6] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			inputEP:    core.NoPosition,
			sides: [2]core.SideState{
				defaultSides[0],
				{KingPosition: core.E8, CastlingRights: core.CastlingRights{QueenSide: true}},
			},
			snapshot: core.Snapshot{
				Move: core.Move{
					Type: core.NORMAL, Piece: core.Piece{Type: core.ROOK, Color: core.BLACK},
					From: core.H8, To: core.H6,
				},
				PreviousSides:           defaultSides,
				PreviousEnPassantTarget: core.NoPosition,
			},
			expectAt: map[core.Position]core.Square{
				core.H8: core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK}),
			},
			expectEmpty: []core.Position{core.H6},
			expectEP:    core.NoPosition,
			expectSides: defaultSides,
		},

		// ==================== Captures ====================
		{
			name: "capture undone restores both pieces",
			setupBoard: func(b *core.Board) {
				b[core.D5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
			},
			sideToMove: core.BLACK,
			inputEP:    core.NoPosition,
			sides:      defaultSides,
			snapshot: core.Snapshot{
				Move: core.Move{
					Type: core.NORMAL, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE},
					From: core.E4, To: core.D5,
					HasCapture: true, Captured: core.Piece{Type: core.PAWN, Color: core.BLACK},
				},
				PreviousSides:           defaultSides,
				PreviousEnPassantTarget: core.NoPosition,
			},
			expectAt: map[core.Position]core.Square{
				core.E4: core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE}),
				core.D5: core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK}),
			},
			expectEP:    core.NoPosition,
			expectSides: defaultSides,
		},
		{
			name: "capturing rook on A8 undone restores rook and rights",
			setupBoard: func(b *core.Board) {
				b[core.A8] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
			},
			sideToMove: core.BLACK,
			inputEP:    core.NoPosition,
			sides: [2]core.SideState{
				defaultSides[0],
				{KingPosition: core.E8, CastlingRights: core.CastlingRights{KingSide: true}},
			},
			snapshot: core.Snapshot{
				Move: core.Move{
					Type: core.NORMAL, Piece: core.Piece{Type: core.BISHOP, Color: core.WHITE},
					From: core.A6, To: core.A8,
					HasCapture: true, Captured: core.Piece{Type: core.ROOK, Color: core.BLACK},
				},
				PreviousSides:           defaultSides,
				PreviousEnPassantTarget: core.NoPosition,
			},
			expectAt: map[core.Position]core.Square{
				core.A6: core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE}),
				core.A8: core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK}),
			},
			expectEP:    core.NoPosition,
			expectSides: defaultSides,
		},
		{
			name: "capturing non-rook on A file undone restores pawn",
			setupBoard: func(b *core.Board) {
				b[core.A6] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
			},
			sideToMove: core.BLACK,
			inputEP:    core.NoPosition,
			sides:      defaultSides,
			snapshot: core.Snapshot{
				Move: core.Move{
					Type: core.NORMAL, Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE},
					From: core.B5, To: core.A6,
					HasCapture: true, Captured: core.Piece{Type: core.PAWN, Color: core.BLACK},
				},
				PreviousSides:           defaultSides,
				PreviousEnPassantTarget: core.NoPosition,
			},
			expectAt: map[core.Position]core.Square{
				core.B5: core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE}),
				core.A6: core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK}),
			},
			expectEP:    core.NoPosition,
			expectSides: defaultSides,
		},

		// ==================== En Passant ====================
		{
			name: "white en passant undone",
			setupBoard: func(b *core.Board) {
				b[core.E6] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
			},
			sideToMove: core.BLACK,
			inputEP:    core.NoPosition,
			sides:      defaultSides,
			snapshot: core.Snapshot{
				Move: core.Move{
					Type: core.EN_PASSANT, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE},
					From: core.D5, To: core.E6,
					HasCapture: true, Captured: core.Piece{Type: core.PAWN, Color: core.BLACK},
				},
				PreviousSides:           defaultSides,
				PreviousEnPassantTarget: core.E6,
			},
			expectAt: map[core.Position]core.Square{
				core.D5: core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE}),
				core.E5: core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK}),
			},
			expectEmpty: []core.Position{core.E6},
			expectEP:    core.E6,
			expectSides: defaultSides,
		},
		{
			name: "black en passant undone",
			setupBoard: func(b *core.Board) {
				b[core.E3] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			inputEP:    core.NoPosition,
			sides:      defaultSides,
			snapshot: core.Snapshot{
				Move: core.Move{
					Type: core.EN_PASSANT, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK},
					From: core.D4, To: core.E3,
					HasCapture: true, Captured: core.Piece{Type: core.PAWN, Color: core.WHITE},
				},
				PreviousSides:           defaultSides,
				PreviousEnPassantTarget: core.E3,
			},
			expectAt: map[core.Position]core.Square{
				core.D4: core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK}),
				core.E4: core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE}),
			},
			expectEmpty: []core.Position{core.E3},
			expectEP:    core.E3,
			expectSides: defaultSides,
		},
		{
			name: "en passant on A file undone restores pawn not rights",
			setupBoard: func(b *core.Board) {
				b[core.A6] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
			},
			sideToMove: core.BLACK,
			inputEP:    core.NoPosition,
			sides:      defaultSides,
			snapshot: core.Snapshot{
				Move: core.Move{
					Type: core.EN_PASSANT, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE},
					From: core.B5, To: core.A6,
					HasCapture: true, Captured: core.Piece{Type: core.PAWN, Color: core.BLACK},
				},
				PreviousSides:           defaultSides,
				PreviousEnPassantTarget: core.A6,
			},
			expectAt: map[core.Position]core.Square{
				core.B5: core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE}),
				core.A5: core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK}),
			},
			expectEmpty: []core.Position{core.A6},
			expectEP:    core.A6,
			expectSides: defaultSides,
		},

		// ==================== Promotion ====================
		{
			name: "promotion to queen undone restores pawn",
			setupBoard: func(b *core.Board) {
				b[core.E8] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.WHITE})
			},
			sideToMove: core.BLACK,
			inputEP:    core.NoPosition,
			sides:      defaultSides,
			snapshot: core.Snapshot{
				Move: core.Move{
					Type: core.PROMOTION, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE},
					From: core.E7, To: core.E8,
					PromoteTo: core.QUEEN,
				},
				PreviousSides:           defaultSides,
				PreviousEnPassantTarget: core.NoPosition,
			},
			expectAt: map[core.Position]core.Square{
				core.E7: core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE}),
			},
			expectEmpty: []core.Position{core.E8},
			expectEP:    core.NoPosition,
			expectSides: defaultSides,
		},
		{
			name: "promotion to knight undone restores pawn",
			setupBoard: func(b *core.Board) {
				b[core.D1] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			inputEP:    core.NoPosition,
			sides:      defaultSides,
			snapshot: core.Snapshot{
				Move: core.Move{
					Type: core.PROMOTION, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK},
					From: core.D2, To: core.D1,
					PromoteTo: core.KNIGHT,
				},
				PreviousSides:           defaultSides,
				PreviousEnPassantTarget: core.NoPosition,
			},
			expectAt: map[core.Position]core.Square{
				core.D2: core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK}),
			},
			expectEmpty: []core.Position{core.D1},
			expectEP:    core.NoPosition,
			expectSides: defaultSides,
		},
		{
			name: "promotion with capture undone restores pawn and captured",
			setupBoard: func(b *core.Board) {
				b[core.H8] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.WHITE})
			},
			sideToMove: core.BLACK,
			inputEP:    core.NoPosition,
			sides: [2]core.SideState{
				defaultSides[0],
				{KingPosition: core.E8, CastlingRights: core.CastlingRights{QueenSide: true}},
			},
			snapshot: core.Snapshot{
				Move: core.Move{
					Type: core.PROMOTION, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE},
					From: core.G7, To: core.H8,
					PromoteTo:  core.QUEEN,
					HasCapture: true, Captured: core.Piece{Type: core.ROOK, Color: core.BLACK},
				},
				PreviousSides:           defaultSides,
				PreviousEnPassantTarget: core.NoPosition,
			},
			expectAt: map[core.Position]core.Square{
				core.G7: core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE}),
				core.H8: core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK}),
			},
			expectEP:    core.NoPosition,
			expectSides: defaultSides,
		},

		// ==================== Castling ====================
		{
			name: "white king-side castling undone",
			setupBoard: func(b *core.Board) {
				b[core.G1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
				b[core.F1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			sideToMove: core.BLACK,
			inputEP:    core.NoPosition,
			sides: [2]core.SideState{
				{KingPosition: core.G1, CastlingRights: core.CastlingRights{}},
				defaultSides[1],
			},
			snapshot: core.Snapshot{
				Move: core.Move{
					Type: core.CASTLING, Piece: core.Piece{Type: core.KING, Color: core.WHITE},
					From: core.E1, To: core.G1,
				},
				PreviousSides:           defaultSides,
				PreviousEnPassantTarget: core.NoPosition,
			},
			expectAt: map[core.Position]core.Square{
				core.E1: core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE}),
				core.H1: core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE}),
			},
			expectEmpty: []core.Position{core.G1, core.F1},
			expectEP:    core.NoPosition,
			expectSides: defaultSides,
		},
		{
			name: "white queen-side castling undone",
			setupBoard: func(b *core.Board) {
				b[core.C1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
				b[core.D1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			sideToMove: core.BLACK,
			inputEP:    core.NoPosition,
			sides: [2]core.SideState{
				{KingPosition: core.C1, CastlingRights: core.CastlingRights{}},
				defaultSides[1],
			},
			snapshot: core.Snapshot{
				Move: core.Move{
					Type: core.CASTLING, Piece: core.Piece{Type: core.KING, Color: core.WHITE},
					From: core.E1, To: core.C1,
				},
				PreviousSides:           defaultSides,
				PreviousEnPassantTarget: core.NoPosition,
			},
			expectAt: map[core.Position]core.Square{
				core.E1: core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE}),
				core.A1: core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE}),
			},
			expectEmpty: []core.Position{core.C1, core.D1},
			expectEP:    core.NoPosition,
			expectSides: defaultSides,
		},
		{
			name: "black king-side castling undone",
			setupBoard: func(b *core.Board) {
				b[core.G8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
				b[core.F8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			inputEP:    core.NoPosition,
			sides: [2]core.SideState{
				defaultSides[0],
				{KingPosition: core.G8, CastlingRights: core.CastlingRights{}},
			},
			snapshot: core.Snapshot{
				Move: core.Move{
					Type: core.CASTLING, Piece: core.Piece{Type: core.KING, Color: core.BLACK},
					From: core.E8, To: core.G8,
				},
				PreviousSides:           defaultSides,
				PreviousEnPassantTarget: core.NoPosition,
			},
			expectAt: map[core.Position]core.Square{
				core.E8: core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK}),
				core.H8: core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK}),
			},
			expectEmpty: []core.Position{core.G8, core.F8},
			expectEP:    core.NoPosition,
			expectSides: defaultSides,
		},
		{
			name: "black queen-side castling undone",
			setupBoard: func(b *core.Board) {
				b[core.C8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
				b[core.D8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			inputEP:    core.NoPosition,
			sides: [2]core.SideState{
				defaultSides[0],
				{KingPosition: core.C8, CastlingRights: core.CastlingRights{}},
			},
			snapshot: core.Snapshot{
				Move: core.Move{
					Type: core.CASTLING, Piece: core.Piece{Type: core.KING, Color: core.BLACK},
					From: core.E8, To: core.C8,
				},
				PreviousSides:           defaultSides,
				PreviousEnPassantTarget: core.NoPosition,
			},
			expectAt: map[core.Position]core.Square{
				core.E8: core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK}),
				core.A8: core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK}),
			},
			expectEmpty: []core.Position{core.C8, core.D8},
			expectEP:    core.NoPosition,
			expectSides: defaultSides,
		},

		// ==================== Double pawn push / EP target restoration ====================
		{
			name: "double pawn push undone restores no-EP state",
			setupBoard: func(b *core.Board) {
				b[core.E4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
			},
			sideToMove: core.BLACK,
			inputEP:    core.E3, // EP target set after the double push
			sides:      defaultSides,
			snapshot: core.Snapshot{
				Move: core.Move{
					Type: core.NORMAL, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE},
					From: core.E2, To: core.E4,
				},
				PreviousSides:           defaultSides,
				PreviousEnPassantTarget: core.NoPosition, // EP was NoPosition before this move
			},
			expectAt: map[core.Position]core.Square{
				core.E2: core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE}),
			},
			expectEmpty: []core.Position{core.E4},
			expectEP:    core.NoPosition,
			expectSides: defaultSides,
		},
		{
			name: "non-pawn move undone restores previous EP target",
			setupBoard: func(b *core.Board) {
				b[core.C3] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
			},
			sideToMove: core.BLACK,
			inputEP:    core.NoPosition, // EP was cleared by this knight move
			sides:      defaultSides,
			snapshot: core.Snapshot{
				Move: core.Move{
					Type: core.NORMAL, Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE},
					From: core.B1, To: core.C3,
				},
				PreviousSides:           defaultSides,
				PreviousEnPassantTarget: core.E3, // EP existed before this knight move
			},
			expectAt: map[core.Position]core.Square{
				core.B1: core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE}),
			},
			expectEmpty: []core.Position{core.C3},
			expectEP:    core.E3,
			expectSides: defaultSides,
		},
	}

	engine := NewDefaultEngine()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var board core.Board
			tt.setupBoard(&board)

			ctx := &core.TurnContext{
				MoveContext: core.MoveContext{
					BoardContext:    core.BoardContext{Board: &board},
					SideToMove:      tt.sideToMove,
					Sides:           tt.sides,
					EnPassantTarget: tt.inputEP,
				},
			}

			engine.Undo(ctx, tt.snapshot)

			// Undo must not touch SideToMove.
			if ctx.SideToMove != tt.sideToMove {
				t.Errorf("SideToMove = %v, want %v (Undo should not touch side)", ctx.SideToMove, tt.sideToMove)
			}

			for pos, expected := range tt.expectAt {
				got := ctx.Board[pos]
				if got != expected {
					t.Errorf("board[%v] = %v, want %v", pos, got, expected)
				}
			}

			for _, pos := range tt.expectEmpty {
				if ctx.Board[pos].IsOccupied() {
					t.Errorf("board[%v] should be empty, got %v", pos, ctx.Board[pos])
				}
			}

			if ctx.EnPassantTarget != tt.expectEP {
				t.Errorf("EnPassantTarget = %v, want %v", ctx.EnPassantTarget, tt.expectEP)
			}

			if ctx.Sides != tt.expectSides {
				t.Errorf("Sides = %+v, want %+v", ctx.Sides, tt.expectSides)
			}
		})
	}
}

// TestApplyThenUndo verifies that Apply followed by Undo returns the context
// to its original state. This catches asymmetry bugs between Apply and Undo
// that the pure Undo test would miss.
func TestApplyThenUndo(t *testing.T) {
	defaultSides := [2]core.SideState{
		{KingPosition: core.E1, CastlingRights: core.CastlingRights{KingSide: true, QueenSide: true}},
		{KingPosition: core.E8, CastlingRights: core.CastlingRights{KingSide: true, QueenSide: true}},
	}

	tests := []struct {
		name       string
		setupBoard func(*core.Board)
		sideToMove core.PieceColor
		inputEP    core.Position
		sides      [2]core.SideState
		move       core.Move
	}{
		// ==================== Normal moves ====================
		{
			name: "normal knight move round trip",
			setupBoard: func(b *core.Board) {
				b[core.B1] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			move: core.Move{
				Type: core.NORMAL, Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE},
				From: core.B1, To: core.C3,
			},
		},
		{
			name: "normal king move round trip",
			setupBoard: func(b *core.Board) {
				b[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			move: core.Move{
				Type: core.NORMAL, Piece: core.Piece{Type: core.KING, Color: core.WHITE},
				From: core.E1, To: core.F1,
			},
		},
		{
			name: "rook move from A1 round trip",
			setupBoard: func(b *core.Board) {
				b[core.A1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			move: core.Move{
				Type: core.NORMAL, Piece: core.Piece{Type: core.ROOK, Color: core.WHITE},
				From: core.A1, To: core.A3,
			},
		},
		{
			name: "rook move from non-home file round trip",
			setupBoard: func(b *core.Board) {
				b[core.C3] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			move: core.Move{
				Type: core.NORMAL, Piece: core.Piece{Type: core.ROOK, Color: core.WHITE},
				From: core.C3, To: core.C5,
			},
		},
		{
			name: "black rook move from H8 round trip",
			setupBoard: func(b *core.Board) {
				b[core.H8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
			},
			sideToMove: core.BLACK,
			sides:      defaultSides,
			move: core.Move{
				Type: core.NORMAL, Piece: core.Piece{Type: core.ROOK, Color: core.BLACK},
				From: core.H8, To: core.H6,
			},
		},
		{
			name: "black rook move from A8 round trip",
			setupBoard: func(b *core.Board) {
				b[core.A8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
			},
			sideToMove: core.BLACK,
			sides:      defaultSides,
			move: core.Move{
				Type: core.NORMAL, Piece: core.Piece{Type: core.ROOK, Color: core.BLACK},
				From: core.A8, To: core.A6,
			},
		},

		// ==================== Captures ====================
		{
			name: "capture round trip",
			setupBoard: func(b *core.Board) {
				b[core.E4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
				b[core.D5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			move: core.Move{
				Type: core.NORMAL, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE},
				From: core.E4, To: core.D5,
				HasCapture: true, Captured: core.Piece{Type: core.PAWN, Color: core.BLACK},
			},
		},
		{
			name: "capture rook on A8 round trip",
			setupBoard: func(b *core.Board) {
				b[core.A6] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
				b[core.A8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			move: core.Move{
				Type: core.NORMAL, Piece: core.Piece{Type: core.BISHOP, Color: core.WHITE},
				From: core.A6, To: core.A8,
				HasCapture: true, Captured: core.Piece{Type: core.ROOK, Color: core.BLACK},
			},
		},
		{
			name: "capture rook on H1 round trip",
			setupBoard: func(b *core.Board) {
				b[core.H3] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.BLACK})
				b[core.H1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			sideToMove: core.BLACK,
			sides:      defaultSides,
			move: core.Move{
				Type: core.NORMAL, Piece: core.Piece{Type: core.BISHOP, Color: core.BLACK},
				From: core.H3, To: core.H1,
				HasCapture: true, Captured: core.Piece{Type: core.ROOK, Color: core.WHITE},
			},
		},
		{
			name: "capture non-rook on A file round trip",
			setupBoard: func(b *core.Board) {
				b[core.B5] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
				b[core.A6] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			move: core.Move{
				Type: core.NORMAL, Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE},
				From: core.B5, To: core.A6,
				HasCapture: true, Captured: core.Piece{Type: core.PAWN, Color: core.BLACK},
			},
		},

		// ==================== En Passant ====================
		{
			name: "white en passant round trip",
			setupBoard: func(b *core.Board) {
				b[core.D5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
				b[core.E5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			inputEP:    core.E6,
			sides:      defaultSides,
			move: core.Move{
				Type: core.EN_PASSANT, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE},
				From: core.D5, To: core.E6,
				HasCapture: true, Captured: core.Piece{Type: core.PAWN, Color: core.BLACK},
			},
		},
		{
			name: "black en passant round trip",
			setupBoard: func(b *core.Board) {
				b[core.E4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
				b[core.D4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
			},
			sideToMove: core.BLACK,
			inputEP:    core.E3,
			sides:      defaultSides,
			move: core.Move{
				Type: core.EN_PASSANT, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK},
				From: core.D4, To: core.E3,
				HasCapture: true, Captured: core.Piece{Type: core.PAWN, Color: core.WHITE},
			},
		},
		{
			name: "en passant on A file round trip",
			setupBoard: func(b *core.Board) {
				b[core.B5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
				b[core.A5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			inputEP:    core.A6,
			sides:      defaultSides,
			move: core.Move{
				Type: core.EN_PASSANT, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE},
				From: core.B5, To: core.A6,
				HasCapture: true, Captured: core.Piece{Type: core.PAWN, Color: core.BLACK},
			},
		},

		// ==================== Promotion ====================
		{
			name: "promotion to queen round trip",
			setupBoard: func(b *core.Board) {
				b[core.E7] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			move: core.Move{
				Type: core.PROMOTION, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE},
				From: core.E7, To: core.E8,
				PromoteTo: core.QUEEN,
			},
		},
		{
			name: "promotion to knight round trip",
			setupBoard: func(b *core.Board) {
				b[core.D2] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
			},
			sideToMove: core.BLACK,
			sides:      defaultSides,
			move: core.Move{
				Type: core.PROMOTION, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK},
				From: core.D2, To: core.D1,
				PromoteTo: core.KNIGHT,
			},
		},
		{
			name: "promotion with capture round trip",
			setupBoard: func(b *core.Board) {
				b[core.G7] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
				b[core.H8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			move: core.Move{
				Type: core.PROMOTION, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE},
				From: core.G7, To: core.H8,
				PromoteTo:  core.QUEEN,
				HasCapture: true, Captured: core.Piece{Type: core.ROOK, Color: core.BLACK},
			},
		},

		// ==================== Castling ====================
		{
			name: "white king-side castling round trip",
			setupBoard: func(b *core.Board) {
				b[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
				b[core.H1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			move: core.Move{
				Type: core.CASTLING, Piece: core.Piece{Type: core.KING, Color: core.WHITE},
				From: core.E1, To: core.G1,
			},
		},
		{
			name: "white queen-side castling round trip",
			setupBoard: func(b *core.Board) {
				b[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
				b[core.A1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			move: core.Move{
				Type: core.CASTLING, Piece: core.Piece{Type: core.KING, Color: core.WHITE},
				From: core.E1, To: core.C1,
			},
		},
		{
			name: "black king-side castling round trip",
			setupBoard: func(b *core.Board) {
				b[core.E8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
				b[core.H8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
			},
			sideToMove: core.BLACK,
			sides:      defaultSides,
			move: core.Move{
				Type: core.CASTLING, Piece: core.Piece{Type: core.KING, Color: core.BLACK},
				From: core.E8, To: core.G8,
			},
		},
		{
			name: "black queen-side castling round trip",
			setupBoard: func(b *core.Board) {
				b[core.E8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
				b[core.A8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
			},
			sideToMove: core.BLACK,
			sides:      defaultSides,
			move: core.Move{
				Type: core.CASTLING, Piece: core.Piece{Type: core.KING, Color: core.BLACK},
				From: core.E8, To: core.C8,
			},
		},

		// ==================== Double pawn push ====================
		{
			name: "white double pawn push round trip",
			setupBoard: func(b *core.Board) {
				b[core.E2] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			move: core.Move{
				Type: core.NORMAL, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE},
				From: core.E2, To: core.E4,
			},
		},
		{
			name: "black double pawn push round trip",
			setupBoard: func(b *core.Board) {
				b[core.D7] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
			},
			sideToMove: core.BLACK,
			sides:      defaultSides,
			move: core.Move{
				Type: core.NORMAL, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK},
				From: core.D7, To: core.D5,
			},
		},
		{
			name: "white double pawn push from A file round trip",
			setupBoard: func(b *core.Board) {
				b[core.A2] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			move: core.Move{
				Type: core.NORMAL, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE},
				From: core.A2, To: core.A4,
			},
		},
		{
			name: "non-pawn move with existing EP target round trip",
			setupBoard: func(b *core.Board) {
				b[core.B1] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
			},
			sideToMove: core.WHITE,
			inputEP:    core.E3,
			sides:      defaultSides,
			move: core.Move{
				Type: core.NORMAL, Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE},
				From: core.B1, To: core.C3,
			},
		},
	}

	engine := NewDefaultEngine()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Build the original context.
			var board core.Board
			tt.setupBoard(&board)

			ctx := &core.TurnContext{
				MoveContext: core.MoveContext{
					BoardContext:    core.BoardContext{Board: &board},
					SideToMove:      tt.sideToMove,
					Sides:           tt.sides,
					EnPassantTarget: tt.inputEP,
				},
			}

			// Snapshot the original state for later comparison.
			originalBoard := *ctx.Board
			originalSideToMove := ctx.SideToMove
			originalSides := ctx.Sides
			originalEP := ctx.EnPassantTarget
			originalHalfMoveClock := ctx.HalfMoveClock
			originalFullMoveNumber := ctx.FullMoveNumber

			// Apply then Undo.
			snap := engine.Apply(ctx, tt.move)
			engine.Undo(ctx, snap)

			// Compare every board square.
			for i := 0; i < 64; i++ {
				pos := core.Position(i)
				got := ctx.Board[pos]
				want := originalBoard[pos]
				if got != want {
					t.Errorf("board[%v] = %v, want %v (after round trip)", pos, got, want)
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
		})
	}
}
