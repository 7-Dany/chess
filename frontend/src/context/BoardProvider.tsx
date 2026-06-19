import React, { useCallback, useReducer } from "react";
import { BoardContext, boardReducer, initialState } from "./useBoard";
import type { Board, Mark } from "../chess/board";
import type { Position } from "../chess/position";

export function BoardProvider({ children }: { children: React.ReactNode }) {
  const [state, dispatch] = useReducer(boardReducer, initialState);

  const selectPosition = useCallback((pos: Position) => {
    dispatch({ type: "SELECT_POSITION", payload: pos });
  }, []);

  const deselect = useCallback(() => {
    dispatch({ type: "DESELECT" });
  }, []);

  const setMark = useCallback((position: Position, mark: Mark) => {
    dispatch({ type: "SET_MARK", payload: { position, mark } });
  }, []);

  const clearMarks = useCallback(() => {
    dispatch({ type: "CLEAR_MARKS" });
  }, []);

  const setBoard = useCallback((board: Board) => {
    dispatch({ type: "SET_BOARD", payload: board });
  }, []);

  return (
    <BoardContext.Provider
      value={{ state, selectPosition, deselect, setMark, clearMarks, setBoard }}
    >
      {children}
    </BoardContext.Provider>
  );
}
