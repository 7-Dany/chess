import { useEffect, useRef } from "react";
import { useControls } from "leva";
import { useBoard } from "../../context/useBoard";
import { NO_POSITION, positionLabel } from "../../chess/position";
import type { Mark } from "../../chess/board";

/** Mark options for the dropdown. */
const MARK_OPTIONS: Mark[] = ["none", "selected", "legal-move", "capture"];

/**
 * Registers a "Board State" section in the debug panel.
 *
 * Shows the currently selected square (set by clicking the 3D board)
 * and its mark. Changing the mark dropdown instantly updates that
 * square. "Deselect" clears the selection. "Clear All Marks" resets
 * every square.
 */
export function useBoardStateControls() {
  const { state, setMark, deselect, clearMarks } = useBoard();

  // Refs so Leva's one-time closures always see current values.
  // Leva's factory function runs once; its onChange/button closures
  // would capture stale state without these.
  const selectedPosRef = useRef(state.selectedPosition);
  const setMarkRef = useRef(setMark);
  const deselectRef = useRef(deselect);
  const clearMarksRef = useRef(clearMarks);

  useEffect(() => {
    selectedPosRef.current = state.selectedPosition;
    setMarkRef.current = setMark;
    deselectRef.current = deselect;
    clearMarksRef.current = clearMarks;
  });

  const hasSelection = state.selectedPosition !== NO_POSITION;

  // Read the current mark of the selected square so we can seed the dropdown.
  const currentMark: Mark = hasSelection
    ? state.board[state.selectedPosition].mark
    : "none";

  const [, set] = useControls("Board State", () => ({
    selected: {
      label: "Selected Square",
      value: hasSelection ? positionLabel(state.selectedPosition) : "-",
      editable: false,
    },
    mark: {
      label: "Mark",
      value: currentMark as Mark,
      options: MARK_OPTIONS,
      onChange: (val: Mark) => {
        const pos = selectedPosRef.current;
        if (pos !== NO_POSITION) {
          setMarkRef.current(pos, val);
        }
      },
      transient: true,
    },
  }));

  // Sync the read-only fields when board state changes
  useEffect(() => {
    set({
      selected: hasSelection ? positionLabel(state.selectedPosition) : "-",
      mark: currentMark,
    });
  }, [state.selectedPosition, currentMark, hasSelection, set]);
}
