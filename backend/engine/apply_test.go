package engine

import (
	"testing"

	"github.com/7-Dany/chess/core"
)

func TestApply(t *testing.T) {
	defaultSides := [2]core.SideState{
		{KingPosition: core.E1, CastlingRights: core.CastlingRights{KingSide: true, QueenSide: true}},
		{KingPosition: core.E8, CastlingRights: core.CastlingRights{KingSide: true, QueenSide: true}},
	}

	tests := []struct {
		name        string
		setupBoard  func(*core.Board)
		sideToMove  core.PieceColor
		sides       [2]core.SideState
		inputEP     core.Position
		move        core.Move
		expectAt    map[core.Position]core.Square
		expectEmpty []core.Position
		expectEP    core.Position
		expectSides [2]core.SideState
		checkSnap   bool
		expectSnap  core.Snapshot
	}{
		// ==================== Normal moves ====================
		{
			name: "normal knight move relocates piece",
			setupBoard: func(b *core.Board) {
				b[core.B1] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			move: core.Move{
				Type: core.NORMAL, Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE},
				From: core.B1, To: core.C3,
			},
			expectAt: map[core.Position]core.Square{
				core.C3: core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE}),
			},
			expectEmpty: []core.Position{core.B1},
			expectEP:    core.NoPosition,
			expectSides: defaultSides,
		},
		{
			name: "normal king move updates position and clears castling rights",
			setupBoard: func(b *core.Board) {
				b[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			move: core.Move{
				Type: core.NORMAL, Piece: core.Piece{Type: core.KING, Color: core.WHITE},
				From: core.E1, To: core.F1,
			},
			expectAt: map[core.Position]core.Square{
				core.F1: core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE}),
			},
			expectEmpty: []core.Position{core.E1},
			expectEP:    core.NoPosition,
			expectSides: [2]core.SideState{
				{KingPosition: core.F1, CastlingRights: core.CastlingRights{}},
				defaultSides[1],
			},
		},
		{
			name: "rook move from A file clears queen-side right only",
			setupBoard: func(b *core.Board) {
				b[core.A1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			move: core.Move{
				Type: core.NORMAL, Piece: core.Piece{Type: core.ROOK, Color: core.WHITE},
				From: core.A1, To: core.A3,
			},
			expectAt: map[core.Position]core.Square{
				core.A3: core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE}),
			},
			expectEmpty: []core.Position{core.A1},
			expectEP:    core.NoPosition,
			expectSides: [2]core.SideState{
				{KingPosition: core.E1, CastlingRights: core.CastlingRights{KingSide: true}},
				defaultSides[1],
			},
		},
		{
			name: "rook move from H file clears king-side right only",
			setupBoard: func(b *core.Board) {
				b[core.H1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			move: core.Move{
				Type: core.NORMAL, Piece: core.Piece{Type: core.ROOK, Color: core.WHITE},
				From: core.H1, To: core.H3,
			},
			expectAt: map[core.Position]core.Square{
				core.H3: core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE}),
			},
			expectEmpty: []core.Position{core.H1},
			expectEP:    core.NoPosition,
			expectSides: [2]core.SideState{
				{KingPosition: core.E1, CastlingRights: core.CastlingRights{QueenSide: true}},
				defaultSides[1],
			},
		},
		{
			name: "rook move from non-home file preserves all rights",
			setupBoard: func(b *core.Board) {
				b[core.C3] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			move: core.Move{
				Type: core.NORMAL, Piece: core.Piece{Type: core.ROOK, Color: core.WHITE},
				From: core.C3, To: core.C5,
			},
			expectAt: map[core.Position]core.Square{
				core.C5: core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE}),
			},
			expectEmpty: []core.Position{core.C3},
			expectEP:    core.NoPosition,
			expectSides: defaultSides,
		},
		{
			name: "black rook move from A8 clears queen-side right",
			setupBoard: func(b *core.Board) {
				b[core.A8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
			},
			sideToMove: core.BLACK,
			sides:      defaultSides,
			move: core.Move{
				Type: core.NORMAL, Piece: core.Piece{Type: core.ROOK, Color: core.BLACK},
				From: core.A8, To: core.A6,
			},
			expectAt: map[core.Position]core.Square{
				core.A6: core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK}),
			},
			expectEmpty: []core.Position{core.A8},
			expectEP:    core.NoPosition,
			expectSides: [2]core.SideState{
				defaultSides[0],
				{KingPosition: core.E8, CastlingRights: core.CastlingRights{KingSide: true}},
			},
		},
		{
			name: "black rook move from H8 clears king-side right",
			setupBoard: func(b *core.Board) {
				b[core.H8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
			},
			sideToMove: core.BLACK,
			sides:      defaultSides,
			move: core.Move{
				Type: core.NORMAL, Piece: core.Piece{Type: core.ROOK, Color: core.BLACK},
				From: core.H8, To: core.H6,
			},
			expectAt: map[core.Position]core.Square{
				core.H6: core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK}),
			},
			expectEmpty: []core.Position{core.H8},
			expectEP:    core.NoPosition,
			expectSides: [2]core.SideState{
				defaultSides[0],
				{KingPosition: core.E8, CastlingRights: core.CastlingRights{QueenSide: true}},
			},
		},

		// ==================== Captures ====================
		{
			name: "capture replaces piece on destination",
			setupBoard: func(b *core.Board) {
				b[core.E4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
				b[core.D5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			move: core.Move{
				Type: core.NORMAL, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE},
				From: core.E4, To: core.D5,
				HasCapture: true, Captured: core.Piece{Type: core.PAWN, Color: core.BLACK},
			},
			expectAt: map[core.Position]core.Square{
				core.D5: core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE}),
			},
			expectEmpty: []core.Position{core.E4},
			expectEP:    core.NoPosition,
			expectSides: defaultSides,
		},
		{
			name: "capturing rook on A8 clears opponent queen-side right",
			setupBoard: func(b *core.Board) {
				b[core.A6] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE})
				b[core.A8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			move: core.Move{
				Type: core.NORMAL, Piece: core.Piece{Type: core.BISHOP, Color: core.WHITE},
				From: core.A6, To: core.A8,
				HasCapture: true, Captured: core.Piece{Type: core.ROOK, Color: core.BLACK},
			},
			expectAt: map[core.Position]core.Square{
				core.A8: core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.WHITE}),
			},
			expectEmpty: []core.Position{core.A6},
			expectEP:    core.NoPosition,
			expectSides: [2]core.SideState{
				defaultSides[0],
				{KingPosition: core.E8, CastlingRights: core.CastlingRights{KingSide: true}},
			},
		},
		{
			name: "capturing rook on H1 clears opponent king-side right",
			setupBoard: func(b *core.Board) {
				b[core.H3] = core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.BLACK})
				b[core.H1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			sideToMove: core.BLACK,
			sides:      defaultSides,
			move: core.Move{
				Type: core.NORMAL, Piece: core.Piece{Type: core.BISHOP, Color: core.BLACK},
				From: core.H3, To: core.H1,
				HasCapture: true, Captured: core.Piece{Type: core.ROOK, Color: core.WHITE},
			},
			expectAt: map[core.Position]core.Square{
				core.H1: core.NewSquare(core.Piece{Type: core.BISHOP, Color: core.BLACK}),
			},
			expectEmpty: []core.Position{core.H3},
			expectEP:    core.NoPosition,
			expectSides: [2]core.SideState{
				{KingPosition: core.E1, CastlingRights: core.CastlingRights{QueenSide: true}},
				defaultSides[1],
			},
		},
		{
			name: "capturing non-rook on A file does not clear rights",
			setupBoard: func(b *core.Board) {
				b[core.B5] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
				b[core.A6] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			move: core.Move{
				Type: core.NORMAL, Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE},
				From: core.B5, To: core.A6,
				HasCapture: true, Captured: core.Piece{Type: core.PAWN, Color: core.BLACK},
			},
			expectAt: map[core.Position]core.Square{
				core.A6: core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE}),
			},
			expectEmpty: []core.Position{core.B5},
			expectEP:    core.NoPosition,
			expectSides: defaultSides,
		},

		// ==================== En Passant ====================
		{
			name: "white en passant capture removes captured pawn",
			setupBoard: func(b *core.Board) {
				b[core.D5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
				b[core.E5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			inputEP:    core.E6,
			sides:      defaultSides,
			move: core.Move{
				Type: core.EN_PASSANT, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE},
				From: core.D5, To: core.E6,
				HasCapture: true, Captured: core.Piece{Type: core.PAWN, Color: core.BLACK},
			},
			expectAt: map[core.Position]core.Square{
				core.E6: core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE}),
			},
			expectEmpty: []core.Position{core.D5, core.E5},
			expectEP:    core.NoPosition,
			expectSides: defaultSides,
		},
		{
			name: "black en passant capture removes captured pawn",
			setupBoard: func(b *core.Board) {
				b[core.E4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
				b[core.D4] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
			},
			sideToMove: core.BLACK,
			inputEP:    core.E3,
			sides:      defaultSides,
			move: core.Move{
				Type: core.EN_PASSANT, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK},
				From: core.D4, To: core.E3,
				HasCapture: true, Captured: core.Piece{Type: core.PAWN, Color: core.WHITE},
			},
			expectAt: map[core.Position]core.Square{
				core.E3: core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK}),
			},
			expectEmpty: []core.Position{core.D4, core.E4},
			expectEP:    core.NoPosition,
			expectSides: defaultSides,
		},
		{
			name: "en passant on A file does not clear queen-side right",
			setupBoard: func(b *core.Board) {
				b[core.B5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
				b[core.A5] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			inputEP:    core.A6,
			sides:      defaultSides,
			move: core.Move{
				Type: core.EN_PASSANT, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE},
				From: core.B5, To: core.A6,
				HasCapture: true, Captured: core.Piece{Type: core.PAWN, Color: core.BLACK},
			},
			expectAt: map[core.Position]core.Square{
				core.A6: core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE}),
			},
			expectEmpty: []core.Position{core.B5, core.A5},
			expectEP:    core.NoPosition,
			expectSides: defaultSides,
		},

		// ==================== Promotion ====================
		{
			name: "white pawn promotes to queen on rank 8",
			setupBoard: func(b *core.Board) {
				b[core.E7] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			move: core.Move{
				Type: core.PROMOTION, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE},
				From: core.E7, To: core.E8,
				PromoteTo: core.QUEEN,
			},
			expectAt: map[core.Position]core.Square{
				core.E8: core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.WHITE}),
			},
			expectEmpty: []core.Position{core.E7},
			expectEP:    core.NoPosition,
			expectSides: defaultSides,
		},
		{
			name: "black pawn promotes to knight on rank 1",
			setupBoard: func(b *core.Board) {
				b[core.D2] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
			},
			sideToMove: core.BLACK,
			sides:      defaultSides,
			move: core.Move{
				Type: core.PROMOTION, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK},
				From: core.D2, To: core.D1,
				PromoteTo: core.KNIGHT,
			},
			expectAt: map[core.Position]core.Square{
				core.D1: core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.BLACK}),
			},
			expectEmpty: []core.Position{core.D2},
			expectEP:    core.NoPosition,
			expectSides: defaultSides,
		},
		{
			name: "promotion with capture on non-home file preserves rights",
			setupBoard: func(b *core.Board) {
				b[core.E7] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
				b[core.D8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			move: core.Move{
				Type: core.PROMOTION, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE},
				From: core.E7, To: core.D8,
				PromoteTo:  core.QUEEN,
				HasCapture: true, Captured: core.Piece{Type: core.ROOK, Color: core.BLACK},
			},
			expectAt: map[core.Position]core.Square{
				core.D8: core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.WHITE}),
			},
			expectEmpty: []core.Position{core.E7},
			expectEP:    core.NoPosition,
			expectSides: defaultSides,
		},
		{
			name: "promotion capture on H8 clears opponent king-side right",
			setupBoard: func(b *core.Board) {
				b[core.G7] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
				b[core.H8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			move: core.Move{
				Type: core.PROMOTION, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE},
				From: core.G7, To: core.H8,
				PromoteTo:  core.QUEEN,
				HasCapture: true, Captured: core.Piece{Type: core.ROOK, Color: core.BLACK},
			},
			expectAt: map[core.Position]core.Square{
				core.H8: core.NewSquare(core.Piece{Type: core.QUEEN, Color: core.WHITE}),
			},
			expectEmpty: []core.Position{core.G7},
			expectEP:    core.NoPosition,
			expectSides: [2]core.SideState{
				defaultSides[0],
				{KingPosition: core.E8, CastlingRights: core.CastlingRights{QueenSide: true}},
			},
		},

		// ==================== Castling ====================
		{
			name: "white king-side castling",
			setupBoard: func(b *core.Board) {
				b[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
				b[core.H1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			move: core.Move{
				Type: core.CASTLING, Piece: core.Piece{Type: core.KING, Color: core.WHITE},
				From: core.E1, To: core.G1,
			},
			expectAt: map[core.Position]core.Square{
				core.G1: core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE}),
				core.F1: core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE}),
			},
			expectEmpty: []core.Position{core.E1, core.H1},
			expectEP:    core.NoPosition,
			expectSides: [2]core.SideState{
				{KingPosition: core.G1, CastlingRights: core.CastlingRights{}},
				defaultSides[1],
			},
		},
		{
			name: "white queen-side castling",
			setupBoard: func(b *core.Board) {
				b[core.E1] = core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE})
				b[core.A1] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			move: core.Move{
				Type: core.CASTLING, Piece: core.Piece{Type: core.KING, Color: core.WHITE},
				From: core.E1, To: core.C1,
			},
			expectAt: map[core.Position]core.Square{
				core.C1: core.NewSquare(core.Piece{Type: core.KING, Color: core.WHITE}),
				core.D1: core.NewSquare(core.Piece{Type: core.ROOK, Color: core.WHITE}),
			},
			expectEmpty: []core.Position{core.E1, core.A1},
			expectEP:    core.NoPosition,
			expectSides: [2]core.SideState{
				{KingPosition: core.C1, CastlingRights: core.CastlingRights{}},
				defaultSides[1],
			},
		},
		{
			name: "black king-side castling",
			setupBoard: func(b *core.Board) {
				b[core.E8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
				b[core.H8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
			},
			sideToMove: core.BLACK,
			sides:      defaultSides,
			move: core.Move{
				Type: core.CASTLING, Piece: core.Piece{Type: core.KING, Color: core.BLACK},
				From: core.E8, To: core.G8,
			},
			expectAt: map[core.Position]core.Square{
				core.G8: core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK}),
				core.F8: core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK}),
			},
			expectEmpty: []core.Position{core.E8, core.H8},
			expectEP:    core.NoPosition,
			expectSides: [2]core.SideState{
				defaultSides[0],
				{KingPosition: core.G8, CastlingRights: core.CastlingRights{}},
			},
		},
		{
			name: "black queen-side castling",
			setupBoard: func(b *core.Board) {
				b[core.E8] = core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK})
				b[core.A8] = core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK})
			},
			sideToMove: core.BLACK,
			sides:      defaultSides,
			move: core.Move{
				Type: core.CASTLING, Piece: core.Piece{Type: core.KING, Color: core.BLACK},
				From: core.E8, To: core.C8,
			},
			expectAt: map[core.Position]core.Square{
				core.C8: core.NewSquare(core.Piece{Type: core.KING, Color: core.BLACK}),
				core.D8: core.NewSquare(core.Piece{Type: core.ROOK, Color: core.BLACK}),
			},
			expectEmpty: []core.Position{core.E8, core.A8},
			expectEP:    core.NoPosition,
			expectSides: [2]core.SideState{
				defaultSides[0],
				{KingPosition: core.C8, CastlingRights: core.CastlingRights{}},
			},
		},

		// ==================== Double pawn push ====================
		{
			name: "white pawn double push sets en passant target",
			setupBoard: func(b *core.Board) {
				b[core.E2] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			move: core.Move{
				Type: core.NORMAL, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE},
				From: core.E2, To: core.E4,
			},
			expectAt: map[core.Position]core.Square{
				core.E4: core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE}),
			},
			expectEmpty: []core.Position{core.E2},
			expectEP:    core.E3,
			expectSides: defaultSides,
		},
		{
			name: "black pawn double push sets en passant target",
			setupBoard: func(b *core.Board) {
				b[core.D7] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK})
			},
			sideToMove: core.BLACK,
			sides:      defaultSides,
			move: core.Move{
				Type: core.NORMAL, Piece: core.Piece{Type: core.PAWN, Color: core.BLACK},
				From: core.D7, To: core.D5,
			},
			expectAt: map[core.Position]core.Square{
				core.D5: core.NewSquare(core.Piece{Type: core.PAWN, Color: core.BLACK}),
			},
			expectEmpty: []core.Position{core.D7},
			expectEP:    core.D6,
			expectSides: defaultSides,
		},
		{
			name: "white pawn double push from A file sets en passant on A3",
			setupBoard: func(b *core.Board) {
				b[core.A2] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			move: core.Move{
				Type: core.NORMAL, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE},
				From: core.A2, To: core.A4,
			},
			expectAt: map[core.Position]core.Square{
				core.A4: core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE}),
			},
			expectEmpty: []core.Position{core.A2},
			expectEP:    core.A3,
			expectSides: defaultSides,
		},
		{
			name: "single pawn push clears en passant target",
			setupBoard: func(b *core.Board) {
				b[core.E2] = core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE})
			},
			sideToMove: core.WHITE,
			sides:      defaultSides,
			move: core.Move{
				Type: core.NORMAL, Piece: core.Piece{Type: core.PAWN, Color: core.WHITE},
				From: core.E2, To: core.E3,
			},
			expectAt: map[core.Position]core.Square{
				core.E3: core.NewSquare(core.Piece{Type: core.PAWN, Color: core.WHITE}),
			},
			expectEmpty: []core.Position{core.E2},
			expectEP:    core.NoPosition,
			expectSides: defaultSides,
		},
		{
			name: "non-pawn move clears previous en passant target",
			setupBoard: func(b *core.Board) {
				b[core.B1] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
			},
			sideToMove: core.WHITE,
			inputEP:    core.E3,
			sides:      defaultSides,
			move: core.Move{
				Type: core.NORMAL, Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE},
				From: core.B1, To: core.C3,
			},
			expectAt: map[core.Position]core.Square{
				core.C3: core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE}),
			},
			expectEmpty: []core.Position{core.B1},
			expectEP:    core.NoPosition,
			expectSides: defaultSides,
		},

		// ==================== Snapshot ====================
		{
			name: "snapshot captures previous sides and en passant target",
			setupBoard: func(b *core.Board) {
				b[core.B1] = core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE})
			},
			sideToMove: core.WHITE,
			inputEP:    core.E3,
			sides:      defaultSides,
			move: core.Move{
				Type: core.NORMAL, Piece: core.Piece{Type: core.KNIGHT, Color: core.WHITE},
				From: core.B1, To: core.C3,
			},
			expectAt: map[core.Position]core.Square{
				core.C3: core.NewSquare(core.Piece{Type: core.KNIGHT, Color: core.WHITE}),
			},
			expectEmpty: []core.Position{core.B1},
			expectEP:    core.NoPosition,
			expectSides: defaultSides,
			checkSnap:   true,
			expectSnap: core.Snapshot{
				PreviousSides:           defaultSides,
				PreviousEnPassantTarget: core.E3,
			},
		},
	}

	engine := NewDefaultEngine()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var board core.Board
			tt.setupBoard(&board)

			ctx := &core.TurnContext{
				MoveContext: core.MoveContext{
					BoardContext:    core.BoardContext{Board: &board},
					SideToMove:      tt.sideToMove,
					Sides:           tt.sides,
					EnPassantTarget: tt.inputEP,
				},
			}

			snap := engine.Apply(ctx, tt.move)

			// Apply no longer touches SideToMove — verify it's unchanged.
			if ctx.SideToMove != tt.sideToMove {
				t.Errorf("SideToMove = %v, want %v (Apply should not flip side)", ctx.SideToMove, tt.sideToMove)
			}

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

			if ctx.EnPassantTarget != tt.expectEP {
				t.Errorf("EnPassantTarget = %v, want %v", ctx.EnPassantTarget, tt.expectEP)
			}

			if ctx.Sides != tt.expectSides {
				t.Errorf("Sides = %+v, want %+v", ctx.Sides, tt.expectSides)
			}

			// Snapshot check — only when explicitly requested via checkSnap.
			if tt.checkSnap {
				if snap.PreviousSides != tt.expectSnap.PreviousSides {
					t.Errorf("snap.PreviousSides = %+v, want %+v", snap.PreviousSides, tt.expectSnap.PreviousSides)
				}
				if snap.PreviousEnPassantTarget != tt.expectSnap.PreviousEnPassantTarget {
					t.Errorf("snap.PreviousEnPassantTarget = %v, want %v",
						snap.PreviousEnPassantTarget, tt.expectSnap.PreviousEnPassantTarget)
				}
			}
		})
	}
}
