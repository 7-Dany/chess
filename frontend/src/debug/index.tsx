import { useState } from "react";
import { Leva } from "leva";
import { Bug } from "lucide-react";
import { Button } from "@/components/ui/button";
import { DebugPanel } from "./DebugPanel";
import { useTheme } from "../theme/useTheme";
import { LEVA_THEME } from "../theme/leva";

export function DebugTools() {
  const [open, setOpen] = useState(false);

  const { theme } = useTheme();

  if (import.meta.env.PROD) return null;

  return (
    <>
      {/* key={theme}: force remount so Leva picks up the new theme object */}
      <Leva key={theme} hidden={!open} theme={LEVA_THEME[theme]} />
      <DebugPanel />
      <Button
        onClick={() => setOpen((o) => !o)}
        variant="outline"
        size="icon"
        className="fixed bottom-4 right-4 z-50"
      >
        <Bug className="size-4" />
      </Button>
    </>
  );
}
