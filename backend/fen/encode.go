package fen

import (
	"strconv"

	"github.com/7-Dany/chess/core"
)

// Encode serializes ctx into a standard FEN string.
func (FEN) Encode(ctx *core.TurnContext) string {
	var buf [92]byte
	b := buf[:0]

	b = encodePiecePlacement(b, ctx.Board)
	b = append(b, ' ')
	b = encodeSideToMove(b, ctx.SideToMove)
	b = append(b, ' ')
	b = encodeCastlingRights(b, ctx.Sides)
	b = append(b, ' ')
	b = encodeEnPassantTarget(b, ctx.EnPassantTarget)
	b = append(b, ' ')
	b = encodeHalfMoveClock(b, ctx.HalfMoveClock)
	b = append(b, ' ')
	b = encodeFullMoveNumber(b, ctx.FullMoveNumber)

	return string(b)
}

// encodePiecePlacement emits the first FEN field: one rank per row from rank 8
// down to rank 1, pieces as their letter, consecutive empty squares collapsed
// into a digit, ranks separated by '/'.
func encodePiecePlacement(b []byte, board *core.Board) []byte {
	for r := 7; r >= 0; r-- {
		empty := 0

		for f := range 8 {
			position := core.NewPosition(core.File(f), core.Rank(r))
			square := board[position]

			if square.IsOccupied() {
				if empty > 0 {
					b = append(b, byte('0'+empty))
				}

				empty = 0
				b = append(b, square.Piece().Char())
				continue
			}

			empty++
		}

		if empty > 0 {
			b = append(b, byte('0'+empty))
		}

		if r > 0 {
			b = append(b, '/')
		}
	}

	return b
}

// encodeSideToMove emits 'w' or 'b'.
func encodeSideToMove(b []byte, sideToMove core.PieceColor) []byte {
	switch sideToMove {
	case core.WHITE:
		return append(b, 'w')
	case core.BLACK:
		return append(b, 'b')
	}
	return b
}

// encodeCastlingRights emits up to four letters (KQkq) or '-' if none.
func encodeCastlingRights(b []byte, sides [2]core.SideState) []byte {
	original := len(b)

	whiteState := sides[core.WHITE]
	if whiteState.CanCastleKingSide {
		b = append(b, 'K')
	}
	if whiteState.CanCastleQueenSide {
		b = append(b, 'Q')
	}

	blackState := sides[core.BLACK]
	if blackState.CanCastleKingSide {
		b = append(b, 'k')
	}
	if blackState.CanCastleQueenSide {
		b = append(b, 'q')
	}

	if original == len(b) {
		return append(b, '-')
	}

	return b
}

// encodeEnPassantTarget emits the target square in algebraic notation (e.g. "e3") or '-'.
func encodeEnPassantTarget(b []byte, enPassantTarget core.Position) []byte {
	if enPassantTarget == core.NoPosition {
		return append(b, '-')
	}

	file := byte(enPassantTarget.File() + 'a')
	rank := byte(enPassantTarget.Rank() + '1')

	return append(b, file, rank)
}

// encodeHalfMoveClock emits the halfmove clock as a decimal number.
func encodeHalfMoveClock(b []byte, halfMoveClock uint16) []byte {
	return strconv.AppendUint(b, uint64(halfMoveClock), 10)
}

// encodeFullMoveNumber emits the fullmove number as a decimal number.
func encodeFullMoveNumber(b []byte, fullMoveNumber uint16) []byte {
	return strconv.AppendUint(b, uint64(fullMoveNumber), 10)
}
