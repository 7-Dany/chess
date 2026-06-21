package core

import "fmt"

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

// ParsePiece converts a standard chess piece letter (as used in FEN, PGN,
// and UCI) into a Piece. Uppercase letters are White, lowercase are Black.
func ParsePiece(ch byte) (Piece, error) {
	color := WHITE
	if ch >= 'a' {
		color = BLACK
		ch -= 'a' - 'A' // fold to uppercase
	}

	var pieceType PieceType
	switch ch {
	case 'P':
		pieceType = PAWN
	case 'N':
		pieceType = KNIGHT
	case 'B':
		pieceType = BISHOP
	case 'R':
		pieceType = ROOK
	case 'Q':
		pieceType = QUEEN
	case 'K':
		pieceType = KING
	default:
		return Piece{}, fmt.Errorf("core: unknown piece letter %q", ch)
	}

	return Piece{Type: pieceType, Color: color}, nil
}

// Char returns the standard chess piece letter for piece
// uppercase for White, lowercase for Black. The inverse of ParsePiece.
func (p Piece) Char() byte {
	letters := [6]byte{'P', 'N', 'B', 'R', 'Q', 'K'}
	ch := letters[p.Type]
	if p.Color == BLACK {
		ch += 'a' - 'A'
	}
	return ch
}
