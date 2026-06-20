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

// KingStartRank returns the rank where this color's king starts.
// White: RANK_1, Black: RANK_8. Same rank where rooks start, so it's
// the rank whose A/H squares guard castling rights.
func (c PieceColor) KingStartRank() Rank {
	if c == WHITE {
		return RANK_1
	}
	return RANK_8
}

// Opponent returns the opposite color.
func (c PieceColor) Opponent() PieceColor {
	return 1 - c
}

type Piece struct {
	Type  PieceType
	Color PieceColor
}
