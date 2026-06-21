// Package fen parses and serializes Forsyth-Edwards Notation (FEN) strings.
//
// FEN describes a chess position in a single line of ASCII text. It has six
// space-separated fields:
//
//	<piece-placement> <side-to-move> <castling-rights> <en-passant-target> <halfmove-clock> <fullmove-number>
//
// Example: rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1
package fen

import "github.com/7-Dany/chess/core"

// FENParser converts between FEN strings and TurnContext values.
// Decode and Encode are inverses: Encode(Decode(s)) should round-trip
// to the same FEN string.
type FENParser interface {
	// Decode parses a FEN string into a TurnContext ready for the engine.
	// Returns an error if any of the six fields is missing or malformed.
	Decode(str string, ctx *core.TurnContext) error

	// Encode serializes a TurnContext back into a standard FEN string.
	Encode(ctx *core.TurnContext) string
}

type FEN struct{}

var defaultFEN = FEN{}

func GetDefaultFenParser() FEN {
	return defaultFEN
}
