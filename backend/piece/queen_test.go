package piece

import (
	"testing"

	"github.com/7-Dany/chess/core"
)

func TestQueenIsAttacking(t *testing.T) {
	tests := []struct {
		name       string
		setupBoard func(*core.Board)
		color      core.PieceColor
		target     core.Position
		want       bool
	}{
		// ==================== Diagonal attacks (bishop-like) ====================
		{
			name: "queen on up-right diagonal attacks target",
			setupBoard: func(b *core.Board) {
				b[core.H7] = core.Square{Piece: core.Piece{Type: core.QUEEN, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true,
		},
		{
			name: "queen on up-left diagonal attacks target",
			setupBoard: func(b *core.Board) {
				b[core.A8] = core.Square{Piece: core.Piece{Type: core.QUEEN, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true,
		},
		{
			name: "queen on down-right diagonal attacks target",
			setupBoard: func(b *core.Board) {
				b[core.H1] = core.Square{Piece: core.Piece{Type: core.QUEEN, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true,
		},
		{
			name: "queen on down-left diagonal attacks target",
			setupBoard: func(b *core.Board) {
				b[core.B1] = core.Square{Piece: core.Piece{Type: core.QUEEN, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true,
		},
		{
			name: "queen adjacent diagonally attacks (distance 1)",
			setupBoard: func(b *core.Board) {
				b[core.D3] = core.Square{Piece: core.Piece{Type: core.QUEEN, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true,
		},
		{
			name: "queen at maximum diagonal distance attacks (distance 7)",
			setupBoard: func(b *core.Board) {
				b[core.H8] = core.Square{Piece: core.Piece{Type: core.QUEEN, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.A1,
			want:   true,
		},

		// ==================== Orthogonal attacks (rook-like) ====================
		{
			name: "queen above target on same file attacks",
			setupBoard: func(b *core.Board) {
				b[core.E8] = core.Square{Piece: core.Piece{Type: core.QUEEN, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true,
		},
		{
			name: "queen below target on same file attacks",
			setupBoard: func(b *core.Board) {
				b[core.E1] = core.Square{Piece: core.Piece{Type: core.QUEEN, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true,
		},
		{
			name: "queen to the right on same rank attacks",
			setupBoard: func(b *core.Board) {
				b[core.H4] = core.Square{Piece: core.Piece{Type: core.QUEEN, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true,
		},
		{
			name: "queen to the left on same rank attacks",
			setupBoard: func(b *core.Board) {
				b[core.A4] = core.Square{Piece: core.Piece{Type: core.QUEEN, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true,
		},
		{
			name: "queen adjacent orthogonally attacks (distance 1)",
			setupBoard: func(b *core.Board) {
				b[core.D4] = core.Square{Piece: core.Piece{Type: core.QUEEN, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true,
		},

		// ==================== Not on any attack line ====================
		{
			name: "queen off-diagonal and off-orthogonal does not attack",
			setupBoard: func(b *core.Board) {
				b[core.F6] = core.Square{Piece: core.Piece{Type: core.QUEEN, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "knight-move away does not attack",
			setupBoard: func(b *core.Board) {
				b[core.D6] = core.Square{Piece: core.Piece{Type: core.QUEEN, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "queen two files and two ranks away on different diagonal does not attack",
			setupBoard: func(b *core.Board) {
				b[core.F6] = core.Square{Piece: core.Piece{Type: core.QUEEN, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},

		// ==================== Blockers (diagonal rays) ====================
		{
			name: "friendly piece blocks diagonal attack",
			setupBoard: func(b *core.Board) {
				b[core.H7] = core.Square{Piece: core.Piece{Type: core.QUEEN, Color: core.WHITE}, Occupied: true}
				b[core.F5] = core.Square{Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "enemy piece blocks diagonal attack",
			setupBoard: func(b *core.Board) {
				b[core.H7] = core.Square{Piece: core.Piece{Type: core.QUEEN, Color: core.WHITE}, Occupied: true}
				b[core.F5] = core.Square{Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "blocker adjacent to target blocks diagonal attack",
			setupBoard: func(b *core.Board) {
				b[core.H7] = core.Square{Piece: core.Piece{Type: core.QUEEN, Color: core.WHITE}, Occupied: true}
				b[core.F5] = core.Square{Piece: core.Piece{Type: core.ROOK, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},

		// ==================== Blockers (orthogonal rays) ====================
		{
			name: "friendly piece blocks orthogonal attack on file",
			setupBoard: func(b *core.Board) {
				b[core.E8] = core.Square{Piece: core.Piece{Type: core.QUEEN, Color: core.WHITE}, Occupied: true}
				b[core.E6] = core.Square{Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "enemy piece blocks orthogonal attack on rank",
			setupBoard: func(b *core.Board) {
				b[core.H4] = core.Square{Piece: core.Piece{Type: core.QUEEN, Color: core.WHITE}, Occupied: true}
				b[core.F4] = core.Square{Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "blocker adjacent to target blocks orthogonal attack",
			setupBoard: func(b *core.Board) {
				b[core.H4] = core.Square{Piece: core.Piece{Type: core.QUEEN, Color: core.WHITE}, Occupied: true}
				b[core.F4] = core.Square{Piece: core.Piece{Type: core.ROOK, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},

		// ==================== Color filtering ====================
		{
			name: "black queen, asking for white attackers",
			setupBoard: func(b *core.Board) {
				b[core.H7] = core.Square{Piece: core.Piece{Type: core.QUEEN, Color: core.BLACK}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "black queen, asking for black attackers",
			setupBoard: func(b *core.Board) {
				b[core.H7] = core.Square{Piece: core.Piece{Type: core.QUEEN, Color: core.BLACK}, Occupied: true}
			},
			color:  core.BLACK,
			target: core.E4,
			want:   true,
		},

		// ==================== Non-queen pieces on attack lines (clean separation) ====================
		{
			name: "bishop on diagonal does NOT trigger queen attack",
			setupBoard: func(b *core.Board) {
				b[core.H7] = core.Square{Piece: core.Piece{Type: core.BISHOP, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "rook on orthogonal does NOT trigger queen attack",
			setupBoard: func(b *core.Board) {
				b[core.E8] = core.Square{Piece: core.Piece{Type: core.ROOK, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "knight on attack line does not trigger queen attack",
			setupBoard: func(b *core.Board) {
				b[core.E5] = core.Square{Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "king on attack line does not trigger queen attack",
			setupBoard: func(b *core.Board) {
				b[core.E5] = core.Square{Piece: core.Piece{Type: core.KING, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "pawn on attack line does not trigger queen attack",
			setupBoard: func(b *core.Board) {
				b[core.E5] = core.Square{Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},

		// ==================== Edge cases ====================
		{
			name:       "empty board returns false",
			setupBoard: func(b *core.Board) {},
			color:      core.WHITE,
			target:     core.E4,
			want:       false,
		},
		{
			name: "queen on target square itself does not attack",
			setupBoard: func(b *core.Board) {
				b[core.E4] = core.Square{Piece: core.Piece{Type: core.QUEEN, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},

		// ==================== Corners ====================
		{
			name: "target on A1, queen on H8 attacks diagonally",
			setupBoard: func(b *core.Board) {
				b[core.H8] = core.Square{Piece: core.Piece{Type: core.QUEEN, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.A1,
			want:   true,
		},
		{
			name: "target on A1, queen on A8 attacks on file",
			setupBoard: func(b *core.Board) {
				b[core.A8] = core.Square{Piece: core.Piece{Type: core.QUEEN, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.A1,
			want:   true,
		},
		{
			name: "target on A1, queen on H1 attacks on rank",
			setupBoard: func(b *core.Board) {
				b[core.H1] = core.Square{Piece: core.Piece{Type: core.QUEEN, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.A1,
			want:   true,
		},
		{
			name: "target on H8, queen on A1 attacks diagonally",
			setupBoard: func(b *core.Board) {
				b[core.A1] = core.Square{Piece: core.Piece{Type: core.QUEEN, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.H8,
			want:   true,
		},
		{
			name: "target on A8, queen on H1 attacks diagonally",
			setupBoard: func(b *core.Board) {
				b[core.H1] = core.Square{Piece: core.Piece{Type: core.QUEEN, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.A8,
			want:   true,
		},
		{
			name: "target on H1, queen on A8 attacks diagonally",
			setupBoard: func(b *core.Board) {
				b[core.A8] = core.Square{Piece: core.Piece{Type: core.QUEEN, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.H1,
			want:   true,
		},

		// ==================== Multiple queens ====================
		{
			name: "multiple queens, one attacks diagonally",
			setupBoard: func(b *core.Board) {
				b[core.A1] = core.Square{Piece: core.Piece{Type: core.QUEEN, Color: core.WHITE}, Occupied: true}
				b[core.H7] = core.Square{Piece: core.Piece{Type: core.QUEEN, Color: core.WHITE}, Occupied: true}
				b[core.A8] = core.Square{Piece: core.Piece{Type: core.QUEEN, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true, // H7 queen attacks via diagonal
		},
		{
			name: "multiple queens, one attacks orthogonally",
			setupBoard: func(b *core.Board) {
				b[core.A1] = core.Square{Piece: core.Piece{Type: core.QUEEN, Color: core.WHITE}, Occupied: true}
				b[core.H8] = core.Square{Piece: core.Piece{Type: core.QUEEN, Color: core.WHITE}, Occupied: true}
				b[core.E8] = core.Square{Piece: core.Piece{Type: core.QUEEN, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true, // E8 queen attacks via file
		},
		{
			name: "multiple enemy queens, none matching color",
			setupBoard: func(b *core.Board) {
				b[core.H7] = core.Square{Piece: core.Piece{Type: core.QUEEN, Color: core.BLACK}, Occupied: true}
				b[core.E8] = core.Square{Piece: core.Piece{Type: core.QUEEN, Color: core.BLACK}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "mixed-color queens, only matching color counts",
			setupBoard: func(b *core.Board) {
				b[core.H7] = core.Square{Piece: core.Piece{Type: core.QUEEN, Color: core.BLACK}, Occupied: true}
				b[core.E8] = core.Square{Piece: core.Piece{Type: core.QUEEN, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true, // E8 (white) attacks
		},

		// ==================== Blocker is queen of wrong color ====================
		{
			name: "enemy queen on ray blocks but doesn't count as attacker",
			setupBoard: func(b *core.Board) {
				b[core.H7] = core.Square{Piece: core.Piece{Type: core.QUEEN, Color: core.WHITE}, Occupied: true}
				b[core.G6] = core.Square{Piece: core.Piece{Type: core.QUEEN, Color: core.BLACK}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false, // black queen at G6 blocks the white queen at H7
		},
		{
			name: "enemy queen as blocker, asking for black",
			setupBoard: func(b *core.Board) {
				b[core.H7] = core.Square{Piece: core.Piece{Type: core.QUEEN, Color: core.WHITE}, Occupied: true}
				b[core.G6] = core.Square{Piece: core.Piece{Type: core.QUEEN, Color: core.BLACK}, Occupied: true}
			},
			color:  core.BLACK,
			target: core.E4,
			want:   true, // black queen at G6 attacks E4 via diagonal
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var board core.Board
			tt.setupBoard(&board)
			ctx := core.BoardContext{Board: &board}

			queen := Queen{}
			got := queen.IsAttacking(tt.color, tt.target, ctx)

			if got != tt.want {
				t.Errorf("IsAttacking(%v, %v) = %v, want %v", tt.color, tt.target, got, tt.want)
			}
		})
	}
}

// test all attack moves for queen.
func TestQueenAttacks(t *testing.T) {
	tests := []struct {
		name       string
		from       core.Position
		setupBoard func(*core.Board)
		expected   []core.Position
	}{
		{
			name:       "center D4 empty board — 27 attacks (14 orthogonal + 13 diagonal)",
			from:       core.D4,
			setupBoard: func(b *core.Board) {},
			expected: []core.Position{
				// Orthogonal
				core.A4, core.B4, core.C4, core.E4, core.F4, core.G4, core.H4,
				core.D1, core.D2, core.D3, core.D5, core.D6, core.D7, core.D8,
				// Diagonals
				core.E5, core.F6, core.G7, core.H8, // NE
				core.E3, core.F2, core.G1, // SE
				core.C5, core.B6, core.A7, // NW
				core.C3, core.B2, core.A1, // SW
			},
		},
		{
			name:       "corner A1 empty board — 21 attacks (14 + 7)",
			from:       core.A1,
			setupBoard: func(b *core.Board) {},
			expected: []core.Position{
				// Orthogonal
				core.B1, core.C1, core.D1, core.E1, core.F1, core.G1, core.H1,
				core.A2, core.A3, core.A4, core.A5, core.A6, core.A7, core.A8,
				// Diagonal — only NE (one direction from corner)
				core.B2, core.C3, core.D4, core.E5, core.F6, core.G7, core.H8,
			},
		},
		{
			name:       "corner H1 empty board — 21 attacks",
			from:       core.H1,
			setupBoard: func(b *core.Board) {},
			expected: []core.Position{
				core.A1, core.B1, core.C1, core.D1, core.E1, core.F1, core.G1,
				core.H2, core.H3, core.H4, core.H5, core.H6, core.H7, core.H8,
				core.G2, core.F3, core.E4, core.D5, core.C6, core.B7, core.A8,
			},
		},
		{
			name:       "corner A8 empty board — 21 attacks",
			from:       core.A8,
			setupBoard: func(b *core.Board) {},
			expected: []core.Position{
				core.B8, core.C8, core.D8, core.E8, core.F8, core.G8, core.H8,
				core.A1, core.A2, core.A3, core.A4, core.A5, core.A6, core.A7,
				core.B7, core.C6, core.D5, core.E4, core.F3, core.G2, core.H1,
			},
		},
		{
			name:       "corner H8 empty board — 21 attacks",
			from:       core.H8,
			setupBoard: func(b *core.Board) {},
			expected: []core.Position{
				core.A8, core.B8, core.C8, core.D8, core.E8, core.F8, core.G8,
				core.H1, core.H2, core.H3, core.H4, core.H5, core.H6, core.H7,
				core.G7, core.F6, core.E5, core.D4, core.C3, core.B2, core.A1,
			},
		},
		{
			name:       "edge A4 empty board — 23 attacks (14 + 9)",
			from:       core.A4,
			setupBoard: func(b *core.Board) {},
			expected: []core.Position{
				// Orthogonal — same as rook on A4 (14 squares)
				core.B4, core.C4, core.D4, core.E4, core.F4, core.G4, core.H4,
				core.A1, core.A2, core.A3, core.A5, core.A6, core.A7, core.A8,
				// Diagonals — only 2 directions from file A
				// NE: B5, C6, D7, E8 (4)
				core.B5, core.C6, core.D7, core.E8,
				// SE: B3, C2, D1 (3)
				core.B3, core.C2, core.D1,
				// NW and SW blocked by file A — 0 squares each
			},
		},
		{
			name:       "edge D1 empty board — 23 attacks (14 + 9)",
			from:       core.D1,
			setupBoard: func(b *core.Board) {},
			expected: []core.Position{
				// Orthogonal
				core.A1, core.B1, core.C1, core.E1, core.F1, core.G1, core.H1,
				core.D2, core.D3, core.D4, core.D5, core.D6, core.D7, core.D8,
				// Diagonals — only 2 upward directions from rank 1
				// NE: E2, F3, G4, H5 (4)
				core.E2, core.F3, core.G4, core.H5,
				// NW: C2, B3, A4 (3)
				core.C2, core.B3, core.A4,
				// SE and SW blocked by rank 1
			},
		},
		{
			name: "center D4 blocked on orthogonal Up and diagonal NE",
			from: core.D4,
			setupBoard: func(b *core.Board) {
				b[core.D6] = core.Square{Occupied: true} // blocks Up ray
				b[core.F6] = core.Square{Occupied: true} // blocks NE ray
			},
			expected: []core.Position{
				// Orthogonal — Up stops at D6 (inclusive)
				core.A4, core.B4, core.C4, core.E4, core.F4, core.G4, core.H4,
				core.D1, core.D2, core.D3, core.D5, core.D6,
				// Diagonal NE stops at F6 (inclusive)
				core.E5, core.F6,
				// Other diagonals unchanged
				core.E3, core.F2, core.G1,
				core.C5, core.B6, core.A7,
				core.C3, core.B2, core.A1,
			},
		},
		{
			name: "center D4 trapped — all 8 directions blocked by adjacent pieces",
			from: core.D4,
			setupBoard: func(b *core.Board) {
				b[core.D5] = core.Square{Occupied: true} // Up
				b[core.D3] = core.Square{Occupied: true} // Down
				b[core.E4] = core.Square{Occupied: true} // Right
				b[core.C4] = core.Square{Occupied: true} // Left
				b[core.E5] = core.Square{Occupied: true} // NE
				b[core.E3] = core.Square{Occupied: true} // SE
				b[core.C5] = core.Square{Occupied: true} // NW
				b[core.C3] = core.Square{Occupied: true} // SW
			},
			expected: []core.Position{
				core.D5, core.D3, core.E4, core.C4,
				core.E5, core.E3, core.C5, core.C3,
			},
		},
	}

	queen := Queen{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var board core.Board
			tt.setupBoard(&board)
			ctx := core.BoardContext{Board: &board}

			got := queen.Attacks(tt.from, ctx)

			if len(got) != len(tt.expected) {
				t.Fatalf("got %d attacks, want %d — got %v", len(got), len(tt.expected), got)
			}

			gotSet := make(map[core.Position]struct{}, len(got))
			for _, p := range got {
				gotSet[p] = struct{}{}
			}

			for _, p := range tt.expected {
				if _, ok := gotSet[p]; !ok {
					t.Errorf("missing expected position %v", p)
				}
			}
		})
	}
}

// test all queen pseudolegal moves
func TestQueenPseudoLegalMoves(t *testing.T) {
	tests := []struct {
		name        string
		from        core.Position
		sideToMove  core.PieceColor
		setupBoard  func(*core.Board)
		expectedTos []core.Position
	}{
		{
			name:       "empty board center D4 — 27 moves",
			from:       core.D4,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {},
			expectedTos: []core.Position{
				// Orthogonal
				core.A4, core.B4, core.C4, core.E4, core.F4, core.G4, core.H4,
				core.D1, core.D2, core.D3, core.D5, core.D6, core.D7, core.D8,
				// Diagonal
				core.E5, core.F6, core.G7, core.H8,
				core.E3, core.F2, core.G1,
				core.C5, core.B6, core.A7,
				core.C3, core.B2, core.A1,
			},
		},
		{
			name:        "empty board corner A1 — 21 moves",
			from:        core.A1,
			sideToMove:  core.WHITE,
			setupBoard:  func(b *core.Board) {},
			expectedTos: []core.Position{core.B1, core.C1, core.D1, core.E1, core.F1, core.G1, core.H1, core.A2, core.A3, core.A4, core.A5, core.A6, core.A7, core.A8, core.B2, core.C3, core.D4, core.E5, core.F6, core.G7, core.H8},
		},
		{
			name:        "empty board corner H8 — 21 moves",
			from:        core.H8,
			sideToMove:  core.WHITE,
			setupBoard:  func(b *core.Board) {},
			expectedTos: []core.Position{core.A8, core.B8, core.C8, core.D8, core.E8, core.F8, core.G8, core.H1, core.H2, core.H3, core.H4, core.H5, core.H6, core.H7, core.G7, core.F6, core.E5, core.D4, core.C3, core.B2, core.A1},
		},
		{
			name:       "can capture enemy on orthogonal and diagonal in same turn",
			from:       core.D4,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.F4] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}} // orthogonal Right
				b[core.F6] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}} // diagonal NE
			},
			expectedTos: []core.Position{
				// Orthogonal — Right stops at F4 (inclusive)
				core.E4, core.F4,
				// Orthogonal — other directions unchanged
				core.A4, core.B4, core.C4,
				core.D1, core.D2, core.D3, core.D5, core.D6, core.D7, core.D8,
				// Diagonal NE stops at F6 (inclusive)
				core.E5, core.F6,
				// Other diagonals unchanged
				core.E3, core.F2, core.G1,
				core.C5, core.B6, core.A7,
				core.C3, core.B2, core.A1,
			},
		},
		{
			name:       "cannot move to friendly occupied square (orthogonal and diagonal)",
			from:       core.D4,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.F4] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}} // orthogonal
				b[core.F6] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}} // diagonal
			},
			expectedTos: []core.Position{
				// Orthogonal Right stops before F4 (exclusive)
				core.E4,
				// Orthogonal other directions unchanged
				core.A4, core.B4, core.C4,
				core.D1, core.D2, core.D3, core.D5, core.D6, core.D7, core.D8,
				// Diagonal NE stops before F6 (exclusive)
				core.E5,
				// Other diagonals unchanged
				core.E3, core.F2, core.G1,
				core.C5, core.B6, core.A7,
				core.C3, core.B2, core.A1,
			},
		},
		{
			name:       "friendly blocks orthogonal; diagonal still open",
			from:       core.D4,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.E4] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}}
				b[core.F4] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}}
			},
			expectedTos: []core.Position{
				// Orthogonal Right fully blocked by E4
				// Orthogonal other directions unchanged
				core.A4, core.B4, core.C4,
				core.D1, core.D2, core.D3, core.D5, core.D6, core.D7, core.D8,
				// Diagonals unchanged
				core.E5, core.F6, core.G7, core.H8,
				core.E3, core.F2, core.G1,
				core.C5, core.B6, core.A7,
				core.C3, core.B2, core.A1,
			},
		},
		{
			name:       "diagonal blocked; orthogonal still open",
			from:       core.D4,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.E5] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}}
				b[core.F6] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}}
			},
			expectedTos: []core.Position{
				// Orthogonal unchanged
				core.A4, core.B4, core.C4, core.E4, core.F4, core.G4, core.H4,
				core.D1, core.D2, core.D3, core.D5, core.D6, core.D7, core.D8,
				// Diagonal NE fully blocked by E5
				// Other diagonals unchanged
				core.E3, core.F2, core.G1,
				core.C5, core.B6, core.A7,
				core.C3, core.B2, core.A1,
			},
		},
		{
			name:       "enemy capturable on diagonal; friendly behind is unreachable",
			from:       core.D4,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.F6] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}}
				b[core.H8] = core.Square{Occupied: true, Piece: core.Piece{Type: core.ROOK, Color: core.WHITE}}
			},
			expectedTos: []core.Position{
				// Diagonal NE stops at F6 (inclusive) — H8 unreachable
				core.E5, core.F6,
				// Orthogonal unchanged
				core.A4, core.B4, core.C4, core.E4, core.F4, core.G4, core.H4,
				core.D1, core.D2, core.D3, core.D5, core.D6, core.D7, core.D8,
				// Other diagonals unchanged
				core.E3, core.F2, core.G1,
				core.C5, core.B6, core.A7,
				core.C3, core.B2, core.A1,
			},
		},
		{
			name:       "enemy capturable on orthogonal; friendly behind is unreachable",
			from:       core.D4,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.F4] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}}
				b[core.H4] = core.Square{Occupied: true, Piece: core.Piece{Type: core.ROOK, Color: core.WHITE}}
			},
			expectedTos: []core.Position{
				// Orthogonal Right stops at F4 (inclusive) — H4 unreachable
				core.E4, core.F4,
				// Orthogonal other directions unchanged
				core.A4, core.B4, core.C4,
				core.D1, core.D2, core.D3, core.D5, core.D6, core.D7, core.D8,
				// Diagonals unchanged
				core.E5, core.F6, core.G7, core.H8,
				core.E3, core.F2, core.G1,
				core.C5, core.B6, core.A7,
				core.C3, core.B2, core.A1,
			},
		},
		{
			name:       "all 8 directions blocked by own pieces yields no moves",
			from:       core.D4,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.D5] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}}
				b[core.D3] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}}
				b[core.E4] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}}
				b[core.C4] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}}
				b[core.E5] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}}
				b[core.E3] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}}
				b[core.C5] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}}
				b[core.C3] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}}
			},
			expectedTos: []core.Position{},
		},
		{
			name:       "all 8 directions blocked by enemy pieces (all capturable)",
			from:       core.D4,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.D5] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}}
				b[core.D3] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}}
				b[core.E4] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}}
				b[core.C4] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}}
				b[core.E5] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}}
				b[core.E3] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}}
				b[core.C5] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}}
				b[core.C3] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}}
			},
			expectedTos: []core.Position{core.D5, core.D3, core.E4, core.C4, core.E5, core.E3, core.C5, core.C3},
		},
		{
			name:       "black queen treats white piece as enemy (included)",
			from:       core.D4,
			sideToMove: core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.F4] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}}
				b[core.F6] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}}
			},
			expectedTos: []core.Position{
				// Orthogonal Right stops at F4 (inclusive)
				core.E4, core.F4,
				// Orthogonal other directions unchanged
				core.A4, core.B4, core.C4,
				core.D1, core.D2, core.D3, core.D5, core.D6, core.D7, core.D8,
				// Diagonal NE stops at F6 (inclusive)
				core.E5, core.F6,
				// Other diagonals unchanged
				core.E3, core.F2, core.G1,
				core.C5, core.B6, core.A7,
				core.C3, core.B2, core.A1,
			},
		},
		{
			name:       "black queen treats black piece as own (excluded)",
			from:       core.D4,
			sideToMove: core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.F4] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}}
				b[core.F6] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}}
			},
			expectedTos: []core.Position{
				// Orthogonal Right stops before F4 (exclusive)
				core.E4,
				// Orthogonal other directions unchanged
				core.A4, core.B4, core.C4,
				core.D1, core.D2, core.D3, core.D5, core.D6, core.D7, core.D8,
				// Diagonal NE stops before F6 (exclusive)
				core.E5,
				// Other diagonals unchanged
				core.E3, core.F2, core.G1,
				core.C5, core.B6, core.A7,
				core.C3, core.B2, core.A1,
			},
		},
		{
			name:       "captures carry the exact enemy piece sitting on destination",
			from:       core.D4,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				// Different enemy piece types — orthogonal captures
				b[core.D6] = core.Square{Occupied: true, Piece: core.Piece{Type: core.QUEEN, Color: core.BLACK}}  // Up
				b[core.D2] = core.Square{Occupied: true, Piece: core.Piece{Type: core.ROOK, Color: core.BLACK}}   // Down
				b[core.F4] = core.Square{Occupied: true, Piece: core.Piece{Type: core.BISHOP, Color: core.BLACK}} // Right
				b[core.B4] = core.Square{Occupied: true, Piece: core.Piece{Type: core.KNIGHT, Color: core.BLACK}} // Left
				// Diagonal captures
				b[core.F6] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}} // NE
				b[core.F2] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}} // SE
				b[core.B6] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}} // NW
				b[core.B2] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}} // SW
			},
			expectedTos: []core.Position{
				// Up: D5, D6 (capture)
				core.D5, core.D6,
				// Down: D3, D2 (capture)
				core.D3, core.D2,
				// Right: E4, F4 (capture)
				core.E4, core.F4,
				// Left: C4, B4 (capture)
				core.C4, core.B4,
				// NE: E5, F6 (capture)
				core.E5, core.F6,
				// SE: E3, F2 (capture)
				core.E3, core.F2,
				// NW: C5, B6 (capture)
				core.C5, core.B6,
				// SW: C3, B2 (capture)
				core.C3, core.B2,
			},
		},
		{
			name:       "queen on edge A4 — orthogonal + 2 diagonals only",
			from:       core.A4,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {},
			expectedTos: []core.Position{
				// Orthogonal (14)
				core.B4, core.C4, core.D4, core.E4, core.F4, core.G4, core.H4,
				core.A1, core.A2, core.A3, core.A5, core.A6, core.A7, core.A8,
				// Diagonals (7) — NE + SE only
				core.B5, core.C6, core.D7, core.E8,
				core.B3, core.C2, core.D1,
			},
		},
		{
			name:       "queen on edge D1 — orthogonal + 2 diagonals only",
			from:       core.D1,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {},
			expectedTos: []core.Position{
				// Orthogonal (14)
				core.A1, core.B1, core.C1, core.E1, core.F1, core.G1, core.H1,
				core.D2, core.D3, core.D4, core.D5, core.D6, core.D7, core.D8,
				// Diagonals (7) — NE + NW only
				core.E2, core.F3, core.G4, core.H5,
				core.C2, core.B3, core.A4,
			},
		},
	}

	queen := Queen{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var board core.Board
			tt.setupBoard(&board)

			ctx := core.MoveContext{
				BoardContext: core.BoardContext{Board: &board},
				SideToMove:   tt.sideToMove,
			}

			got := queen.PseudoLegalMoves(tt.from, ctx)

			if len(got) != len(tt.expectedTos) {
				t.Fatalf("got %d moves, want %d — got %v", len(got), len(tt.expectedTos), got)
			}

			expectedMover := core.Piece{Type: core.QUEEN, Color: tt.sideToMove}

			// build expected captures map: destination -> captured piece
			expectedCaptures := map[core.Position]core.Piece{}
			for _, pos := range tt.expectedTos {
				sq := board[pos]
				if sq.Occupied && sq.Piece.Color != tt.sideToMove {
					expectedCaptures[pos] = sq.Piece
				}
			}

			gotSet := make(map[core.Position]struct{}, len(got))
			for _, move := range got {
				// every queen move is NORMAL
				if move.Type != core.NORMAL {
					t.Errorf("expected move type NORMAL, got %v", move.Type)
				}
				// every move must carry its source
				if move.From != tt.from {
					t.Errorf("move to %v: From=%v, want %v", move.To, move.From, tt.from)
				}
				// every move must carry the moving piece
				if move.Piece != expectedMover {
					t.Errorf("move to %v: Piece=%v, want %v", move.To, move.Piece, expectedMover)
				}

				// capture info must match what's on the destination square
				wantCapture, shouldCapture := expectedCaptures[move.To]
				if shouldCapture {
					if !move.HasCapture {
						t.Errorf("move to %v: HasCapture=false, want true", move.To)
					} else if move.Captured != wantCapture {
						t.Errorf("move to %v: Captured=%v, want %v", move.To, move.Captured, wantCapture)
					}
				} else {
					if move.HasCapture {
						t.Errorf("move to %v: HasCapture=true, want false", move.To)
					}
				}

				gotSet[move.To] = struct{}{}
			}

			for _, pos := range tt.expectedTos {
				if _, ok := gotSet[pos]; !ok {
					t.Errorf("missing expected position %v", pos)
				}
			}
		})
	}
}
