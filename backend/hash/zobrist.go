package hash

import (
	"math/rand/v2"

	"github.com/7-Dany/chess/core"
)

type Zobrist struct {
	piecePosition [12][64]uint64 // [12]piece * [64]position
	castling      [4]uint64      // [0, 1] -> White Castling Rights, [2, 3] -> Black Castling Rights
	enPassant     [8]uint64      // 8 different en passant targets positions
	sideToMove    uint64         // current side xor in for black, xor out for white
}

func newZobrist() Zobrist {
	var z Zobrist
	rng := rand.New(rand.NewPCG(0x9e3779b97f4a7c15, 0))

	for piece := range 12 {
		for pos := range 64 {
			z.piecePosition[piece][pos] = rng.Uint64()
		}
	}

	for i := range 4 {
		z.castling[i] = rng.Uint64()
	}

	for file := range 8 {
		z.enPassant[file] = rng.Uint64()
	}

	z.sideToMove = rng.Uint64()

	return z
}

var defaultZobrist = newZobrist()

func GetDefaultHasher() Zobrist {
	return defaultZobrist
}

func (z Zobrist) InitHash(ctx *core.TurnContext) uint64 {
	var hash uint64

	// hash board
	for i, square := range ctx.Board {
		if square.IsOccupied() {
			hash ^= z.piecePosition[square-1][i]
		}
	}

	// hash castling rights
	whiteSide := ctx.Sides[core.WHITE]
	if whiteSide.CanCastleKingSide {
		hash ^= z.castling[0]
	}
	if whiteSide.CanCastleQueenSide {
		hash ^= z.castling[1]
	}

	blackSide := ctx.Sides[core.BLACK]
	if blackSide.CanCastleKingSide {
		hash ^= z.castling[2]
	}
	if blackSide.CanCastleQueenSide {
		hash ^= z.castling[3]
	}

	// if en passant target available
	if ctx.EnPassantTarget != core.NoPosition {
		file := ctx.EnPassantTarget.File()
		hash ^= z.enPassant[file]
	}

	// if black turn
	if ctx.SideToMove == core.BLACK {
		hash ^= z.sideToMove
	}

	return hash
}

func (z Zobrist) Hash(current uint64, move core.MoveHash) uint64 {
	hash := current

	switch move.Type {
	case core.NORMAL:
		hash = z.hashNormal(hash, move)
	case core.PROMOTION:
		hash = z.hashPromotion(hash, move)
	case core.CASTLING:
		hash = z.hashCastling(hash, move)
	case core.EN_PASSANT:
		hash = z.hashEnPassantCapture(hash, move)
	}

	hash = z.hashEnCastlingRights(hash, move)

	hash = z.hashEnPassantTarget(hash, move)

	hash ^= z.sideToMove

	return hash
}

func (z *Zobrist) hashNormal(hash uint64, move core.MoveHash) uint64 {
	pieceSquare := core.NewSquare(move.Piece) - 1

	// hash the movement
	hash ^= z.piecePosition[pieceSquare][move.From]
	hash ^= z.piecePosition[pieceSquare][move.To]

	// if any piece captured hash it out
	if move.HasCapture {
		capturedSquare := core.NewSquare(move.Captured) - 1
		hash ^= z.piecePosition[capturedSquare][move.To]
	}

	return hash
}

func (z Zobrist) hashPromotion(hash uint64, move core.MoveHash) uint64 {
	pawnSquare := core.NewSquare(move.Piece) - 1
	promotedSquare := core.NewSquare(core.Piece{Type: move.PromoteTo, Color: move.Piece.Color}) - 1

	// hash the movement
	hash ^= z.piecePosition[pawnSquare][move.From]
	hash ^= z.piecePosition[promotedSquare][move.To]

	// if any piece captured, hash it out
	if move.HasCapture {
		capturedSquare := core.NewSquare(move.Captured) - 1
		hash ^= z.piecePosition[capturedSquare][move.To]
	}

	return hash
}

func (z Zobrist) hashCastling(hash uint64, move core.MoveHash) uint64 {
	kingSquare := core.NewSquare(move.Piece) - 1
	rookSquare := core.NewSquare(core.Piece{Type: core.ROOK, Color: move.Piece.Color}) - 1
	rookFrom, rookTo := move.CastlingRookPositions()

	// move king
	hash ^= z.piecePosition[kingSquare][move.From]
	hash ^= z.piecePosition[kingSquare][move.To]

	// move rook
	hash ^= z.piecePosition[rookSquare][rookFrom]
	hash ^= z.piecePosition[rookSquare][rookTo]

	return hash
}

func (z Zobrist) hashEnPassantCapture(hash uint64, move core.MoveHash) uint64 {
	pieceSquare := core.NewSquare(move.Piece) - 1

	// move piece
	hash ^= z.piecePosition[pieceSquare][move.From]
	hash ^= z.piecePosition[pieceSquare][move.To]

	// capture the pawn before that position
	capturedPawnPosition := move.EnPassantCapturedPosition()
	capturedPawnSquare := core.NewSquare(move.Captured) - 1
	hash ^= z.piecePosition[capturedPawnSquare][capturedPawnPosition]

	return hash
}

func (z Zobrist) hashEnCastlingRights(hash uint64, move core.MoveHash) uint64 {
	// new castling rights for white side
	nWCR := move.NewSides[core.WHITE]
	// previous castling rights for white side
	pWCR := move.PreviousSides[core.WHITE]

	if pWCR.CanCastleKingSide != nWCR.CanCastleKingSide {
		hash ^= z.castling[0]
	}
	if pWCR.CanCastleQueenSide != nWCR.CanCastleQueenSide {
		hash ^= z.castling[1]
	}

	// new castling rights for black side
	nBCR := move.NewSides[core.BLACK]
	// previous castling rights for white side
	pBCR := move.PreviousSides[core.BLACK]

	if pBCR.CanCastleKingSide != nBCR.CanCastleKingSide {
		hash ^= z.castling[2]
	}
	if pBCR.CanCastleQueenSide != nBCR.CanCastleQueenSide {
		hash ^= z.castling[3]
	}

	return hash
}

func (z Zobrist) hashEnPassantTarget(hash uint64, move core.MoveHash) uint64 {
	if move.PreviousEnPassantTarget != core.NoPosition {
		file := move.PreviousEnPassantTarget.File()
		hash ^= z.enPassant[file]
	}

	if move.IsDoublePawnPush() {
		file := move.EnPassantTarget().File()
		hash ^= z.enPassant[file]
	}

	return hash
}
