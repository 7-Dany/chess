package core

type PieceType uint8

const (
	PAWN PieceType = iota
	KNIGHT
	BISHOP
	ROOK
	QUEEN
	KING
)

type PieceColor uint8

const (
	WHITE PieceColor = iota
	BLACK
)

// Opponent returns the opposite color.
func (c PieceColor) Opponent() PieceColor {
	return 1 - c
}

type Piece struct {
	Type  PieceType
	Color PieceColor
}
