package piece

import (
	"testing"

	"github.com/7-Dany/chess/core"
)

func TestRookIsAttacking(t *testing.T) {
	tests := []struct {
		name       string
		setupBoard func(*core.Board)
		color      core.PieceColor
		target     core.Position
		want       bool
	}{
		// ==================== Direct attacks (all 4 orthogonal directions) ====================
		{
			name: "rook above target on same file attacks",
			setupBoard: func(b *core.Board) {
				b[core.E8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true,
		},
		{
			name: "rook below target on same file attacks",
			setupBoard: func(b *core.Board) {
				b[core.E1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true,
		},
		{
			name: "rook to the right on same rank attacks",
			setupBoard: func(b *core.Board) {
				b[core.H4] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true,
		},
		{
			name: "rook to the left on same rank attacks",
			setupBoard: func(b *core.Board) {
				b[core.A4] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true,
		},
		{
			name: "rook adjacent (distance 1) attacks",
			setupBoard: func(b *core.Board) {
				b[core.E5] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true,
		},
		{
			name: "rook at maximum distance attacks (distance 7)",
			setupBoard: func(b *core.Board) {
				b[core.E8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.E1,
			want:   true,
		},

		// ==================== Not on attack line ====================
		{
			name: "rook on diagonal does not attack",
			setupBoard: func(b *core.Board) {
				b[core.H7] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "rook off-file and off-rank does not attack",
			setupBoard: func(b *core.Board) {
				b[core.F5] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "rook two files and one rank away (knight-move) does not attack",
			setupBoard: func(b *core.Board) {
				b[core.G5] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "rook one file and one rank away (diagonal) does not attack",
			setupBoard: func(b *core.Board) {
				b[core.D5] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},

		// ==================== Blockers ====================
		{
			name: "friendly piece blocks attack on file",
			setupBoard: func(b *core.Board) {
				b[core.E8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
				b[core.E6] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "enemy piece blocks attack on file",
			setupBoard: func(b *core.Board) {
				b[core.E8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
				b[core.E6] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "friendly piece blocks attack on rank",
			setupBoard: func(b *core.Board) {
				b[core.H4] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
				b[core.F4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "enemy piece blocks attack on rank",
			setupBoard: func(b *core.Board) {
				b[core.H4] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
				b[core.F4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "blocker adjacent to target blocks attack",
			setupBoard: func(b *core.Board) {
				b[core.E8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
				b[core.E5] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true, // E5 rook attacks (closer one), E8 is blocked by E5
		},
		{
			name: "blocker adjacent to rook blocks attack",
			setupBoard: func(b *core.Board) {
				b[core.E8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
				b[core.E7] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "piece behind target does not block (target between rook and piece)",
			setupBoard: func(b *core.Board) {
				b[core.E8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
				b[core.E2] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true,
		},
		{
			name: "two rooks on same file, closer one blocks farther one",
			setupBoard: func(b *core.Board) {
				b[core.E8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
				b[core.E6] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true, // E6 rook attacks (closer one)
		},

		// ==================== Color filtering ====================
		{
			name: "black rook, asking for white attackers",
			setupBoard: func(b *core.Board) {
				b[core.E8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "black rook, asking for black attackers",
			setupBoard: func(b *core.Board) {
				b[core.E8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
			},
			color:  core.BLACK,
			target: core.E4,
			want:   true,
		},

		// ==================== Non-rook pieces on attack line (clean separation) ====================
		{
			name: "queen on file does NOT trigger rook attack",
			setupBoard: func(b *core.Board) {
				b[core.E8] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "queen on rank does NOT trigger rook attack",
			setupBoard: func(b *core.Board) {
				b[core.H4] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "bishop on file does not trigger rook attack",
			setupBoard: func(b *core.Board) {
				b[core.E8] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "knight on file does not trigger rook attack",
			setupBoard: func(b *core.Board) {
				b[core.E8] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "king on file does not trigger rook attack",
			setupBoard: func(b *core.Board) {
				b[core.E5] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "pawn on file does not trigger rook attack",
			setupBoard: func(b *core.Board) {
				b[core.E5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
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
			name: "rook on target square itself does not attack",
			setupBoard: func(b *core.Board) {
				b[core.E4] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},

		// ==================== Corners ====================
		{
			name: "target on A1, rook on A8 attacks on file",
			setupBoard: func(b *core.Board) {
				b[core.A8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.A1,
			want:   true,
		},
		{
			name: "target on A1, rook on H1 attacks on rank",
			setupBoard: func(b *core.Board) {
				b[core.H1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.A1,
			want:   true,
		},
		{
			name: "target on H8, rook on H1 attacks on file",
			setupBoard: func(b *core.Board) {
				b[core.H1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.H8,
			want:   true,
		},
		{
			name: "target on H8, rook on A8 attacks on rank",
			setupBoard: func(b *core.Board) {
				b[core.A8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.H8,
			want:   true,
		},
		{
			name: "target on A8, rook on A1 attacks on file",
			setupBoard: func(b *core.Board) {
				b[core.A1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.A8,
			want:   true,
		},
		{
			name: "target on H1, rook on H8 attacks on file",
			setupBoard: func(b *core.Board) {
				b[core.H8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.H1,
			want:   true,
		},
		{
			name: "target on A1, rook on H8 (diagonal) does not attack",
			setupBoard: func(b *core.Board) {
				b[core.H8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.A1,
			want:   false,
		},

		// ==================== Edge: target on edge file/rank ====================
		{
			name: "target on A4, rook on A8 attacks on file",
			setupBoard: func(b *core.Board) {
				b[core.A8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.A4,
			want:   true,
		},
		{
			name: "target on A4, rook on H4 attacks on rank",
			setupBoard: func(b *core.Board) {
				b[core.H4] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.A4,
			want:   true,
		},

		// ==================== Multiple rooks ====================
		{
			name: "multiple rooks, one attacks on file",
			setupBoard: func(b *core.Board) {
				b[core.A1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
				b[core.E8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
				b[core.H1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true, // E8 attacks via file
		},
		{
			name: "multiple rooks, one attacks on rank",
			setupBoard: func(b *core.Board) {
				b[core.A1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
				b[core.H4] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
				b[core.A8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true, // H4 attacks via rank
		},
		{
			name: "multiple enemy rooks, none matching color",
			setupBoard: func(b *core.Board) {
				b[core.E8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
				b[core.H4] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "mixed-color rooks, only matching color counts",
			setupBoard: func(b *core.Board) {
				b[core.E8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
				b[core.H4] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true, // H4 (white) attacks
		},

		// ==================== Blocker is rook of wrong color ====================
		{
			name: "enemy rook on ray blocks but doesn't count as attacker",
			setupBoard: func(b *core.Board) {
				b[core.E8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
				b[core.E6] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false, // black rook at E6 blocks the white rook at E8
		},
		{
			name: "enemy rook as blocker, asking for black",
			setupBoard: func(b *core.Board) {
				b[core.E8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
				b[core.E6] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
			},
			color:  core.BLACK,
			target: core.E4,
			want:   true, // black rook at E6 attacks E4
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var board core.Board
			tt.setupBoard(&board)
			ctx := core.BoardContext{Board: &board}

			rook := Rook{}
			got := rook.IsAttacking(tt.color, tt.target, ctx)

			if got != tt.want {
				t.Errorf("IsAttacking(%v, %v) = %v, want %v", tt.color, tt.target, got, tt.want)
			}
		})
	}
}

// test all attack moves for rook.
func TestRookAttacks(t *testing.T) {
	tests := []struct {
		name       string
		from       core.Position
		setupBoard func(*core.Board)
		expected   []core.Position
	}{
		{
			name:       "center D4 empty board",
			from:       core.D4,
			setupBoard: func(b *core.Board) {},
			expected: []core.Position{
				// Horizontal
				core.A4, core.B4, core.C4, core.E4, core.F4, core.G4, core.H4,
				// Vertical
				core.D1, core.D2, core.D3, core.D5, core.D6, core.D7, core.D8,
			},
		},
		{
			name:       "corner A1 empty board",
			from:       core.A1,
			setupBoard: func(b *core.Board) {},
			expected: []core.Position{
				// Horizontal
				core.B1, core.C1, core.D1, core.E1, core.F1, core.G1, core.H1,
				// Vertical
				core.A2, core.A3, core.A4, core.A5, core.A6, core.A7, core.A8,
			},
		},
		{
			name:       "corner H1 empty board",
			from:       core.H1,
			setupBoard: func(b *core.Board) {},
			expected: []core.Position{
				core.A1, core.B1, core.C1, core.D1, core.E1, core.F1, core.G1,
				core.H2, core.H3, core.H4, core.H5, core.H6, core.H7, core.H8,
			},
		},
		{
			name:       "corner A8 empty board",
			from:       core.A8,
			setupBoard: func(b *core.Board) {},
			expected: []core.Position{
				core.B8, core.C8, core.D8, core.E8, core.F8, core.G8, core.H8,
				core.A1, core.A2, core.A3, core.A4, core.A5, core.A6, core.A7,
			},
		},
		{
			name:       "corner H8 empty board",
			from:       core.H8,
			setupBoard: func(b *core.Board) {},
			expected: []core.Position{
				core.A8, core.B8, core.C8, core.D8, core.E8, core.F8, core.G8,
				core.H1, core.H2, core.H3, core.H4, core.H5, core.H6, core.H7,
			},
		},
		{
			name:       "edge A4 empty board",
			from:       core.A4,
			setupBoard: func(b *core.Board) {},
			expected: []core.Position{
				// Horizontal (only rightward — A is the leftmost file)
				core.B4, core.C4, core.D4, core.E4, core.F4, core.G4, core.H4,
				// Vertical
				core.A1, core.A2, core.A3, core.A5, core.A6, core.A7, core.A8,
			},
		},
		{
			name:       "edge D1 empty board",
			from:       core.D1,
			setupBoard: func(b *core.Board) {},
			expected: []core.Position{
				// Horizontal
				core.A1, core.B1, core.C1, core.E1, core.F1, core.G1, core.H1,
				// Vertical (only upward — 1 is the bottom rank)
				core.D2, core.D3, core.D4, core.D5, core.D6, core.D7, core.D8,
			},
		},
		{
			name: "center D4 blocked on Up at D6",
			from: core.D4,
			setupBoard: func(b *core.Board) {
				b[core.D6] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
			},
			expected: []core.Position{
				// Horizontal (unchanged)
				core.A4, core.B4, core.C4, core.E4, core.F4, core.G4, core.H4,
				// Vertical — Up stops at D6 (inclusive), Down unchanged
				core.D1, core.D2, core.D3, core.D5, core.D6,
			},
		},
		{
			name: "center D4 blocked on all four directions",
			from: core.D4,
			setupBoard: func(b *core.Board) {
				b[core.D5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
				b[core.D3] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
				b[core.E4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
				b[core.C4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
			},
			expected: []core.Position{core.D5, core.D3, core.E4, core.C4},
		},
		{
			name: "corner A1 blocked at C1 and A3",
			from: core.A1,
			setupBoard: func(b *core.Board) {
				b[core.C1] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
				b[core.A3] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
			},
			expected: []core.Position{
				// Horizontal stops at C1
				core.B1, core.C1,
				// Vertical stops at A3
				core.A2, core.A3,
			},
		},
	}

	rook := Rook{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var board core.Board
			tt.setupBoard(&board)
			ctx := core.BoardContext{Board: &board}

			got := rook.Attacks(tt.from, ctx)

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

// test all rook pseudolegal moves
func TestRookPseudoLegalMoves(t *testing.T) {
	tests := []struct {
		name        string
		from        core.Position
		sideToMove  core.PieceColor
		setupBoard  func(*core.Board)
		expectedTos []core.Position
	}{
		{
			name:       "empty board center D4",
			from:       core.D4,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {},
			expectedTos: []core.Position{
				core.A4, core.B4, core.C4, core.E4, core.F4, core.G4, core.H4,
				core.D1, core.D2, core.D3, core.D5, core.D6, core.D7, core.D8,
			},
		},
		{
			name:        "empty board corner A1",
			from:        core.A1,
			sideToMove:  core.WHITE,
			setupBoard:  func(b *core.Board) {},
			expectedTos: []core.Position{core.B1, core.C1, core.D1, core.E1, core.F1, core.G1, core.H1, core.A2, core.A3, core.A4, core.A5, core.A6, core.A7, core.A8},
		},
		{
			name:        "empty board corner H8",
			from:        core.H8,
			sideToMove:  core.WHITE,
			setupBoard:  func(b *core.Board) {},
			expectedTos: []core.Position{core.A8, core.B8, core.C8, core.D8, core.E8, core.F8, core.G8, core.H1, core.H2, core.H3, core.H4, core.H5, core.H6, core.H7},
		},
		{
			name:       "can capture enemy piece (stops ray)",
			from:       core.D4,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.F4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
			},
			expectedTos: []core.Position{
				// Horizontal right stops at F4 (inclusive)
				core.E4, core.F4,
				// Horizontal left unchanged
				core.A4, core.B4, core.C4,
				// Vertical unchanged
				core.D1, core.D2, core.D3, core.D5, core.D6, core.D7, core.D8,
			},
		},
		{
			name:       "cannot move to friendly occupied square",
			from:       core.D4,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.F4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
			},
			expectedTos: []core.Position{
				// Horizontal right stops BEFORE F4 (exclusive)
				core.E4,
				// Horizontal left unchanged
				core.A4, core.B4, core.C4,
				// Vertical unchanged
				core.D1, core.D2, core.D3, core.D5, core.D6, core.D7, core.D8,
			},
		},
		{
			name:       "friendly blocks slide path; enemy behind is unreachable",
			from:       core.D4,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.E4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
				b[core.F4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
			},
			expectedTos: []core.Position{
				// Horizontal right stops before E4 (exclusive)
				// Horizontal left unchanged
				core.A4, core.B4, core.C4,
				// Vertical unchanged
				core.D1, core.D2, core.D3, core.D5, core.D6, core.D7, core.D8,
			},
		},
		{
			name:       "enemy blocks slide path but is capturable",
			from:       core.D4,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.F4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
				b[core.H4] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
			},
			expectedTos: []core.Position{
				// Horizontal right stops at F4 (inclusive) — H4 unreachable
				core.E4, core.F4,
				// Horizontal left unchanged
				core.A4, core.B4, core.C4,
				// Vertical unchanged
				core.D1, core.D2, core.D3, core.D5, core.D6, core.D7, core.D8,
			},
		},
		{
			name:       "enemy capturable; friendly behind is unreachable",
			from:       core.D4,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.F4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
				b[core.H4] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			expectedTos: []core.Position{
				// Same as "enemy blocks slide path" — H4 unreachable regardless of color
				core.E4, core.F4,
				core.A4, core.B4, core.C4,
				core.D1, core.D2, core.D3, core.D5, core.D6, core.D7, core.D8,
			},
		},
		{
			name:       "mixed blockers in all four directions",
			from:       core.D4,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				// Up: enemy at D6 (capturable, stops ray)
				b[core.D6] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
				// Down: friendly at D2 (blocks ray, not capturable)
				b[core.D2] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
				// Right: enemy at F4 (capturable, stops ray)
				b[core.F4] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.BLACK})
				// Left: friendly at B4 (blocks ray, not capturable)
				b[core.B4] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
			},
			expectedTos: []core.Position{
				// Up: D5, D6 (inclusive of capturable)
				core.D5, core.D6,
				// Down: D3 only (D2 blocks, exclusive)
				core.D3,
				// Right: E4, F4 (inclusive)
				core.E4, core.F4,
				// Left: C4 only (B4 blocks, exclusive)
				core.C4,
			},
		},
		{
			name:       "all directions blocked by own pieces yields no moves",
			from:       core.D4,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.D5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
				b[core.D3] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
				b[core.E4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
				b[core.C4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
			},
			expectedTos: []core.Position{},
		},
		{
			name:       "all four directions blocked by enemy pieces (all capturable)",
			from:       core.D4,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.D5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
				b[core.D3] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
				b[core.E4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
				b[core.C4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
			},
			expectedTos: []core.Position{core.D5, core.D3, core.E4, core.C4},
		},
		{
			name:       "black rook treats white piece as enemy (included)",
			from:       core.D4,
			sideToMove: core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.F4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
			},
			expectedTos: []core.Position{
				core.E4, core.F4,
				core.A4, core.B4, core.C4,
				core.D1, core.D2, core.D3, core.D5, core.D6, core.D7, core.D8,
			},
		},
		{
			name:       "black rook treats black piece as own (excluded)",
			from:       core.D4,
			sideToMove: core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.F4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
			},
			expectedTos: []core.Position{
				core.E4,
				core.A4, core.B4, core.C4,
				core.D1, core.D2, core.D3, core.D5, core.D6, core.D7, core.D8,
			},
		},
		{
			name:       "captures carry the exact enemy piece sitting on destination",
			from:       core.D4,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				// Different enemy piece types in each direction
				b[core.D6] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.BLACK})  // Up
				b[core.D2] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})   // Down
				b[core.F4] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.BLACK}) // Right
				b[core.B4] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.BLACK}) // Left
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
			},
		},
		{
			name:       "rook on edge A4 — only rightward horizontal moves",
			from:       core.A4,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {},
			expectedTos: []core.Position{
				core.B4, core.C4, core.D4, core.E4, core.F4, core.G4, core.H4,
				core.A1, core.A2, core.A3, core.A5, core.A6, core.A7, core.A8,
			},
		},
		{
			name:       "rook on edge D1 — only upward vertical moves",
			from:       core.D1,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {},
			expectedTos: []core.Position{
				core.A1, core.B1, core.C1, core.E1, core.F1, core.G1, core.H1,
				core.D2, core.D3, core.D4, core.D5, core.D6, core.D7, core.D8,
			},
		},
	}

	rook := Rook{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var board core.Board
			tt.setupBoard(&board)

			ctx := core.MoveContext{
				BoardContext: core.BoardContext{Board: &board},
				SideToMove:   tt.sideToMove,
			}

			got := rook.PseudoLegalMoves(tt.from, ctx)

			if len(got) != len(tt.expectedTos) {
				t.Fatalf("got %d moves, want %d — got %v", len(got), len(tt.expectedTos), got)
			}

			expectedMover := core.Piece{Type: core.ROOK, Color: tt.sideToMove}

			// build expected captures map: destination -> captured piece
			expectedCaptures := map[core.Position]core.Piece{}
			for _, pos := range tt.expectedTos {
				sq := board[pos]
				if sq.IsOccupied() && sq.Color() != tt.sideToMove {
					expectedCaptures[pos] = sq.Piece()
				}
			}

			gotSet := make(map[core.Position]struct{}, len(got))
			for _, move := range got {
				// every rook move is NORMAL
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
