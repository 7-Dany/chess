package engine

import (
	"testing"

	"github.com/7-Dany/chess/core"
)

func TestGetPseudoLegalMoves(t *testing.T) {
	defaultSides := [2]core.SideState{
		{KingPosition: core.E1, CanCastleKingSide: true, CanCastleQueenSide: true},
		{KingPosition: core.E8, CanCastleKingSide: true, CanCastleQueenSide: true},
	}

	tests := []struct {
		name         string
		setupBoard   func(*core.Board)
		sideToMove   core.PieceColor
		position     core.Position
		wantEmpty    bool
		wantNonEmpty bool
	}{
		{
			name: "empty square returns no moves",
			setupBoard: func(b *core.Board) {
				b[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
			},
			sideToMove: core.WHITE,
			position:   core.E4,
			wantEmpty:  true,
		},
		{
			name: "enemy piece returns no moves (out of turn)",
			setupBoard: func(b *core.Board) {
				b[core.E8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			position:   core.E8,
			wantEmpty:  true,
		},
		{
			name: "own piece delegates to piece implementation",
			setupBoard: func(b *core.Board) {
				b[core.B1] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
			},
			sideToMove:   core.WHITE,
			position:     core.B1,
			wantNonEmpty: true,
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
					Sides:        defaultSides,
				},
			}
			moves := engine.GetPseudoLegalMoves(tt.position, ctx)
			if tt.wantEmpty && len(moves) != 0 {
				t.Errorf("expected empty, got %d moves", len(moves))
			}
			if tt.wantNonEmpty && len(moves) == 0 {
				t.Errorf("expected non-empty, got 0 moves")
			}
		})
	}
}

func TestCastlingMoves(t *testing.T) {
	defaultSides := [2]core.SideState{
		{KingPosition: core.E1, CanCastleKingSide: true, CanCastleQueenSide: true},
		{KingPosition: core.E8, CanCastleKingSide: true, CanCastleQueenSide: true},
	}
	kingSideOnly := [2]core.SideState{
		{KingPosition: core.E1, CanCastleKingSide: true},
		{KingPosition: core.E8, CanCastleKingSide: true, CanCastleQueenSide: true},
	}
	queenSideOnly := [2]core.SideState{
		{KingPosition: core.E1, CanCastleQueenSide: true},
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
		position   core.Position
		wantCount  int
		checks     []moveCheck
	}{
		{
			name: "king not on E file returns no castling",
			setupBoard: func(b *core.Board) {
				b[core.D1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
				b[core.A1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
				b[core.H1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			sideToMove: core.WHITE,
			sides: [2]core.SideState{
				{KingPosition: core.D1, CanCastleKingSide: true, CanCastleQueenSide: true},
				defaultSides[1],
			},
			position:  core.D1,
			wantCount: 0,
		},
		{
			name: "king in check returns no castling",
			setupBoard: func(b *core.Board) {
				b[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
				b[core.A1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
				b[core.H1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
				b[core.E8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			position:   core.E1,
			wantCount:  0,
		},
		{
			name: "both sides available returns two castling moves",
			setupBoard: func(b *core.Board) {
				b[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
				b[core.A1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
				b[core.H1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
				b[core.E8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			position:   core.E1,
			wantCount:  2,
			checks: []moveCheck{
				{core.E1, core.G1, core.CASTLING, true},
				{core.E1, core.C1, core.CASTLING, true},
			},
		},
		{
			name: "only king-side right returns one castling move",
			setupBoard: func(b *core.Board) {
				b[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
				b[core.H1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
				b[core.E8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			sides:      kingSideOnly,
			position:   core.E1,
			wantCount:  1,
			checks: []moveCheck{
				{core.E1, core.G1, core.CASTLING, true},
				{core.E1, core.C1, core.CASTLING, false},
			},
		},
		{
			name: "only queen-side right returns one castling move",
			setupBoard: func(b *core.Board) {
				b[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
				b[core.A1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
				b[core.E8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			sides:      queenSideOnly,
			position:   core.E1,
			wantCount:  1,
			checks: []moveCheck{
				{core.E1, core.C1, core.CASTLING, true},
				{core.E1, core.G1, core.CASTLING, false},
			},
		},
		{
			name: "black king-side castling on rank 8",
			setupBoard: func(b *core.Board) {
				b[core.E8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
				b[core.H8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
				b[core.A8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
				b[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
			},
			sideToMove: core.BLACK,
			sides:      defaultSides,
			position:   core.E8,
			wantCount:  2,
			checks: []moveCheck{
				{core.E8, core.G8, core.CASTLING, true},
				{core.E8, core.C8, core.CASTLING, true},
			},
		},
		{
			name: "both sides blocked returns no castling",
			setupBoard: func(b *core.Board) {
				b[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
				b[core.A1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
				b[core.H1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
				b[core.F1] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
				b[core.D1] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.WHITE})
				b[core.E8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			position:   core.E1,
			wantCount:  0,
		},
		{
			name: "queen-side castling removed when B1 occupied",
			setupBoard: func(b *core.Board) {
				b[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
				b[core.A1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
				b[core.H1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
				b[core.B1] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
				b[core.E8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			position:   core.E1,
			wantCount:  1,
			checks: []moveCheck{
				{core.E1, core.G1, core.CASTLING, true},
				{core.E1, core.C1, core.CASTLING, false},
			},
		},
		{
			name: "queen-side castling removed when enemy on C1",
			setupBoard: func(b *core.Board) {
				b[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
				b[core.A1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
				b[core.H1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
				b[core.C1] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.BLACK})
				b[core.E8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			position:   core.E1,
			wantCount:  1,
			checks: []moveCheck{
				{core.E1, core.G1, core.CASTLING, true},
				{core.E1, core.C1, core.CASTLING, false},
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
					BoardContext: core.BoardContext{Board: &board},
					SideToMove:   tt.sideToMove,
					Sides:        tt.sides,
				},
			}
			moves := engine.castlingMoves(nil, tt.position, ctx)
			if len(moves) != tt.wantCount {
				t.Errorf("count = %d, want %d", len(moves), tt.wantCount)
			}
			checkMoves(moves, tt.checks)
		})
	}
}

func TestCanCastleKingSide(t *testing.T) {
	defaultSides := [2]core.SideState{
		{KingPosition: core.E1, CanCastleKingSide: true, CanCastleQueenSide: true},
		{KingPosition: core.E8, CanCastleKingSide: true, CanCastleQueenSide: true},
	}
	noKingSide := [2]core.SideState{
		{KingPosition: core.E1, CanCastleQueenSide: true},
		{KingPosition: core.E8, CanCastleKingSide: true, CanCastleQueenSide: true},
	}

	tests := []struct {
		name       string
		setupBoard func(*core.Board)
		sideToMove core.PieceColor
		sides      [2]core.SideState
		rank       core.Rank
		want       bool
	}{
		{
			name:       "no king-side right",
			setupBoard: func(b *core.Board) {},
			sideToMove: core.WHITE,
			sides:      noKingSide,
			rank:       core.RANK_1,
			want:       false,
		},
		{
			name: "F1 occupied by own piece",
			setupBoard: func(b *core.Board) {
				b[core.F1] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			rank:       core.RANK_1,
			want:       false,
		},
		{
			name: "G1 occupied by own piece",
			setupBoard: func(b *core.Board) {
				b[core.G1] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			rank:       core.RANK_1,
			want:       false,
		},
		{
			name: "F1 occupied by enemy piece",
			setupBoard: func(b *core.Board) {
				b[core.F1] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			rank:       core.RANK_1,
			want:       false,
		},
		{
			name: "F1 attacked by enemy rook",
			setupBoard: func(b *core.Board) {
				b[core.F8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			rank:       core.RANK_1,
			want:       false,
		},
		{
			name: "G1 attacked by enemy rook",
			setupBoard: func(b *core.Board) {
				b[core.G8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			rank:       core.RANK_1,
			want:       false,
		},
		{
			name: "F1 attacked by enemy bishop",
			setupBoard: func(b *core.Board) {
				b[core.A6] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			rank:       core.RANK_1,
			want:       false,
		},
		{
			name:       "all clear white rank 1",
			setupBoard: func(b *core.Board) {},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			rank:       core.RANK_1,
			want:       true,
		},
		{
			name: "F8 occupied by own piece (black rank 8)",
			setupBoard: func(b *core.Board) {
				b[core.F8] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.BLACK})
			},
			sideToMove: core.BLACK,
			sides:      defaultSides,
			rank:       core.RANK_8,
			want:       false,
		},
		{
			name: "F8 attacked by white rook (black rank 8)",
			setupBoard: func(b *core.Board) {
				b[core.F1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			sideToMove: core.BLACK,
			sides:      defaultSides,
			rank:       core.RANK_8,
			want:       false,
		},
		{
			name:       "all clear black rank 8",
			setupBoard: func(b *core.Board) {},
			sideToMove: core.BLACK,
			sides:      defaultSides,
			rank:       core.RANK_8,
			want:       true,
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
			got := engine.canCastleKingSide(tt.rank, ctx)
			if got != tt.want {
				t.Errorf("canCastleKingSide = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCanCastleQueenSide(t *testing.T) {
	defaultSides := [2]core.SideState{
		{KingPosition: core.E1, CanCastleKingSide: true, CanCastleQueenSide: true},
		{KingPosition: core.E8, CanCastleKingSide: true, CanCastleQueenSide: true},
	}
	noQueenSide := [2]core.SideState{
		{KingPosition: core.E1, CanCastleKingSide: true},
		{KingPosition: core.E8, CanCastleKingSide: true, CanCastleQueenSide: true},
	}

	tests := []struct {
		name       string
		setupBoard func(*core.Board)
		sideToMove core.PieceColor
		sides      [2]core.SideState
		rank       core.Rank
		want       bool
	}{
		{
			name:       "no queen-side right",
			setupBoard: func(b *core.Board) {},
			sideToMove: core.WHITE,
			sides:      noQueenSide,
			rank:       core.RANK_1,
			want:       false,
		},
		{
			name: "B1 occupied by own piece",
			setupBoard: func(b *core.Board) {
				b[core.B1] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			rank:       core.RANK_1,
			want:       false,
		},
		{
			name: "C1 occupied by own piece",
			setupBoard: func(b *core.Board) {
				b[core.C1] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			rank:       core.RANK_1,
			want:       false,
		},
		{
			name: "D1 occupied by own piece",
			setupBoard: func(b *core.Board) {
				b[core.D1] = core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.WHITE})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			rank:       core.RANK_1,
			want:       false,
		},
		{
			name: "C1 occupied by enemy piece blocks castling",
			setupBoard: func(b *core.Board) {
				b[core.C1] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			rank:       core.RANK_1,
			want:       false,
		},
		{
			name: "D1 attacked by enemy rook",
			setupBoard: func(b *core.Board) {
				b[core.D8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			rank:       core.RANK_1,
			want:       false,
		},
		{
			name: "C1 attacked by enemy rook",
			setupBoard: func(b *core.Board) {
				b[core.C8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			rank:       core.RANK_1,
			want:       false,
		},
		{
			name: "D1 attacked by enemy bishop",
			setupBoard: func(b *core.Board) {
				b[core.A4] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			rank:       core.RANK_1,
			want:       false,
		},
		{
			name: "B1 attacked but C1 and D1 safe — castling allowed",
			setupBoard: func(b *core.Board) {
				b[core.B8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			rank:       core.RANK_1,
			want:       true,
		},
		{
			name:       "all clear white rank 1",
			setupBoard: func(b *core.Board) {},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			rank:       core.RANK_1,
			want:       true,
		},
		{
			name: "B8 occupied by own piece (black rank 8)",
			setupBoard: func(b *core.Board) {
				b[core.B8] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.BLACK})
			},
			sideToMove: core.BLACK,
			sides:      defaultSides,
			rank:       core.RANK_8,
			want:       false,
		},
		{
			name: "D8 attacked by white rook (black rank 8)",
			setupBoard: func(b *core.Board) {
				b[core.D1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			sideToMove: core.BLACK,
			sides:      defaultSides,
			rank:       core.RANK_8,
			want:       false,
		},
		{
			name: "B1 attacked but C8 and D8 safe — black castling allowed",
			setupBoard: func(b *core.Board) {
				b[core.B1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			sideToMove: core.BLACK,
			sides:      defaultSides,
			rank:       core.RANK_8,
			want:       true,
		},
		{
			name:       "all clear black rank 8",
			setupBoard: func(b *core.Board) {},
			sideToMove: core.BLACK,
			sides:      defaultSides,
			rank:       core.RANK_8,
			want:       true,
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
			got := engine.canCastleQueenSide(tt.rank, ctx)
			if got != tt.want {
				t.Errorf("canCastleQueenSide = %v, want %v", got, tt.want)
			}
		})
	}
}
