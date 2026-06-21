# Benchmarks

Run with:

```
go test ./... -bench=. -benchmem
```

## Headline

**Every engine method is allocation-free.** All 44 benchmarks report
`0 B/op, 0 allocs/op`. The numbers below are from a local run; per-op times
will vary by machine, but the allocation column is exact and machine-
independent — it is `0` for every benchmark, every time.

## Results

Every row shows `0 B/op` and `0 allocs/op`. The `ns/op` column is the
per-operation time on the local machine.

### GetPseudoLegalMoves — pseudo-legal move generation per piece

Generates moves for a single piece without the king-safety filter. The
buffer is stack-allocated; passing `buf[:0]` each iteration means no
allocation.

| Benchmark | ns/op | B/op | allocs/op |
|---|---|---|---|
| GetPseudoLegalMoves Pawn | 30.30 | 0 | 0 |
| GetPseudoLegalMoves Knight | 38.92 | 0 | 0 |
| GetPseudoLegalMoves Bishop | 24.54 | 0 | 0 |
| GetPseudoLegalMoves Rook | 22.35 | 0 | 0 |
| GetPseudoLegalMoves Queen | 32.41 | 0 | 0 |
| GetPseudoLegalMoves King | 110.2 | 0 | 0 |
| GetPseudoLegalMoves Kiwipete Knight | 73.96 | 0 | 0 |
| GetPseudoLegalMoves Kiwipete Queen | 97.68 | 0 | 0 |

### GetLegalMoves — legal move generation per piece

Generates pseudo-legal moves AND filters them for king safety. Each
pseudo-move is applied, the king is checked, and the move is undone. This is
the real per-piece cost during search.

| Benchmark | ns/op | B/op | allocs/op |
|---|---|---|---|
| GetLegalMoves Pawn | 248.7 | 0 | 0 |
| GetLegalMoves Knight | 246.8 | 0 | 0 |
| GetLegalMoves Bishop | 28.22 | 0 | 0 |
| GetLegalMoves Rook | 27.57 | 0 | 0 |
| GetLegalMoves Queen | 39.41 | 0 | 0 |
| GetLegalMoves King | 118.4 | 0 | 0 |
| GetLegalMoves Kiwipete Knight | 863.3 | 0 | 0 |
| GetLegalMoves Kiwipete Queen | 1120 | 0 | 0 |
| GetLegalMoves Kiwipete King | 892.4 | 0 | 0 |

The starting-position pieces (Pawn, Knight) are slow because every pseudo-
move must be filtered — a pawn push or knight move rarely resolves check, so
nearly all of them are tried and most are kept. The sliders (Bishop, Rook,
Queen) are fast on the starting position because they're hemmed in by pawns
and generate 0–1 moves. The Kiwipete cases have many pseudo-moves, so the
king-safety filter runs many apply/undo cycles — the dominant cost during
search.

### HasAnyLegalMoves — checkmate / stalemate detection

Returns true as soon as one legal move is found; returns false only after
trying every pseudo-move of every piece. The "false" cases (Fool's Mate,
Stalemate) are the worst case.

| Benchmark | ns/op | B/op | allocs/op |
|---|---|---|---|
| HasAnyLegalMoves Start | 151.5 | 0 | 0 |
| HasAnyLegalMoves Fool's Mate | 2298 | 0 | 0 |
| HasAnyLegalMoves Stalemate | 351.0 | 0 | 0 |
| HasAnyLegalMoves Kiwipete | 149.7 | 0 | 0 |
| HasAnyLegalMoves Endgame | 257.0 | 0 | 0 |

### IsSquareAttacked — attack detection

Scans from the target outward: 3 leaper checks (knight, king, pawn) + 8
slider rays (4 diagonal for bishop/queen, 4 orthogonal for rook/queen).
Never allocates.

| Benchmark | ns/op | B/op | allocs/op |
|---|---|---|---|
| IsSquareAttacked Empty | 116.5 | 0 | 0 |
| IsSquareAttacked Start | 71.95 | 0 | 0 |
| IsSquareAttacked Kiwipete | 83.39 | 0 | 0 |
| IsSquareAttacked Attacked | 95.06 | 0 | 0 |
| IsSquareAttacked Corner | 79.52 | 0 | 0 |

### Apply — applying a move to the board

Each iteration copies the base context (so Apply has a fresh board to
mutate) and applies the move. The `Copy()` is included in the measurement;
Apply itself is 0-alloc.

| Benchmark | ns/op | B/op | allocs/op |
|---|---|---|---|
| Apply Normal pawn push | 25.03 | 0 | 0 |
| Apply Normal knight move | 23.56 | 0 | 0 |
| Apply Capture | 24.66 | 0 | 0 |
| Apply En passant | 23.55 | 0 | 0 |
| Apply Castling (king side) | 24.29 | 0 | 0 |
| Apply Castling (queen side) | 24.06 | 0 | 0 |
| Apply Promotion | 25.14 | 0 | 0 |
| Apply Double pawn push | 24.96 | 0 | 0 |

### Undo — reversing an Apply

Each iteration copies the base context, applies the move, then undoes it.
The measurement includes Apply + Copy (unavoidable setup) + Undo. To see the
Undo cost in isolation, compare these numbers against the matching Apply
benchmark above — the difference is roughly the Undo cost.

| Benchmark | ns/op | B/op | allocs/op |
|---|---|---|---|
| Undo Normal pawn push | 39.81 | 0 | 0 |
| Undo Normal knight move | 38.28 | 0 | 0 |
| Undo Capture | 40.71 | 0 | 0 |
| Undo En passant | 40.70 | 0 | 0 |
| Undo Castling (king side) | 42.27 | 0 | 0 |
| Undo Promotion | 41.37 | 0 | 0 |

### Apply + Undo round-trip — the real search hot-path cost

During a search, every node applies a move, recurses, then undoes. The
round-trip cost (Apply + Undo, excluding the copy) is what matters for
nodes-per-second.

| Benchmark | ns/op | B/op | allocs/op |
|---|---|---|---|
| Apply + Undo Pawn push | 40.15 | 0 | 0 |
| Apply + Undo Knight move | 38.37 | 0 | 0 |
| Apply + Undo Castling | 43.25 | 0 | 0 |

## Takeaway

- **Allocations: fully eliminated.** `GetPseudoLegalMoves`, `GetLegalMoves`,
  and `HasAnyLegalMoves` previously allocated 1–16 times per call; they now
  allocate zero. `Apply`, `Undo`, `IsSquareAttacked` were already
  allocation-free and stay that way.
- **Move generation is fast.** Pseudo-legal generation for a single piece
  takes 22–110 ns; the king is the slowest (it checks castling eligibility,
  which calls `IsSquareAttacked` up to 5 times). Legal move generation is
  dominated by the apply/undo cycle (one `IsSquareAttacked` per pseudo-move).
- **`DefaultEngine` is stateless** (`struct{}`): all state lives in the
  `TurnContext` passed to each method, so a single engine instance is safe to
  share across goroutines.

## Buffer-size contract

Callers of `GetPseudoLegalMoves` / `GetLegalMoves` / `PseudoLegalMoves` /
`Attacks` must pass a buffer of at least (== `core.MAX_MOVES` == 32) to guarantee no allocation. The largest
single-piece move set is a queen on an open board (27 moves); a king with
both castling moves adds at most 2. 32 covers both with headroom.
