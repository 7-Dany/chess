package engine

import (
	"slices"

	"github.com/7-Dany/chess/core"
)

func (e *DefaultEngine) GetPseudoLegalMoves(position core.Position, ctx core.TurnContext) []core.Move {
	square := ctx.Board[position]
	if !square.IsOccupiedBy(ctx.SideToMove) {
		return []core.Move{}
	}

	piece := e.pieces.GetPiece(square.Piece.Type)
	moves := piece.PseudoLegalMoves(position, ctx.MoveContext)

	// add castling moves
	if ctx.Sides[ctx.SideToMove].KingPosition == position {
		moves = append(moves, e.castlingMoves(position, ctx)...)
	}

	return moves
}

// return king castling moves, if rights are valid
func (e *DefaultEngine) castlingMoves(kingPosition core.Position, ctx core.TurnContext) []core.Move {
	current := ctx.SideToMove
	enemy := current.Opponent()

	// King must be on its home file. If CastlingRights are set but the king
	// is elsewhere, state is corrupt — bail.
	if kingPosition.File() != core.FILE_E {
		return nil
	}

	// if the king is in check, return
	if e.IsSquareAttacked(kingPosition, enemy, ctx) {
		return nil
	}

	moves := make([]core.Move, 0, 2)
	king := core.Piece{Type: core.KING, Color: current}
	rank := kingPosition.Rank()

	if e.canCastleKingSide(rank, ctx) {
		moves = append(moves, core.Move{
			Type:  core.CASTLING,
			Piece: king,
			From:  kingPosition,
			To:    core.NewPosition(core.FILE_G, rank),
		})
	}

	if e.canCastleQueenSide(rank, ctx) {
		moves = append(moves, core.Move{
			Type:  core.CASTLING,
			Piece: king,
			From:  kingPosition,
			To:    core.NewPosition(core.FILE_C, rank),
		})
	}

	return slices.Clip(moves)
}

// canCastleKingSide return true, if rights allow the king to castle from king side
func (e *DefaultEngine) canCastleKingSide(rank core.Rank, ctx core.TurnContext) bool {
	if !ctx.Sides[ctx.SideToMove].CastlingRights.KingSide {
		return false
	}

	enemy := ctx.SideToMove.Opponent()

	path := core.NewPosition(core.FILE_F, rank)
	dest := core.NewPosition(core.FILE_G, rank)

	if ctx.Board[path].Occupied || ctx.Board[dest].Occupied {
		return false
	}

	if e.IsSquareAttacked(path, enemy, ctx) || e.IsSquareAttacked(dest, enemy, ctx) {
		return false
	}

	return true
}

// canCastleQueenSide return true, if rights allow the king to castle from queen side
func (e *DefaultEngine) canCastleQueenSide(rank core.Rank, ctx core.TurnContext) bool {
	if !ctx.Sides[ctx.SideToMove].CastlingRights.QueenSide {
		return false
	}

	enemy := ctx.SideToMove.Opponent()

	path := core.NewPosition(core.FILE_D, rank)
	dest := core.NewPosition(core.FILE_C, rank)
	between := core.NewPosition(core.FILE_B, rank)

	if ctx.Board[path].Occupied || ctx.Board[dest].Occupied || ctx.Board[between].Occupied {
		return false
	}

	if e.IsSquareAttacked(path, enemy, ctx) || e.IsSquareAttacked(dest, enemy, ctx) {
		return false
	}

	return true
}
