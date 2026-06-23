package core

type MoveHash struct {
	Move
	PreviousSides           [2]SideState
	PreviousEnPassantTarget Position
	NewSides                [2]SideState
}

func NewMoveHash(snap Snapshot, ctx TurnContext) MoveHash {
	return MoveHash{
		Move:                    snap.Move,
		PreviousSides:           snap.PreviousSides,
		PreviousEnPassantTarget: snap.PreviousEnPassantTarget,
		NewSides:                ctx.Sides,
	}
}
