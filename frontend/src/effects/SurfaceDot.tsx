import { memo, useRef } from "react";
import { useFrame } from "@react-three/fiber";
import * as THREE from "three";
import { MARK_EFFECTS } from "../theme/effects";
import { livePulse, useElapsedTime } from "./pulse";

/** Square surface thickness — dot sits above this. */
const SQUARE_HEIGHT = 0.2;

const TORUS_RADIAL_SEGMENTS = 12;
const TORUS_TUBULAR_SEGMENTS = 24;

export const SurfaceDot = memo(function SurfaceDot() {
  const meshRef = useRef<THREE.Mesh>(null);
  const materialRef = useRef<THREE.MeshStandardMaterial>(null);
  const getElapsed = useElapsedTime();

  /** Previous geometry params — to detect changes and rebuild geometry. */
  const prevGeo = useRef({ radius: 0, thickness: 0 });

  useFrame(() => {
    const mesh = meshRef.current;
    const mat = materialRef.current;
    if (!mesh || !mat) return;

    const config = MARK_EFFECTS["legal-move"];
    const dot = config.dot!;
    const pulse = livePulse("legal-move", getElapsed());

    // — Material properties (every frame) —
    mat.emissiveIntensity = dot.emissiveIntensity * pulse;
    mat.opacity = dot.opacity * pulse;
    mat.color.set(dot.color);
    mat.emissive.set(dot.color);

    // — Position (every frame — thickness might change) —
    mesh.position.y = SQUARE_HEIGHT / 2 + dot.thickness / 2;

    // — Geometry rebuild (only when shape params change) —
    if (
      prevGeo.current.radius !== dot.radius ||
      prevGeo.current.thickness !== dot.thickness
    ) {
      const oldGeo = mesh.geometry as THREE.TorusGeometry;
      oldGeo.dispose();
      mesh.geometry = new THREE.TorusGeometry(
        dot.radius,
        dot.thickness,
        TORUS_RADIAL_SEGMENTS,
        TORUS_TUBULAR_SEGMENTS,
      );
      prevGeo.current = { radius: dot.radius, thickness: dot.thickness };
    }
  });

  // Initial values for the first render
  const config = MARK_EFFECTS["legal-move"].dot!;
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
