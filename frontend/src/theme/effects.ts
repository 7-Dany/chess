import type { Mark } from "../chess/board";
import type { MarkEffectConfig } from "../effects/types";

/**
 * Per-mark visual effect settings.
 *
 * **Mutable** — the debug panel writes to these objects directly.
 * useFrame callbacks read them every frame, so changes are instant.
 */
export const MARK_EFFECTS: Record<Mark, MarkEffectConfig> = {
  none: {
    emissive: "#000000",
    emissiveIntensity: 0,
    pulseSpeed: 0,
    pulseRange: 0,
  },
  selected: {
    emissive: "#0ea5e9",
    emissiveIntensity: 0.35,
    pulseSpeed: 2,
    pulseRange: 0.4,
    cone: {
      color: "#38bdf8",
      height: 0.5,
      baseRadius: 0.5,
      topRadius: 0.7,
      baseOpacity: 0.2,
    },
  },
  "legal-move": {
    emissive: "#84cc16",
    emissiveIntensity: 0.15,
    pulseSpeed: 0.7,
    pulseRange: 0.4,
    dot: {
      color: "#a3e635",
      radius: 0.4,
      thickness: 0.01,
      opacity: 0.7,
      emissiveIntensity: 1.0,
    },
  },
  capture: {
    emissive: "#ef4444",
    emissiveIntensity: 0.35,
    pulseSpeed: 3,
    pulseRange: 0.5,
    ring: {
      color: "#f87171",
      radius: 0.38,
      thickness: 0.035,
      opacity: 0.8,
      emissiveIntensity: 0.6,
    },
  },
};
