// src/App.tsx
import WorkerCard from "./components/WorkerCard";
import { WorkerMock } from "./mock/Worker";

export default function App() {
  return (
    <div style={{ padding: "2rem", display: "flex", gap: "1rem", flexWrap: "wrap" }}>
      {WorkerMock.map((worker) => (
        <WorkerCard key={worker.id} worker={worker} />
      ))}
    </div>
    
  );
}