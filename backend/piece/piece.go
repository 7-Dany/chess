package piece

import "github.com/7-Dany/chess/core"

// Piece represents the behavior of a single chess piece type.
//
// Implementations are stateless value types (e.g. Pawn{}, Bishop{}).
// All position-specific and board-specific information is passed in
// through the method arguments, so the same instance can be reused
// across the entire program.
type Piece interface {
	// IsAttacking reports whether any piece of this type and color
	// attacks target.
	//
	// Scans from target outward using this piece's attack geometry —
	// the right tool for "is my king in check?" or "is this square
	// defended?" when the caller knows target but not candidate sources.
	//
	// The receiver is a stateless piece-type: Knight.IsAttacking(WHITE,
	// E4, ctx) = "is E4 attacked by a white knight?".
	//
	// Special cases:
	//   - Pawns scan the two squares "behind" target relative to
	//     attacker color (white below, black above).
	//   - Sliders check the first blocker on each ray; queen is covered
	//     by the bishop and rook scans, not dispatched separately.
	IsAttacking(color core.PieceColor, target core.Position, ctx core.BoardContext) bool

	// Attacks returns every square this piece threatens from the given
	// position, given the current board layout.
	//
	// Includes squares occupied by friendly pieces — a piece "attacks" a
	// square even if it can't move there. Essential for check detection.
	//
	// Special cases:
	//   - Pawns attack diagonally, not forward.
	//   - Castling is NOT an attack.
	//   - Sliders stop at the first occupied square but include it.
	//
	// Color dependency: geometry is color-independent for all pieces
	// except pawn, which reads color from ctx.Board[from] — so `from`
	// must be occupied when calling Attacks on a pawn.
	Attacks(from core.Position, ctx core.BoardContext) []core.Position

	// PseudoLegalMoves returns all move options this piece has from the
	// given position, respecting movement rules, board edges, blockers,
	// and capture eligibility — but NOT filtering for king safety.
	//
	// Moves that leave the moving side's king in check are still
	// returned; the Engine layer applies that filter.
	//
	// Special cases:
	//   - Pawn moves include single push, double push, diagonal
	//     captures, en passant, and promotions.
	//   - Castling is NOT included — the Engine adds it.
	PseudoLegalMoves(from core.Position, ctx core.MoveContext) []core.Move
}

// PieceProvider maps each PieceType to its stateless Piece implementation.
// It is constructed once and shared across all queries.
type PieceProvider struct {
	pieces [6]Piece
}

// NewPieceProvider creates a PieceProvider with the six standard piece
// implementations. Each piece is a zero-size value type, so the map
// holds only six pointers with no per-query allocation.
func newPieceProvider() *PieceProvider {
	return &PieceProvider{
		pieces: [6]Piece{Pawn{}, Knight{}, Bishop{}, Rook{}, Queen{}, King{}},
	}
}

// GetPiece returns the Piece implementation for the given PieceType.
// Callers should not need to construct Piece values directly.
func (p *PieceProvider) GetPiece(pt core.PieceType) Piece {
	return p.pieces[pt]
}

var provider = newPieceProvider()

func DefaultProvider() *PieceProvider {
	return provider
}
