package engine

import (
	"testing"

	"github.com/7-Dany/chess/core"
)

func TestIsSquareAttacked(t *testing.T) {
	tests := []struct {
		name          string
		position      core.Position
		attackerColor core.PieceColor
		sideToMove    core.PieceColor
		setupBoard    func(*core.Board)
		expected      bool
	}{
		// ==================== Knight ====================
		{
			name:          "knight attacks from valid L-shape",
			position:      core.E5,
			attackerColor: core.WHITE,
			sideToMove:    core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.D7] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
			},
			expected: true,
		},
		{
			name:          "knight does not attack from adjacent square",
			position:      core.E5,
			attackerColor: core.WHITE,
			sideToMove:    core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.D6] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
			},
			expected: false,
		},
		{
			name:          "knight of wrong color ignored",
			position:      core.E5,
			attackerColor: core.WHITE,
			sideToMove:    core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.D7] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.BLACK})
			},
			expected: false,
		},
		{
			name:          "knight attacks from corner H8 to F7",
			position:      core.F7,
			attackerColor: core.WHITE,
			sideToMove:    core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.H8] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
			},
			expected: true,
		},

		// ==================== Bishop ====================
		{
			name:          "bishop attacks diagonally long range",
			position:      core.E5,
			attackerColor: core.WHITE,
			sideToMove:    core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.B2] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
			},
			expected: true,
		},
		{
			name:          "bishop attacks diagonally short range",
			position:      core.E5,
			attackerColor: core.WHITE,
			sideToMove:    core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.D4] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
			},
			expected: true,
		},
		{
			name:          "bishop blocked by own piece no other attacker",
			position:      core.E5,
			attackerColor: core.WHITE,
			sideToMove:    core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.B2] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
				b[core.D4] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
			},
			expected: false,
		},
		{
			name:          "bishop blocked by enemy piece",
			position:      core.E5,
			attackerColor: core.WHITE,
			sideToMove:    core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.B2] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
				b[core.D4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
			},
			expected: false,
		},
		{
			name:          "bishop does not attack orthogonally",
			position:      core.E5,
			attackerColor: core.WHITE,
			sideToMove:    core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.E3] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
			},
			expected: false,
		},

		// ==================== Rook ====================
		{
			name:          "rook attacks along file",
			position:      core.E5,
			attackerColor: core.WHITE,
			sideToMove:    core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.E2] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			expected: true,
		},
		{
			name:          "rook attacks along rank",
			position:      core.E5,
			attackerColor: core.WHITE,
			sideToMove:    core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.A5] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			expected: true,
		},
		{
			name:          "rook blocked by own piece no other attacker",
			position:      core.E5,
			attackerColor: core.WHITE,
			sideToMove:    core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.E2] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
				b[core.E4] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
			},
			expected: false,
		},
		{
			name:          "rook blocked by enemy piece",
			position:      core.E5,
			attackerColor: core.WHITE,
			sideToMove:    core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.E2] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
				b[core.E4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
			},
			expected: false,
		},
		{
			name:          "rook does not attack diagonally",
			position:      core.E5,
			attackerColor: core.WHITE,
			sideToMove:    core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.D4] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			expected: false,
		},

		// ==================== Queen (via bishop and rook paths) ====================
		{
			name:          "queen attacks diagonally caught by bishop path",
			position:      core.E5,
			attackerColor: core.WHITE,
			sideToMove:    core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.B2] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.WHITE})
			},
			expected: true,
		},
		{
			name:          "queen attacks orthogonally caught by rook path",
			position:      core.E5,
			attackerColor: core.WHITE,
			sideToMove:    core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.E2] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.WHITE})
			},
			expected: true,
		},
		{
			name:          "queen blocked diagonally no other attacker",
			position:      core.E5,
			attackerColor: core.WHITE,
			sideToMove:    core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.B2] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.WHITE})
				b[core.D4] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
			},
			expected: false,
		},
		{
			name:          "queen blocked orthogonally by enemy piece",
			position:      core.E5,
			attackerColor: core.WHITE,
			sideToMove:    core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.E2] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.WHITE})
				b[core.E4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
			},
			expected: false,
		},

		// ==================== King ====================
		{
			name:          "king attacks adjacent square",
			position:      core.E5,
			attackerColor: core.WHITE,
			sideToMove:    core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.E6] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
			},
			expected: true,
		},
		{
			name:          "king attacks diagonally adjacent",
			position:      core.E5,
			attackerColor: core.WHITE,
			sideToMove:    core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.F6] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
			},
			expected: true,
		},
		{
			name:          "king does not attack two squares away",
			position:      core.E5,
			attackerColor: core.WHITE,
			sideToMove:    core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.E7] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
			},
			expected: false,
		},
		{
			name:          "king of wrong color ignored",
			position:      core.E5,
			attackerColor: core.WHITE,
			sideToMove:    core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.E6] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
			},
			expected: false,
		},

		// ==================== White Pawn ====================
		{
			name:          "white pawn attacks from below-left",
			position:      core.E5,
			attackerColor: core.WHITE,
			sideToMove:    core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.D4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
			},
			expected: true,
		},
		{
			name:          "white pawn attacks from below-right",
			position:      core.E5,
			attackerColor: core.WHITE,
			sideToMove:    core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.F4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
			},
			expected: true,
		},
		{
			name:          "white pawn does not attack from above",
			position:      core.E5,
			attackerColor: core.WHITE,
			sideToMove:    core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.D6] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
			},
			expected: false,
		},
		{
			name:          "white pawn does not attack forward same file",
			position:      core.E5,
			attackerColor: core.WHITE,
			sideToMove:    core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.E4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
			},
			expected: false,
		},

		// ==================== Black Pawn ====================
		{
			name:          "black pawn attacks from above-left",
			position:      core.E5,
			attackerColor: core.BLACK,
			sideToMove:    core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.D6] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
			},
			expected: true,
		},
		{
			name:          "black pawn attacks from above-right",
			position:      core.E5,
			attackerColor: core.BLACK,
			sideToMove:    core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.F6] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
			},
			expected: true,
		},
		{
			name:          "black pawn does not attack from below",
			position:      core.E5,
			attackerColor: core.BLACK,
			sideToMove:    core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.D4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
			},
			expected: false,
		},
		{
			name:          "black pawn does not attack forward same file",
			position:      core.E5,
			attackerColor: core.BLACK,
			sideToMove:    core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.E6] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
			},
			expected: false,
		},

		// ==================== Pawn edge cases ====================
		{
			name:          "white pawn on rank 1 cannot attack from below impossible",
			position:      core.E1,
			attackerColor: core.WHITE,
			sideToMove:    core.BLACK,
			setupBoard:    func(b *core.Board) {},
			expected:      false,
		},
		{
			name:          "black pawn on rank 8 cannot attack from above impossible",
			position:      core.E8,
			attackerColor: core.BLACK,
			sideToMove:    core.WHITE,
			setupBoard:    func(b *core.Board) {},
			expected:      false,
		},
		{
			name:          "black pawn can attack rank 1 from above",
			position:      core.E1,
			attackerColor: core.BLACK,
			sideToMove:    core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.D2] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
			},
			expected: true,
		},
		{
			name:          "white pawn can attack rank 8 from below",
			position:      core.E8,
			attackerColor: core.WHITE,
			sideToMove:    core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.D7] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
			},
			expected: true,
		},
		{
			name:          "A file pawn only has right diagonal",
			position:      core.A5,
			attackerColor: core.WHITE,
			sideToMove:    core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.B4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
			},
			expected: true,
		},
		{
			name:          "A file has no left diagonal for pawn",
			position:      core.A5,
			attackerColor: core.WHITE,
			sideToMove:    core.BLACK,
			setupBoard:    func(b *core.Board) {},
			expected:      false,
		},
		{
			name:          "H file pawn only has left diagonal",
			position:      core.H5,
			attackerColor: core.WHITE,
			sideToMove:    core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.G4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
			},
			expected: true,
		},
		{
			name:          "H file has no right diagonal for pawn",
			position:      core.H5,
			attackerColor: core.WHITE,
			sideToMove:    core.BLACK,
			setupBoard:    func(b *core.Board) {},
			expected:      false,
		},
		{
			name:          "pawn of wrong color on diagonal ignored",
			position:      core.E5,
			attackerColor: core.WHITE,
			sideToMove:    core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.D4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
			},
			expected: false,
		},
		{
			name:          "non-pawn on pawn diagonal does not trigger pawn attack",
			position:      core.E5,
			attackerColor: core.WHITE,
			sideToMove:    core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.D4] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			expected: false,
		},

		// ==================== Sliding piece stacking ====================
		{
			name:          "rook behind rook on same rank first one attacks",
			position:      core.E5,
			attackerColor: core.WHITE,
			sideToMove:    core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.D5] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
				b[core.C5] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			expected: true,
		},
		{
			name:          "queen behind enemy blocker on file is blocked",
			position:      core.E5,
			attackerColor: core.WHITE,
			sideToMove:    core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.E3] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
				b[core.E2] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.WHITE})
			},
			expected: false,
		},
		{
			name:          "queen behind enemy blocker on diagonal is blocked",
			position:      core.E5,
			attackerColor: core.WHITE,
			sideToMove:    core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.D4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
				b[core.C3] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.WHITE})
			},
			expected: false,
		},

		// ==================== Corner positions ====================
		{
			name:          "rook attacks along full file from A1 to A8",
			position:      core.A8,
			attackerColor: core.WHITE,
			sideToMove:    core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.A1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			expected: true,
		},
		{
			name:          "bishop attacks along full diagonal from A1 to H8",
			position:      core.H8,
			attackerColor: core.WHITE,
			sideToMove:    core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.A1] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
			},
			expected: true,
		},

		// ==================== Multiple attackers ====================
		{
			name:          "multiple attackers of same color",
			position:      core.E5,
			attackerColor: core.WHITE,
			sideToMove:    core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.D7] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
				b[core.E2] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			expected: true,
		},
		{
			name:          "only friendly pieces not attacked",
			position:      core.E5,
			attackerColor: core.WHITE,
			sideToMove:    core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.D7] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.BLACK})
				b[core.E2] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
			},
			expected: false,
		},

		// ==================== Empty board ====================
		{
			name:          "completely empty board",
			position:      core.E5,
			attackerColor: core.WHITE,
			sideToMove:    core.BLACK,
			setupBoard:    func(b *core.Board) {},
			expected:      false,
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
				},
			}

			got := engine.IsSquareAttacked(tt.position, tt.attackerColor, ctx)
			if got != tt.expected {
				t.Errorf("IsSquareAttacked(%v, %v) = %v, want %v", tt.position, tt.attackerColor, got, tt.expected)
			}
		})
	}
}
