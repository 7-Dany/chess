import type { Leva } from "leva";

type LevaCustomTheme = NonNullable<Parameters<typeof Leva>[0]["theme"]>;

export const LEVA_THEME: Record<"dark" | "light", LevaCustomTheme> = {
  dark: {
    colors: {
      elevation1: "#1a1a1a",
      elevation2: "#242424",
      elevation3: "#2e2e2e",
      accent1: "#e8e8e8",
      accent2: "#c8c8c8",
      accent3: "#aaaaaa",
      highlight1: "#f9f9f9",
      highlight2: "#aaaaaa",
      highlight3: "#111111",
      vivid1: "#000000",
      folderTextColor: "#f9f9f9",
      folderWidgetColor: "#f9f9f9",
    },
    fontWeights: { label: "500", folder: "600", button: "700" },
    radii: { xs: "3px", sm: "4px", lg: "8px" },
    fonts: {
      mono: '"JetBrains Mono", "Fira Code", monospace',
      sans: '"Inter", system-ui, sans-serif',
    },
    fontSizes: { root: "11px" },
    sizes: {
      rootWidth: "300px",
      controlWidth: "150px",
      titleBarHeight: "38px",
    },
    space: { sm: "6px", md: "10px", rowGap: "6px", colGap: "6px" },
  },
  light: {
    colors: {
      elevation1: "#f5f5f5",
      elevation2: "#ffffff",
      elevation3: "#e8e8e8",
      accent1: "#282828",
      accent2: "#444444",
      accent3: "#666666",
      highlight1: "#1a1a1a",
      highlight2: "#555555",
      highlight3: "#f0f0f0",
      vivid1: "#ffffff",
      folderTextColor: "#1a1a1a",
      folderWidgetColor: "#1a1a1a",
    },
    fontWeights: { label: "500", folder: "600", button: "700" },
    radii: { xs: "3px", sm: "4px", lg: "8px" },
    fonts: {
      mono: '"JetBrains Mono", "Fira Code", monospace',
      sans: '"Inter", system-ui, sans-serif',
    },
    fontSizes: { root: "11px" },
    sizes: {
      rootWidth: "300px",
      controlWidth: "150px",
      titleBarHeight: "38px",
    },
    space: { sm: "6px", md: "10px", rowGap: "6px", colGap: "6px" },
  },
};
