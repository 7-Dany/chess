import { memo, useRef, useMemo } from "react";
import { useFrame } from "@react-three/fiber";
import * as THREE from "three";
import { MARK_EFFECTS } from "../theme/effects";
import { livePulse, useElapsedTime } from "./pulse";

/** Square surface thickness — cone starts above this. */
const SQUARE_HEIGHT = 0.2;

/**
 * 4-segment cylinder = square cross-section.
 * Matches the board squares.
 */
const CONE_SEGMENTS = 4;

/**
 * CylinderGeometry's "radius" is circumradius (center → vertex).
 * For a square with half-side `s`, circumradius = `s × √2`.
 */
const CIRCUMRADIUS_FACTOR = Math.SQRT2;

/** Rotation to align square faces with the board squares. */
const CONE_ALIGNMENT_ROTATION = Math.PI / 4;

const vertexShader = /* glsl */ `
  varying float vNormY;
  varying float vEdgeDist;

  void main() {
    vNormY = uv.y;
    vEdgeDist = abs(uv.x - 0.5) * 2.0;
    gl_Position = projectionMatrix * modelViewMatrix * vec4(position, 1.0);
  }
`;

const fragmentShader = /* glsl */ `
  uniform vec3 uColor;
  uniform float uBaseOpacity;
  uniform float uPulse;

  varying float vNormY;
  varying float vEdgeDist;

  void main() {
    float bottomRamp = smoothstep(0.0, 0.15, vNormY);
    float topFade = 1.0 - pow(vNormY, 1.4);
    float edgeSoftness = 1.0 - 0.3 * pow(vEdgeDist, 4.0);
    float alpha = uBaseOpacity * bottomRamp * topFade * edgeSoftness * uPulse;
    gl_FragColor = vec4(uColor, alpha);
  }
`;

export const GlowCone = memo(function GlowCone() {
  const meshRef = useRef<THREE.Mesh>(null);
  const materialRef = useRef<THREE.ShaderMaterial>(null);
  const getElapsed = useElapsedTime();

  /** Previous geometry params — to detect changes and rebuild geometry. */
  const prevGeo = useRef({ height: 0, baseCR: 0, topCR: 0 });

  // Initial uniform values for the <shaderMaterial> JSX prop.
  // After the material is created, useFrame updates the material's
  // ACTUAL uniforms via materialRef.current.uniforms — NOT via this
  // object. (R3F shallow-merges uniform values, so the material's
  // inner uniform objects are separate references.)
  const initialUniforms = useMemo(
    () => ({
      uColor: { value: new THREE.Color(MARK_EFFECTS.selected.cone!.color) },
      uBaseOpacity: { value: MARK_EFFECTS.selected.cone!.baseOpacity },
      uPulse: { value: 1.0 },
    }),
    [],
  );

  useFrame(() => {
    const mesh = meshRef.current;
    const mat = materialRef.current;
    if (!mesh || !mat) return;

    const cone = MARK_EFFECTS.selected.cone!;
    const pulse = livePulse("selected", getElapsed());

    // — Material uniforms via the material's actual uniform store —
    // MUST use mat.uniforms, NOT the closure initialUniforms.
    // R3F's applyProps creates separate uniform objects inside the
    // material; mutating the closure copy has no effect on rendering.
    mat.uniforms.uPulse.value = pulse;
    mat.uniforms.uBaseOpacity.value = cone.baseOpacity;
    mat.uniforms.uColor.value.set(cone.color);

    // — Position (every frame — height might change) —
    mesh.position.y = SQUARE_HEIGHT / 2 + cone.height / 2;

    // — Geometry rebuild (only when shape params change) —
    const baseCR = cone.baseRadius * CIRCUMRADIUS_FACTOR;
    const topCR = cone.topRadius * CIRCUMRADIUS_FACTOR;

    if (
      prevGeo.current.height !== cone.height ||
      prevGeo.current.baseCR !== baseCR ||
      prevGeo.current.topCR !== topCR
    ) {
      const oldGeo = mesh.geometry as THREE.CylinderGeometry;
      oldGeo.dispose();
      mesh.geometry = new THREE.CylinderGeometry(
        topCR,
        baseCR,
        cone.height,
        CONE_SEGMENTS,
        1,
        true,
      );
      prevGeo.current = { height: cone.height, baseCR, topCR };
    }
  });

  // Initial values for the first render
  const config = MARK_EFFECTS.selected.cone!;
  const y = SQUARE_HEIGHT / 2 + config.height / 2;
  const baseCR = config.baseRadius * CIRCUMRADIUS_FACTOR;
  const topCR = config.topRadius * CIRCUMRADIUS_FACTOR;

  // Seed prevGeo so we don't rebuild on the very first frame
  prevGeo.current = { height: config.height, baseCR, topCR };

  return (
    <mesh
      ref={meshRef}
      position={[0, y, 0]}
      rotation={[0, CONE_ALIGNMENT_ROTATION, 0]}
    >
      <cylinderGeometry
        args={[topCR, baseCR, config.height, CONE_SEGMENTS, 1, true]}
      />
      <shaderMaterial
        ref={materialRef}
        vertexShader={vertexShader}
        fragmentShader={fragmentShader}
        uniforms={initialUniforms}
        transparent
        blending={THREE.AdditiveBlending}
        depthWrite={false}
        side={THREE.DoubleSide}
      />
    </mesh>
  );
});
