import type { Mark } from "@/chess/board";
import { useBoardStateControls } from "./board/useBoardStateControls";
import { useMarkEffectControls } from "./effect/useMarkEffectControls";
import { useAppearanceDebugControls } from "./theme/useAppearanceControls";
import { useBoard } from "@/context/useBoard";
import { NO_POSITION } from "@/chess/position";

/**
 * Registers Leva debug controls for board and theme state.
 * Renders nothing — only side effects into the Leva panel.
 * Must be mounted inside BoardProvider and ThemeProvider.
 */
export function DebugPanel() {
  const { state } = useBoard();

  // Derive the active mark from the selected square
  const hasSelection = state.selectedPosition !== NO_POSITION;
  const activeMark: Mark = hasSelection
    ? state.board[state.selectedPosition].mark
    : "none";

  // Appearance (always visible)
  useAppearanceDebugControls();

  // Board state (always visible — shows selection + mark dropdown)
  useBoardStateControls();

  // Effects (only visible for the active mark)
  useMarkEffectControls(activeMark);

  return null;
}
