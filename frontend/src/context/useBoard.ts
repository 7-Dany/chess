import { createContext, use } from "react";
import {
  createInitialBoard,
  normaliseBackendBoard,
  type Board,
  type Mark,
} from "../chess/board";
import { NO_POSITION, type Position } from "../chess/position";

export interface BoardState {
  board: Board;
  selectedPosition: Position;
  enPassantTarget: Position;
}

export const initialState: BoardState = {
  board: createInitialBoard(),
  selectedPosition: NO_POSITION,
  enPassantTarget: NO_POSITION,
};

export type BoardAction =
  | { type: "SELECT_POSITION"; payload: Position }
  | { type: "DESELECT" }
  | { type: "SET_MARK"; payload: { position: Position; mark: Mark } }
  | { type: "CLEAR_MARKS" }
  | { type: "SET_BOARD"; payload: Board };

export function boardReducer(
  state: BoardState,
  action: BoardAction,
): BoardState {
  switch (action.type) {
    case "SELECT_POSITION": {
      const clickedPos = action.payload;

      // Toggle: clicking the same square deselects it
      if (state.selectedPosition === clickedPos) {
        return {
          ...state,
          selectedPosition: NO_POSITION,
          board: state.board.map((sq, i) =>
            i === clickedPos ? { ...sq, mark: "none" as Mark } : sq,
          ),
        };
      }

      // Select the new square: clear previous mark, set new one
      const prevSelected = state.selectedPosition;
      return {
        ...state,
        selectedPosition: clickedPos,
        board: state.board.map((sq, i) => {
          if (i === clickedPos) {
            return { ...sq, mark: "selected" as Mark };
          }
          if (i === prevSelected && prevSelected !== NO_POSITION) {
            return { ...sq, mark: "none" as Mark };
          }
          return sq;
        }),
      };
    }

    case "DESELECT": {
      const prevSelected = state.selectedPosition;
      if (prevSelected === NO_POSITION) return state;

      return {
        ...state,
        selectedPosition: NO_POSITION,
        board: state.board.map((sq, i) =>
          i === prevSelected ? { ...sq, mark: "none" as Mark } : sq,
        ),
      };
    }

    case "SET_MARK": {
      const { position, mark } = action.payload;
      return {
        ...state,
        board: state.board.map((sq, i) =>
          i === position ? { ...sq, mark } : sq,
        ),
      };
    }

    case "CLEAR_MARKS": {
      return {
        ...state,
        selectedPosition: NO_POSITION,
        board: state.board.map((sq) =>
          sq.mark !== "none" ? { ...sq, mark: "none" as Mark } : sq,
        ),
      };
    }

    case "SET_BOARD": {
      // The backend payload lacks `mark`. Normalise it so every
      // square has `mark: "none"`. Also reset selection state.
      return {
        ...state,
        selectedPosition: NO_POSITION,
        enPassantTarget: NO_POSITION,
        board: normaliseBackendBoard(action.payload),
      };
    }

    default:
      return state;
  }
}

export interface BoardContextValue {
  state: BoardState;
  selectPosition: (pos: Position) => void;
  deselect: () => void;
  setMark: (position: Position, mark: Mark) => void;
  clearMarks: () => void;
  setBoard: (board: Board) => void;
}

export const BoardContext = createContext<BoardContextValue | null>(null);

export function useBoard(): BoardContextValue {
  const context = use(BoardContext);
  if (!context) {
    throw new Error("useBoard must be used inside BoardProvider");
  }
  return context;
}
