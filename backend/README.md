# Chess Backend

A chess engine written in Go, built for correctness and clear separation of concerns.

## Project Structure

```
backend/
├── core/      # Domain types: board, position, piece, move, game state, snapshots
├── engine/    # Move generation, legality filtering, apply/undo, attack detection
├── piece/     # Per-piece movement and attack rules (pawn, knight, bishop, rook, queen, king)
├── fen/       # FEN string parsing and serialization
├── hash/      # Incremental Zobrist hashing for position identity
├── history/   # Move history stack (undo support)
├── tracker/   # Position-occurrence counter (threefold repetition detection)
├── rules/     # Game termination rules (checkmate, stalemate, draws)
├── game/      # Orchestrator: coordinates all subsystems behind a single Chess API
├── testutil/  # Shared test helpers (context builders, board assertions, move assertions)
└── main.go
```

## Packages

**`core`** — The shared domain layer. Defines every type the other packages build on:
board coordinates (`File`, `Rank`, `Position`), piece identity (`PieceType`, `PieceColor`,
`Piece`), board representation (`Square`, `Board`), move description (`Move`, `MoveType`),
and the full game-state context (`TurnContext`, `SideState`, `Snapshot`). Has no
chess-logic dependencies — nothing here generates or validates moves — so all other
packages can import it without creating cycles.

**`piece`** — Movement and attack rules for each of the six piece types. Every
implementation (`Pawn`, `Knight`, `Bishop`, `Rook`, `Queen`, `King`) is a stateless
zero-value struct that satisfies the `Piece` interface. The three methods are:
- `PseudoLegalMoves` — all moves from a square, respecting geometry and blockers but not king safety.
- `Attacks` — every square threatened from a position (including friendly-occupied squares).
- `IsAttacking` — reverse scan: "does any piece of this type and color attack this square?".

All three append into a caller-owned buffer, so move generation never allocates
(see `core.MAX_MOVES`).

**`engine`** — The rules layer. Combines `core` types with `piece` logic to provide:
- `GetPseudoLegalMoves` / `GetLegalMoves` / `GetAllLegalMoves` — single-piece and whole-side move generation.
- `IsLegalMove` — reports whether a specific move is legal; used by the game orchestrator for validation.
- `HasAnyLegalMoves` — early-exit scan used for checkmate and stalemate detection.
- `IsSquareAttacked` — reverse-ray scan across all piece types for check detection and castling validation.
- `Apply` / `Undo` — mutate a `TurnContext` in place and revert it using a `Snapshot`.

All methods are allocation-free. `GetPseudoLegalMoves` and `GetLegalMoves` take a
caller-owned `[]core.Move` buffer. `DefaultEngine` is stateless so a single instance
is safe to share across goroutines.

**`fen`** — Parses and serializes Forsyth-Edwards Notation. `Decode` fills a
caller-owned `TurnContext` from a FEN string; `Encode` serializes one back.
Decode and Encode are inverses: `Encode(Decode(s))` round-trips to the same string.

**`hash`** — Incremental Zobrist hashing. `InitHash` computes a full 64-bit hash
from scratch (call once after FEN decode); `Hash` updates it incrementally after
each `Apply` or before each `Undo`. Because XOR is self-inverse, the same call
works in both directions. The random table is seeded from a fixed PCG source, so
hashes are stable across process restarts.

**`history`** — A stack of `core.Snapshot` values representing the moves played
so far. `MemoryStore` is the default in-memory implementation; the `HistoryStore`
interface leaves room for persistent backends (Redis, database) in the future.

**`tracker`** — Counts how many times each position hash has been seen. Used to
detect threefold repetition. `Record` is called after each `Apply` (once the hash
is updated); `Undo` is called before each move is reversed (before the hash reverts).

**`rules`** — Pure, stateless evaluator for all game-ending conditions. Every method
receives only the data it actually needs — no `*Chess` pointer. `GetGameResult` is
the single method to call after every move; it short-circuits from cheapest to most
expensive:
- `IsFiftyMoveRule` — O(1), single integer comparison on `HalfMoveClock`.
- `IsThreefoldRepetition` — O(1), single map lookup via the `Tracker`.
- `IsInsufficientMaterial` — O(64), one board scan; bails immediately on any pawn, rook, or queen.
- Checkmate / Stalemate — `HasAnyLegalMoves` called once; `IsSquareAttacked` on the king determines which.

`DefaultRules` is the standard implementation. The `Rules` interface allows alternative
implementations (e.g. custom draw rules) without touching the orchestrator.

**`game`** — The top-level orchestrator. `Chess` owns the mutable game state
(`TurnContext`, Zobrist hash) and coordinates all subsystems behind a clean API.
All subsystems have sensible defaults and can be overridden via functional options:

```go
// standard starting position, all defaults
g, err := game.New()

// custom FEN, persistent history
g, err := game.New(
    game.WithFEN("rnbqkbnr/pp1ppppp/..."),
    game.WithHistory(myRedisStore),
)
```

The two core methods enforce a strict call ordering:
- `MakeMove` — validates legality, applies the move, updates hash, records position, pushes to history.
- `UndoMove` — pops history, reverts tracker and hash (before the board), then reverts the board.

Interleaving direct `engine.Apply` / `engine.Undo` calls with `MakeMove` / `UndoMove`
will desync the hash, tracker, and history from the board.

**`testutil`** — Test-only helpers shared across suites. Provides `TurnContext`
builders (`NewTurn`, `WithSides`, `WithEnPassantTarget`), `SideState` factories
(`DefaultSides`, `FullWhite`, `FullBlack`), and assertion helpers
(`AssertSquareHas`, `AssertMovePresent`, `AssertMoveCount`, `AssertPositionsMatch`, etc.).

## Usage

### Starting a game

```go
import "github.com/7-Dany/chess/game"

// Standard starting position, all defaults.
g, err := game.New()
if err != nil {
    log.Fatal(err)
}

// Custom position via FEN.
g, err = game.New(game.WithFEN("r1bqkbnr/pppp1ppp/2n5/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R w KQkq - 2 3"))
```

### Making and undoing moves

```go
move := core.Move{
    Piece: core.Piece{Type: core.PAWN, Color: core.WHITE},
    From:  core.E2,
    To:    core.E4,
    Type:  core.NORMAL,
}

if err := g.MakeMove(move); err == game.ErrIllegalMove {
    // move was not legal in the current position
}

// Revert the last move.
if err := g.UndoMove(); err == game.ErrNothingToUndo {
    // no moves to undo
}
```

### Querying legal moves

```go
// All legal moves for the piece on E2.
moves := g.LegalMoves(core.E2)

// Inspect the current board state.
ctx := g.TurnContext()
fmt.Println(ctx.SideToMove) // WHITE or BLACK
fmt.Println(ctx.Board[core.E4]) // the square on E4
```

### Checking game result

```go
// Call after every MakeMove.
result := g.GameResult()

switch result.Status {
case core.InProgress:
    // game continues
case core.CheckMate:
    fmt.Printf("%v wins\n", result.Winner)
case core.Draw:
    switch result.DrawReason {
    case core.Stalemate:
        fmt.Println("draw by stalemate")
    case core.ThreefoldRepetition:
        fmt.Println("draw by threefold repetition")
    case core.FiftyMoveRule:
        fmt.Println("draw by fifty-move rule")
    case core.InsufficientMaterial:
        fmt.Println("draw by insufficient material")
    }
}
```

### Saving and restoring state

```go
// Serialize to a flat struct (e.g. for persistence or transmission).
state := g.State()

// Restore into a new Chess instance.
g2, _ := game.New()
g2.LoadState(state)
// Note: move history and position tracker are not part of ChessState
// and are reset on LoadState. Use WithHistory to plug in a persistent store
// if you need full undo support after a restore.
```

### Plugging in custom subsystems

```go
// All subsystems are optional — provide only what you want to override.
g, err := game.New(
    game.WithHistory(myRedisStore),   // persistent move history
    game.WithTracker(myTracker),      // custom repetition tracker
    game.WithEngine(myEngine),        // alternative move generator
)
```

## Correctness — Perft Validation

The engine is validated against the six standard perft positions from the
[Chess Programming Wiki](https://www.chessprogramming.org/Perft_Results).
Perft (performance test) recursively counts leaf nodes in the full move tree
to a given depth and compares against trusted reference values. A matching
count proves the move generator is correct.

All six positions pass at depth 1–3, and positions 1, 3, and 4 pass at depth 4:

| Position | Depth 1 | Depth 2 | Depth 3 | Depth 4 |
|---|---|---|---|---|
| Start | 20 | 400 | 8,902 | 197,281 |
| Kiwipete | 48 | 2,039 | 97,862 | — |
| Position 3 | 14 | 191 | 2,812 | 43,238 |
| Position 4 | 6 | 264 | 9,467 | 422,333 |
| Position 5 | 44 | 1,486 | 62,379 | — |
| Position 6 | 46 | 2,079 | 89,890 | — |

These positions collectively stress every edge case: en passant (including
discovered-check en passant), promotion, promotion-captures, castling
through/into/out-of check, pins, and double-checks.

```
go test ./engine -run TestPerft -v
```

## Current State

`core`, `piece`, `engine`, `fen`, `hash`, `history`, `tracker`, `rules`, and `game` are complete with tests.
`main.go` is a placeholder — no HTTP API or AI search yet.

All engine methods are allocation-free; see [BENCHMARKS.md](./BENCHMARKS.md).

## Benchmarks

See [BENCHMARKS.md](./BENCHMARKS.md).
