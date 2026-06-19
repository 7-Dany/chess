import { memo } from "react";
import type { Mark } from "../chess/board";
import { GlowCone } from "./GlowCone";
import { SurfaceDot } from "./SurfaceDot";
import { EdgeRing } from "./EdgeRing";

interface MarkEffectProps {
  mark: Mark;
}

export const MarkEffect = memo(function MarkEffect({ mark }: MarkEffectProps) {
  if (mark === "none") return null;

  switch (mark) {
    case "selected":
      return <GlowCone />;
    case "legal-move":
      return <SurfaceDot />;
    case "capture":
      return <EdgeRing />;
  }
});
