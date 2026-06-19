/**
 * The type of a chess piece.
 */
export type PieceType =
  | "PAWN"
  | "KNIGHT"
  | "BISHOP"
  | "ROOK"
  | "QUEEN"
  | "KING";

/**
 * The color of a chess piece.
 */
export type PieceColor = "WHITE" | "BLACK";

/**
 * A chess piece with a type and color.
 */
export interface Piece {
  readonly type: PieceType;
  readonly color: PieceColor;
}

/**
 * All piece types ordered weakest → strongest.
 */
export const PIECE_TYPE_ORDER: readonly PieceType[] = [
  "PAWN",
  "KNIGHT",
  "BISHOP",
  "ROOK",
  "QUEEN",
  "KING",
];

/**
 * Piece types a pawn can promote to.
 */
export const PROMOTABLE_PIECES: readonly PieceType[] = [
  "KNIGHT",
  "BISHOP",
  "ROOK",
  "QUEEN",
];

/**
 * Creates a Piece. Prefer this over object literals for consistency.
 *
 * @example
 * const whiteKing = createPiece('KING', 'WHITE');
 */
export function createPiece(type: PieceType, color: PieceColor): Piece {
  return { type, color };
}

/**
 * Returns the opposite color.
 */
export function oppositeColor(color: PieceColor): PieceColor {
  return color === "WHITE" ? "BLACK" : "WHITE";
}

/**
 * Returns true if the piece is white.
 */
export function isWhite(piece: Piece): boolean {
  return piece.color === "WHITE";
}

/**
 * Returns true if the piece is black.
 */
export function isBlack(piece: Piece): boolean {
  return piece.color === "BLACK";
}

/**
 * Returns true if two pieces belong to opposing sides.
 * Used to determine if a move is a capture.
 */
export function isOpponent(a: Piece, b: Piece): boolean {
  return a.color !== b.color;
}

/**
 * Returns true if two pieces belong to the same side.
 * Used to prevent moving onto a friendly piece.
 */
export function isFriendly(a: Piece, b: Piece): boolean {
  return a.color === b.color;
}
