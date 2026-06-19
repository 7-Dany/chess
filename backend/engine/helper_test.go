package engine

import (
	"testing"

	"github.com/7-Dany/chess/core"
)

func TestIsDoublePawnPush(t *testing.T) {
	tests := []struct {
		name     string
		move     core.Move
		expected bool
	}{
		// ==================== Valid double pushes ====================
		{
			name:     "white pawn double push E2 to E4",
			move:     core.Move{Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}, From: core.NewPosition(core.FILE_E, core.RANK_2), To: core.NewPosition(core.FILE_E, core.RANK_4)},
			expected: true,
		},
		{
			name:     "black pawn double push D7 to D5",
			move:     core.Move{Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}, From: core.NewPosition(core.FILE_D, core.RANK_7), To: core.NewPosition(core.FILE_D, core.RANK_5)},
			expected: true,
		},
		{
			name:     "white pawn double push A2 to A4 (edge file)",
			move:     core.Move{Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}, From: core.NewPosition(core.FILE_A, core.RANK_2), To: core.NewPosition(core.FILE_A, core.RANK_4)},
			expected: true,
		},
		{
			name:     "black pawn double push H7 to H5 (edge file)",
			move:     core.Move{Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}, From: core.NewPosition(core.FILE_H, core.RANK_7), To: core.NewPosition(core.FILE_H, core.RANK_5)},
			expected: true,
		},

		// ==================== Single pushes ====================
		{
			name:     "white pawn single push E2 to E3",
			move:     core.Move{Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}, From: core.NewPosition(core.FILE_E, core.RANK_2), To: core.NewPosition(core.FILE_E, core.RANK_3)},
			expected: false,
		},
		{
			name:     "black pawn single push D7 to D6",
			move:     core.Move{Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}, From: core.NewPosition(core.FILE_D, core.RANK_7), To: core.NewPosition(core.FILE_D, core.RANK_6)},
			expected: false,
		},

		// ==================== Non-pawn pieces ====================
		{
			name:     "knight is never a double pawn push",
			move:     core.Move{Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE}, From: core.NewPosition(core.FILE_E, core.RANK_2), To: core.NewPosition(core.FILE_E, core.RANK_4)},
			expected: false,
		},
		{
			name:     "rook is never a double pawn push",
			move:     core.Move{Piece: core.Piece{Type: core.ROOK, Color: core.BLACK}, From: core.NewPosition(core.FILE_A, core.RANK_7), To: core.NewPosition(core.FILE_A, core.RANK_5)},
			expected: false,
		},
		{
			name:     "bishop is never a double pawn push",
			move:     core.Move{Piece: core.Piece{Type: core.BISHOP, Color: core.WHITE}, From: core.NewPosition(core.FILE_C, core.RANK_1), To: core.NewPosition(core.FILE_F, core.RANK_4)},
			expected: false,
		},
		{
			name:     "king is never a double pawn push",
			move:     core.Move{Piece: core.Piece{Type: core.KING, Color: core.WHITE}, From: core.NewPosition(core.FILE_E, core.RANK_1), To: core.NewPosition(core.FILE_E, core.RANK_3)},
			expected: false,
		},

		// ==================== Wrong direction (defensive cases) ====================
		{
			name:     "white pawn moving backward rank 7 to rank 5",
			move:     core.Move{Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}, From: core.NewPosition(core.FILE_E, core.RANK_7), To: core.NewPosition(core.FILE_E, core.RANK_5)},
			expected: false,
		},
		{
			name:     "black pawn moving backward rank 2 to rank 4",
			move:     core.Move{Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}, From: core.NewPosition(core.FILE_E, core.RANK_2), To: core.NewPosition(core.FILE_E, core.RANK_4)},
			expected: false,
		},
		{
			name:     "white pawn moving backward rank 8 to rank 6",
			move:     core.Move{Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}, From: core.NewPosition(core.FILE_E, core.RANK_8), To: core.NewPosition(core.FILE_E, core.RANK_6)},
			expected: false,
		},
		{
			name:     "black pawn moving backward rank 1 to rank 3",
			move:     core.Move{Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}, From: core.NewPosition(core.FILE_E, core.RANK_1), To: core.NewPosition(core.FILE_E, core.RANK_3)},
			expected: false,
		},

		// ==================== Not a 2-rank move ====================
		{
			name:     "white pawn rank 3 to rank 5 (forward 2 but not from start)",
			move:     core.Move{Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}, From: core.NewPosition(core.FILE_E, core.RANK_3), To: core.NewPosition(core.FILE_E, core.RANK_5)},
			expected: true,
		},
		{
			name:     "black pawn rank 6 to rank 4 (forward 2 but not from start)",
			move:     core.Move{Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}, From: core.NewPosition(core.FILE_E, core.RANK_6), To: core.NewPosition(core.FILE_E, core.RANK_4)},
			expected: true,
		},
		{
			name:     "white pawn rank 4 to rank 5 (single push mid-board)",
			move:     core.Move{Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}, From: core.NewPosition(core.FILE_E, core.RANK_4), To: core.NewPosition(core.FILE_E, core.RANK_5)},
			expected: false,
		},
		{
			name:     "black pawn rank 5 to rank 4 (single push mid-board)",
			move:     core.Move{Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}, From: core.NewPosition(core.FILE_E, core.RANK_5), To: core.NewPosition(core.FILE_E, core.RANK_4)},
			expected: false,
		},
		{
			name:     "white pawn 3-rank move (impossible, should be false)",
			move:     core.Move{Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}, From: core.NewPosition(core.FILE_E, core.RANK_2), To: core.NewPosition(core.FILE_E, core.RANK_5)},
			expected: false,
		},
		{
			name:     "black pawn 3-rank move (impossible, should be false)",
			move:     core.Move{Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}, From: core.NewPosition(core.FILE_E, core.RANK_7), To: core.NewPosition(core.FILE_E, core.RANK_4)},
			expected: false,
		},

		// ==================== Edge files ====================
		{
			name:     "white pawn A2 to A4 forward 2",
			move:     core.Move{Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}, From: core.NewPosition(core.FILE_A, core.RANK_2), To: core.NewPosition(core.FILE_A, core.RANK_4)},
			expected: true,
		},
		{
			name:     "black pawn H7 to H5 forward 2",
			move:     core.Move{Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}, From: core.NewPosition(core.FILE_H, core.RANK_7), To: core.NewPosition(core.FILE_H, core.RANK_5)},
			expected: true,
		},

		// ==================== Same square (degenerate) ====================
		{
			name:     "white pawn same square",
			move:     core.Move{Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}, From: core.NewPosition(core.FILE_E, core.RANK_4), To: core.NewPosition(core.FILE_E, core.RANK_4)},
			expected: false,
		},
		{
			name:     "black pawn same square",
			move:     core.Move{Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}, From: core.NewPosition(core.FILE_E, core.RANK_4), To: core.NewPosition(core.FILE_E, core.RANK_4)},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isDoublePawnPush(tt.move)
			if got != tt.expected {
				t.Errorf("isDoublePawnPush() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestEnPassantTarget(t *testing.T) {
	tests := []struct {
		name     string
		move     core.Move
		expected core.Position
	}{
		{
			name:     "white double push E2->E4 target is E3",
			move:     core.Move{Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}, From: core.NewPosition(core.FILE_E, core.RANK_2), To: core.NewPosition(core.FILE_E, core.RANK_4)},
			expected: core.NewPosition(core.FILE_E, core.RANK_3),
		},
		{
			name:     "black double push D7->D5 target is D6",
			move:     core.Move{Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}, From: core.NewPosition(core.FILE_D, core.RANK_7), To: core.NewPosition(core.FILE_D, core.RANK_5)},
			expected: core.NewPosition(core.FILE_D, core.RANK_6),
		},
		{
			name:     "white double push A2->A4 target is A3",
			move:     core.Move{Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}, From: core.NewPosition(core.FILE_A, core.RANK_2), To: core.NewPosition(core.FILE_A, core.RANK_4)},
			expected: core.NewPosition(core.FILE_A, core.RANK_3),
		},
		{
			name:     "black double push H7->H5 target is H6",
			move:     core.Move{Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}, From: core.NewPosition(core.FILE_H, core.RANK_7), To: core.NewPosition(core.FILE_H, core.RANK_5)},
			expected: core.NewPosition(core.FILE_H, core.RANK_6),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := enPassantTarget(tt.move)
			if got != tt.expected {
				t.Errorf("enPassantTarget() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestMoveRook(t *testing.T) {
	tests := []struct {
		name        string
		setupBoard  func(*core.Board)
		rank        core.Rank
		from        core.File
		to          core.File
		expectAt    map[core.Position]core.Square
		expectEmpty []core.Position
	}{
		{
			name: "white king-side rook H1 to F1",
			setupBoard: func(b *core.Board) {
				b[core.NewPosition(core.FILE_H, core.RANK_1)] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			rank: core.RANK_1,
			from: core.FILE_H,
			to:   core.FILE_F,
			expectAt: map[core.Position]core.Square{
				core.NewPosition(core.FILE_F, core.RANK_1): core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE}),
			},
			expectEmpty: []core.Position{core.NewPosition(core.FILE_H, core.RANK_1)},
		},
		{
			name: "white queen-side rook A1 to D1",
			setupBoard: func(b *core.Board) {
				b[core.NewPosition(core.FILE_A, core.RANK_1)] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			rank: core.RANK_1,
			from: core.FILE_A,
			to:   core.FILE_D,
			expectAt: map[core.Position]core.Square{
				core.NewPosition(core.FILE_D, core.RANK_1): core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE}),
			},
			expectEmpty: []core.Position{core.NewPosition(core.FILE_A, core.RANK_1)},
		},
		{
			name: "black king-side rook H8 to F8",
			setupBoard: func(b *core.Board) {
				b[core.NewPosition(core.FILE_H, core.RANK_8)] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
			},
			rank: core.RANK_8,
			from: core.FILE_H,
			to:   core.FILE_F,
			expectAt: map[core.Position]core.Square{
				core.NewPosition(core.FILE_F, core.RANK_8): core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK}),
			},
			expectEmpty: []core.Position{core.NewPosition(core.FILE_H, core.RANK_8)},
		},
		{
			name: "black queen-side rook A8 to D8",
			setupBoard: func(b *core.Board) {
				b[core.NewPosition(core.FILE_A, core.RANK_8)] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
			},
			rank: core.RANK_8,
			from: core.FILE_A,
			to:   core.FILE_D,
			expectAt: map[core.Position]core.Square{
				core.NewPosition(core.FILE_D, core.RANK_8): core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK}),
			},
			expectEmpty: []core.Position{core.NewPosition(core.FILE_A, core.RANK_8)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var board core.Board
			tt.setupBoard(&board)

			ctx := &core.TurnContext{
				MoveContext: core.MoveContext{
					BoardContext: core.BoardContext{Board: &board},
				},
			}

			moveRook(ctx, tt.rank, tt.from, tt.to)

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
		})
	}
}

func TestClearCastlingRightByFile(t *testing.T) {
	tests := []struct {
		name        string
		color       core.PieceColor
		file        core.File
		inputSides  [2]core.SideState
		expectSides [2]core.SideState
	}{
		{
			name:  "FILE_A clears queen-side right for white",
			color: core.WHITE,
			file:  core.FILE_A,
			inputSides: [2]core.SideState{
				{KingPosition: core.E1, CastlingRights: core.CastlingRights{KingSide: true, QueenSide: true}},
				{KingPosition: core.E8, CastlingRights: core.CastlingRights{KingSide: true, QueenSide: true}},
			},
			expectSides: [2]core.SideState{
				{KingPosition: core.E1, CastlingRights: core.CastlingRights{KingSide: true}},
				{KingPosition: core.E8, CastlingRights: core.CastlingRights{KingSide: true, QueenSide: true}},
			},
		},
		{
			name:  "FILE_H clears king-side right for white",
			color: core.WHITE,
			file:  core.FILE_H,
			inputSides: [2]core.SideState{
				{KingPosition: core.E1, CastlingRights: core.CastlingRights{KingSide: true, QueenSide: true}},
				{KingPosition: core.E8, CastlingRights: core.CastlingRights{KingSide: true, QueenSide: true}},
			},
			expectSides: [2]core.SideState{
				{KingPosition: core.E1, CastlingRights: core.CastlingRights{QueenSide: true}},
				{KingPosition: core.E8, CastlingRights: core.CastlingRights{KingSide: true, QueenSide: true}},
			},
		},
		{
			name:  "FILE_A clears queen-side right for black",
			color: core.BLACK,
			file:  core.FILE_A,
			inputSides: [2]core.SideState{
				{KingPosition: core.E1, CastlingRights: core.CastlingRights{KingSide: true, QueenSide: true}},
				{KingPosition: core.E8, CastlingRights: core.CastlingRights{KingSide: true, QueenSide: true}},
			},
			expectSides: [2]core.SideState{
				{KingPosition: core.E1, CastlingRights: core.CastlingRights{KingSide: true, QueenSide: true}},
				{KingPosition: core.E8, CastlingRights: core.CastlingRights{KingSide: true}},
			},
		},
		{
			name:  "FILE_H clears king-side right for black",
			color: core.BLACK,
			file:  core.FILE_H,
			inputSides: [2]core.SideState{
				{KingPosition: core.E1, CastlingRights: core.CastlingRights{KingSide: true, QueenSide: true}},
				{KingPosition: core.E8, CastlingRights: core.CastlingRights{KingSide: true, QueenSide: true}},
			},
			expectSides: [2]core.SideState{
				{KingPosition: core.E1, CastlingRights: core.CastlingRights{KingSide: true, QueenSide: true}},
				{KingPosition: core.E8, CastlingRights: core.CastlingRights{QueenSide: true}},
			},
		},
		{
			name:  "non-home file does nothing",
			color: core.WHITE,
			file:  core.FILE_D,
			inputSides: [2]core.SideState{
				{KingPosition: core.E1, CastlingRights: core.CastlingRights{KingSide: true, QueenSide: true}},
				{KingPosition: core.E8, CastlingRights: core.CastlingRights{KingSide: true, QueenSide: true}},
			},
			expectSides: [2]core.SideState{
				{KingPosition: core.E1, CastlingRights: core.CastlingRights{KingSide: true, QueenSide: true}},
				{KingPosition: core.E8, CastlingRights: core.CastlingRights{KingSide: true, QueenSide: true}},
			},
		},
		{
			name:  "already cleared right stays cleared",
			color: core.WHITE,
			file:  core.FILE_A,
			inputSides: [2]core.SideState{
				{KingPosition: core.E1, CastlingRights: core.CastlingRights{KingSide: true}},
				{KingPosition: core.E8, CastlingRights: core.CastlingRights{KingSide: true, QueenSide: true}},
			},
			expectSides: [2]core.SideState{
				{KingPosition: core.E1, CastlingRights: core.CastlingRights{KingSide: true}},
				{KingPosition: core.E8, CastlingRights: core.CastlingRights{KingSide: true, QueenSide: true}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var board core.Board
			ctx := &core.TurnContext{
				MoveContext: core.MoveContext{
					BoardContext: core.BoardContext{Board: &board},
					Sides:        tt.inputSides,
				},
			}

			clearCastlingRightByFile(ctx, tt.color, tt.file)

			if ctx.Sides != tt.expectSides {
				t.Errorf("Sides = %+v, want %+v", ctx.Sides, tt.expectSides)
			}
		})
	}
}
