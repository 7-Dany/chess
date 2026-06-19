import {
  createPiece,
  type Piece,
  type PieceColor,
  type PieceType,
} from "./piece";
import type { Position } from "./position";

/** Square marks */
export type Mark = "none" | "selected" | "legal-move" | "capture";

/**
 * A single square on the board.
 *
 * `mark` is owned by the UI layer (the reducer in useBoard.ts)
 * and defaults to "none". It is **not** part of the backend Board
 * payload, so SET_BOARD must inject it when building state.board.
 */
export interface Square {
  readonly piece: Piece | null;
  readonly occupied: boolean;
  readonly mark: Mark;
}

/**
 * The board as a flat array of 64 squares.
 *
 * Index formula: file * 8 + rank (same as Position)
 */
export type Board = readonly Square[];

/**
 * Returns the square at a given position.
 */
export function squareAt(board: Board, pos: Position): Square {
  return board[pos];
}

/**
 * Returns the piece at a position, or null if unoccupied.
 */
export function pieceAt(board: Board, pos: Position): Piece | null {
  const square = board[pos];
  return square.occupied ? square.piece : null;
}

/**
 * Returns true if the square at a position is empty.
 */
export function isEmpty(board: Board, pos: Position): boolean {
  return !board[pos].occupied;
}

/**
 * Creates an empty square (no piece, no mark).
 */
function emptySquare(): Square {
  return { piece: null, occupied: false, mark: "none" };
}

/**
 * Creates a square occupied by the given piece.
 */
function occupiedSquare(type: PieceType, color: PieceColor): Square {
  return { piece: createPiece(type, color), occupied: true, mark: "none" };
}

/**
 * Returns the standard chess starting position as a 64-element Board.
 *
 * Layout (file × rank):
 *   Rank 1 (index 0–7):   White back rank  R N B Q K B N R
 *   Rank 2 (index 8–15):  White pawns
 *   Ranks 3–6:            Empty
 *   Rank 7 (index 48–55): Black pawns
 *   Rank 8 (index 56–63): Black back rank  R N B Q K B N R
 */
export function createInitialBoard(): Board {
  const squares: Square[] = Array.from({ length: 64 }, () => emptySquare());

  // White back rank (rank 1 = index 0..7 for file A..H)
  const whiteBackRank: PieceType[] = [
    "ROOK",
    "KNIGHT",
    "BISHOP",
    "QUEEN",
    "KING",
    "BISHOP",
    "KNIGHT",
    "ROOK",
  ];
  for (let file = 0; file < 8; file++) {
    squares[file * 8 + 0] = occupiedSquare(whiteBackRank[file], "WHITE");
  }

  // White pawns (rank 2)
  for (let file = 0; file < 8; file++) {
    squares[file * 8 + 1] = occupiedSquare("PAWN", "WHITE");
  }

  // Black pawns (rank 7)
  for (let file = 0; file < 8; file++) {
    squares[file * 8 + 6] = occupiedSquare("PAWN", "BLACK");
  }

  // Black back rank (rank 8)
  const blackBackRank: PieceType[] = [
    "ROOK",
    "KNIGHT",
    "BISHOP",
    "QUEEN",
    "KING",
    "BISHOP",
    "KNIGHT",
    "ROOK",
  ];
  for (let file = 0; file < 8; file++) {
    squares[file * 8 + 7] = occupiedSquare(blackBackRank[file], "BLACK");
  }

  return squares;
}

/**
 * Normalises a backend Board payload (which lacks `mark`) into
 * the full Square shape used by the UI.
 *
 * Every square gets `mark: "none"` so the reducer can later
 * swap individual squares to "selected" / "legal-move" / "capture"
 * with immutable updates.
 */
export function normaliseBackendBoard(
  raw: readonly { piece: Piece | null; occupied: boolean }[],
): Board {
  return raw.map((sq) => ({
    ...sq,
    mark: "none" as Mark,
  }));
}
