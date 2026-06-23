// Package core defines the fundamental chess domain types and the helper
// methods that operate on them. It has no chess-logic dependencies — no move
// generation, no rule enforcement, no I/O — so every other package can import
// it without creating import cycles.
//
// The package covers:
//
//   - Board coordinates: File, Rank, and the combined Position type, each with
//     bounds-checked arithmetic and parse/format helpers.
//   - Piece identity: PieceType (PAWN…KING), PieceColor (WHITE/BLACK), and the
//     Piece struct, with ASCII letter parsing (ParsePiece / Char).
//   - Board representation: Square (a one-byte cell that packs piece type and
//     color) and the 64-element Board array that uses it.
//   - Move description: Move, MoveType (NORMAL, CASTLING, EN_PASSANT, PROMOTION),
//     and the derived queries (IsDoublePawnPush, EnPassantTarget, etc.).
//   - Game state: SideState (king position + castling rights), MoveContext,
//     TurnContext (the complete mutable position handed to the engine), and
//     Snapshot (the pre-move state needed to undo a move).
//   - Game outcome: GameStatus, DrawReason, and GameResult.
//   - Hash primitives: MoveHash (the before/after context needed for incremental
//     Zobrist updates).
package core
