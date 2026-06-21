package core

import (
	"fmt"
	"strings"
)

type SideState struct {
	KingPosition       Position
	CanCastleKingSide  bool
	CanCastleQueenSide bool
}

// ClearCastlingRights revokes both rights — called when the king moves or castles.
func (s *SideState) ClearCastlingRights() {
	s.CanCastleKingSide = false
	s.CanCastleQueenSide = false
}

// ClearCastlingRight revokes the single right tied to a rook's home file.
// Used when that rook moves, or is captured on its home square.
func (s *SideState) ClearCastlingRight(file File) {
	switch file {
	case FILE_A:
		s.CanCastleQueenSide = false
	case FILE_H:
		s.CanCastleKingSide = false
	}
}

type BoardContext struct {
	Board *Board
}

type MoveContext struct {
	BoardContext
	SideToMove      PieceColor
	Sides           [2]SideState
	EnPassantTarget Position
}

// ForfeitCastlingRight clears the single castling right — if any — forfeited
// by this move. A rook moving from its back rank, or a rook captured on its
// back rank, each forfeits one right. King moves and castling forfeit all rights
func (ctx *MoveContext) ForfeitCastlingRight(move Move) {
	// if the piece is rook and its position not from start forfeit the castle right for that file
	if move.Piece.Type == ROOK && move.From.Rank() == move.Piece.Color.KingStartRank() {
		ctx.Sides[move.Piece.Color].ClearCastlingRight(move.From.File())
		return
	}
	// if there is a rook capture forfeit the other side castle rights
	if move.HasCapture && move.Captured.Type == ROOK &&
		move.To.Rank() == move.Captured.Color.KingStartRank() {
		ctx.Sides[move.Captured.Color].ClearCastlingRight(move.To.File())
	}
}

// SetEnPassantTarget sets the en passant target square if this move is a
// double pawn push, otherwise clears it. Called after every move.
func (ctx *MoveContext) SetEnPassantTarget(move Move) {
	if move.IsDoublePawnPush() {
		ctx.EnPassantTarget = move.EnPassantTarget()
	} else {
		ctx.EnPassantTarget = NoPosition
	}
}

type ClockContext struct {
	HalfMoveClock  uint16
	FullMoveNumber uint16
}

type TurnContext struct {
	MoveContext
	ClockContext
}

// Reset zeroes the context and attaches a fresh Board. Call this before
// filling in a TurnContext — the zero value has a nil Board pointer.
func (ctx *TurnContext) Reset() {
	*ctx = TurnContext{}
	ctx.Board = new(Board)
}

// Copy returns a deep copy of the TurnContext with its own independent Board.
// Mutating the copy's board does not affect the original.
func (ctx *TurnContext) Copy() TurnContext {
	board := *ctx.Board
	return TurnContext{
		MoveContext: MoveContext{
			BoardContext:    BoardContext{Board: &board},
			SideToMove:      ctx.SideToMove,
			Sides:           ctx.Sides,
			EnPassantTarget: ctx.EnPassantTarget,
		},
		ClockContext: ctx.ClockContext,
	}
}

// Return a snapshot of current turn context values
func (ctx *TurnContext) Snapshot(move Move) Snapshot {
	return Snapshot{
		Move:                    move,
		PreviousSides:           ctx.Sides,
		PreviousEnPassantTarget: ctx.EnPassantTarget,
	}
}

func (ctx *TurnContext) String() string {
	var sb strings.Builder

	sb.WriteString(ctx.Board.String())
	sb.WriteByte('\n')

	if ctx.SideToMove == WHITE {
		sb.WriteString("Side:       White\n")
	} else {
		sb.WriteString("Side:       Black\n")
	}

	castling := ""
	if ctx.Sides[WHITE].CanCastleKingSide {
		castling += "K"
	}
	if ctx.Sides[WHITE].CanCastleQueenSide {
		castling += "Q"
	}
	if ctx.Sides[BLACK].CanCastleKingSide {
		castling += "k"
	}
	if ctx.Sides[BLACK].CanCastleQueenSide {
		castling += "q"
	}
	if castling == "" {
		castling = "-"
	}
	fmt.Fprintf(&sb, "Castling:   %s\n", castling)
	fmt.Fprintf(&sb, "En passant: %s\n", ctx.EnPassantTarget)
	fmt.Fprintf(&sb, "Clocks:     %d / %d\n", ctx.HalfMoveClock, ctx.FullMoveNumber)

	return sb.String()
}

type ChessState struct {
	Board           Board
	SideToMove      PieceColor
	Sides           [2]SideState
	EnPassantTarget Position
	HalfMoveClock   uint16
	FullMoveNumber  uint16
	Hash            uint64
}
