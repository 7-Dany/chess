import { memo, useRef } from "react";
import { useFrame } from "@react-three/fiber";
import * as THREE from "three";
import { MARK_EFFECTS } from "../theme/effects";
import { livePulse, useElapsedTime } from "./pulse";

/** Square surface thickness — ring sits above this. */
const SQUARE_HEIGHT = 0.2;

const TORUS_RADIAL_SEGMENTS = 12;
const TORUS_TUBULAR_SEGMENTS = 32;

export const EdgeRing = memo(function EdgeRing() {
  const meshRef = useRef<THREE.Mesh>(null);
  const materialRef = useRef<THREE.MeshStandardMaterial>(null);
  const getElapsed = useElapsedTime();

  /** Previous geometry params — to detect changes and rebuild geometry. */
  const prevGeo = useRef({ radius: 0, thickness: 0 });

  useFrame(() => {
    const mesh = meshRef.current;
    const mat = materialRef.current;
    if (!mesh || !mat) return;

    const config = MARK_EFFECTS.capture;
    const ring = config.ring!;
    const pulse = livePulse("capture", getElapsed());

    // — Material properties (every frame) —
    mat.emissiveIntensity = ring.emissiveIntensity * pulse;
    mat.opacity = ring.opacity * pulse;
    mat.color.set(ring.color);
    mat.emissive.set(ring.color);

    // — Position (every frame — thickness might change) —
    mesh.position.y = SQUARE_HEIGHT / 2 + ring.thickness / 2;

    // — Geometry rebuild (only when shape params change) —
    if (
      prevGeo.current.radius !== ring.radius ||
      prevGeo.current.thickness !== ring.thickness
    ) {
      const oldGeo = mesh.geometry as THREE.TorusGeometry;
      oldGeo.dispose();
      mesh.geometry = new THREE.TorusGeometry(
        ring.radius,
        ring.thickness,
        TORUS_RADIAL_SEGMENTS,
        TORUS_TUBULAR_SEGMENTS,
      );
      prevGeo.current = { radius: ring.radius, thickness: ring.thickness };
    }
  });

  // Initial values for the first render
  const config = MARK_EFFECTS.capture.ring!;
  const y = SQUARE_HEIGHT / 2 + config.thickness / 2;

  // Seed prevGeo so we don't rebuild on the very first frame
  prevGeo.current = { radius: config.radius, thickness: config.thickness };

  return (
    <mesh ref={meshRef} position={[0, y, 0]} rotation={[Math.PI / 2, 0, 0]}>
      <torusGeometry
        args={[
          config.radius,
          config.thickness,
          TORUS_RADIAL_SEGMENTS,
          TORUS_TUBULAR_SEGMENTS,
        ]}
      />
      <meshStandardMaterial
        ref={materialRef}
        color={config.color}
        emissive={config.color}
        emissiveIntensity={config.emissiveIntensity}
        transparent
        opacity={config.opacity}
        depthWrite={false}
        toneMapped={false}
      />
    </mesh>
  );
});
