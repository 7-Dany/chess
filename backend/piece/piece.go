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
	Attacks(attacks []core.Position, from core.Position, ctx core.BoardContext) []core.Position

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
	PseudoLegalMoves(moves []core.Move, from core.Position, ctx core.MoveContext) []core.Move
}

const MAX_MOVES uint8 = 32

type Pieces struct {
	pawn   Pawn
	knight Knight
	bishop Bishop
	rook   Rook
	queen  Queen
	king   King
}

var defaultPieces = Pieces{}

func GetDefaultPieces() Pieces { return defaultPieces }

// Pawn returns the concrete Pawn instance. The return type is Pawn (not Piece),
// so callers get static method dispatch.
func (p Pieces) Pawn() Pawn { return p.pawn }

// Knight returns the concrete Knight instance.
func (p Pieces) Knight() Knight { return p.knight }

// Bishop returns the concrete Bishop instance.
func (p Pieces) Bishop() Bishop { return p.bishop }

// Rook returns the concrete Rook instance.
func (p Pieces) Rook() Rook { return p.rook }

// Queen returns the concrete Queen instance.
func (p Pieces) Queen() Queen { return p.queen }

// King returns the concrete King instance.
func (p Pieces) King() King { return p.king }
