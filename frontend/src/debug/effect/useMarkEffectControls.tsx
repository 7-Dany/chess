import { useEffect, useRef } from "react";
import { useControls } from "leva";
import { MARK_EFFECTS } from "../../theme/effects";
import type { Mark } from "../../chess/board";

/**
 * Registers "Mark Effects" sections in the debug panel.
 *
 * Only the folder matching the currently selected square's mark is
 * visible — all fields in other folders get `render: () => false`.
 * This keeps the panel clean and focused on what's actually on screen.
 *
 * All three useControls calls always run (Rules of Hooks), but their
 * fields are shown/hidden via Leva's per-field `render` function.
 *
 * Because Leva's factory function `() => ({...})` only runs once,
 * `render` closures would capture a stale `activeMark`. We avoid this
 * with a ref so the closures always read the latest value.
 *
 * However, Leva only re-evaluates `render` functions when its store
 * updates for that folder. A hidden `_activeMark` trigger field is
 * included in each subfolder; updating it via `set()` when the mark
 * changes forces Leva to re-evaluate every `render` in that folder.
 *
 * Values are written directly to MARK_EFFECTS (mutable), so changes
 * appear instantly in the 3D scene — no re-render needed.
 *
 * @param activeMark The mark of the currently selected square, or
 *                   "none" if nothing is selected. Effect folders are
 *                   only shown when their mark matches.
 */
export function useMarkEffectControls(activeMark: Mark) {
  // Refs so Leva's one-time closures always see current values.
  // All ref writes happen inside useEffect (React 19 forbids
  // writing ref.current during render).
  const activeMarkRef = useRef(activeMark);
  const setSelectedRef = useRef<(v: Record<string, unknown>) => void>(() => {});
  const setLegalMoveRef = useRef<(v: Record<string, unknown>) => void>(
    () => {},
  );
  const setCaptureRef = useRef<(v: Record<string, unknown>) => void>(() => {});

  // Render functions read from the ref — called by Leva, not during React render
  const renderSelected = () => activeMarkRef.current === "selected";
  const renderLegalMove = () => activeMarkRef.current === "legal-move";
  const renderCapture = () => activeMarkRef.current === "capture";

  /* ── selected (GlowCone) ─────────────────────────────────── */

  const [, setSelected] = useControls("Mark Selected", () => ({
    pulseSpeed: {
      label: "Speed (Hz)",
      value: MARK_EFFECTS.selected.pulseSpeed,
      min: 0,
      max: 10,
      step: 0.1,
      onChange: (v: number) => {
        MARK_EFFECTS.selected.pulseSpeed = v;
      },
      render: renderSelected,
      transient: true,
    },
    pulseRange: {
      label: "Range",
      value: MARK_EFFECTS.selected.pulseRange,
      min: 0,
      max: 1,
      step: 0.01,
      onChange: (v: number) => {
        MARK_EFFECTS.selected.pulseRange = v;
      },
      render: renderSelected,
      transient: true,
    },
    emissiveIntensity: {
      label: "Emissive",
      value: MARK_EFFECTS.selected.emissiveIntensity,
      min: 0,
      max: 2,
      step: 0.01,
      onChange: (v: number) => {
        MARK_EFFECTS.selected.emissiveIntensity = v;
      },
      render: renderSelected,
      transient: true,
    },
    coneHeight: {
      label: "Cone Height",
      value: MARK_EFFECTS.selected.cone!.height,
      min: 0.1,
      max: 5,
      step: 0.05,
      onChange: (v: number) => {
        MARK_EFFECTS.selected.cone!.height = v;
      },
      render: renderSelected,
      transient: true,
    },
    coneBaseRadius: {
      label: "Cone Base R",
      value: MARK_EFFECTS.selected.cone!.baseRadius,
      min: 0,
      max: 2,
      step: 0.01,
      onChange: (v: number) => {
        MARK_EFFECTS.selected.cone!.baseRadius = v;
      },
      render: renderSelected,
      transient: true,
    },
    coneTopRadius: {
      label: "Cone Top R",
      value: MARK_EFFECTS.selected.cone!.topRadius,
      min: 0,
      max: 2,
      step: 0.01,
      onChange: (v: number) => {
        MARK_EFFECTS.selected.cone!.topRadius = v;
      },
      render: renderSelected,
      transient: true,
    },
    coneBaseOpacity: {
      label: "Cone Base Opacity",
      value: MARK_EFFECTS.selected.cone!.baseOpacity,
      min: 0,
      max: 1,
      step: 0.01,
      onChange: (v: number) => {
        MARK_EFFECTS.selected.cone!.baseOpacity = v;
      },
      render: renderSelected,
      transient: true,
    },
  }));

  /* ── legal-move (SurfaceDot) ──────────────────────────────── */

  const [, setLegalMove] = useControls("Mark Legal-move", () => ({
    pulseSpeed: {
      label: "Speed (Hz)",
      value: MARK_EFFECTS["legal-move"].pulseSpeed,
      min: 0,
      max: 10,
      step: 0.1,
      onChange: (v: number) => {
        MARK_EFFECTS["legal-move"].pulseSpeed = v;
      },
      render: renderLegalMove,
      transient: true,
    },
    pulseRange: {
      label: "Range",
      value: MARK_EFFECTS["legal-move"].pulseRange,
      min: 0,
      max: 1,
      step: 0.01,
      onChange: (v: number) => {
        MARK_EFFECTS["legal-move"].pulseRange = v;
      },
      render: renderLegalMove,
      transient: true,
    },
    emissiveIntensity: {
      label: "Emissive",
      value: MARK_EFFECTS["legal-move"].emissiveIntensity,
      min: 0,
      max: 2,
      step: 0.01,
      onChange: (v: number) => {
        MARK_EFFECTS["legal-move"].emissiveIntensity = v;
      },
      render: renderLegalMove,
      transient: true,
    },
    dotRadius: {
      label: "Dot Radius",
      value: MARK_EFFECTS["legal-move"].dot!.radius,
      min: 0.01,
      max: 0.5,
      step: 0.01,
      onChange: (v: number) => {
        MARK_EFFECTS["legal-move"].dot!.radius = v;
      },
      render: renderLegalMove,
      transient: true,
    },
    dotThickness: {
      label: "Dot Thickness",
      value: MARK_EFFECTS["legal-move"].dot!.thickness,
      min: 0.01,
      max: 0.2,
      step: 0.005,
      onChange: (v: number) => {
        MARK_EFFECTS["legal-move"].dot!.thickness = v;
      },
      render: renderLegalMove,
      transient: true,
    },
    dotOpacity: {
      label: "Dot Opacity",
      value: MARK_EFFECTS["legal-move"].dot!.opacity,
      min: 0,
      max: 1,
      step: 0.01,
      onChange: (v: number) => {
        MARK_EFFECTS["legal-move"].dot!.opacity = v;
      },
      render: renderLegalMove,
      transient: true,
    },
    dotEmissive: {
      label: "Dot Emissive",
      value: MARK_EFFECTS["legal-move"].dot!.emissiveIntensity,
      min: 0,
      max: 2,
      step: 0.01,
      onChange: (v: number) => {
        MARK_EFFECTS["legal-move"].dot!.emissiveIntensity = v;
      },
      render: renderLegalMove,
      transient: true,
    },
  }));

  /* ── capture (EdgeRing) ───────────────────────────────────── */

  const [, setCapture] = useControls("Mark Capture", () => ({
    pulseSpeed: {
      label: "Speed (Hz)",
      value: MARK_EFFECTS.capture.pulseSpeed,
      min: 0,
      max: 10,
      step: 0.1,
      onChange: (v: number) => {
        MARK_EFFECTS.capture.pulseSpeed = v;
      },
      render: renderCapture,
      transient: true,
    },
    pulseRange: {
      label: "Range",
      value: MARK_EFFECTS.capture.pulseRange,
      min: 0,
      max: 1,
      step: 0.01,
      onChange: (v: number) => {
        MARK_EFFECTS.capture.pulseRange = v;
      },
      render: renderCapture,
      transient: true,
    },
    emissiveIntensity: {
      label: "Emissive",
      value: MARK_EFFECTS.capture.emissiveIntensity,
      min: 0,
      max: 2,
      step: 0.01,
      onChange: (v: number) => {
        MARK_EFFECTS.capture.emissiveIntensity = v;
      },
      render: renderCapture,
      transient: true,
    },
    ringRadius: {
      label: "Ring Radius",
      value: MARK_EFFECTS.capture.ring!.radius,
      min: 0.1,
      max: 0.7,
      step: 0.01,
      onChange: (v: number) => {
        MARK_EFFECTS.capture.ring!.radius = v;
      },
      render: renderCapture,
      transient: true,
    },
    ringThickness: {
      label: "Ring Thickness",
      value: MARK_EFFECTS.capture.ring!.thickness,
      min: 0.005,
      max: 0.1,
      step: 0.005,
      onChange: (v: number) => {
        MARK_EFFECTS.capture.ring!.thickness = v;
      },
      render: renderCapture,
      transient: true,
    },
    ringOpacity: {
      label: "Ring Opacity",
      value: MARK_EFFECTS.capture.ring!.opacity,
      min: 0,
      max: 1,
      step: 0.01,
      onChange: (v: number) => {
        MARK_EFFECTS.capture.ring!.opacity = v;
      },
      render: renderCapture,
      transient: true,
    },
    ringEmissive: {
      label: "Ring Emissive",
      value: MARK_EFFECTS.capture.ring!.emissiveIntensity,
      min: 0,
      max: 2,
      step: 0.01,
      onChange: (v: number) => {
        MARK_EFFECTS.capture.ring!.emissiveIntensity = v;
      },
      render: renderCapture,
      transient: true,
    },
  }));

  /* ── Update all refs (MUST be declared BEFORE the sync effect) ── */

  useEffect(() => {
    activeMarkRef.current = activeMark;
    setSelectedRef.current = setSelected;
    setLegalMoveRef.current = setLegalMove;
    setCaptureRef.current = setCapture;
  });

  /* ── Sync MARK_EFFECTS → Leva when activeMark changes ────── */

  useEffect(() => {
    setSelectedRef.current({
      pulseSpeed: MARK_EFFECTS.selected.pulseSpeed,
      pulseRange: MARK_EFFECTS.selected.pulseRange,
      emissiveIntensity: MARK_EFFECTS.selected.emissiveIntensity,
      coneHeight: MARK_EFFECTS.selected.cone!.height,
      coneBaseRadius: MARK_EFFECTS.selected.cone!.baseRadius,
      coneTopRadius: MARK_EFFECTS.selected.cone!.topRadius,
      coneBaseOpacity: MARK_EFFECTS.selected.cone!.baseOpacity,
    });
    setLegalMoveRef.current({
      pulseSpeed: MARK_EFFECTS["legal-move"].pulseSpeed,
      pulseRange: MARK_EFFECTS["legal-move"].pulseRange,
      emissiveIntensity: MARK_EFFECTS["legal-move"].emissiveIntensity,
      dotRadius: MARK_EFFECTS["legal-move"].dot!.radius,
      dotThickness: MARK_EFFECTS["legal-move"].dot!.thickness,
      dotOpacity: MARK_EFFECTS["legal-move"].dot!.opacity,
      dotEmissive: MARK_EFFECTS["legal-move"].dot!.emissiveIntensity,
    });
    setCaptureRef.current({
      pulseSpeed: MARK_EFFECTS.capture.pulseSpeed,
      pulseRange: MARK_EFFECTS.capture.pulseRange,
      emissiveIntensity: MARK_EFFECTS.capture.emissiveIntensity,
      ringRadius: MARK_EFFECTS.capture.ring!.radius,
      ringThickness: MARK_EFFECTS.capture.ring!.thickness,
      ringOpacity: MARK_EFFECTS.capture.ring!.opacity,
      ringEmissive: MARK_EFFECTS.capture.ring!.emissiveIntensity,
    });
  }, [activeMark]);
}
