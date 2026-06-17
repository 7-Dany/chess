package core

type GameStatus uint8

const (
	InProgress GameStatus = iota
	CheckMate
	Draw
)

type DrawReason uint8

const (
	NoDrawReason DrawReason = iota
	Stalemate
	ThreefoldRepetition
	FiftyMoveRule
	InsufficientMaterial
)

type GameResult struct {
	Status     GameStatus
	Winner     PieceColor
	HasWinner  bool
	DrawReason DrawReason
}
