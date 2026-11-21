import type { Worker } from "../../models/Worker";
import WorkerCard from "../WorkerCard";
import styles from "./WorkerGrid.module.css";

export interface WorkerGridProps {
  workers: Worker[];
  onWorkerClick?: (worker: Worker) => void;
  activeConfirmationWorkerId?: string | null;
  onRequestConfirmation?: (workerId: string) => void;
  onClearConfirmation?: () => void;
  onTrigger?: (workerId: string) => Promise<void> | void;
  onClearFilters?: () => void;
  hasActiveFilters?: boolean;
}

export default function WorkerGrid({
  workers,
  onWorkerClick,
  activeConfirmationWorkerId,
  onRequestConfirmation,
  onClearConfirmation,
  onTrigger,
  onClearFilters,
  hasActiveFilters,
}: WorkerGridProps) {
  if (workers.length === 0) {
    return (
      <div className={styles.emptyState}>
        <svg className={styles.emptyIcon} viewBox="0 0 24 24" fill="none" stroke="currentColor">
          <circle cx="11" cy="11" r="8" />
          <path d="m21 21-4.3-4.3" />
          <path d="M8 11h6" />
        </svg>
        <p className={styles.emptyTitle}>No workers found</p>
        <p className={styles.emptyDescription}>
          No workers match the current criteria or no services are registered.
        </p>
        {hasActiveFilters && onClearFilters && (
          <button type="button" className={styles.clearFiltersBtn} onClick={onClearFilters}>
            Clear active filters
          </button>
        )}
      </div>
    );
  }

  return (
    <section className={styles.gridContainer} aria-label="Workers grid">
      <div className={styles.grid}>
        {workers.map((worker) => (
          <div key={worker.id} className={styles.cardWrapper}>
            <WorkerCard
              worker={worker}
              onSelect={onWorkerClick}
              activeConfirmationWorkerId={activeConfirmationWorkerId}
              onRequestConfirmation={onRequestConfirmation}
              onClearConfirmation={onClearConfirmation}
              onTrigger={onTrigger}
            />
          </div>
        ))}
      </div>
    </section>
  );
}
