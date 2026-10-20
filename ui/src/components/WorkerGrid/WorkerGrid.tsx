import type { Worker } from "../../models/Worker";
import WorkerCard from "../WorkerCard";
import styles from "./WorkerGrid.module.css";

export interface WorkerGridProps {
  workers: Worker[];
  onWorkerClick?: (worker: Worker) => void;
}

export default function WorkerGrid({ workers, onWorkerClick }: WorkerGridProps) {
  if (workers.length === 0) {
    return (
      <div className={styles.emptyState}>
        <svg
          className={styles.emptyIcon}
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
        >
          <circle cx="11" cy="11" r="8" />
          <path d="m21 21-4.3-4.3" />
          <path d="M8 11h6" />
        </svg>
        <p className={styles.emptyTitle}>No workers found</p>
        <p className={styles.emptyDescription}>
          No workers match the current search criteria or no services are registered.
        </p>
      </div>
    );
  }

  return (
    <section className={styles.gridContainer} aria-label="Workers grid">
      <div className={styles.grid}>
        {workers.map((worker) => (
          <div key={worker.id} className={styles.cardWrapper}>
            <WorkerCard worker={worker} onSelect={onWorkerClick} />
          </div>
        ))}
      </div>
    </section>
  );
}
