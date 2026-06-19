/** Glow cone settings (used by "selected" mark). */
export interface ConeConfig {
  /** Base colour of the cone glow (CSS hex string, e.g. "#38bdf8"). */
  color: string;
  /** Total height of the cone in world units. Higher = taller beacon. */
  height: number;
  /**
   * Radius at the cone's base (bottom), in world units.
   * This is the "half-side" of the square cross-section — the engine
   * multiplies by √2 to get the circumradius CylinderGeometry needs.
   */
  baseRadius: number;
  /**
   * Radius at the cone's top, in world units.
   * Same half-side → circumradius conversion as baseRadius.
   * Set larger than baseRadius for an outward-flaring shape;
   * set smaller for a tapered/pointed shape.
   */
  topRadius: number;
  /**
   * Master opacity of the entire cone (0–1).
   * This is multiplied by the pulse value and a per-fragment fade
   * (bottom ramp, top fade, edge softness) to produce the final alpha.
   * Increase for a more opaque/solid glow; decrease for a ghostly look.
   */
  baseOpacity: number;
}

/** Surface dot settings (used by "legal-move" mark). */
export interface DotConfig {
  /** Colour of the dot (CSS hex string, e.g. "#a3e635"). */
  color: string;
  /** Radius of the dot ring (distance from centre to the tube centre), in world units. */
  radius: number;
  /** Thickness of the dot ring's tube, in world units. */
  thickness: number;
  /**
   * Master opacity of the dot (0–1).
   * Multiplied by the pulse value each frame.
   * Higher = more visible; lower = more transparent.
   */
  opacity: number;
  /**
   * How strongly the dot emits its own light (0–2).
   * This is the MeshStandardMaterial emissiveIntensity, scaled by pulse.
   * Higher values make the dot glow more vividly even in shadow.
   */
  emissiveIntensity: number;
}

/** Edge ring settings (used by "capture" mark). */
export interface RingConfig {
  /** Colour of the ring (CSS hex string, e.g. "#f87171"). */
  color: string;
  /** Radius of the ring (distance from centre to the tube centre), in world units. */
  radius: number;
  /** Thickness of the ring's tube, in world units. */
  thickness: number;
  /**
   * Master opacity of the ring (0–1).
   * Multiplied by the pulse value each frame.
   * Higher = more visible; lower = more transparent.
   */
  opacity: number;
  /**
   * How strongly the ring emits its own light (0–2).
   * This is the MeshStandardMaterial emissiveIntensity, scaled by pulse.
   * Higher values make the ring glow more vividly even in shadow.
   */
  emissiveIntensity: number;
}

/**
 * Per-mark effect configuration.
 *
 * Each mark has:
 *   - Surface emissive (colour + intensity on the square mesh)
 *   - Pulse (speed in Hz, range as fraction of base)
 *   - One optional 3D overlay (cone / dot / ring)
 */
export interface MarkEffectConfig {
  /**
   * Emissive colour applied to the square's surface mesh (CSS hex string).
   * This tints the square itself — separate from the 3D overlay colour.
   * Use a dark/zero value ("#000000") for no surface tint.
   */
  emissive: string;
  /**
   * Emissive intensity on the square surface (0–2).
   * Controls how brightly the square glows with the emissive colour.
   * Multiplied by the pulse value each frame for breathing animation.
   */
  emissiveIntensity: number;
  /**
   * Pulse oscillation speed in Hz (cycles per second).
   * 0 = no animation (steady state).
   * Typical range: 1–3 for a gentle breathing effect.
   */
  pulseSpeed: number;
  /**
   * How much the effect "breathes" — the ±percentage around the base value.
   * 0   = steady, no animation at all.
   * 0.2 = subtle ±20% breathing (fades to 80%, brightens to 120%).
   * 0.5 = strong ±50% breathing (fades to 50%, brightens to 150%).
   */
  pulseRange: number;
  /** Glow cone overlay (only present on "selected" mark). */
  cone?: ConeConfig;
  /** Surface dot overlay (only present on "legal-move" mark). */
  dot?: DotConfig;
  /** Edge ring overlay (only present on "capture" mark). */
  ring?: RingConfig;
}
