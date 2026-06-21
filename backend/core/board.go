package core

import (
	"fmt"
	"slices"
	"strings"
)

// Square is a packed board cell.
//
//	0      empty
//	1-6    white PAWN..KING (PieceType + 1)
//	7-12   black PAWN..KING (PieceType + 7)
//
// The zero value, EmptySquare, means empty — so a freshly zeroed Board
// starts all-empty with no initialization needed.
type Square uint8

// EmptySquare is the zero value of Square — an unoccupied cell.
const EmptySquare Square = 0

// NewSquare packs a Piece into its Square encoding.
func NewSquare(p Piece) Square {
	return Square(uint8(p.Color)*6 + uint8(p.Type) + 1)
}

// IsEmpty reports whether the square holds no piece.
func (s Square) IsEmpty() bool {
	return s == EmptySquare
}

// IsOccupied reports whether the square holds a piece.
func (s Square) IsOccupied() bool {
	return s != EmptySquare
}

// Type returns the piece type on this square. Only meaningful when
// the square is occupied.
func (s Square) Type() PieceType {
	return PieceType((uint8(s) - 1) % 6)
}

// Color returns the piece color on this square. Only meaningful when
// the square is occupied.
func (s Square) Color() PieceColor {
	return PieceColor((uint8(s) - 1) / 6)
}

// Piece unpacks the square back into a Piece value. Only meaningful
// when the square is occupied.
func (s Square) Piece() Piece {
	return Piece{Type: s.Type(), Color: s.Color()}
}

// IsOccupiedBy reports whether the square is occupied by a piece
// of the given color.
func (s Square) IsOccupiedBy(color PieceColor) bool {
	return s.IsOccupied() && s.Color() == color
}

// IsOccupiedByAny reports whether the square is occupied by a piece
// of the given color and one of the given types.
func (s Square) IsOccupiedByAny(color PieceColor, types ...PieceType) bool {
	if !s.IsOccupiedBy(color) {
		return false
	}
	return slices.Contains(types, s.Type())
}

// Board is the 8x8 grid of squares. At one byte per square.
type Board [64]Square

// IsOccupied reports whether pos holds a piece. Use this for a single
// one-shot check; if you need more than one fact about the square
// (type, color, the piece itself), index the board directly into a
// local variable instead of repeating Board.<Method>(pos) calls.
func (b *Board) IsOccupied(pos Position) bool {
	return b[pos].IsOccupied()
}

// Place puts p on pos, overwriting whatever was there.
func (b *Board) Place(pos Position, p Piece) {
	b[pos] = NewSquare(p)
}

// Clear empties pos.
func (b *Board) Clear(pos Position) {
	b[pos] = EmptySquare
}

// Move relocates whatever occupies from to to, emptying from. For
// moves where the piece itself changes (promotion), use Clear + Place
// instead — Move assumes the mover is unchanged.
func (b *Board) Move(from, to Position) {
	b[to] = b[from]
	b[from] = EmptySquare
}

func (b *Board) String() string {
	var sb strings.Builder
	for rank := int(RANK_8); rank >= int(RANK_1); rank-- {
		fmt.Fprintf(&sb, "%d  ", rank+1)
		for file := FILE_A; file <= FILE_H; file++ {
			pos := NewPosition(file, Rank(rank))
			sq := b[pos]
			if sq.IsEmpty() {
				sb.WriteByte('.')
			} else {
				sb.WriteByte(sq.Piece().Char())
			}
			if file != FILE_H {
				sb.WriteByte(' ')
			}
		}
		sb.WriteByte('\n')
	}
	sb.WriteString("   a b c d e f g h\n")

	return sb.String()
}
