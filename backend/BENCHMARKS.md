# Benchmarks

Measured on an i5-4590 CPU @ 3.30GHz. Run with:

```
go test ./... -bench=. -benchmem
```

## Results

| Category | Benchmark | ns/op | B/op | allocs/op |
|---|---|---|---|---|
| Pseudo-legal moves | Pawn | 63.21 | 48 | 1 |
| Pseudo-legal moves | Knight | 83.21 | 80 | 1 |
| Pseudo-legal moves | Bishop | 74.88 | 128 | 1 |
| Pseudo-legal moves | Rook | 75.62 | 128 | 1 |
| Pseudo-legal moves | Queen | 112.5 | 256 | 1 |
| Pseudo-legal moves | King | 182.0 | 96 | 1 |
| Pseudo-legal moves | Kiwipete Knight | 117.1 | 80 | 1 |
| Pseudo-legal moves | Kiwipete Queen | 183.9 | 256 | 1 |
| Legal moves | Pawn | 305.1 | 48 | 1 |
| Legal moves | Knight | 318.9 | 80 | 1 |
| Legal moves | Bishop | 78.03 | 128 | 1 |
| Legal moves | Rook | 78.06 | 128 | 1 |
| Legal moves | Queen | 122.2 | 256 | 1 |
| Legal moves | King | 191.5 | 96 | 1 |
| Legal moves | Kiwipete Knight | 1004 | 80 | 1 |
| Legal moves | Kiwipete Queen | 1281 | 256 | 1 |
| Legal moves | Kiwipete King | 1067 | 96 | 1 |
| Has any legal moves | Start position | 269.6 | 176 | 2 |
| Has any legal moves | Fool's Mate | 3494 | 1408 | 16 |
| Has any legal moves | Stalemate | 462.3 | 96 | 1 |
| Has any legal moves | Kiwipete | 228.8 | 128 | 1 |
| Has any legal moves | Endgame | 321.0 | 128 | 1 |
| Square attacked | Empty board | 119.7 | 0 | 0 |
| Square attacked | Start position | 78.70 | 0 | 0 |
| Square attacked | Kiwipete | 91.53 | 0 | 0 |
| Square attacked | Attacked square | 97.17 | 0 | 0 |
| Square attacked | Corner square | 83.28 | 0 | 0 |
| Apply | Normal pawn push | 41.41 | 0 | 0 |
| Apply | Normal knight move | 38.37 | 0 | 0 |
| Apply | Capture | 40.74 | 0 | 0 |
| Apply | En passant | 41.87 | 0 | 0 |
| Apply | Castling (king side) | 27.33 | 0 | 0 |
| Apply | Castling (queen side) | 27.58 | 0 | 0 |
| Apply | Promotion | 23.87 | 0 | 0 |
| Apply | Double pawn push | 42.01 | 0 | 0 |
| Undo | Normal pawn push | 42.50 | 0 | 0 |
| Undo | Normal knight move | 43.60 | 0 | 0 |
| Undo | Capture | 43.74 | 0 | 0 |
| Undo | En passant | 42.11 | 0 | 0 |
| Undo | Castling (king side) | 48.42 | 0 | 0 |
| Undo | Promotion | 42.63 | 0 | 0 |
| Apply + Undo round-trip | Pawn push | 42.17 | 0 | 0 |
| Apply + Undo round-trip | Knight move | 42.90 | 0 | 0 |
| Apply + Undo round-trip | Castling | 47.11 | 0 | 0 |

## Takeaway

Move generation (`GetPseudoLegalMoves` / `GetLegalMoves`) now allocates a single slice per call across the board (previously 1-5 allocs depending on piece type), and is 35-71% faster. `HasAnyLegalMoves` and `IsSquareAttacked` improved similarly across all positions tested.

`Apply`, `Undo`, and `IsSquareAttacked` remain fully allocation-free (0 B/op, 0 allocs/op). `Undo` and castling/promotion `Apply` cases got faster, but plain `Apply` move types (pawn push, knight move, capture, en passant, double pawn push) regressed 7-19% in ns/op despite no new allocations — worth investigating if `Apply` throughput matters for the search path.
