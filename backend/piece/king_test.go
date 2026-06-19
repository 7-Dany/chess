package piece

import (
	"testing"

	"github.com/7-Dany/chess/core"
)

func TestKingIsAttacking(t *testing.T) {
	tests := []struct {
		name       string
		setupBoard func(*core.Board)
		color      core.PieceColor
		target     core.Position
		want       bool
	}{
		// ==================== Direct attacks ====================
		{
			name: "king attacks target above (same file)",
			setupBoard: func(b *core.Board) {
				b[core.E5] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true,
		},
		{
			name: "king attacks target below (same file)",
			setupBoard: func(b *core.Board) {
				b[core.E3] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true,
		},
		{
			name: "king attacks target to the right (same rank)",
			setupBoard: func(b *core.Board) {
				b[core.F4] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true,
		},
		{
			name: "king attacks target to the left (same rank)",
			setupBoard: func(b *core.Board) {
				b[core.D4] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true,
		},
		{
			name: "king attacks target up-right diagonal",
			setupBoard: func(b *core.Board) {
				b[core.F5] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true,
		},
		{
			name: "king attacks target up-left diagonal",
			setupBoard: func(b *core.Board) {
				b[core.D5] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true,
		},
		{
			name: "king attacks target down-right diagonal",
			setupBoard: func(b *core.Board) {
				b[core.F3] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true,
		},
		{
			name: "king attacks target down-left diagonal",
			setupBoard: func(b *core.Board) {
				b[core.D3] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true,
		},

		// ==================== Not adjacent ====================
		{
			name: "king two squares away on same file does not attack",
			setupBoard: func(b *core.Board) {
				b[core.E6] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "king two squares away on same rank does not attack",
			setupBoard: func(b *core.Board) {
				b[core.G4] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "king two squares away diagonally does not attack",
			setupBoard: func(b *core.Board) {
				b[core.G6] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "knight-move away does not attack",
			setupBoard: func(b *core.Board) {
				b[core.D6] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},

		// ==================== Color filtering ====================
		{
			name: "black king, asking for white attackers",
			setupBoard: func(b *core.Board) {
				b[core.E5] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "black king, asking for black attackers",
			setupBoard: func(b *core.Board) {
				b[core.E5] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
			},
			color:  core.BLACK,
			target: core.E4,
			want:   true,
		},

		// ==================== Non-king pieces adjacent ====================
		{
			name: "queen adjacent does NOT trigger king attack (clean separation)",
			setupBoard: func(b *core.Board) {
				b[core.E5] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "rook adjacent does not trigger king attack",
			setupBoard: func(b *core.Board) {
				b[core.E5] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "bishop adjacent does not trigger king attack",
			setupBoard: func(b *core.Board) {
				b[core.F5] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "knight adjacent does not trigger king attack",
			setupBoard: func(b *core.Board) {
				b[core.E5] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "pawn adjacent does not trigger king attack",
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
			name: "king on target square itself does not attack",
			setupBoard: func(b *core.Board) {
				b[core.E4] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},

		// ==================== Corners (only 3 valid adjacencies) ====================
		{
			name: "target on A1, king attacks from B2",
			setupBoard: func(b *core.Board) {
				b[core.B2] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.A1,
			want:   true,
		},
		{
			name: "target on A1, king attacks from A2",
			setupBoard: func(b *core.Board) {
				b[core.A2] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.A1,
			want:   true,
		},
		{
			name: "target on A1, king attacks from B1",
			setupBoard: func(b *core.Board) {
				b[core.B1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.A1,
			want:   true,
		},
		{
			name: "target on A1, king on D4 does not attack",
			setupBoard: func(b *core.Board) {
				b[core.D4] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.A1,
			want:   false,
		},
		{
			name: "target on H8, king attacks from G7",
			setupBoard: func(b *core.Board) {
				b[core.G7] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.H8,
			want:   true,
		},
		{
			name: "target on H1, king attacks from G2",
			setupBoard: func(b *core.Board) {
				b[core.G2] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.H1,
			want:   true,
		},
		{
			name: "target on A8, king attacks from B7",
			setupBoard: func(b *core.Board) {
				b[core.B7] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.A8,
			want:   true,
		},

		// ==================== Edge: target next to corner ====================
		{
			name: "target on B1, king attacks from A2",
			setupBoard: func(b *core.Board) {
				b[core.A2] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.B1,
			want:   true,
		},
		{
			name: "target on B1, king attacks from C2",
			setupBoard: func(b *core.Board) {
				b[core.C2] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.B1,
			want:   true,
		},
		{
			name: "target on B1, king on A1 does not attack (same square impossible, but verifies scan)",
			setupBoard: func(b *core.Board) {
				b[core.A1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.B1,
			want:   true,
		},

		// ==================== Multiple kings (impossible but tests iteration) ====================
		{
			name: "multiple friendly pieces adjacent, one is king",
			setupBoard: func(b *core.Board) {
				b[core.D3] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
				b[core.D4] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
				b[core.D5] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
				b[core.E5] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true, // D4 king attacks
		},
		{
			name: "enemy pieces adjacent, none matching",
			setupBoard: func(b *core.Board) {
				b[core.D3] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
				b[core.D4] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
				b[core.D5] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.BLACK})
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false, // All enemies; we want white attackers
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var board core.Board
			tt.setupBoard(&board)
			ctx := core.BoardContext{Board: &board}

			king := King{}
			got := king.IsAttacking(tt.color, tt.target, ctx)

			if got != tt.want {
				t.Errorf("IsAttacking(%v, %v) = %v, want %v", tt.color, tt.target, got, tt.want)
			}
		})
	}
}

// test all attack moves for king.
func TestKingAttacks(t *testing.T) {
	tests := []struct {
		name     string
		from     core.Position
		expected []core.Position
	}{
		{
			name: "center D4 has 8 attacks",
			from: core.D4,
			expected: []core.Position{
				core.D5, core.D3,
				core.E4, core.C4,
				core.E5, core.E3,
				core.C5, core.C3,
			},
		},
		{
			name: "corner A1 has 3 attacks",
			from: core.A1,
			expected: []core.Position{
				core.A2, core.B1, core.B2,
			},
		},
		{
			name: "corner H1 has 3 attacks",
			from: core.H1,
			expected: []core.Position{
				core.H2, core.G1, core.G2,
			},
		},
		{
			name: "corner A8 has 3 attacks",
			from: core.A8,
			expected: []core.Position{
				core.A7, core.B8, core.B7,
			},
		},
		{
			name: "corner H8 has 3 attacks",
			from: core.H8,
			expected: []core.Position{
				core.H7, core.G8, core.G7,
			},
		},
		{
			name: "edge A4 has 5 attacks",
			from: core.A4,
			expected: []core.Position{
				core.A5, core.A3,
				core.B4, core.B5, core.B3,
			},
		},
		{
			name: "edge H4 has 5 attacks",
			from: core.H4,
			expected: []core.Position{
				core.H5, core.H3,
				core.G4, core.G5, core.G3,
			},
		},
		{
			name: "edge D1 has 5 attacks",
			from: core.D1,
			expected: []core.Position{
				core.D2,
				core.E1, core.C1,
				core.E2, core.C2,
			},
		},
		{
			name: "edge D8 has 5 attacks",
			from: core.D8,
			expected: []core.Position{
				core.D7,
				core.E8, core.C8,
				core.E7, core.C7,
			},
		},
	}

	king := King{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var board core.Board
			got := king.Attacks(tt.from, core.BoardContext{Board: &board})

			if len(got) != len(tt.expected) {
				t.Fatalf("got %d attacks, want %d", len(got), len(tt.expected))
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

// test all king pseudolegal moves
func TestKingPseudoLegalMoves(t *testing.T) {
	tests := []struct {
		name        string
		from        core.Position
		sideToMove  core.PieceColor
		setupBoard  func(*core.Board)
		expectedTos []core.Position
	}{
		{
			name:       "empty board center D4 has 8 moves",
			from:       core.D4,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {},
			expectedTos: []core.Position{
				core.D5, core.D3,
				core.E4, core.C4,
				core.E5, core.E3,
				core.C5, core.C3,
			},
		},
		{
			name:        "empty board corner A1 has 3 moves",
			from:        core.A1,
			sideToMove:  core.WHITE,
			setupBoard:  func(b *core.Board) {},
			expectedTos: []core.Position{core.A2, core.B1, core.B2},
		},
		{
			name:        "empty board corner H8 has 3 moves",
			from:        core.H8,
			sideToMove:  core.WHITE,
			setupBoard:  func(b *core.Board) {},
			expectedTos: []core.Position{core.H7, core.G8, core.G7},
		},
		{
			name:        "empty board edge A4 has 5 moves",
			from:        core.A4,
			sideToMove:  core.WHITE,
			setupBoard:  func(b *core.Board) {},
			expectedTos: []core.Position{core.A5, core.A3, core.B4, core.B5, core.B3},
		},
		{
			name:       "own piece on attack square is excluded",
			from:       core.D4,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.D5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
			},
			expectedTos: []core.Position{
				core.D3,
				core.E4, core.C4,
				core.E5, core.E3,
				core.C5, core.C3,
			},
		},
		{
			name:       "enemy piece on attack square is included (capture)",
			from:       core.D4,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.D5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
			},
			expectedTos: []core.Position{
				core.D5, core.D3,
				core.E4, core.C4,
				core.E5, core.E3,
				core.C5, core.C3,
			},
		},
		{
			name:       "captures carry the exact enemy piece sitting on destination",
			from:       core.D4,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.D5] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.BLACK})
				b[core.E5] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
			},
			expectedTos: []core.Position{
				core.D5, core.D3,
				core.E4, core.C4,
				core.E5, core.E3,
				core.C5, core.C3,
			},
		},
		{
			name:       "mix of own and enemy blockers",
			from:       core.D4,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.D5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
				b[core.C5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
				b[core.E5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
			},
			expectedTos: []core.Position{
				core.D3,
				core.E4, core.C4,
				core.E5, core.E3,
				core.C3,
			},
		},
		{
			name:       "all attack squares blocked by own pieces yields no moves",
			from:       core.D4,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				for _, pos := range []core.Position{core.D5, core.D3, core.E4, core.C4, core.E5, core.E3, core.C5, core.C3} {
					b[pos] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
				}
			},
			expectedTos: []core.Position{},
		},
		{
			name:       "black king treats white piece as enemy (included)",
			from:       core.D4,
			sideToMove: core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.D5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
			},
			expectedTos: []core.Position{
				core.D5, core.D3,
				core.E4, core.C4,
				core.E5, core.E3,
				core.C5, core.C3,
			},
		},
		{
			name:       "black king treats black piece as own (excluded)",
			from:       core.D4,
			sideToMove: core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.D5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
			},
			expectedTos: []core.Position{
				core.D3,
				core.E4, core.C4,
				core.E5, core.E3,
				core.C5, core.C3,
			},
		},
		{
			name:       "does not filter squares attacked by enemy pieces",
			from:       core.D4,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.D8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
			},
			expectedTos: []core.Position{
				core.D5, core.D3,
				core.E4, core.C4,
				core.E5, core.E3,
				core.C5, core.C3,
			},
		},
	}

	king := King{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var board core.Board
			tt.setupBoard(&board)

			ctx := core.MoveContext{
				BoardContext: core.BoardContext{Board: &board},
				SideToMove:   tt.sideToMove,
			}

			got := king.PseudoLegalMoves(tt.from, ctx)

			if len(got) != len(tt.expectedTos) {
				t.Fatalf("got %d moves, want %d", len(got), len(tt.expectedTos))
			}

			expectedMover := core.Piece{Type: core.KING, Color: tt.sideToMove}

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
				// every king move is NORMAL (castling lives in GetLegalMoves)
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
