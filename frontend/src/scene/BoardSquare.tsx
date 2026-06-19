import { memo, useRef } from "react";
import { useFrame } from "@react-three/fiber";
import { rankOf, fileOf, type Position } from "../chess/position";
import type { Square } from "../chess/board";
import { getSquareColor } from "../theme/board";
import { MARK_EFFECTS } from "../theme/effects";
import { useTheme } from "../theme/useTheme";
import { MarkEffect } from "../effects/MarkEffect";
import type { MeshStandardMaterial } from "three";
import { useElapsedTime } from "@/effects/pulse";

/** Width and depth of each square in world units. */
const SQUARE_SIZE = 1;

/** Height (thickness) of each square in world units. */
const SQUARE_HEIGHT = 0.2;

/** Offset to center the 8×8 board at the scene origin. */
const BOARD_OFFSET = 3.5;

interface BoardSquareProps {
  position: Position;
  square: Square;
  onSelect: (position: Position) => void;
}

/**
 * A single 3D square on the chess board.
 * Pure/presentational: its mark state and click behavior are
 * provided by the parent (BoardScene), so it has no board-context
 * dependency and can be memoized.
 */
export const BoardSquare = memo(function BoardSquare({
  position,
  square,
  onSelect,
}: BoardSquareProps) {
  const { theme } = useTheme();
  const getElapsed = useElapsedTime();
  const materialRef = useRef<MeshStandardMaterial>(null);

  // Light if file + rank parity is even, same formula as chess board pattern
  const isLight = (fileOf(position) + rankOf(position)) % 2 === 0;

  // Derive colour from theme, square shade, and the square's own mark
  const color = getSquareColor(isLight, square.mark, theme);

  // Read live config from MARK_EFFECTS (mutable — debug panel can change it)
  const config = MARK_EFFECTS[square.mark];

  // Animate emissive intensity per-frame from live config
  useFrame(() => {
    if (!materialRef.current || config.pulseSpeed === 0) return;

    const t = getElapsed();
    const oscillation = Math.sin(t * config.pulseSpeed * Math.PI * 2);
    materialRef.current.emissiveIntensity =
      config.emissiveIntensity * (1 + oscillation * config.pulseRange);
  });

  // Offset by 3.5 to center the board at the scene origin
  const x = fileOf(position) - BOARD_OFFSET;
  const y = 0;
  const z = rankOf(position) - BOARD_OFFSET;

  return (
    <group position={[x, y, z]}>
      {/* The square surface */}
      <mesh
        onClick={(e) => {
          e.stopPropagation();
          onSelect(position);
        }}
      >
        <boxGeometry args={[SQUARE_SIZE, SQUARE_HEIGHT, SQUARE_SIZE]} />
        <meshStandardMaterial
          ref={materialRef}
          color={color}
          emissive={config.emissive}
          emissiveIntensity={config.emissiveIntensity}
          toneMapped={false}
        />
      </mesh>

      {/* Mark-specific visual effect (cone / dot / ring) */}
      <MarkEffect mark={square.mark} />
    </group>
  );
});
