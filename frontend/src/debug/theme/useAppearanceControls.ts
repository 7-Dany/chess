import { useTheme } from "@/theme/useTheme";
import { button, useControls } from "leva";
import { useEffect, useRef } from "react";

export function useAppearanceDebugControls() {
  const { toggleTheme } = useTheme();

  const toggleThemeRef = useRef(toggleTheme);
  useEffect(() => {
    toggleThemeRef.current = toggleTheme;
  });

  useControls("Appearance", () => ({
    "Toggle Theme": button(() => toggleThemeRef.current()),
  }));
}
