package engine

import (
	"testing"

	"github.com/7-Dany/chess/core"
)

func TestGetLegalMoves(t *testing.T) {
	defaultSides := [2]core.SideState{
		{KingPosition: core.E1, CanCastleKingSide: true, CanCastleQueenSide: true},
		{KingPosition: core.E8, CanCastleKingSide: true, CanCastleQueenSide: true},
	}
	noRights := [2]core.SideState{
		{KingPosition: core.E1},
		{KingPosition: core.E8, CanCastleKingSide: true, CanCastleQueenSide: true},
	}

	type moveCheck struct {
		from     core.Position
		to       core.Position
		moveType core.MoveType
		want     bool
	}

	hasMove := func(moves []core.Move, from, to core.Position, moveType core.MoveType) bool {
		for _, m := range moves {
			if m.From == from && m.To == to && m.Type == moveType {
				return true
			}
		}
		return false
	}

	checkMoves := func(moves []core.Move, checks []moveCheck) {
		t.Helper()
		for _, c := range checks {
			got := hasMove(moves, c.from, c.to, c.moveType)
			if c.want && !got {
				t.Errorf("expected move %v->%v (%v) present", c.from, c.to, c.moveType)
			}
			if !c.want && got {
				t.Errorf("expected move %v->%v (%v) absent", c.from, c.to, c.moveType)
			}
		}
	}

	tests := []struct {
		name       string
		setupBoard func(*core.Board)
		sideToMove core.PieceColor
		sides      [2]core.SideState
		inputEP    core.Position
		position   core.Position
		wantCount  int
		checks     []moveCheck
	}{
		// ==================== King safety filter ====================
		{
			name: "pinned rook can only move along pin lines",
			setupBoard: func(b *core.Board) {
				b[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
				b[core.E2] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
				b[core.E8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			position:   core.E2,
			wantCount:  6,
			checks: []moveCheck{
				{core.E2, core.E3, core.NORMAL, true},
				{core.E2, core.E8, core.NORMAL, true},
				{core.E2, core.D2, core.NORMAL, false},
				{core.E2, core.F2, core.NORMAL, false},
			},
		},
		{
			name: "pinned bishop has no legal moves",
			setupBoard: func(b *core.Board) {
				b[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
				b[core.E2] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
				b[core.E8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
				b[core.H8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			sides: [2]core.SideState{
				{KingPosition: core.E1},
				{KingPosition: core.H8},
			},
			position:  core.E2,
			wantCount: 0,
		},
		{
			name: "king in check on E-file must escape sideways",
			setupBoard: func(b *core.Board) {
				b[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
				b[core.E8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			position:   core.E1,
			wantCount:  4,
			checks: []moveCheck{
				{core.E1, core.D1, core.NORMAL, true},
				{core.E1, core.D2, core.NORMAL, true},
				{core.E1, core.F2, core.NORMAL, true},
				{core.E1, core.F1, core.NORMAL, true},
				{core.E1, core.E2, core.NORMAL, false},
				{core.E1, core.G1, core.CASTLING, false},
			},
		},
		{
			name: "king in check on rank 1 cannot stay on rank",
			setupBoard: func(b *core.Board) {
				b[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
				b[core.A1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			position:   core.E1,
			wantCount:  3,
			checks: []moveCheck{
				{core.E1, core.D2, core.NORMAL, true},
				{core.E1, core.E2, core.NORMAL, true},
				{core.E1, core.F2, core.NORMAL, true},
				{core.E1, core.D1, core.NORMAL, false},
				{core.E1, core.F1, core.NORMAL, false},
			},
		},
		{
			name: "king can capture undefended checker",
			setupBoard: func(b *core.Board) {
				b[core.E4] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
				b[core.E5] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
				b[core.E8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			sides: [2]core.SideState{
				{KingPosition: core.E4},
				{KingPosition: core.E8},
			},
			position:  core.E4,
			wantCount: 5,
			checks: []moveCheck{
				{core.E4, core.D3, core.NORMAL, true},
				{core.E4, core.D4, core.NORMAL, true},
				{core.E4, core.F3, core.NORMAL, true},
				{core.E4, core.F4, core.NORMAL, true},
				{core.E4, core.E5, core.NORMAL, true},
				{core.E4, core.D5, core.NORMAL, false},
				{core.E4, core.E3, core.NORMAL, false},
				{core.E4, core.F5, core.NORMAL, false},
			},
		},
		{
			name: "king cannot capture defended checker",
			setupBoard: func(b *core.Board) {
				b[core.E4] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
				b[core.E5] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
				b[core.D6] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
				b[core.E8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			sides: [2]core.SideState{
				{KingPosition: core.E4},
				{KingPosition: core.E8},
			},
			position:  core.E4,
			wantCount: 4,
			checks: []moveCheck{
				{core.E4, core.D3, core.NORMAL, true},
				{core.E4, core.D4, core.NORMAL, true},
				{core.E4, core.F3, core.NORMAL, true},
				{core.E4, core.F4, core.NORMAL, true},
				{core.E4, core.E5, core.NORMAL, false},
				{core.E4, core.D5, core.NORMAL, false},
				{core.E4, core.E3, core.NORMAL, false},
				{core.E4, core.F5, core.NORMAL, false},
			},
		},
		{
			name: "king cannot capture defended adjacent piece",
			setupBoard: func(b *core.Board) {
				b[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
				b[core.D2] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.BLACK})
				b[core.C3] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
				b[core.E8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			sides: [2]core.SideState{
				{KingPosition: core.E1},
				{KingPosition: core.E8},
			},
			position:  core.E1,
			wantCount: 3,
			checks: []moveCheck{
				{core.E1, core.D1, core.NORMAL, true},
				{core.E1, core.E2, core.NORMAL, true},
				{core.E1, core.F2, core.NORMAL, true},
				{core.E1, core.D2, core.NORMAL, false},
				{core.E1, core.F1, core.NORMAL, false},
			},
		},
		{
			name: "king cannot move adjacent to enemy king",
			setupBoard: func(b *core.Board) {
				b[core.E4] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
				b[core.E6] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			sides: [2]core.SideState{
				{KingPosition: core.E4},
				{KingPosition: core.E6},
			},
			position:  core.E4,
			wantCount: 5,
			checks: []moveCheck{
				{core.E4, core.D3, core.NORMAL, true},
				{core.E4, core.D4, core.NORMAL, true},
				{core.E4, core.E3, core.NORMAL, true},
				{core.E4, core.F3, core.NORMAL, true},
				{core.E4, core.F4, core.NORMAL, true},
				{core.E4, core.D5, core.NORMAL, false},
				{core.E4, core.E5, core.NORMAL, false},
				{core.E4, core.F5, core.NORMAL, false},
			},
		},
		{
			name: "knight can block check by interposing",
			setupBoard: func(b *core.Board) {
				b[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
				b[core.D6] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
				b[core.E8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
				b[core.H8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			sides: [2]core.SideState{
				{KingPosition: core.E1},
				{KingPosition: core.H8},
			},
			position:  core.D6,
			wantCount: 2,
			checks: []moveCheck{
				{core.D6, core.E4, core.NORMAL, true},
				{core.D6, core.E8, core.NORMAL, true},
				{core.D6, core.B5, core.NORMAL, false},
				{core.D6, core.F7, core.NORMAL, false},
				{core.D6, core.C8, core.NORMAL, false},
			},
		},
		{
			name: "knight can capture checker or interpose to resolve check",
			setupBoard: func(b *core.Board) {
				b[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
				b[core.C7] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
				b[core.E8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
				b[core.H8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			sides: [2]core.SideState{
				{KingPosition: core.E1},
				{KingPosition: core.H8},
			},
			position:  core.C7,
			wantCount: 2,
			checks: []moveCheck{
				{core.C7, core.E8, core.NORMAL, true}, // capture the checker
				{core.C7, core.E6, core.NORMAL, true}, // interpose on E-file
				{core.C7, core.A6, core.NORMAL, false},
				{core.C7, core.D5, core.NORMAL, false},
			},
		},
		{
			name: "en passant exposing king on the rank is illegal",
			setupBoard: func(b *core.Board) {
				b[core.H5] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
				b[core.F5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
				b[core.E5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
				b[core.A5] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
				b[core.E8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			sides: [2]core.SideState{
				{KingPosition: core.H5},
				{KingPosition: core.E8},
			},
			inputEP:   core.E6,
			position:  core.F5,
			wantCount: 1,
			checks: []moveCheck{
				{core.F5, core.F6, core.NORMAL, true},
				{core.F5, core.E6, core.EN_PASSANT, false},
			},
		},
		{
			name: "promotion push blocked by check is filtered",
			setupBoard: func(b *core.Board) {
				b[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
				b[core.D7] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
				b[core.E8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
				b[core.H8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			sides: [2]core.SideState{
				{KingPosition: core.E1},
				{KingPosition: core.H8},
			},
			position:  core.D7,
			wantCount: 4,
			checks: []moveCheck{
				{core.D7, core.D8, core.PROMOTION, false},
				{core.D7, core.E8, core.PROMOTION, true},
			},
		},

		// ==================== Castling integration ====================
		{
			name: "castling available when all conditions met",
			setupBoard: func(b *core.Board) {
				b[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
				b[core.A1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
				b[core.H1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
				b[core.E8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			position:   core.E1,
			wantCount:  7,
			checks: []moveCheck{
				{core.E1, core.D1, core.NORMAL, true},
				{core.E1, core.F1, core.NORMAL, true},
				{core.E1, core.G1, core.CASTLING, true},
				{core.E1, core.C1, core.CASTLING, true},
			},
		},
		{
			name: "castling removed when king in check",
			setupBoard: func(b *core.Board) {
				b[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
				b[core.A1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
				b[core.H1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
				b[core.E8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			position:   core.E1,
			wantCount:  4,
			checks: []moveCheck{
				{core.E1, core.G1, core.CASTLING, false},
				{core.E1, core.C1, core.CASTLING, false},
			},
		},
		{
			name: "king-side castling removed when F1 occupied",
			setupBoard: func(b *core.Board) {
				b[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
				b[core.A1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
				b[core.H1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
				b[core.F1] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
				b[core.E8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			position:   core.E1,
			wantCount:  5,
			checks: []moveCheck{
				{core.E1, core.G1, core.CASTLING, false},
				{core.E1, core.C1, core.CASTLING, true},
				{core.E1, core.F1, core.NORMAL, false},
			},
		},
		{
			name: "king-side castling removed when F1 attacked",
			setupBoard: func(b *core.Board) {
				b[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
				b[core.A1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
				b[core.H1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
				b[core.F8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
				b[core.E8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			position:   core.E1,
			wantCount:  4,
			checks: []moveCheck{
				{core.E1, core.G1, core.CASTLING, false},
				{core.E1, core.C1, core.CASTLING, true},
				{core.E1, core.F1, core.NORMAL, false},
				{core.E1, core.F2, core.NORMAL, false},
				{core.E1, core.D1, core.NORMAL, true},
			},
		},
		{
			name: "no castling when rights lost",
			setupBoard: func(b *core.Board) {
				b[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
				b[core.A1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
				b[core.H1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
				b[core.E8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			sides:      noRights,
			position:   core.E1,
			wantCount:  5,
			checks: []moveCheck{
				{core.E1, core.G1, core.CASTLING, false},
				{core.E1, core.C1, core.CASTLING, false},
			},
		},
		{
			name: "non-king piece never gets castling moves",
			setupBoard: func(b *core.Board) {
				b[core.B1] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
				b[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
				b[core.E8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			position:   core.B1,
			wantCount:  3,
			checks: []moveCheck{
				{core.B1, core.A3, core.NORMAL, true},
				{core.B1, core.C3, core.NORMAL, true},
				{core.B1, core.D2, core.NORMAL, true},
				{core.B1, core.G1, core.CASTLING, false},
			},
		},
	}

	engine := NewDefaultEngine()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var board core.Board
			tt.setupBoard(&board)
			ctx := core.TurnContext{
				MoveContext: core.MoveContext{
					BoardContext:    core.BoardContext{Board: &board},
					SideToMove:      tt.sideToMove,
					Sides:           tt.sides,
					EnPassantTarget: tt.inputEP,
				},
			}
			moves := engine.GetLegalMoves(tt.position, ctx)
			if tt.wantCount >= 0 && len(moves) != tt.wantCount {
				t.Errorf("count = %d, want %d", len(moves), tt.wantCount)
			}
			checkMoves(moves, tt.checks)
		})
	}
}

func TestHasAnyLegalMoves(t *testing.T) {
	defaultSides := [2]core.SideState{
		{KingPosition: core.E1, CanCastleKingSide: true, CanCastleQueenSide: true},
		{KingPosition: core.E8, CanCastleKingSide: true, CanCastleQueenSide: true},
	}

	tests := []struct {
		name       string
		setupBoard func(*core.Board)
		sideToMove core.PieceColor
		sides      [2]core.SideState
		want       bool
	}{
		{
			name: "side with moves returns true",
			setupBoard: func(b *core.Board) {
				b[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
				b[core.B1] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
				b[core.E8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			want:       true,
		},
		{
			name: "checkmate returns false",
			setupBoard: func(b *core.Board) {
				b[core.H8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
				b[core.G7] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.WHITE})
				b[core.F6] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
			},
			sideToMove: core.BLACK,
			sides: [2]core.SideState{
				{KingPosition: core.F6},
				{KingPosition: core.H8},
			},
			want: false,
		},
		{
			name: "stalemate returns false",
			setupBoard: func(b *core.Board) {
				b[core.A1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
				b[core.B3] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.WHITE})
				b[core.C2] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
			},
			sideToMove: core.BLACK,
			sides: [2]core.SideState{
				{KingPosition: core.C2},
				{KingPosition: core.A1},
			},
			want: false,
		},
		{
			name: "only checks side-to-move pieces (white stalemated, black has moves)",
			setupBoard: func(b *core.Board) {
				b[core.A1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
				b[core.B3] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.BLACK})
				b[core.C2] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			sides: [2]core.SideState{
				{KingPosition: core.A1},
				{KingPosition: core.C2},
			},
			want: false,
		},
		{
			name: "same board, black to move has moves",
			setupBoard: func(b *core.Board) {
				b[core.A1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
				b[core.B3] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.BLACK})
				b[core.C2] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
			},
			sideToMove: core.BLACK,
			sides: [2]core.SideState{
				{KingPosition: core.A1},
				{KingPosition: core.C2},
			},
			want: true,
		},
		{
			name: "no pieces for side to move returns false",
			setupBoard: func(b *core.Board) {
				b[core.E8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			sides: [2]core.SideState{
				{KingPosition: core.NoPosition},
				{KingPosition: core.E8},
			},
			want: false,
		},
		{
			name: "first piece blocked, second piece has moves",
			setupBoard: func(b *core.Board) {
				b[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
				b[core.A2] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
				b[core.A3] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
				b[core.B1] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
				b[core.E8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			want:       true,
		},
		{
			name: "pinned piece has no moves but other piece does",
			setupBoard: func(b *core.Board) {
				b[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
				b[core.E2] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
				b[core.B1] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
				b[core.E8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
				b[core.H8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			sides: [2]core.SideState{
				{KingPosition: core.E1},
				{KingPosition: core.H8},
			},
			want: true,
		},
	}

	engine := NewDefaultEngine()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var board core.Board
			tt.setupBoard(&board)
			ctx := core.TurnContext{
				MoveContext: core.MoveContext{
					BoardContext: core.BoardContext{Board: &board},
					SideToMove:   tt.sideToMove,
					Sides:        tt.sides,
				},
			}
			got := engine.HasAnyLegalMoves(ctx)
			if got != tt.want {
				t.Errorf("HasAnyLegalMoves = %v, want %v", got, tt.want)
			}
		})
	}
}
