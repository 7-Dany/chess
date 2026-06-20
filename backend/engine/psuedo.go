package engine

import (
	"github.com/7-Dany/chess/core"
)

func (e *DefaultEngine) GetPseudoLegalMoves(position core.Position, ctx core.TurnContext) []core.Move {
	square := ctx.Board[position]
	if !square.IsOccupiedBy(ctx.SideToMove) {
		return nil
	}

	piece := e.pieces.GetPiece(square.Type())
	moves := piece.PseudoLegalMoves(position, ctx.MoveContext)

	// add castling moves
	if ctx.Sides[ctx.SideToMove].KingPosition == position {
		moves = e.castlingMoves(moves, position, ctx)
	}

	return moves
}

// return king castling moves, if rights are valid
func (e *DefaultEngine) castlingMoves(moves []core.Move, kingPosition core.Position, ctx core.TurnContext) []core.Move {
	current := ctx.SideToMove
	enemy := current.Opponent()

	// King must be on its home file. If CastlingRights are set but the king
	// is elsewhere, state is corrupt — bail.
	if kingPosition.File() != core.FILE_E {
		return moves
	}

	// if the king is in check, return
	if e.IsSquareAttacked(kingPosition, enemy, ctx.BoardContext) {
		return moves
	}

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

	return moves
}

// canCastleKingSide return true, if rights allow the king to castle from king side
func (e *DefaultEngine) canCastleKingSide(rank core.Rank, ctx core.TurnContext) bool {
	if !ctx.Sides[ctx.SideToMove].CanCastleKingSide {
		return false
	}

	enemy := ctx.SideToMove.Opponent()

	path := core.NewPosition(core.FILE_F, rank)
	dest := core.NewPosition(core.FILE_G, rank)

	if ctx.Board[path].IsOccupied() || ctx.Board[dest].IsOccupied() {
		return false
	}

	if e.IsSquareAttacked(path, enemy, ctx.BoardContext) || e.IsSquareAttacked(dest, enemy, ctx.BoardContext) {
		return false
	}

	return true
}

// canCastleQueenSide return true, if rights allow the king to castle from queen side
func (e *DefaultEngine) canCastleQueenSide(rank core.Rank, ctx core.TurnContext) bool {
	if !ctx.Sides[ctx.SideToMove].CanCastleQueenSide {
		return false
	}

	enemy := ctx.SideToMove.Opponent()

	path := core.NewPosition(core.FILE_D, rank)
	dest := core.NewPosition(core.FILE_C, rank)
	between := core.NewPosition(core.FILE_B, rank)

	if ctx.Board[path].IsOccupied() || ctx.Board[dest].IsOccupied() || ctx.Board[between].IsOccupied() {
		return false
	}

	if e.IsSquareAttacked(path, enemy, ctx.BoardContext) || e.IsSquareAttacked(dest, enemy, ctx.BoardContext) {
		return false
	}

	return true
}
