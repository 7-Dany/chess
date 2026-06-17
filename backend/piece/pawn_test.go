package piece

import (
	"testing"

	"github.com/7-Dany/chess/core"
)

func TestPawnIsAttacking(t *testing.T) {
	tests := []struct {
		name       string
		setupBoard func(*core.Board)
		color      core.PieceColor
		target     core.Position
		want       bool
	}{
		// ==================== White pawn attacks (from below) ====================
		// Target E4 is attacked by a white pawn on D3 or F3.
		{
			name: "white pawn attacks from down-left (D3 attacks E4)",
			setupBoard: func(b *core.Board) {
				b[core.D3] = core.Square{Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true,
		},
		{
			name: "white pawn attacks from down-right (F3 attacks E4)",
			setupBoard: func(b *core.Board) {
				b[core.F3] = core.Square{Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true,
		},

		// ==================== Black pawn attacks (from above) ====================
		// Target E4 is attacked by a black pawn on D5 or F5.
		{
			name: "black pawn attacks from up-left (D5 attacks E4)",
			setupBoard: func(b *core.Board) {
				b[core.D5] = core.Square{Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}, Occupied: true}
			},
			color:  core.BLACK,
			target: core.E4,
			want:   true,
		},
		{
			name: "black pawn attacks from up-right (F5 attacks E4)",
			setupBoard: func(b *core.Board) {
				b[core.F5] = core.Square{Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}, Occupied: true}
			},
			color:  core.BLACK,
			target: core.E4,
			want:   true,
		},

		// ==================== Not on attack diagonal ====================
		{
			name: "white pawn directly below target (same file) does NOT attack",
			setupBoard: func(b *core.Board) {
				b[core.E3] = core.Square{Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false, // pawns push forward but don't attack forward
		},
		{
			name: "black pawn directly above target (same file) does NOT attack",
			setupBoard: func(b *core.Board) {
				b[core.E5] = core.Square{Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}, Occupied: true}
			},
			color:  core.BLACK,
			target: core.E4,
			want:   false,
		},
		{
			name: "white pawn two ranks below target does not attack",
			setupBoard: func(b *core.Board) {
				b[core.D2] = core.Square{Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false, // wrong rank (two below instead of one)
		},
		{
			name: "white pawn above target does not attack (wrong direction)",
			setupBoard: func(b *core.Board) {
				b[core.D5] = core.Square{Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false, // white pawns attack up, not down
		},
		{
			name: "black pawn below target does not attack (wrong direction)",
			setupBoard: func(b *core.Board) {
				b[core.F3] = core.Square{Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}, Occupied: true}
			},
			color:  core.BLACK,
			target: core.E4,
			want:   false, // black pawns attack down, not up
		},
		{
			name: "white pawn on same rank does not attack",
			setupBoard: func(b *core.Board) {
				b[core.D4] = core.Square{Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},

		// ==================== Color filtering ====================
		{
			name: "black pawn on white-attack square, asking for white",
			setupBoard: func(b *core.Board) {
				b[core.D3] = core.Square{Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false, // wrong color even though on the geometric attack square
		},
		{
			name: "white pawn on black-attack square, asking for black",
			setupBoard: func(b *core.Board) {
				b[core.D5] = core.Square{Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}, Occupied: true}
			},
			color:  core.BLACK,
			target: core.E4,
			want:   false,
		},
		{
			name: "black pawn on white-attack square, asking for black",
			setupBoard: func(b *core.Board) {
				b[core.D3] = core.Square{Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}, Occupied: true}
			},
			color:  core.BLACK,
			target: core.E4,
			want:   false, // D3 is one rank BELOW target — black pawns don't attack downward
		},
		{
			name: "white pawn on black-attack square, asking for white",
			setupBoard: func(b *core.Board) {
				b[core.D5] = core.Square{Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false, // D5 is one rank ABOVE target — white pawns don't attack downward
		},

		// ==================== Non-pawn on attack diagonal (clean separation) ====================
		{
			name: "queen on attack diagonal does NOT trigger pawn attack",
			setupBoard: func(b *core.Board) {
				b[core.D3] = core.Square{Piece: core.Piece{Type: core.QUEEN, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "bishop on attack diagonal does not trigger pawn attack",
			setupBoard: func(b *core.Board) {
				b[core.D3] = core.Square{Piece: core.Piece{Type: core.BISHOP, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "knight on attack diagonal does not trigger pawn attack",
			setupBoard: func(b *core.Board) {
				b[core.F3] = core.Square{Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "rook on attack diagonal does not trigger pawn attack",
			setupBoard: func(b *core.Board) {
				b[core.F3] = core.Square{Piece: core.Piece{Type: core.ROOK, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},
		{
			name: "king on attack diagonal does not trigger pawn attack",
			setupBoard: func(b *core.Board) {
				b[core.D3] = core.Square{Piece: core.Piece{Type: core.KING, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false,
		},

		// ==================== Edge ranks (pawns can't attack from off-board) ====================
		{
			name: "target on rank 1 cannot be attacked by white pawn (no rank below)",
			setupBoard: func(b *core.Board) {
				b[core.D1] = core.Square{Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E1,
			want:   false,
		},
		{
			name: "target on rank 1 CAN be attacked by black pawn from rank 2",
			setupBoard: func(b *core.Board) {
				b[core.D2] = core.Square{Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}, Occupied: true}
			},
			color:  core.BLACK,
			target: core.E1,
			want:   true,
		},
		{
			name: "target on rank 8 cannot be attacked by black pawn (no rank above)",
			setupBoard: func(b *core.Board) {
				b[core.D8] = core.Square{Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}, Occupied: true}
			},
			color:  core.BLACK,
			target: core.E8,
			want:   false,
		},
		{
			name: "target on rank 8 CAN be attacked by white pawn from rank 7",
			setupBoard: func(b *core.Board) {
				b[core.D7] = core.Square{Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E8,
			want:   true,
		},

		// ==================== Edge files (only one valid diagonal) ====================
		{
			name: "target on A4, white pawn attacks from B3 (only right diagonal valid)",
			setupBoard: func(b *core.Board) {
				b[core.B3] = core.Square{Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.A4,
			want:   true,
		},
		{
			name: "target on A4, no white pawn attacks from off-board left diagonal",
			setupBoard: func(b *core.Board) {
				b[core.B3] = core.Square{Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.A4,
			want:   true, // B3 attacks; the off-board left diagonal is correctly skipped
		},
		{
			name: "target on H4, white pawn attacks from G3 (only left diagonal valid)",
			setupBoard: func(b *core.Board) {
				b[core.G3] = core.Square{Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.H4,
			want:   true,
		},
		{
			name: "target on A4, black pawn attacks from B5",
			setupBoard: func(b *core.Board) {
				b[core.B5] = core.Square{Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}, Occupied: true}
			},
			color:  core.BLACK,
			target: core.A4,
			want:   true,
		},
		{
			name: "target on H4, black pawn attacks from G5",
			setupBoard: func(b *core.Board) {
				b[core.G5] = core.Square{Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}, Occupied: true}
			},
			color:  core.BLACK,
			target: core.H4,
			want:   true,
		},

		// ==================== Corners ====================
		{
			name: "target on A1, no white pawn attack possible (rank 1 + A file)",
			setupBoard: func(b *core.Board) {
				b[core.B1] = core.Square{Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.A1,
			want:   false, // B1 is on rank 1, can't attack A1 (white pawns attack upward, so attacker must be on rank 0 — off-board)
		},
		{
			name: "target on A1, black pawn attacks from B2",
			setupBoard: func(b *core.Board) {
				b[core.B2] = core.Square{Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}, Occupied: true}
			},
			color:  core.BLACK,
			target: core.A1,
			want:   true,
		},
		{
			name: "target on H8, no black pawn attack possible (rank 8 + H file)",
			setupBoard: func(b *core.Board) {
				b[core.G8] = core.Square{Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}, Occupied: true}
			},
			color:  core.BLACK,
			target: core.H8,
			want:   false,
		},
		{
			name: "target on H8, white pawn attacks from G7",
			setupBoard: func(b *core.Board) {
				b[core.G7] = core.Square{Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.H8,
			want:   true,
		},
		{
			name: "target on A8, white pawn attacks from B7",
			setupBoard: func(b *core.Board) {
				b[core.B7] = core.Square{Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.A8,
			want:   true,
		},
		{
			name: "target on H1, black pawn attacks from G2",
			setupBoard: func(b *core.Board) {
				b[core.G2] = core.Square{Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}, Occupied: true}
			},
			color:  core.BLACK,
			target: core.H1,
			want:   true,
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
			name: "pawn on target square itself does not attack",
			setupBoard: func(b *core.Board) {
				b[core.E4] = core.Square{Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false, // scan checks D3/F3, never E4 itself
		},

		// ==================== Multiple pawns ====================
		{
			name: "two white pawns on both attack diagonals",
			setupBoard: func(b *core.Board) {
				b[core.D3] = core.Square{Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}, Occupied: true}
				b[core.F3] = core.Square{Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true,
		},
		{
			name: "white and black pawns on opposite diagonals, asking for white",
			setupBoard: func(b *core.Board) {
				b[core.D3] = core.Square{Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}, Occupied: true}
				b[core.D5] = core.Square{Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   true, // D3 (white) attacks
		},
		{
			name: "white and black pawns on opposite diagonals, asking for black",
			setupBoard: func(b *core.Board) {
				b[core.D3] = core.Square{Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}, Occupied: true}
				b[core.D5] = core.Square{Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}, Occupied: true}
			},
			color:  core.BLACK,
			target: core.E4,
			want:   true, // D5 (black) attacks
		},
		{
			name: "pawns on wrong-color attack squares don't match",
			setupBoard: func(b *core.Board) {
				b[core.D3] = core.Square{Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}, Occupied: true}
				b[core.D5] = core.Square{Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}, Occupied: true}
			},
			color:  core.WHITE,
			target: core.E4,
			want:   false, // D3 black (wrong color), D5 white (wrong direction for white)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var board core.Board
			tt.setupBoard(&board)
			ctx := core.BoardContext{Board: &board}

			pawn := Pawn{}
			got := pawn.IsAttacking(tt.color, tt.target, ctx)

			if got != tt.want {
				t.Errorf("IsAttacking(%v, %v) = %v, want %v", tt.color, tt.target, got, tt.want)
			}
		})
	}
}

// test pawn attacks
func TestPawnAttacks(t *testing.T) {
	p := Pawn{}

	newCtx := func(from core.Position, color core.PieceColor) core.BoardContext {
		board := &core.Board{}
		board[from] = core.Square{Piece: core.Piece{Color: color}, Occupied: true}
		return core.BoardContext{Board: board}
	}

	tests := []struct {
		name  string
		from  core.Position
		color core.PieceColor
		want  []core.Position
	}{
		// Middle of the board
		{
			name:  "white pawn in the middle attacks two squares",
			from:  core.E4,
			color: core.WHITE,
			want:  []core.Position{core.F5, core.D5},
		},
		{
			name:  "black pawn in the middle attacks two squares",
			from:  core.E5,
			color: core.BLACK,
			want:  []core.Position{core.F4, core.D4},
		},

		// A file — no left attack
		{
			name:  "white pawn on A file attacks one square only",
			from:  core.A4,
			color: core.WHITE,
			want:  []core.Position{core.B5},
		},
		{
			name:  "black pawn on A file attacks one square only",
			from:  core.A5,
			color: core.BLACK,
			want:  []core.Position{core.B4},
		},

		// H file — no right attack
		{
			name:  "white pawn on H file attacks one square only",
			from:  core.H4,
			color: core.WHITE,
			want:  []core.Position{core.G5},
		},
		{
			name:  "black pawn on H file attacks one square only",
			from:  core.H5,
			color: core.BLACK,
			want:  []core.Position{core.G4},
		},

		// Last rank — rank guard triggers, no attacks
		{
			name:  "white pawn on last rank returns no attacks",
			from:  core.E8,
			color: core.WHITE,
			want:  []core.Position{},
		},
		{
			name:  "black pawn on last rank returns no attacks",
			from:  core.E1,
			color: core.BLACK,
			want:  []core.Position{},
		},

		// Behind starting rank — still produces valid attacks
		{
			name:  "white pawn on rank 1 attacks normally",
			from:  core.E1,
			color: core.WHITE,
			want:  []core.Position{core.F2, core.D2},
		},
		{
			name:  "black pawn on rank 8 attacks normally",
			from:  core.E8,
			color: core.BLACK,
			want:  []core.Position{core.F7, core.D7},
		},

		// Corners — combines file and last rank edges
		{
			name:  "white pawn on A8 corner returns no attacks",
			from:  core.A8,
			color: core.WHITE,
			want:  []core.Position{},
		},
		{
			name:  "white pawn on H8 corner returns no attacks",
			from:  core.H8,
			color: core.WHITE,
			want:  []core.Position{},
		},
		{
			name:  "black pawn on A1 corner returns no attacks",
			from:  core.A1,
			color: core.BLACK,
			want:  []core.Position{},
		},
		{
			name:  "black pawn on H1 corner returns no attacks",
			from:  core.H1,
			color: core.BLACK,
			want:  []core.Position{},
		},

		// Corners — file edge but not last rank
		{
			name:  "white pawn on A1 corner attacks one square",
			from:  core.A1,
			color: core.WHITE,
			want:  []core.Position{core.B2},
		},
		{
			name:  "white pawn on H1 corner attacks one square",
			from:  core.H1,
			color: core.WHITE,
			want:  []core.Position{core.G2},
		},
		{
			name:  "black pawn on A8 corner attacks one square",
			from:  core.A8,
			color: core.BLACK,
			want:  []core.Position{core.B7},
		},
		{
			name:  "black pawn on H8 corner attacks one square",
			from:  core.H8,
			color: core.BLACK,
			want:  []core.Position{core.G7},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := p.Attacks(tt.from, newCtx(tt.from, tt.color))

			if len(got) != len(tt.want) {
				t.Fatalf("len: got %d, want %d — got %v, want %v", len(got), len(tt.want), got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("[%d]: got %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

// test all pawn pseudolegal moves
func TestPawnPseudoLegalMoves(t *testing.T) {
	promotionPieces := []core.PieceType{core.QUEEN, core.ROOK, core.BISHOP, core.KNIGHT}

	// assertPromotions verifies that all four promotion options exist for a given target square.
	assertPromotions := func(t *testing.T, got []core.Move, to core.Position) {
		t.Helper()
		found := make(map[core.PieceType]bool)
		for _, move := range got {
			if move.To == to && move.Type == core.PROMOTION {
				found[move.PromoteTo] = true
			}
		}
		for _, pt := range promotionPieces {
			if !found[pt] {
				t.Errorf("missing promotion to %v at %v", pt, to)
			}
		}
	}

	// assertMoveType verifies that the move to a given square carries the expected type.
	assertMoveType := func(t *testing.T, got []core.Move, to core.Position, want core.MoveType) {
		t.Helper()
		for _, move := range got {
			if move.To == to {
				if move.Type != want {
					t.Errorf("move to %v: got type %v, want %v", to, move.Type, want)
				}
				return
			}
		}
		t.Errorf("no move found to %v", to)
	}

	tests := []struct {
		name            string
		from            core.Position
		sideToMove      core.PieceColor
		setupBoard      func(*core.Board)
		enPassantTarget core.Position
		// expectedTos lists every distinct destination square expected.
		// Promotion squares appear once here even though they yield four Moves.
		expectedTos []core.Position
		// promotionTos lists squares that must yield all four promotion options.
		promotionTos []core.Position
		// enPassantTos lists squares that must carry MoveType EN_PASSANT.
		enPassantTos []core.Position
	}{
		// ── single push ──────────────────────────────────────────────────────────

		{
			name:        "white single push from non-starting rank",
			from:        core.E4,
			sideToMove:  core.WHITE,
			setupBoard:  func(b *core.Board) {},
			expectedTos: []core.Position{core.E5},
		},
		{
			name:        "black single push from non-starting rank",
			from:        core.E5,
			sideToMove:  core.BLACK,
			setupBoard:  func(b *core.Board) {},
			expectedTos: []core.Position{core.E4},
		},
		{
			name:       "white single push blocked by own piece",
			from:       core.E4,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.E5] = core.Square{Occupied: true, Piece: core.Piece{Type: core.ROOK, Color: core.WHITE}}
			},
			expectedTos: []core.Position{},
		},
		{
			name:       "white single push blocked by enemy piece",
			from:       core.E4,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.E5] = core.Square{Occupied: true, Piece: core.Piece{Type: core.ROOK, Color: core.BLACK}}
			},
			expectedTos: []core.Position{},
		},
		{
			name:       "black single push blocked by own piece",
			from:       core.E5,
			sideToMove: core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.E4] = core.Square{Occupied: true, Piece: core.Piece{Type: core.ROOK, Color: core.BLACK}}
			},
			expectedTos: []core.Position{},
		},
		{
			name:       "black single push blocked by enemy piece",
			from:       core.E5,
			sideToMove: core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.E4] = core.Square{Occupied: true, Piece: core.Piece{Type: core.ROOK, Color: core.WHITE}}
			},
			expectedTos: []core.Position{},
		},

		// ── double push ──────────────────────────────────────────────────────────

		{
			name:        "white double push from starting rank",
			from:        core.E2,
			sideToMove:  core.WHITE,
			setupBoard:  func(b *core.Board) {},
			expectedTos: []core.Position{core.E3, core.E4},
		},
		{
			name:        "black double push from starting rank",
			from:        core.E7,
			sideToMove:  core.BLACK,
			setupBoard:  func(b *core.Board) {},
			expectedTos: []core.Position{core.E6, core.E5},
		},
		{
			name:        "double push not available from non-starting rank",
			from:        core.E3,
			sideToMove:  core.WHITE,
			setupBoard:  func(b *core.Board) {},
			expectedTos: []core.Position{core.E4},
		},
		{
			name:       "white double push blocked by own piece on first square",
			from:       core.E2,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.E3] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}}
			},
			expectedTos: []core.Position{},
		},
		{
			name:       "white double push blocked by enemy piece on first square",
			from:       core.E2,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.E3] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}}
			},
			expectedTos: []core.Position{},
		},
		{
			name:       "white double push blocked by own piece on second square",
			from:       core.E2,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.E4] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}}
			},
			expectedTos: []core.Position{core.E3},
		},
		{
			name:       "white double push blocked by enemy piece on second square",
			from:       core.E2,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.E4] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}}
			},
			expectedTos: []core.Position{core.E3},
		},
		{
			name:       "black double push blocked by own piece on first square",
			from:       core.E7,
			sideToMove: core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.E6] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}}
			},
			expectedTos: []core.Position{},
		},
		{
			name:       "black double push blocked by enemy piece on first square",
			from:       core.E7,
			sideToMove: core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.E6] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}}
			},
			expectedTos: []core.Position{},
		},
		{
			name:       "black double push blocked by own piece on second square",
			from:       core.E7,
			sideToMove: core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.E5] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}}
			},
			expectedTos: []core.Position{core.E6},
		},
		{
			name:       "black double push blocked by enemy piece on second square",
			from:       core.E7,
			sideToMove: core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.E5] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}}
			},
			expectedTos: []core.Position{core.E6},
		},

		// ── diagonal captures ────────────────────────────────────────────────────

		{
			name:       "white captures enemy on both diagonals",
			from:       core.E4,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.D5] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}}
				b[core.F5] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}}
			},
			expectedTos: []core.Position{core.E5, core.D5, core.F5},
		},
		{
			name:       "white cannot capture own piece on diagonal",
			from:       core.E4,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.D5] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}}
				b[core.F5] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}}
			},
			expectedTos: []core.Position{core.E5},
		},
		{
			name:        "white does not capture on empty diagonal squares",
			from:        core.E4,
			sideToMove:  core.WHITE,
			setupBoard:  func(b *core.Board) {},
			expectedTos: []core.Position{core.E5},
		},
		{
			name:       "black captures enemy on both diagonals",
			from:       core.E5,
			sideToMove: core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.D4] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}}
				b[core.F4] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}}
			},
			expectedTos: []core.Position{core.E4, core.D4, core.F4},
		},
		{
			name:       "black cannot capture own piece on diagonal",
			from:       core.E5,
			sideToMove: core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.D4] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}}
				b[core.F4] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}}
			},
			expectedTos: []core.Position{core.E4},
		},
		{
			name:        "black does not capture on empty diagonal squares",
			from:        core.E5,
			sideToMove:  core.BLACK,
			setupBoard:  func(b *core.Board) {},
			expectedTos: []core.Position{core.E4},
		},

		// ── A/H file edge cases ──────────────────────────────────────────────────

		{
			name:       "white pawn on A file captures right only",
			from:       core.A4,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.B5] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}}
			},
			expectedTos: []core.Position{core.A5, core.B5},
		},
		{
			name:       "white pawn on H file captures left only",
			from:       core.H4,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.G5] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}}
			},
			expectedTos: []core.Position{core.H5, core.G5},
		},
		{
			name:       "black pawn on A file captures right only",
			from:       core.A5,
			sideToMove: core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.B4] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}}
			},
			expectedTos: []core.Position{core.A4, core.B4},
		},
		{
			name:       "black pawn on H file captures left only",
			from:       core.H5,
			sideToMove: core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.G4] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}}
			},
			expectedTos: []core.Position{core.H4, core.G4},
		},

		// ── en passant ───────────────────────────────────────────────────────────

		{
			name:            "white en passant to the right",
			from:            core.E5,
			sideToMove:      core.WHITE,
			setupBoard:      func(b *core.Board) {},
			enPassantTarget: core.F6,
			expectedTos:     []core.Position{core.E6, core.F6},
			enPassantTos:    []core.Position{core.F6},
		},
		{
			name:            "white en passant to the left",
			from:            core.E5,
			sideToMove:      core.WHITE,
			setupBoard:      func(b *core.Board) {},
			enPassantTarget: core.D6,
			expectedTos:     []core.Position{core.E6, core.D6},
			enPassantTos:    []core.Position{core.D6},
		},
		{
			name:            "black en passant to the right",
			from:            core.E4,
			sideToMove:      core.BLACK,
			setupBoard:      func(b *core.Board) {},
			enPassantTarget: core.F3,
			expectedTos:     []core.Position{core.E3, core.F3},
			enPassantTos:    []core.Position{core.F3},
		},
		{
			name:            "black en passant to the left",
			from:            core.E4,
			sideToMove:      core.BLACK,
			setupBoard:      func(b *core.Board) {},
			enPassantTarget: core.D3,
			expectedTos:     []core.Position{core.E3, core.D3},
			enPassantTos:    []core.Position{core.D3},
		},
		{
			name:            "en passant not generated when target is not an attack square",
			from:            core.E5,
			sideToMove:      core.WHITE,
			setupBoard:      func(b *core.Board) {},
			enPassantTarget: core.G6,
			expectedTos:     []core.Position{core.E6},
		},
		{
			// enPassantTarget is left as the zero value (A1). A1 is rank 1, which
			// is never a pawn attack square for any pawn — so the EP check fails
			// the same way "target is not an attack square" does. This test is
			// about "no EP target set", not about A1 specifically.
			name:            "en passant not generated when no target is set",
			from:            core.E5,
			sideToMove:      core.WHITE,
			setupBoard:      func(b *core.Board) {},
			enPassantTarget: core.Position(0),
			expectedTos:     []core.Position{core.E6},
		},
		{
			// En passant and normal capture coexist: one diagonal has an enemy
			// piece (regular capture) and the other diagonal is the en passant target.
			name:       "white en passant on one side and normal capture on the other",
			from:       core.E5,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.D6] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}}
			},
			enPassantTarget: core.F6,
			expectedTos:     []core.Position{core.E6, core.D6, core.F6},
			enPassantTos:    []core.Position{core.F6},
		},
		{
			name:       "black en passant on one side and normal capture on the other",
			from:       core.E4,
			sideToMove: core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.D3] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}}
			},
			enPassantTarget: core.F3,
			expectedTos:     []core.Position{core.E3, core.D3, core.F3},
			enPassantTos:    []core.Position{core.F3},
		},
		{
			// A-file pawn has only one attack square; en passant target is that square.
			name:            "white en passant on A file (single attack square)",
			from:            core.A5,
			sideToMove:      core.WHITE,
			setupBoard:      func(b *core.Board) {},
			enPassantTarget: core.B6,
			expectedTos:     []core.Position{core.A6, core.B6},
			enPassantTos:    []core.Position{core.B6},
		},
		{
			// H-file pawn has only one attack square; en passant target is that square.
			name:            "white en passant on H file (single attack square)",
			from:            core.H5,
			sideToMove:      core.WHITE,
			setupBoard:      func(b *core.Board) {},
			enPassantTarget: core.G6,
			expectedTos:     []core.Position{core.H6, core.G6},
			enPassantTos:    []core.Position{core.G6},
		},
		{
			// Symmetric to white-on-A: black A-file pawn has only one attack square.
			name:            "black en passant on A file (single attack square)",
			from:            core.A4,
			sideToMove:      core.BLACK,
			setupBoard:      func(b *core.Board) {},
			enPassantTarget: core.B3,
			expectedTos:     []core.Position{core.A3, core.B3},
			enPassantTos:    []core.Position{core.B3},
		},
		{
			// Symmetric to white-on-H: black H-file pawn has only one attack square.
			name:            "black en passant on H file (single attack square)",
			from:            core.H4,
			sideToMove:      core.BLACK,
			setupBoard:      func(b *core.Board) {},
			enPassantTarget: core.G3,
			expectedTos:     []core.Position{core.H3, core.G3},
			enPassantTos:    []core.Position{core.G3},
		},

		// ── promotion ────────────────────────────────────────────────────────────

		{
			name:         "white promotes via single push",
			from:         core.E7,
			sideToMove:   core.WHITE,
			setupBoard:   func(b *core.Board) {},
			promotionTos: []core.Position{core.E8},
			expectedTos:  []core.Position{core.E8},
		},
		{
			name:       "white promotes via push and both diagonal captures",
			from:       core.E7,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.D8] = core.Square{Occupied: true, Piece: core.Piece{Type: core.ROOK, Color: core.BLACK}}
				b[core.F8] = core.Square{Occupied: true, Piece: core.Piece{Type: core.ROOK, Color: core.BLACK}}
			},
			promotionTos: []core.Position{core.E8, core.D8, core.F8},
			expectedTos:  []core.Position{core.E8, core.D8, core.F8},
		},
		{
			name:       "white promotion push blocked, only captures promote",
			from:       core.E7,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.E8] = core.Square{Occupied: true, Piece: core.Piece{Type: core.ROOK, Color: core.BLACK}}
				b[core.F8] = core.Square{Occupied: true, Piece: core.Piece{Type: core.ROOK, Color: core.BLACK}}
			},
			promotionTos: []core.Position{core.F8},
			expectedTos:  []core.Position{core.F8},
		},
		{
			// Own piece on a promotion diagonal must not generate promotion moves.
			name:       "white cannot promote by capturing own piece",
			from:       core.E7,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.D8] = core.Square{Occupied: true, Piece: core.Piece{Type: core.ROOK, Color: core.WHITE}}
				b[core.F8] = core.Square{Occupied: true, Piece: core.Piece{Type: core.ROOK, Color: core.WHITE}}
			},
			promotionTos: []core.Position{core.E8},
			expectedTos:  []core.Position{core.E8},
		},
		{
			name:         "black promotes via single push",
			from:         core.E2,
			sideToMove:   core.BLACK,
			setupBoard:   func(b *core.Board) {},
			promotionTos: []core.Position{core.E1},
			expectedTos:  []core.Position{core.E1},
		},
		{
			name:       "black promotes via push and both diagonal captures",
			from:       core.E2,
			sideToMove: core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.D1] = core.Square{Occupied: true, Piece: core.Piece{Type: core.ROOK, Color: core.WHITE}}
				b[core.F1] = core.Square{Occupied: true, Piece: core.Piece{Type: core.ROOK, Color: core.WHITE}}
			},
			promotionTos: []core.Position{core.E1, core.D1, core.F1},
			expectedTos:  []core.Position{core.E1, core.D1, core.F1},
		},
		{
			name:       "black promotion push blocked, only captures promote",
			from:       core.E2,
			sideToMove: core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.E1] = core.Square{Occupied: true, Piece: core.Piece{Type: core.ROOK, Color: core.WHITE}}
				b[core.F1] = core.Square{Occupied: true, Piece: core.Piece{Type: core.ROOK, Color: core.WHITE}}
			},
			promotionTos: []core.Position{core.F1},
			expectedTos:  []core.Position{core.F1},
		},
		{
			name:       "black cannot promote by capturing own piece",
			from:       core.E2,
			sideToMove: core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.D1] = core.Square{Occupied: true, Piece: core.Piece{Type: core.ROOK, Color: core.BLACK}}
				b[core.F1] = core.Square{Occupied: true, Piece: core.Piece{Type: core.ROOK, Color: core.BLACK}}
			},
			promotionTos: []core.Position{core.E1},
			expectedTos:  []core.Position{core.E1},
		},
		{
			// A-file pawn: only one diagonal capture is possible on promotion rank.
			name:       "white A-file pawn promotes via push and one diagonal capture",
			from:       core.A7,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.B8] = core.Square{Occupied: true, Piece: core.Piece{Type: core.ROOK, Color: core.BLACK}}
			},
			promotionTos: []core.Position{core.A8, core.B8},
			expectedTos:  []core.Position{core.A8, core.B8},
		},
		{
			// H-file pawn: only one diagonal capture is possible on promotion rank.
			name:       "white H-file pawn promotes via push and one diagonal capture",
			from:       core.H7,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.G8] = core.Square{Occupied: true, Piece: core.Piece{Type: core.ROOK, Color: core.BLACK}}
			},
			promotionTos: []core.Position{core.H8, core.G8},
			expectedTos:  []core.Position{core.H8, core.G8},
		},
		{
			name:       "black A-file pawn promotes via push and one diagonal capture",
			from:       core.A2,
			sideToMove: core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.B1] = core.Square{Occupied: true, Piece: core.Piece{Type: core.ROOK, Color: core.WHITE}}
			},
			promotionTos: []core.Position{core.A1, core.B1},
			expectedTos:  []core.Position{core.A1, core.B1},
		},
		{
			name:       "black H-file pawn promotes via push and one diagonal capture",
			from:       core.H2,
			sideToMove: core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.G1] = core.Square{Occupied: true, Piece: core.Piece{Type: core.ROOK, Color: core.WHITE}}
			},
			promotionTos: []core.Position{core.H1, core.G1},
			expectedTos:  []core.Position{core.H1, core.G1},
		},
		{
			// Promotion path with everything blocked — push square and both diagonals
			// unavailable. Verifies that promote returns an empty slice, not a stray move.
			name:       "white pawn on rank 7 fully blocked yields no moves",
			from:       core.E7,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.E8] = core.Square{Occupied: true, Piece: core.Piece{Type: core.ROOK, Color: core.WHITE}}
				b[core.D8] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}}
				b[core.F8] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}}
			},
			expectedTos: []core.Position{},
		},
		{
			name:       "black pawn on rank 2 fully blocked yields no moves",
			from:       core.E2,
			sideToMove: core.BLACK,
			setupBoard: func(b *core.Board) {
				b[core.E1] = core.Square{Occupied: true, Piece: core.Piece{Type: core.ROOK, Color: core.BLACK}}
				b[core.D1] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}}
				b[core.F1] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}}
			},
			expectedTos: []core.Position{},
		},

		// ── combined ─────────────────────────────────────────────────────────────

		{
			name:       "white on starting rank with double push and both captures",
			from:       core.E2,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.D3] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}}
				b[core.F3] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK}}
			},
			expectedTos: []core.Position{core.E3, core.E4, core.D3, core.F3},
		},
		{
			name:            "white push and en passant both available",
			from:            core.D5,
			sideToMove:      core.WHITE,
			setupBoard:      func(b *core.Board) {},
			enPassantTarget: core.E6,
			expectedTos:     []core.Position{core.D6, core.E6},
			enPassantTos:    []core.Position{core.E6},
		},
		{
			name:       "fully blocked pawn with no captures yields no moves",
			from:       core.E4,
			sideToMove: core.WHITE,
			setupBoard: func(b *core.Board) {
				b[core.E5] = core.Square{Occupied: true, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE}}
			},
			expectedTos: []core.Position{},
		},

		// ── defensive: impossible positions ──────────────────────────────────────
		// These should never occur in real chess (promotion happens on arrival),
		// but the dispatcher guards with `if !ok { return []core.Move{} }`. These
		// tests lock in that contract so a future refactor that removes the guard
		// gets caught immediately.

		{
			name:        "white pawn on last rank returns no moves (defensive)",
			from:        core.E8,
			sideToMove:  core.WHITE,
			setupBoard:  func(b *core.Board) {},
			expectedTos: []core.Position{},
		},
		{
			name:        "black pawn on last rank returns no moves (defensive)",
			from:        core.E1,
			sideToMove:  core.BLACK,
			setupBoard:  func(b *core.Board) {},
			expectedTos: []core.Position{},
		},
	}

	pawn := Pawn{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var board core.Board
			tt.setupBoard(&board)

			ctx := core.MoveContext{
				BoardContext:    core.BoardContext{Board: &board},
				SideToMove:      tt.sideToMove,
				EnPassantTarget: tt.enPassantTarget,
			}

			got := pawn.PseudoLegalMoves(tt.from, ctx)

			// Each promotion square yields 4 Moves (one per piece type).
			// expectedTos lists each destination once, so add 3 extra per promotion square.
			promoSquares := make(map[core.Position]bool, len(tt.promotionTos))
			for _, pos := range tt.promotionTos {
				promoSquares[pos] = true
			}
			expectedCount := len(tt.expectedTos) + len(tt.promotionTos)*3

			if len(got) != expectedCount {
				t.Fatalf("got %d moves, want %d — got %v", len(got), expectedCount, got)
			}

			// All four promotion piece types must be present for every promotion square.
			for _, pos := range tt.promotionTos {
				assertPromotions(t, got, pos)
			}

			// Every declared en passant destination must carry MoveType EN_PASSANT.
			for _, pos := range tt.enPassantTos {
				assertMoveType(t, got, pos, core.EN_PASSANT)
			}

			// Non-promotion, non-en-passant moves must be NORMAL.
			epSquares := make(map[core.Position]bool, len(tt.enPassantTos))
			for _, pos := range tt.enPassantTos {
				epSquares[pos] = true
			}
			for _, move := range got {
				if promoSquares[move.To] || epSquares[move.To] {
					continue
				}
				if move.Type != core.NORMAL {
					t.Errorf("move to %v: got type %v, want NORMAL", move.To, move.Type)
				}
			}

			// ── field validation ─────────────────────────────────────────────────
			// Every move must carry its source, the moving piece, and correct
			// capture info derived from the board state (or the EP target).

			expectedMover := core.Piece{Type: core.PAWN, Color: tt.sideToMove}
			enemy := tt.sideToMove.Opponent()

			// build expected captures map: destination -> captured piece
			// normal captures: the enemy piece on the destination
			// en passant: always a pawn of the enemy color (target square is empty)
			expectedCaptures := map[core.Position]core.Piece{}
			for _, pos := range tt.expectedTos {
				sq := board[pos]
				if sq.Occupied && sq.Piece.Color == enemy {
					expectedCaptures[pos] = sq.Piece
				}
			}
			for _, pos := range tt.enPassantTos {
				expectedCaptures[pos] = core.Piece{Type: core.PAWN, Color: enemy}
			}

			for _, move := range got {
				// every move must carry its source
				if move.From != tt.from {
					t.Errorf("move to %v: From=%v, want %v", move.To, move.From, tt.from)
				}
				// every move must carry the moving piece
				if move.Piece != expectedMover {
					t.Errorf("move to %v: Piece=%v, want %v", move.To, move.Piece, expectedMover)
				}

				// capture info must match expectations
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
			}

			// Every expected destination must appear at least once in the results.
			gotSet := make(map[core.Position]struct{}, len(got))
			for _, move := range got {
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

func TestPawnDirection(t *testing.T) {
	p := Pawn{}

	tests := []struct {
		name      string
		color     core.PieceColor
		wantStep  int8
		wantStart core.Rank
		wantLast  core.Rank
	}{
		{
			name:      "white moves up the board",
			color:     core.WHITE,
			wantStep:  1,
			wantStart: core.RANK_2,
			wantLast:  core.RANK_8,
		},
		{
			name:      "black moves down the board",
			color:     core.BLACK,
			wantStep:  -1,
			wantStart: core.RANK_7,
			wantLast:  core.RANK_1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			step, start, last := p.direction(tt.color)

			if step != tt.wantStep {
				t.Errorf("step: got %d, want %d", step, tt.wantStep)
			}
			if start != tt.wantStart {
				t.Errorf("start: got %v, want %v", start, tt.wantStart)
			}
			if last != tt.wantLast {
				t.Errorf("last: got %v, want %v", last, tt.wantLast)
			}
		})
	}
}
