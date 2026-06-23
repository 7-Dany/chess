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
├── history/   # Move history stack (undo support, threefold-repetition input)
├── tracker/   # Position-occurrence counter (threefold repetition detection)
├── rules/     # Game termination rules (checkmate, stalemate, draws)
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

**`testutil`** — Test-only helpers shared across suites. Provides `TurnContext`
builders (`NewTurn`, `WithSides`, `WithEnPassantTarget`), `SideState` factories
(`DefaultSides`, `FullWhite`, `FullBlack`), and assertion helpers
(`AssertSquareHas`, `AssertMovePresent`, `AssertMoveCount`, `AssertPositionsMatch`, etc.).

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

`core`, `piece`, `engine`, `fen`, `hash`, `history`, `tracker`, and `rules` are complete with tests.
`main.go` is a placeholder — no game orchestration, CLI, or search yet.

All engine methods are allocation-free; see [BENCHMARKS.md](./BENCHMARKS.md).

## Benchmarks

See [BENCHMARKS.md](./BENCHMARKS.md).

