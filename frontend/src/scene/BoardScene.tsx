import { Canvas } from "@react-three/fiber";
import { OrbitControls } from "@react-three/drei";
import { BoardSquare } from "./BoardSquare";
import { useBoard } from "../context/useBoard";
import type { Position } from "../chess/position";

/** Starting position per side. Flips when turn changes. */
const CAMERA_POSITIONS = {
  WHITE: [0, 8, 10] as const,
  BLACK: [0, 8, -10] as const,
} as const;

/** Field of view in degrees. */
const CAMERA_FOV = 50;

/** Zoom limits in world units. */
const CAMERA_MIN_DISTANCE = 5;
const CAMERA_MAX_DISTANCE = 18;

/** Prevents camera from going below the board surface. */
const CAMERA_MAX_POLAR_ANGLE = Math.PI / 2.2;

/** Lights configuration */
const AMBIENT_LIGHT_INTENSITY = 0.8;
const DIRECTIONAL_LIGHT_INTENSITY = 1.6;
const DIRECTIONAL_LIGHT_POSITION = [5, 10, 5] as const;

export function BoardScene() {
  const { state, selectPosition } = useBoard();

  // TODO: get currentTurn from useBoard and switch camera position based on it
  const cameraPosition = CAMERA_POSITIONS.BLACK;

  return (
    <Canvas camera={{ position: cameraPosition, fov: CAMERA_FOV }}>
      {/* Lights */}
      <ambientLight intensity={AMBIENT_LIGHT_INTENSITY} />
      <directionalLight
        position={DIRECTIONAL_LIGHT_POSITION}
        intensity={DIRECTIONAL_LIGHT_INTENSITY}
      />

      {/* Controls */}
      <OrbitControls
        minDistance={CAMERA_MIN_DISTANCE}
        maxDistance={CAMERA_MAX_DISTANCE}
        maxPolarAngle={CAMERA_MAX_POLAR_ANGLE}
      />

      {/* Board squares — iterate the board array directly */}
      {state.board.map((square, position) => (
        <BoardSquare
          key={position}
          position={position as Position}
          square={square}
          onSelect={selectPosition}
        />
      ))}
    </Canvas>
  );
}
