package core

type Snapshot struct {
	Move                    Move
	PreviousSides           [2]SideState
	PreviousEnPassantTarget Position
	PreviousHalfMoveClock   uint16
	PreviousFullMoveNumber  uint16
}
