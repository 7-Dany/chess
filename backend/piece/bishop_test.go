package piece

import (
	"testing"

	"github.com/7-Dany/chess/core"
)

func TestBishopIsAttacking(t *testing.T) {
	tests := []struct {
		name       string
		setupBoard func(*core.Board)
		color      core.PieceColor
		target     core.Position
		want       bool
	}{
		// ==================== Direct attacks (no blockers) ====================
		{
			name: "bishop on up-right diagonal attacks target",
			setupBoard: func(b *core.Board) {
				b[core.H7] = core.Square{Piece: core.Piece{Type: core.BISHOP, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true,
		},
		{
			name: "bishop on up-left diagonal attacks target",
			setupBoard: func(b *core.Board) {
				b[core.A8] = core.Square{Piece: core.Piece{Type: core.BISHOP, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true,
		},
		{
			name: "bishop on down-right diagonal attacks target",
			setupBoard: func(b *core.Board) {
				b[core.H1] = core.Square{Piece: core.Piece{Type: core.BISHOP, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true,
		},
		{
			name: "bishop on down-left diagonal attacks target",
			setupBoard: func(b *core.Board) {
				b[core.B1] = core.Square{Piece: core.Piece{Type: core.BISHOP, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true,
		},
		{
			name: "bishop adjacent to target attacks (distance 1)",
			setupBoard: func(b *core.Board) {
				b[core.D3] = core.Square{Piece: core.Piece{Type: core.BISHOP, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true,
		},
		{
			name: "bishop at maximum diagonal distance attacks (distance 7)",
			setupBoard: func(b *core.Board) {
				b[core.H8] = core.Square{Piece: core.Piece{Type: core.BISHOP, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.A1,
			want:   true,
		},

		// ==================== Not on diagonal ====================
		{
			name: "bishop on same file does not attack",
			setupBoard: func(b *core.Board) {
				b[core.E7] = core.Square{Piece: core.Piece{Type: core.BISHOP, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "bishop on same rank does not attack",
			setupBoard: func(b *core.Board) {
				b[core.H4] = core.Square{Piece: core.Piece{Type: core.BISHOP, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "bishop on non-diagonal square does not attack",
			setupBoard: func(b *core.Board) {
				b[core.F4] = core.Square{Piece: core.Piece{Type: core.BISHOP, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},

		// ==================== Blockers ====================
		{
			name: "friendly piece between bishop and target blocks attack",
			setupBoard: func(b *core.Board) {
				b[core.H7] = core.Square{Piece: core.Piece{Type: core.BISHOP, Color: core.WHITE}, Occupied: true}
				b[core.F5] = core.Square{Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "enemy piece between bishop and target blocks attack",
			setupBoard: func(b *core.Board) {
				b[core.H7] = core.Square{Piece: core.Piece{Type: core.BISHOP, Color: core.WHITE}, Occupied: true}
				b[core.F5] = core.Square{Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "piece behind target does not block (target between bishop and piece)",
			setupBoard: func(b *core.Board) {
				b[core.H7] = core.Square{Piece: core.Piece{Type: core.BISHOP, Color: core.WHITE}, Occupied: true}
				b[core.C2] = core.Square{Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true,
		},
		{
			name: "two bishops on same diagonal, closer one blocks farther one",
			setupBoard: func(b *core.Board) {
				b[core.H7] = core.Square{Piece: core.Piece{Type: core.BISHOP, Color: core.WHITE}, Occupied: true}
				b[core.F5] = core.Square{Piece: core.Piece{Type: core.BISHOP, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true, // F5 bishop attacks target (closer one)
		},

		// ==================== Color filtering ====================
		{
			name: "black bishop, asking for white attackers",
			setupBoard: func(b *core.Board) {
				b[core.H7] = core.Square{Piece: core.Piece{Type: core.BISHOP, Color: core.BLACK}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "black bishop, asking for black attackers",
			setupBoard: func(b *core.Board) {
				b[core.H7] = core.Square{Piece: core.Piece{Type: core.BISHOP, Color: core.BLACK}, Occupied: true}
			},
			color:  core.BLACK,
			target: core.E4,
			want:   true,
		},

		// ==================== Non-bishop pieces on diagonal ====================
		{
			name: "queen on diagonal does NOT trigger bishop attack (clean separation)",
			setupBoard: func(b *core.Board) {
				b[core.H7] = core.Square{Piece: core.Piece{Type: core.QUEEN, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "rook on diagonal does not trigger bishop attack",
			setupBoard: func(b *core.Board) {
				b[core.H7] = core.Square{Piece: core.Piece{Type: core.ROOK, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "pawn on diagonal does not trigger bishop attack",
			setupBoard: func(b *core.Board) {
				b[core.D5] = core.Square{Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "king on diagonal does not trigger bishop attack",
			setupBoard: func(b *core.Board) {
				b[core.F5] = core.Square{Piece: core.Piece{Type: core.KING, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "knight on diagonal does not trigger bishop attack",
			setupBoard: func(b *core.Board) {
				b[core.G6] = core.Square{Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE}, Occupied: true}
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
			name: "target on corner A1, bishop on H8 attacks",
			setupBoard: func(b *core.Board) {
				b[core.H8] = core.Square{Piece: core.Piece{Type: core.BISHOP, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.A1,
			want:   true,
		},
		{
			name: "target on corner H8, bishop on A1 attacks",
			setupBoard: func(b *core.Board) {
				b[core.A1] = core.Square{Piece: core.Piece{Type: core.BISHOP, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.H8,
			want:   true,
		},
		{
			name: "target on corner A8, bishop on H1 attacks",
			setupBoard: func(b *core.Board) {
				b[core.H1] = core.Square{Piece: core.Piece{Type: core.BISHOP, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.A8,
			want:   true,
		},
		{
			name: "target on corner H1, bishop on A8 attacks",
			setupBoard: func(b *core.Board) {
				b[core.A8] = core.Square{Piece: core.Piece{Type: core.BISHOP, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.H1,
			want:   true,
		},
		{
			name: "bishop on target square itself does not attack",
			setupBoard: func(b *core.Board) {
				b[core.E4] = core.Square{Piece: core.Piece{Type: core.BISHOP, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},

		// ==================== Multiple bishops ====================
		{
			name: "multiple bishops, one attacks along open diagonal",
			setupBoard: func(b *core.Board) {
				b[core.A1] = core.Square{Piece: core.Piece{Type: core.BISHOP, Color: core.WHITE}, Occupied: true}
				b[core.H1] = core.Square{Piece: core.Piece{Type: core.BISHOP, Color: core.WHITE}, Occupied: true}
				b[core.A8] = core.Square{Piece: core.Piece{Type: core.BISHOP, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true, // H1 bishop attacks via down-right diagonal
		},
		{
			name: "multiple enemy bishops, none attack (all blocked)",
			setupBoard: func(b *core.Board) {
				b[core.H7] = core.Square{Piece: core.Piece{Type: core.BISHOP, Color: core.WHITE}, Occupied: true}
				b[core.A7] = core.Square{Piece: core.Piece{Type: core.BISHOP, Color: core.WHITE}, Occupied: true}
				b[core.H1] = core.Square{Piece: core.Piece{Type: core.BISHOP, Color: core.WHITE}, Occupied: true}
				b[core.B1] = core.Square{Piece: core.Piece{Type: core.BISHOP, Color: core.WHITE}, Occupied: true}
				// Blockers on all four diagonals
				b[core.F5] = core.Square{Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}, Occupied: true}
				b[core.D5] = core.Square{Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}, Occupied: true}
				b[core.F3] = core.Square{Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}, Occupied: true}
				b[core.D3] = core.Square{Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "mixed-color bishops, only matching color counts",
			setupBoard: func(b *core.Board) {
				b[core.H7] = core.Square{Piece: core.Piece{Type: core.BISHOP, Color: core.BLACK}, Occupied: true}
				b[core.A8] = core.Square{Piece: core.Piece{Type: core.BISHOP, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true,
		},

		// ==================== Blocker just before target ====================
		{
			name: "blocker immediately adjacent to target blocks attack",
			setupBoard: func(b *core.Board) {
				b[core.H7] = core.Square{Piece: core.Piece{Type: core.BISHOP, Color: core.WHITE}, Occupied: true}
				b[core.F5] = core.Square{Piece: core.Piece{Type: core.ROOK, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "blocker immediately adjacent to bishop blocks attack",
			setupBoard: func(b *core.Board) {
				b[core.H7] = core.Square{Piece: core.Piece{Type: core.BISHOP, Color: core.WHITE}, Occupied: true}
				b[core.G6] = core.Square{Piece: core.Piece{Type: core.ROOK, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var board core.Board
			tt.setupBoard(&board)
			ctx := core.BoardContext{Board: &board}

			bishop := Bishop{}
			got := bishop.IsAttacking(tt.color, tt.target, ctx)

			if got != tt.want {
				t.Errorf("IsAttacking(%v, %v) = %v, want %v", tt.color, tt.target, got, tt.want)
			}
		})
	}
}

// test all attack moves for bishop
func TestBishopAttacks(t *testing.T) {
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
				core.E5, core.F6, core.G7, core.H8, // NE
				core.E3, core.F2, core.G1, // SE
				core.C5, core.B6, core.A7, // NW
				core.C3, core.B2, core.A1, // SW
			},
		},
		{
			name:       "corner A1 one diagonal",
			from:       core.A1,
			setupBoard: func(b *core.Board) {},
			expected: []core.Position{
				core.B2, core.C3, core.D4, core.E5, core.F6, core.G7, core.H8,
			},
		},
		{
			name:       "corner H1 one diagonal",
			from:       core.H1,
			setupBoard: func(b *core.Board) {},
			expected: []core.Position{
				core.G2, core.F3, core.E4, core.D5, core.C6, core.B7, core.A8,
			},
		},
		{
			name:       "corner A8 one diagonal",
			from:       core.A8,
			setupBoard: func(b *core.Board) {},
			expected: []core.Position{
				core.B7, core.C6, core.D5, core.E4, core.F3, core.G2, core.H1,
			},
		},
		{
			name:       "corner H8 one diagonal",
			from:       core.H8,
			setupBoard: func(b *core.Board) {},
			expected: []core.Position{
				core.G7, core.F6, core.E5, core.D4, core.C3, core.B2, core.A1,
			},
		},
		{
			name:       "edge A4 two diagonals",
			from:       core.A4,
			setupBoard: func(b *core.Board) {},
			expected: []core.Position{
				core.B5, core.C6, core.D7, core.E8, // NE
				core.B3, core.C2, core.D1, // SE
			},
		},
		{
			name: "center D4 blocked on NE at F6",
			from: core.D4,
			setupBoard: func(b *core.Board) {
				b[core.F6] = core.Square{Occupied: true}
			},
			expected: []core.Position{
				core.E5, core.F6, // NE stops
				core.E3, core.F2, core.G1, // SE
				core.C5, core.B6, core.A7, // NW
				core.C3, core.B2, core.A1, // SW
			},
		},
		{
			name: "center D4 trapped all diagonals blocked",
			from: core.D4,
			setupBoard: func(b *core.Board) {
				b[core.E5] = core.Square{Occupied: true}
				b[core.E3] = core.Square{Occupied: true}
				b[core.C5] = core.Square{Occupied: true}
				b[core.C3] = core.Square{Occupied: true}
			},
			expected: []core.Position{core.E5, core.E3, core.C5, core.C3},
		},
		{
			name: "corner A1 blocked at C3",
			from: core.A1,
			setupBoard: func(b *core.Board) {
				b[core.C3] = core.Square{Occupied: true}
			},
			expected: []core.Position{core.B2, core.C3},
		},
	}

	bishop := Bishop{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var board core.Board
			tt.setupBoard(&board)
			ctx := core.BoardContext{Board: &board}

			got := bishop.Attacks(tt.from, ctx)

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

// test all bishop pseudolegal moves
func TestBishopPseudoLegalMoves(t *testing.T) {
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
				core.E5, core.F6, core.G7, core.H8,
				core.E3, core.F2, core.G1,
				core.C5, core.B6, core.A7,
				core.C3, core.B2, core.A1,
			},
		},
		{
			name:       "can capture enemy piece",
			from:       core.D4,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.F6] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}}
			},
			expectedTos: []core.Position{
				core.E5, core.F6,
				core.E3, core.F2, core.G1,
				core.C5, core.B6, core.A7,
				core.C3, core.B2, core.A1,
			},
		},
		{
			name:       "cannot move to friendly occupied square",
			from:       core.D4,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.F6] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}}
			},
			expectedTos: []core.Position{
				core.E5,
				core.E3, core.F2, core.G1,
				core.C5, core.B6, core.A7,
				core.C3, core.B2, core.A1,
			},
		},
		{
			name:       "friendly blocks slide path; enemy behind is unreachable",
			from:       core.D4,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.E5] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}}
				b[core.F6] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}}
			},
			expectedTos: []core.Position{
				core.E3, core.F2, core.G1,
				core.C5, core.B6, core.A7,
				core.C3, core.B2, core.A1,
			},
		},
		{
			name:       "enemy blocks slide path but is capturable",
			from:       core.D4,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.F6] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}}
				b[core.H8] = core.Square{Occupied: true, Piece: core.Piece{Type: core.ROOK, Color: core.BLACK}}
			},
			expectedTos: []core.Position{
				core.E5, core.F6,
				core.E3, core.F2, core.G1,
				core.C5, core.B6, core.A7,
				core.C3, core.B2, core.A1,
			},
		},
		{
			name:       "enemy capturable; friendly behind is unreachable",
			from:       core.D4,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.F6] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}}
				b[core.H8] = core.Square{Occupied: true, Piece: core.Piece{Type: core.ROOK, Color: core.WHITE}}
			},
			expectedTos: []core.Position{
				core.E5, core.F6,
				core.E3, core.F2, core.G1,
				core.C5, core.B6, core.A7,
				core.C3, core.B2, core.A1,
			},
		},
		{
			name:       "mixed friendly and enemy on all diagonals",
			from:       core.D4,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.E5] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}}
				b[core.E3] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}}
				b[core.C5] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}}
				b[core.C3] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}}
			},
			expectedTos: []core.Position{core.E5, core.C5},
		},
		{
			name:       "all diagonals blocked by own pieces yields no moves",
			from:       core.D4,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.E5] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}}
				b[core.E3] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}}
				b[core.C5] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}}
				b[core.C3] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}}
			},
			expectedTos: []core.Position{},
		},
		{
			name:       "black bishop treats white piece as enemy (included)",
			from:       core.D4,
			sideToMove: core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.F6] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}}
			},
			expectedTos: []core.Position{
				core.E5, core.F6,
				core.E3, core.F2, core.G1,
				core.C5, core.B6, core.A7,
				core.C3, core.B2, core.A1,
			},
		},
		{
			name:       "black bishop treats black piece as own (excluded)",
			from:       core.D4,
			sideToMove: core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.F6] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}}
			},
			expectedTos: []core.Position{
				core.E5,
				core.E3, core.F2, core.G1,
				core.C5, core.B6, core.A7,
				core.C3, core.B2, core.A1,
			},
		},
	}

	bishop := Bishop{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var board core.Board
			tt.setupBoard(&board)

			ctx := core.MoveContext{
				BoardContext: core.BoardContext{Board: &board},
				SideToMove:   tt.sideToMove,
			}

			got := bishop.PseudoLegalMoves(tt.from, ctx)

			if len(got) != len(tt.expectedTos) {
				t.Fatalf("got %d moves, want %d", len(got), len(tt.expectedTos))
			}

			expectedMover := core.Piece{Type: core.BISHOP, Color: tt.sideToMove}

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
				// every bishop move is NORMAL
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
