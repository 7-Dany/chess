package core

type MoveType uint8

const (
	// NORMAL is a regular move or capture made by one piece.
	NORMAL MoveType = iota

	// CASTLING moves the king two squares toward a rook, then moves that rook
	// to the square the king crossed.
	CASTLING

	// EN_PASSANT is a pawn capture of a pawn that just advanced two squares,
	// where the capturing pawn moves diagonally to the skipped square.
	EN_PASSANT

	// PROMOTION replaces a pawn with another piece after it reaches the last rank.
	PROMOTION
)

// MAX_MOVES is the buffer size that guarantees PseudoLegalMoves and Attacks
// never need to grow their backing array, and therefore never allocate.
//
// The largest single-piece move set is a queen on a wide-open board (27
// moves); a king with two castling moves adds at most 2 more. 32 provides
// headroom over both. Pass a buffer of at least this size — e.g.
// `var buf [core.MAX_MOVES]core.Move` on the stack — to keep move generation
// allocation-free.
const MAX_MOVES = 32

type Move struct {
	Piece      Piece
	From       Position
	To         Position
	Type       MoveType
	PromoteTo  PieceType
	Captured   Piece
	HasCapture bool
}

// IsDoublePawnPush reports whether this move is a pawn advancing two squares
// from its starting rank. Used to set the en passant target.
func (m Move) IsDoublePawnPush() bool {
	if m.Piece.Type != PAWN {
		return false
	}
	rankDiff := int(m.To.Rank()) - int(m.From.Rank())
	if m.Piece.Color == WHITE {
		return rankDiff == 2
	}
	return rankDiff == -2
}

// EnPassantTarget returns the square that sits behind a just-pushed pawn —
// the square an enemy pawn could capture into. Only meaningful when
// IsDoublePawnPush() is true.
func (m Move) EnPassantTarget() Position {
	step := int8(1)
	if m.Piece.Color == BLACK {
		step = -1
	}
	rank, _ := m.From.Rank().Add(step)
	return NewPosition(m.To.File(), rank)
}
