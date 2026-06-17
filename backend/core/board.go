package core

import "slices"

type Square struct {
	Piece    Piece
	Occupied bool
}

type Board [64]Square

// IsOccupiedBy reports whether the square is occupied by a piece
// of the given color.
func (s Square) IsOccupiedBy(color PieceColor) bool {
	return s.Occupied && s.Piece.Color == color
}

// IsOccupiedByAny reports whether the square is occupied by a piece
// of the given color and one of the given types.
func (s Square) IsOccupiedByAny(color PieceColor, types ...PieceType) bool {
	if !s.Occupied || s.Piece.Color != color {
		return false
	}
	return slices.Contains(types, s.Piece.Type)
}
