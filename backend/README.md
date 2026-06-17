# Chess Backend

A chess engine written in Go, built for correctness and clear separation of concerns.

## Project Structure

```
backend/
├── core/      # Board, position, piece, move, and game-state primitives
├── engine/    # Move generation, legality checking, apply/undo, attack detection
├── piece/     # Per-piece-type move logic (pawn, knight, bishop, rook, queen, king)
└── main.go
```

## Packages

- **`core`** — Data model (board, position, pieces, moves, game state). Plain value types designed to be cheap to copy.
- **`piece`** — Move rules for each piece type. Stateless implementations shared via a global provider — zero allocation per lookup.
- **`engine`** — The rules layer: generates moves, filters for legality, applies and undoes moves, detects attacks.

## Current state

`core`, `engine`, and `piece` are complete with tests and benchmarks. `main.go` is a placeholder — no game orchestration, CLI, FEN parsing, or search yet.

## Benchmarks

See [BENCHMARKS.md](./BENCHMARKS.md).
