package core

// Board columns as files from A -> H
type File uint8

const (
	FILE_A File = iota
	FILE_B
	FILE_C
	FILE_D
	FILE_E
	FILE_F
	FILE_G
	FILE_H
)

func (f File) Add(value int8) (File, bool) {
	result := int8(f) + value
	if result < 0 || result > 7 {
		return 0, false
	}
	return File(result), true
}

func (f File) String() string {
	return string(rune(f + 'A'))
}

// Board raws as rank from 1 -> 8
type Rank uint8

const (
	RANK_1 Rank = iota
	RANK_2
	RANK_3
	RANK_4
	RANK_5
	RANK_6
	RANK_7
	RANK_8
)

func (r Rank) Add(value int8) (Rank, bool) {
	result := int8(r) + value
	if result < 0 || result > 7 {
		return 0, false
	}
	return Rank(result), true
}

func (r Rank) String() string {
	return string(rune('1' + r))
}

type Position uint8

// Positions from A1 -> H8
const (
	A1 Position = iota
	A2
	A3
	A4
	A5
	A6
	A7
	A8
	B1
	B2
	B3
	B4
	B5
	B6
	B7
	B8
	C1
	C2
	C3
	C4
	C5
	C6
	C7
	C8
	D1
	D2
	D3
	D4
	D5
	D6
	D7
	D8
	E1
	E2
	E3
	E4
	E5
	E6
	E7
	E8
	F1
	F2
	F3
	F4
	F5
	F6
	F7
	F8
	G1
	G2
	G3
	G4
	G5
	G6
	G7
	G8
	H1
	H2
	H3
	H4
	H5
	H6
	H7
	H8
	NoPosition
)

func NewPosition(file File, rank Rank) Position {
	return Position(uint8(file)*8 + uint8(rank))
}

// Board columns as files from A -> H
func (p Position) File() File {
	return File(p / 8)
}

// Board raws as rank from 1 -> 8
func (p Position) Rank() Rank {
	return Rank(p % 8)
}

func (p Position) Index() uint8 {
	return uint8(p)
}

func (p Position) IsValid() bool {
	return p < NoPosition
}

func (p Position) String() string {
	if !p.IsValid() {
		return "-"
	}
	return p.File().String() + p.Rank().String()
}
