/**
 * A board column, A through H.
 */
export type File = 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7;

/**
 * A board row, 1 through 8.
 */
export type Rank = 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7;

/**
 * A board position encoded as a single number (0–63).
 * Formula: file * 8 + rank
 *
 * A1=0, A2=1 ... A8=7, B1=8 ... H8=63
 * NoPosition=64 is the sentinel for "no square".
 */
export type Position = number;

/** Not valid position, used for filtering */
export const NO_POSITION: Position = 64;

/**
 * Named Positions
 */
export const A1 = 0,
  A2 = 1,
  A3 = 2,
  A4 = 3,
  A5 = 4,
  A6 = 5,
  A7 = 6,
  A8 = 7;
export const B1 = 8,
  B2 = 9,
  B3 = 10,
  B4 = 11,
  B5 = 12,
  B6 = 13,
  B7 = 14,
  B8 = 15;
export const C1 = 16,
  C2 = 17,
  C3 = 18,
  C4 = 19,
  C5 = 20,
  C6 = 21,
  C7 = 22,
  C8 = 23;
export const D1 = 24,
  D2 = 25,
  D3 = 26,
  D4 = 27,
  D5 = 28,
  D6 = 29,
  D7 = 30,
  D8 = 31;
export const E1 = 32,
  E2 = 33,
  E3 = 34,
  E4 = 35,
  E5 = 36,
  E6 = 37,
  E7 = 38,
  E8 = 39;
export const F1 = 40,
  F2 = 41,
  F3 = 42,
  F4 = 43,
  F5 = 44,
  F6 = 45,
  F7 = 46,
  F8 = 47;
export const G1 = 48,
  G2 = 49,
  G3 = 50,
  G4 = 51,
  G5 = 52,
  G6 = 53,
  G7 = 54,
  G8 = 55;
export const H1 = 56,
  H2 = 57,
  H3 = 58,
  H4 = 59,
  H5 = 60,
  H6 = 61,
  H7 = 62,
  H8 = 63;

/**
 * File (Columns) A -> H
 */
export const FILE_A = 0,
  FILE_B = 1,
  FILE_C = 2,
  FILE_D = 3,
  FILE_E = 4,
  FILE_F = 5,
  FILE_G = 6,
  FILE_H = 7;

/**
 * Rank (Raws) 1 -> 8
 */
export const RANK_1 = 0,
  RANK_2 = 1,
  RANK_3 = 2,
  RANK_4 = 3,
  RANK_5 = 4,
  RANK_6 = 5,
  RANK_7 = 6,
  RANK_8 = 7;

/**
 * Creates a Position from a file and rank.
 *
 * @example
 * newPosition(FILE_A, RANK_1) // 0  (A1)
 * newPosition(FILE_H, RANK_8) // 63 (H8)
 */
export function newPosition(file: File, rank: Rank): Position {
  return file * 8 + rank;
}

/**
 * Returns the file (column) of a position.
 */
export function fileOf(pos: Position): File {
  return Math.floor(pos / 8) as File;
}

/**
 * Returns the rank (row) of a position.
 */
export function rankOf(pos: Position): Rank {
  return (pos % 8) as Rank;
}

/**
 * Returns true if the position is a valid board square.
 */
export function isValidPosition(pos: Position): boolean {
  return pos >= 0 && pos < NO_POSITION;
}

/**
 * Adds a value to a file. Returns null if the result is off the board.
 *
 * @example
 * addFile(FILE_A, 1)  // 1  (FILE_B)
 * addFile(FILE_H, 1)  // null (off board)
 */
export function addFile(file: File, value: number): File | null {
  const result = file + value;
  if (result < 0 || result > 7) return null;
  return result as File;
}

/**
 * Adds a value to a rank. Returns null if the result is off the board.
 *
 * @example
 * addRank(RANK_1, 2)  // 2  (RANK_3)
 * addRank(RANK_8, 1)  // null (off board)
 */
export function addRank(rank: Rank, value: number): Rank | null {
  const result = rank + value;
  if (result < 0 || result > 7) return null;
  return result as Rank;
}

const FILE_LABELS = ["A", "B", "C", "D", "E", "F", "G", "H"] as const;

/**
 * Returns the string label of a position, e.g. "A1", "H8".
 * Returns "-" for NoPosition.
 */
export function positionLabel(pos: Position): string {
  if (!isValidPosition(pos)) return "-";
  return FILE_LABELS[fileOf(pos)] + (rankOf(pos) + 1);
}
