import { useRef, useCallback } from "react";
import { useFrame } from "@react-three/fiber";
import type { Mark } from "../chess/board";
import { MARK_EFFECTS } from "../theme/effects";

/**
 * Returns a pulse multiplier for a mark at the given time.
 * Oscillates around 1.0 by ±pulseRange.
 *
 * - `pulseRange: 0`  → always returns 1.0 (steady)
 * - `pulseRange: 0.4` → swings between 0.6 and 1.4
 */
export function livePulse(mark: Mark, time: number): number {
  const config = MARK_EFFECTS[mark];
  if (config.pulseSpeed === 0) return 1.0;
  return (
    1.0 + Math.sin(time * config.pulseSpeed * Math.PI * 2) * config.pulseRange
  );
}

/**
 * Hook that tracks elapsed time using `useFrame`'s `delta` parameter,
 * avoiding the deprecated `THREE.Clock`.
 *
 * Returns a getter function that yields the accumulated time in seconds.
 * Call this inside a component that is already within the R3F `<Canvas>`.
 */
export function useElapsedTime(): () => number {
  const elapsed = useRef(0);

  useFrame((_, delta) => {
    elapsed.current += delta;
  });

  return useCallback(() => elapsed.current, []);
}
