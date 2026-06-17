# Benchmarks

Measured on an i5-4590 CPU @ 3.30GHz. Run with:

```
go test ./... -bench=. -benchmem
```

## Results

| Category | Benchmark | ns/op | B/op | allocs/op |
|---|---|---|---|---|
| Pseudo-legal moves | Pawn | 215.5 | 98 | 4 |
| Pseudo-legal moves | Knight | 162.1 | 40 | 2 |
| Pseudo-legal moves | Bishop | 120.0 | 128 | 1 |
| Pseudo-legal moves | Rook | 124.8 | 128 | 1 |
| Pseudo-legal moves | Queen | 209.9 | 256 | 1 |
| Pseudo-legal moves | King | 372.0 | 112 | 3 |
| Pseudo-legal moves | Kiwipete Knight | 210.7 | 88 | 2 |
| Pseudo-legal moves | Kiwipete Queen | 296.8 | 256 | 1 |
| Legal moves | Pawn | 600.0 | 122 | 5 |
| Legal moves | Knight | 537.0 | 64 | 3 |
| Legal moves | Bishop | 136.7 | 128 | 1 |
| Legal moves | Rook | 134.5 | 128 | 1 |
| Legal moves | Queen | 230.0 | 256 | 1 |
| Legal moves | King | 387.0 | 112 | 3 |
| Legal moves | Kiwipete Knight | 1597 | 152 | 3 |
| Legal moves | Kiwipete Queen | 2093 | 352 | 2 |
| Legal moves | Kiwipete King | 1837 | 208 | 5 |
| Has any legal moves | Start position | 592.9 | 226 | 5 |
| Has any legal moves | Fool's Mate | 7149 | 1728 | 43 |
| Has any legal moves | Stalemate | 759.9 | 88 | 2 |
| Has any legal moves | Kiwipete | 383.3 | 128 | 1 |
| Has any legal moves | Endgame | 530.4 | 128 | 1 |
| Square attacked | Empty board | 209.5 | 0 | 0 |
| Square attacked | Start position | 121.3 | 0 | 0 |
| Square attacked | Kiwipete | 140.6 | 0 | 0 |
| Square attacked | Attacked square | 163.9 | 0 | 0 |
| Square attacked | Corner square | 125.1 | 0 | 0 |
| Apply | Normal pawn push | 36.84 | 0 | 0 |
| Apply | Normal knight move | 35.73 | 0 | 0 |
| Apply | Capture | 34.38 | 0 | 0 |
| Apply | En passant | 36.74 | 0 | 0 |
| Apply | Castling (king side) | 42.03 | 0 | 0 |
| Apply | Castling (queen side) | 41.53 | 0 | 0 |
| Apply | Promotion | 35.55 | 0 | 0 |
| Apply | Double pawn push | 38.26 | 0 | 0 |
| Undo | Normal pawn push | 50.60 | 0 | 0 |
| Undo | Normal knight move | 47.97 | 0 | 0 |
| Undo | Capture | 46.69 | 0 | 0 |
| Undo | En passant | 50.87 | 0 | 0 |
| Undo | Castling (king side) | 60.95 | 0 | 0 |
| Undo | Promotion | 47.34 | 0 | 0 |
| Apply + Undo round-trip | Pawn push | 50.21 | 0 | 0 |
| Apply + Undo round-trip | Knight move | 48.55 | 0 | 0 |
| Apply + Undo round-trip | Castling | 62.48 | 0 | 0 |

## Takeaway

`Apply`, `Undo`, and `IsSquareAttacked` are fully allocation-free (0 B/op, 0 allocs/op) — the path that matters most for high-throughput use like perft or search, since it runs millions of times per second.

Move generation (`GetPseudoLegalMoves` / `GetLegalMoves`) still allocates a slice per call; this is the next obvious optimization target, likely via a reusable move buffer instead of fresh slice allocation on every query.
