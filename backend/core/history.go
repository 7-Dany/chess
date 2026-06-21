package core

type Snapshot struct {
	Move                    Move
	PreviousSides           [2]SideState
	PreviousEnPassantTarget Position
}
