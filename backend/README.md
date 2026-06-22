# Chess Backend

A chess engine written in Go, built for correctness and clear separation of concerns.

## Project Structure

```
backend/
├── core/      # Board, position, piece, move, and game-state primitives
├── engine/    # Move generation, legality checking, apply/undo, attack detection
├── piece/     # Per-piece-type move logic (pawn, knight, bishop, rook, queen, king)
├── fen/       # FEN string parsing and serialization
├── testutil/  # Shared test helpers
└── main.go
```

## Packages

- **`core`** — Data model (board, position, pieces, moves, game state). Plain value types designed to be cheap to copy.
- **`piece`** — Move rules for each piece type. Stateless zero-size value types. `PseudoLegalMoves`/`Attacks` append into a caller-owned buffer, so move generation never allocates (see `core.MAX_MOVES`).
- **`engine`** — The rules layer: generates moves, filters for legality, applies and undoes moves, detects attacks. Every method is allocation-free; `GetPseudoLegalMoves`/`GetLegalMoves` take a caller-owned `[]core.Move` buffer, and piece dispatch on the hot path is a concrete type switch (not the `Piece` interface) so escape analysis can keep the buffer on the stack. `DefaultEngine` is stateless (`struct{}`).
- **`fen`** — Parses and serializes Forsyth-Edwards Notation strings. `Decode` fills a caller-owned `TurnContext`; `Encode` serializes one back to a FEN string.

## Correctness — Perft Validation

The engine is validated against the six standard perft positions from the
[Chess Programming Wiki](https://www.chessprogramming.org/Perft_Results).
Perft (performance test) recursively counts the number of leaf nodes in the
full move tree to a given depth, then compares against trusted reference
values. If the count matches, the move generator is correct.

All six positions pass at depth 1-3, and positions 1, 3, and 4 pass at
depth 4:

| Position | Depth 1 | Depth 2 | Depth 3 | Depth 4 |
|---|---|---|---|---|
| Start | 20 | 400 | 8,902 | 197,281 |
| Kiwipete | 48 | 2,039 | 97,862 | — |
| Position 3 | 14 | 191 | 2,812 | 43,238 |
| Position 4 | 6 | 264 | 9,467 | 422,333 |
| Position 5 | 44 | 1,486 | 62,379 | — |
| Position 6 | 46 | 2,079 | 89,890 | — |

Run the perft tests:

```
go test ./engine -run TestPerft -v
```

These positions stress every edge case: en passant (including discovered-check
en passant), promotion, promotion-captures, castling through/into check,
pins, and double-checks.

## Current state

`core`, `engine`, `piece`, and `fen` are complete with tests and benchmarks.
`main.go` is a placeholder — no game orchestration, CLI, or search yet.

All engine methods are allocation-free; see [BENCHMARKS.md](./BENCHMARKS.md).

## Benchmarks

See [BENCHMARKS.md](./BENCHMARKS.md).
