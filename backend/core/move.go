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

type Move struct {
	Piece      Piece
	From       Position
	To         Position
	Type       MoveType
	PromoteTo  PieceType
	Captured   Piece
	HasCapture bool
}
