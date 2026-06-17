package piece

import (
	"testing"

	"github.com/7-Dany/chess/core"
)

func TestKnightIsAttacking(t *testing.T) {
	tests := []struct {
		name       string
		setupBoard func(*core.Board)
		color      core.PieceColor
		target     core.Position
		want       bool
	}{
		// ==================== Direct attacks (all 8 L-shapes from target E4) ====================
		// From E4, knight attacks: C3, C5, D2, D6, F2, F6, G3, G5
		{
			name: "knight attacks from C3 (down-left L1)",
			setupBoard: func(b *core.Board) {
				b[core.C3] = core.Square{Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true,
		},
		{
			name: "knight attacks from C5 (up-left L2)",
			setupBoard: func(b *core.Board) {
				b[core.C5] = core.Square{Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true,
		},
		{
			name: "knight attacks from D2 (down-near L3)",
			setupBoard: func(b *core.Board) {
				b[core.D2] = core.Square{Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true,
		},
		{
			name: "knight attacks from D6 (up-near L4)",
			setupBoard: func(b *core.Board) {
				b[core.D6] = core.Square{Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true,
		},
		{
			name: "knight attacks from F2 (down-near L5)",
			setupBoard: func(b *core.Board) {
				b[core.F2] = core.Square{Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true,
		},
		{
			name: "knight attacks from F6 (up-near L6)",
			setupBoard: func(b *core.Board) {
				b[core.F6] = core.Square{Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true,
		},
		{
			name: "knight attacks from G3 (down-right L7)",
			setupBoard: func(b *core.Board) {
				b[core.G3] = core.Square{Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true,
		},
		{
			name: "knight attacks from G5 (up-right L8)",
			setupBoard: func(b *core.Board) {
				b[core.G5] = core.Square{Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true,
		},

		// ==================== Not on L-shape ====================
		{
			name: "knight adjacent (not L-shape) does not attack",
			setupBoard: func(b *core.Board) {
				b[core.D3] = core.Square{Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "knight on same file does not attack",
			setupBoard: func(b *core.Board) {
				b[core.E6] = core.Square{Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "knight on same rank does not attack",
			setupBoard: func(b *core.Board) {
				b[core.H4] = core.Square{Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "knight on diagonal does not attack",
			setupBoard: func(b *core.Board) {
				b[core.G6] = core.Square{Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "knight on long diagonal does not attack",
			setupBoard: func(b *core.Board) {
				b[core.H7] = core.Square{Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},

		// ==================== Color filtering ====================
		{
			name: "black knight on L-shape, asking for white",
			setupBoard: func(b *core.Board) {
				b[core.D6] = core.Square{Piece: core.Piece{Type: core.KNIGHT, Color: core.BLACK}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "black knight on L-shape, asking for black",
			setupBoard: func(b *core.Board) {
				b[core.D6] = core.Square{Piece: core.Piece{Type: core.KNIGHT, Color: core.BLACK}, Occupied: true}
			},
			color:  core.BLACK,
			target: core.E4,
			want:   true,
		},

		// ==================== Non-knight on L-shape (clean separation) ====================
		{
			name: "queen on L-shape does NOT trigger knight attack",
			setupBoard: func(b *core.Board) {
				b[core.D6] = core.Square{Piece: core.Piece{Type: core.QUEEN, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "rook on L-shape does not trigger knight attack",
			setupBoard: func(b *core.Board) {
				b[core.D6] = core.Square{Piece: core.Piece{Type: core.ROOK, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "bishop on L-shape does not trigger knight attack",
			setupBoard: func(b *core.Board) {
				b[core.D6] = core.Square{Piece: core.Piece{Type: core.BISHOP, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "king on L-shape does not trigger knight attack",
			setupBoard: func(b *core.Board) {
				b[core.D6] = core.Square{Piece: core.Piece{Type: core.KING, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "pawn on L-shape does not trigger knight attack",
			setupBoard: func(b *core.Board) {
				b[core.D6] = core.Square{Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}, Occupied: true}
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
			name: "knight on target square itself does not attack",
			setupBoard: func(b *core.Board) {
				b[core.E4] = core.Square{Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},

		// ==================== Corners ====================
		// Knight on A1 attacks B3 and C2 only.
		{
			name: "target on A1, knight attacks from B3",
			setupBoard: func(b *core.Board) {
				b[core.B3] = core.Square{Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.A1,
			want:   true,
		},
		{
			name: "target on A1, knight attacks from C2",
			setupBoard: func(b *core.Board) {
				b[core.C2] = core.Square{Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.A1,
			want:   true,
		},
		{
			name: "target on A1, knight on H8 does not attack",
			setupBoard: func(b *core.Board) {
				b[core.H8] = core.Square{Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.A1,
			want:   false,
		},
		{
			name: "target on H8, knight attacks from F7",
			setupBoard: func(b *core.Board) {
				b[core.F7] = core.Square{Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.H8,
			want:   true,
		},
		{
			name: "target on H8, knight attacks from G6",
			setupBoard: func(b *core.Board) {
				b[core.G6] = core.Square{Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.H8,
			want:   true,
		},
		{
			name: "target on A8, knight attacks from B6",
			setupBoard: func(b *core.Board) {
				b[core.B6] = core.Square{Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.A8,
			want:   true,
		},
		{
			name: "target on H1, knight attacks from F2",
			setupBoard: func(b *core.Board) {
				b[core.F2] = core.Square{Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.H1,
			want:   true,
		},

		// ==================== Edge: target near edge, fewer valid L-shapes ====================
		{
			name: "target on A4, knight attacks from B6 (only 4 valid L-shapes on edge)",
			setupBoard: func(b *core.Board) {
				b[core.B6] = core.Square{Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.A4,
			want:   true,
		},
		{
			name: "target on A4, knight attacks from C5",
			setupBoard: func(b *core.Board) {
				b[core.C5] = core.Square{Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.A4,
			want:   true,
		},
		{
			name: "target on A4, knight attacks from B2",
			setupBoard: func(b *core.Board) {
				b[core.B2] = core.Square{Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.A4,
			want:   true,
		},
		{
			name: "target on A4, knight attacks from C3",
			setupBoard: func(b *core.Board) {
				b[core.C3] = core.Square{Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.A4,
			want:   true,
		},

		// ==================== Multiple knights ====================
		{
			name: "multiple knights, one on L-shape attacks",
			setupBoard: func(b *core.Board) {
				b[core.A1] = core.Square{Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE}, Occupied: true}
				b[core.D6] = core.Square{Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE}, Occupied: true}
				b[core.H8] = core.Square{Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true, // D6 attacks
		},
		{
			name: "multiple enemy knights, none matching color",
			setupBoard: func(b *core.Board) {
				b[core.D6] = core.Square{Piece: core.Piece{Type: core.KNIGHT, Color: core.BLACK}, Occupied: true}
				b[core.F6] = core.Square{Piece: core.Piece{Type: core.KNIGHT, Color: core.BLACK}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "mixed-color knights, only matching color counts",
			setupBoard: func(b *core.Board) {
				b[core.D6] = core.Square{Piece: core.Piece{Type: core.KNIGHT, Color: core.BLACK}, Occupied: true}
				b[core.F6] = core.Square{Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true, // F6 (white) attacks
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var board core.Board
			tt.setupBoard(&board)
			ctx := core.BoardContext{Board: &board}

			knight := Knight{}
			got := knight.IsAttacking(tt.color, tt.target, ctx)

			if got != tt.want {
				t.Errorf("IsAttacking(%v, %v) = %v, want %v", tt.color, tt.target, got, tt.want)
			}
		})
	}
}

// test all attack moves for knight.
func TestKnightAttacks(t *testing.T) {
	tests := []struct {
		name     string
		from     core.Position
		expected []core.Position
	}{
		{
			name: "center D4 has 8 attacks",
			from: core.D4,
			expected: []core.Position{
				core.E6, core.E2,
				core.C6, core.C2,
				core.F5, core.F3,
				core.B5, core.B3,
			},
		},
		{
			name: "corner A1 has 2 attacks",
			from: core.A1,
			expected: []core.Position{
				core.B3, core.C2,
			},
		},
		{
			name: "corner H1 has 2 attacks",
			from: core.H1,
			expected: []core.Position{
				core.G3, core.F2,
			},
		},
		{
			name: "corner A8 has 2 attacks",
			from: core.A8,
			expected: []core.Position{
				core.B6, core.C7,
			},
		},
		{
			name: "corner H8 has 2 attacks",
			from: core.H8,
			expected: []core.Position{
				core.G6, core.F7,
			},
		},
		{
			name: "edge A4 has 4 attacks",
			from: core.A4,
			expected: []core.Position{
				core.B6, core.B2,
				core.C5, core.C3,
			},
		},
		{
			name: "edge H4 has 4 attacks",
			from: core.H4,
			expected: []core.Position{
				core.G6, core.G2,
				core.F5, core.F3,
			},
		},
		{
			name: "edge D1 has 4 attacks",
			from: core.D1,
			expected: []core.Position{
				core.E3, core.C3,
				core.F2, core.B2,
			},
		},
		{
			name: "edge D8 has 4 attacks",
			from: core.D8,
			expected: []core.Position{
				core.E6, core.C6,
				core.F7, core.B7,
			},
		},
		{
			name: "near-corner B2 has 4 attacks",
			from: core.B2,
			expected: []core.Position{
				core.D3, core.D1,
				core.A4, core.C4,
			},
		},
	}

	knight := Knight{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var board core.Board
			ctx := core.BoardContext{Board: &board}

			got := knight.Attacks(tt.from, ctx)

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

// test all knight pseudolegal moves
func TestKnightPseudoLegalMoves(t *testing.T) {
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
				core.E6, core.E2,
				core.C6, core.C2,
				core.F5, core.F3,
				core.B5, core.B3,
			},
		},
		{
			name:        "empty board corner A1 has 2 moves",
			from:        core.A1,
			sideToMove:  core.WHITE,
			setupBoard:  func(b *core.Board) {},
			expectedTos: []core.Position{core.B3, core.C2},
		},
		{
			name:        "empty board corner H8 has 2 moves",
			from:        core.H8,
			sideToMove:  core.WHITE,
			setupBoard:  func(b *core.Board) {},
			expectedTos: []core.Position{core.G6, core.F7},
		},
		{
			name:        "empty board edge A4 has 4 moves",
			from:        core.A4,
			sideToMove:  core.WHITE,
			setupBoard:  func(b *core.Board) {},
			expectedTos: []core.Position{core.B6, core.B2, core.C5, core.C3},
		},
		{
			name:       "own piece on attack square is excluded",
			from:       core.D4,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.E6] = core.Square{Occupied: true, Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE}}
			},
			expectedTos: []core.Position{
				core.E2,
				core.C6, core.C2,
				core.F5, core.F3,
				core.B5, core.B3,
			},
		},
		{
			name:       "enemy piece on attack square is included (capture)",
			from:       core.D4,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.E6] = core.Square{Occupied: true, Piece: core.Piece{Type: core.KNIGHT, Color: core.BLACK}}
			},
			expectedTos: []core.Position{
				core.E6, core.E2,
				core.C6, core.C2,
				core.F5, core.F3,
				core.B5, core.B3,
			},
		},
		{
			name:       "captures carry the exact enemy piece sitting on destination",
			from:       core.D4,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.E6] = core.Square{Occupied: true, Piece: core.Piece{Type: core.QUEEN, Color: core.BLACK}}
				b[core.F5] = core.Square{Occupied: true, Piece: core.Piece{Type: core.ROOK, Color: core.BLACK}}
				b[core.C2] = core.Square{Occupied: true, Piece: core.Piece{Type: core.BISHOP, Color: core.BLACK}}
			},
			expectedTos: []core.Position{
				core.E6, core.E2,
				core.C6, core.C2,
				core.F5, core.F3,
				core.B5, core.B3,
			},
		},
		{
			name:       "mix of own and enemy blockers",
			from:       core.D4,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.E6] = core.Square{Occupied: true, Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE}}
				b[core.C6] = core.Square{Occupied: true, Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE}}
				b[core.F5] = core.Square{Occupied: true, Piece: core.Piece{Type: core.KNIGHT, Color: core.BLACK}}
			},
			expectedTos: []core.Position{
				core.E2,
				core.C2,
				core.F5,
				core.F3,
				core.B5, core.B3,
			},
		},
		{
			name:       "all attack squares blocked by own pieces yields no moves",
			from:       core.D4,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				for _, pos := range []core.Position{core.E6, core.E2, core.C6, core.C2, core.F5, core.F3, core.B5, core.B3} {
					b[pos] = core.Square{Occupied: true, Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE}}
				}
			},
			expectedTos: []core.Position{},
		},
		{
			name:       "black knight treats white piece as enemy (included)",
			from:       core.D4,
			sideToMove: core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.E6] = core.Square{Occupied: true, Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE}}
			},
			expectedTos: []core.Position{
				core.E6, core.E2,
				core.C6, core.C2,
				core.F5, core.F3,
				core.B5, core.B3,
			},
		},
		{
			name:       "black knight treats black piece as own (excluded)",
			from:       core.D4,
			sideToMove: core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.E6] = core.Square{Occupied: true, Piece: core.Piece{Type: core.KNIGHT, Color: core.BLACK}}
			},
			expectedTos: []core.Position{
				core.E2,
				core.C6, core.C2,
				core.F5, core.F3,
				core.B5, core.B3,
			},
		},
	}

	knight := Knight{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var board core.Board
			tt.setupBoard(&board)

			ctx := core.MoveContext{
				BoardContext: core.BoardContext{Board: &board},
				SideToMove:   tt.sideToMove,
			}

			got := knight.PseudoLegalMoves(tt.from, ctx)

			if len(got) != len(tt.expectedTos) {
				t.Fatalf("got %d moves, want %d", len(got), len(tt.expectedTos))
			}

			expectedMover := core.Piece{Type: core.KNIGHT, Color: tt.sideToMove}

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
				// every knight move is NORMAL
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
