package core

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
// back rank, each forfeits one right. King moves and castling forfeit all
// rights and are handled in applyNormal/applyCastling.
//
// The back-rank guard makes the chess invariant explicit: a rook on its
// color's king-start rank at A/H is either the original or the right is
// already cleared. ClearCastlingRight is a no-op for non-A/H files.
func (ctx *MoveContext) ForfeitCastlingRight(move Move) {
	if move.Piece.Type == ROOK && move.From.Rank() == move.Piece.Color.KingStartRank() {
		ctx.Sides[move.Piece.Color].ClearCastlingRight(move.From.File())
		return
	}
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

type ChessState struct {
	Board           Board
	SideToMove      PieceColor
	Sides           [2]SideState
	EnPassantTarget Position
	HalfMoveClock   uint16
	FullMoveNumber  uint16
	Hash            uint64
}
