import type { Mark } from "../chess/board";

export const BOARD_COLORS = {
  light: {
    lightSquare: "#f0d9b5",
    darkSquare: "#b58863",
    selected: "#0ea5e9",
    legalMove: "#84cc16",
    capture: "#ef4444",
  },
  dark: {
    lightSquare: "#6b7280",
    darkSquare: "#1f2937",
    selected: "#38bdf8",
    legalMove: "#a3e635",
    capture: "#f87171",
  },
} as const;

export function getSquareColor(
  isLight: boolean,
  mark: Mark,
  theme: keyof typeof BOARD_COLORS,
): string {
  const colors = BOARD_COLORS[theme];
  switch (mark) {
    case "selected":
      return colors.selected;
    case "legal-move":
      return colors.legalMove;
    case "capture":
      return colors.capture;
    case "none":
      return isLight ? colors.lightSquare : colors.darkSquare;
  }
}
