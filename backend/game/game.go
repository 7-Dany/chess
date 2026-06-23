// Package game is the top-level chess orchestrator. It owns the mutable game
// state — board position, Zobrist hash, move history — and coordinates all
// subsystems: move generation and execution (engine), position hashing (hash),
// draw/checkmate detection (rules), move history (history), and repetition
// tracking (tracker).
//
// Chess is the single entry point for callers. All subsystems are injected at
// construction via functional options; every subsystem has a sensible default
// so New() with no arguments starts a standard game from the opening position.
//
// # Typical usage
//
//	game, err := game.New()                       // standard starting position
//	game, err := game.New(game.WithFEN("..."))    // custom position
//	game, err := game.New(game.WithHistory(store)) // persistent history
//
// # Sequence contracts
//
// MakeMove and UndoMove maintain a strict ordering of subsystem calls.
// Callers must not interleave them with direct engine Apply/Undo calls —
// doing so would desync the hash, tracker, and history from the board.
package game

import (
	"errors"

	"github.com/7-Dany/chess/core"
	"github.com/7-Dany/chess/engine"
	"github.com/7-Dany/chess/fen"
	"github.com/7-Dany/chess/hash"
	"github.com/7-Dany/chess/history"
	"github.com/7-Dany/chess/rules"
	"github.com/7-Dany/chess/tracker"
)

var (
	// ErrIllegalMove is returned by MakeMove when the given move is not
	// legal in the current position.
	ErrIllegalMove = errors.New("illegal move")

	// ErrNothingToUndo is returned by UndoMove when the move history is empty.
	ErrNothingToUndo = errors.New("nothing to undo")
)

// STARTING_FEN is the standard chess opening position in Forsyth-Edwards Notation.
const STARTING_FEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

// Chess is the top-level game object. It owns the mutable state (ctx, hash)
// and holds references to all stateless subsystems. Create it with New —
// the zero value is not usable.
type Chess struct {
	ctx     core.TurnContext
	hash    uint64
	fen     string
	engine  engine.Engine
	hasher  hash.Hasher
	tracker tracker.Tracker
	rules   rules.Rules
	parser  fen.Parser
	history history.HistoryStore
}

// Option is a functional option for configuring a Chess instance.
// Pass one or more to New to override specific subsystems.
type Option func(*Chess)

// New creates a Chess instance ready to play. All subsystems default to their
// standard implementations; pass Option values to override any of them.
// Returns an error if the FEN string (default or custom) is invalid.
func New(opts ...Option) (*Chess, error) {
	// set defaults
	e := engine.GetDefaultEngine()
	p := fen.GetDefaultFenParser()
	r := rules.GetDefaultRules()
	h := hash.GetDefaultHasher()

	c := &Chess{
		engine:  &e,
		parser:  &p,
		rules:   &r,
		history: history.NewMemoryStore(),
		tracker: tracker.NewPositionTracker(),
		hasher:  &h,
		fen:     STARTING_FEN,
	}

	// apply overrides
	for _, opt := range opts {
		opt(c)
	}

	// bootstrap
	c.ctx.Reset()
	if err := c.parser.Decode(c.fen, &c.ctx); err != nil {
		return nil, err
	}

	c.hash = c.hasher.InitHash(&c.ctx)
	c.tracker.Record(c.hash)

	return c, nil
}

// WithEngine overrides the default move-generation and execution subsystem.
func WithEngine(e engine.Engine) Option {
	return func(c *Chess) { c.engine = e }
}

// WithParser overrides the default FEN parser used to decode the starting position.
func WithParser(p fen.Parser) Option {
	return func(c *Chess) { c.parser = p }
}

// WithRules overrides the default draw and checkmate detection subsystem.
func WithRules(r rules.Rules) Option {
	return func(c *Chess) { c.rules = r }
}

// WithHistory overrides the default in-memory move history store.
// Use this to plug in a persistent store (e.g. Redis, database).
func WithHistory(h history.HistoryStore) Option {
	return func(c *Chess) { c.history = h }
}

// WithTracker overrides the default position tracker used for threefold repetition detection.
func WithTracker(t tracker.Tracker) Option {
	return func(c *Chess) { c.tracker = t }
}

// WithHasher overrides the default Zobrist hasher.
func WithHasher(h hash.Hasher) Option {
	return func(c *Chess) { c.hasher = h }
}

// WithFEN sets a custom starting position. The FEN string is decoded during
// New — an invalid FEN causes New to return an error.
func WithFEN(fen string) Option {
	return func(c *Chess) { c.fen = fen }
}

// State serializes the current game into a ChessState snapshot suitable for
// persistence or transmission. The snapshot can be restored with LoadState.
func (c *Chess) State() core.ChessState {
	return core.ChessState{
		Board:           *c.ctx.Board,
		SideToMove:      c.ctx.SideToMove,
		Sides:           c.ctx.Sides,
		EnPassantTarget: c.ctx.EnPassantTarget,
		HalfMoveClock:   c.ctx.HalfMoveClock,
		FullMoveNumber:  c.ctx.FullMoveNumber,
		Hash:            c.hash,
	}
}

// LoadState restores a previously saved ChessState into this Chess instance.
// Note: history and tracker are not part of ChessState and are not restored.
func (c *Chess) LoadState(state core.ChessState) {
	c.ctx.Board = &state.Board
	c.ctx.SideToMove = state.SideToMove
	c.ctx.Sides = state.Sides
	c.ctx.EnPassantTarget = state.EnPassantTarget
	c.ctx.HalfMoveClock = state.HalfMoveClock
	c.ctx.FullMoveNumber = state.FullMoveNumber
	c.hash = state.Hash
}

// TurnContext returns a copy of the current board position and game state.
// Callers may read from it freely; mutating it does not affect the game.
func (c *Chess) TurnContext() core.TurnContext {
	return c.ctx
}

// Hash returns the Zobrist hash of the current position. The hash is
// maintained incrementally — it is updated on every MakeMove and UndoMove
// without a full board scan.
func (c *Chess) Hash() uint64 {
	return c.hash
}

// LegalMoves returns all legal moves for the piece at from. Returns an empty
// slice if the square is empty, occupied by the opponent, or has no legal moves.
func (c *Chess) LegalMoves(from core.Position) []core.Move {
	var buf [engine.MAX_TOTAL_MOVES]core.Move
	return c.engine.GetLegalMoves(buf[:0], from, c.ctx)
}

// GameResult evaluates all termination conditions and returns the current
// result. Returns InProgress when the game continues. Evaluates cheapest
// conditions first: fifty-move rule, threefold repetition, insufficient
// material, then checkmate/stalemate (most expensive). Call after every move.
func (c *Chess) GameResult() core.GameResult {
	return c.rules.GetGameResult(c.ctx, c.engine, c.tracker, c.hash)
}

// MakeMove validates and applies move to the current position. Returns
// ErrIllegalMove if the move is not legal. On success the board, hash,
// tracker, and history are all updated atomically in the correct order:
// Apply first (so ctx reflects the new position), then hash (which needs
// the post-apply ctx to read the new castling rights), then tracker and
// history.
func (c *Chess) MakeMove(move core.Move) error {
	if !c.engine.IsLegalMove(move, c.ctx) {
		return ErrIllegalMove
	}

	// Apply mutates the board and returns a snapshot for undo.
	snapshot := c.engine.Apply(&c.ctx, move)

	// Flip side to move
	c.ctx.SideToMove = c.ctx.SideToMove.Opponent()

	// Hash after Apply so ctx.Sides reflects the new castling rights.
	moveHash := core.NewMoveHash(snapshot, c.ctx)
	c.hash = c.hasher.Hash(c.hash, moveHash)

	// Record the new position and save the snapshot.
	c.tracker.Record(c.hash)
	c.history.Push(snapshot)

	return nil
}

// UndoMove reverts the last move. Returns ErrNothingToUndo if the history
// is empty. The revert order is the strict inverse of MakeMove: tracker and
// hash are reverted before the board — NewMoveHash needs ctx still in the
// post-move state to read the castling rights that were active when the move
// was hashed. Only then is it safe to call engine.Undo.
func (c *Chess) UndoMove() error {
	snapshot, ok := c.history.Pop()
	if !ok {
		return ErrNothingToUndo
	}

	// Revert tracker and hash before the board — ctx must still be post-move.
	c.tracker.Undo(c.hash)

	moveHash := core.NewMoveHash(snapshot, c.ctx)
	c.hash = c.hasher.Hash(c.hash, moveHash)

	// Now safe to revert the board.
	c.engine.Undo(&c.ctx, snapshot)

	// Flip side to move back
	c.ctx.SideToMove = c.ctx.SideToMove.Opponent()

	return nil
}
