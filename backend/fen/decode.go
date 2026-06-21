package fen

import (
	"fmt"
	"strconv"

	"github.com/7-Dany/chess/core"
)

// Decode parses a FEN string into ctx.
func (FEN) Decode(str string, ctx *core.TurnContext) error {
	ctx.Reset()

	index, err := decodePiecePlacement(str, ctx.Board)
	if err != nil {
		return err
	}

	index, err = decodeSideToMove(str, index, ctx)
	if err != nil {
		return err
	}

	index, err = decodeCastlingRights(str, index, ctx)
	if err != nil {
		return err
	}

	index, err = decodeEnPassantTarget(str, index, ctx)
	if err != nil {
		return err
	}

	index, err = decodeHalfMoveClock(str, index, ctx)
	if err != nil {
		return err
	}

	err = decodeFullMoveNumber(str, index, ctx)
	if err != nil {
		return err
	}

	return nil
}

// decodePiecePlacement parses the first FEN field (piece placement) into board.
// Returns the index of the space after the field (start of the next field).
func decodePiecePlacement(str string, board *core.Board) (int, error) {
	rank, file := uint8(0), uint8(0)

	for i, letter := range str {
		switch {
		case letter == ' ':
			if rank+1 != 8 || file != 8 {
				return 0, fmt.Errorf("fen: expected 8 ranks, got %d", core.Rank(rank+1))
			}

			return i + 1, nil

		case letter == '/':
			if file != 8 {
				return 0, fmt.Errorf("fen: rank %d has %d files, want 8", core.Rank(rank).Reverse(), file)
			}

			file = 0
			rank++
			if rank >= 8 {
				return 0, fmt.Errorf("fen: too many ranks")
			}

		// if letter is number, they are empty squares.
		case letter >= '1' && letter <= '8':
			empty := uint8(letter - '0')
			if (file + empty) > 8 {
				return 0, fmt.Errorf("fen: rank %d overflows 8 files", core.Rank(rank).Reverse())
			}

			file += empty

		default:
			if file >= 8 {
				return 0, fmt.Errorf("fen: rank %d overflows 8 files", core.Rank(rank).Reverse())
			}

			// current letter is piece
			piece, err := core.ParsePiece(byte(letter))
			if err != nil {
				return 0, fmt.Errorf("fen: rank %d: %w", core.Rank(rank).Reverse(), err)
			}

			position := core.NewPosition(core.File(file), core.Rank(rank).Reverse())
			board.Place(position, piece)

			file++
		}
	}

	return 0, fmt.Errorf("fen: piece placement field is incomplete, no space terminator")
}

// decodeSideToMove parses the 'w' or 'b' field.
func decodeSideToMove(str string, index int, ctx *core.TurnContext) (int, error) {
	if index >= len(str) {
		return 0, fmt.Errorf("fen: missing side-to-move field")
	}

	sideToMove := str[index]
	switch sideToMove {
	case 'w':
		ctx.SideToMove = core.WHITE
	case 'b':
		ctx.SideToMove = core.BLACK
	default:
		return 0, fmt.Errorf("fen: invalid sideToMove letter %q, expected w or b", sideToMove)
	}

	return index + 2, nil
}

// decodeCastlingRights parses the KQkq field (or '-' for none).
func decodeCastlingRights(str string, index int, ctx *core.TurnContext) (int, error) {
	if index >= len(str) {
		return 0, fmt.Errorf("fen: missing castling-rights field")
	}

	white := core.WHITE
	black := core.BLACK

	for i := index; i < len(str); i++ {
		letter := str[i]
		switch letter {
		case ' ':
			return i + 1, nil
		case 'k':
			ctx.Sides[black].CanCastleKingSide = true
		case 'q':
			ctx.Sides[black].CanCastleQueenSide = true
		case 'K':
			ctx.Sides[white].CanCastleKingSide = true
		case 'Q':
			ctx.Sides[white].CanCastleQueenSide = true
		case '-':
		// no rights, nothing to set, turn initialized with no castle rights

		default:
			return 0, fmt.Errorf("fen: invalid castle rights letter %v, expected any of [K, Q, k, q, -]", letter)
		}
	}

	return 0, fmt.Errorf("fen: castling-rights field is incomplete, no space terminator")
}

// decodeEnPassantTarget parses the en passant square (e.g. "e3") or '-'.
func decodeEnPassantTarget(str string, index int, ctx *core.TurnContext) (int, error) {
	if index >= len(str) {
		return 0, fmt.Errorf("fen: missing en-passant-target field")
	}

	if str[index] == '-' {
		ctx.EnPassantTarget = core.NoPosition
		return index + 2, nil // skip '-' + space
	}

	if index+1 >= len(str) {
		return 0, fmt.Errorf("fen: en-passant target is too short, expected a file letter and rank digit")
	}

	file, err := core.ParseFile(str[index])
	if err != nil {
		return 0, fmt.Errorf("fen: en-passant target file: %w", err)
	}
	rank, err := core.ParseRank(str[index+1])
	if err != nil {
		return 0, fmt.Errorf("fen: en-passant target rank: %w", err)
	}

	// En passant targets can only be on rank 3 (white just moved) or rank 6 (black just moved).
	if rank != core.RANK_3 && rank != core.RANK_6 {
		return 0, fmt.Errorf("fen: en-passant target rank must be 3 or 6, got %d", rank+1)
	}

	ctx.EnPassantTarget = core.NewPosition(file, rank)
	return index + 3, nil // skip file + rank + space
}

// decodeHalfMoveClock parses the halfmove clock (a decimal number).
func decodeHalfMoveClock(str string, index int, ctx *core.TurnContext) (int, error) {
	start := index
	for index < len(str) && str[index] != ' ' {
		index++
	}
	if index == start {
		return 0, fmt.Errorf("fen: missing halfmove-clock field")
	}

	clock, err := strconv.ParseUint(str[start:index], 10, 16)
	if err != nil {
		return 0, fmt.Errorf("fen: invalid halfmove-clock %q: %w", str[start:index], err)
	}

	ctx.HalfMoveClock = uint16(clock)

	// Skip the space (if present — this might be the last field).
	if index < len(str) {
		index++ // skip ' '
	}

	return index, nil
}

// decodeFullMoveNumber parses the fullmove number (a decimal number).
func decodeFullMoveNumber(str string, index int, ctx *core.TurnContext) error {
	start := index
	for index < len(str) {
		index++
	}
	if index == start {
		return fmt.Errorf("fen: missing fullmove-number field")
	}

	num, err := strconv.ParseUint(str[start:index], 10, 16)
	if err != nil {
		return fmt.Errorf("fen: invalid fullmove-number %q: %w", str[start:index], err)
	}

	ctx.FullMoveNumber = uint16(num)

	return nil
}
