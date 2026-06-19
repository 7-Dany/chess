import { BoardProvider } from "./context/BoardProvider";
import { BoardScene } from "./scene/BoardScene";
import { ThemeProvider } from "./theme/ThemeProvider";
import { DebugTools } from "./debug";

function App() {
  return (
    <ThemeProvider>
      <BoardProvider>
        <div className="relative w-screen h-screen">
          <BoardScene />
          <DebugTools />
        </div>
      </BoardProvider>
    </ThemeProvider>
  );
}

export default App;
