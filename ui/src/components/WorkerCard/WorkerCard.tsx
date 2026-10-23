import { useState, type MouseEvent, type KeyboardEvent } from "react";
import { Blobatar } from "@blobatar/react";
import "blobatar/motion.css";
import { idle, mad, sleepy, thinking } from "blobatar/expression";
import { WorkerStatus, type Worker } from "../../models/Worker";
import { useWorkerTrigger } from "./hooks/useWorkerTrigger";
import { useGaze } from "@blobatar/react/gaze";

import styles from "./WorkerCard.module.css";
import "blobatar/gaze.css";


export interface WorkerCardProps {
  worker: Worker;
  size?: number;
  onSelect?: (worker: Worker) => void;
  onTrigger?: (workerId: string) => Promise<void> | void;
}

const expressions = {
  [WorkerStatus.OPERATIONAL]: idle,
  [WorkerStatus.DEGRADED]: thinking,
  [WorkerStatus.ERROR]: mad,
  [WorkerStatus.PAUSED]: sleepy,
} as const;

const statusClassMap = {
  [WorkerStatus.OPERATIONAL]: styles.statusOperational,
  [WorkerStatus.DEGRADED]: styles.statusDegraded,
  [WorkerStatus.ERROR]: styles.statusError,
  [WorkerStatus.PAUSED]: styles.statusPaused,
};

export default function WorkerCard({
  worker,
  size = 160,
  onSelect,
  onTrigger,
}: WorkerCardProps) {
  const [isHovered, setIsHovered] = useState(false);
  const {
    isOperational,
    isTriggering,
    notice,
    executeTrigger,
    tooltipText,
  } = useWorkerTrigger({ worker, onCustomTrigger: onTrigger });
  const { ref } = useGaze({
  travel: 10,
  lookAt: "pointer",
});
  const handleCardClick = () => {
    onSelect?.(worker);
  };

  const handleKeyDown = (e: KeyboardEvent<HTMLDivElement>) => {
    if (e.key === "Enter" || e.key === " ") {
      e.preventDefault();
      onSelect?.(worker);
    }
  };

  const handleDetailsClick = (e: MouseEvent<HTMLButtonElement>) => {
    e.stopPropagation();
    onSelect?.(worker);
  };

  const handleTriggerClick = (e: MouseEvent<HTMLButtonElement>) => {
    e.stopPropagation();
    void executeTrigger();
  };

  const statusStyle = statusClassMap[worker.status] || styles.statusPaused;

  return (
    <div
      className={styles.card}
      onClick={handleCardClick}
      onKeyDown={handleKeyDown}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
      role="button"
      tabIndex={0}
      aria-label={`Worker ${worker.name}, status: ${worker.status}. Click to view details.`}
    >
      <div className={styles.blobWrapper}>
        <Blobatar
          name={worker.name}
          ref={ref}
          size={size}
          animate={isHovered ? "always" : undefined}
          expression={expressions[worker.status] || idle}
        />
      </div>

      <h2 className={styles.name}>{worker.name}</h2>

      <div className={`${styles.statusBadge} ${statusStyle}`}>
        <span className={styles.statusDot} />
        <span>{worker.status}</span>
      </div>

      {notice && (
        <div
          className={`${styles.notice} ${
            notice.type === "success"
              ? styles.noticeSuccess
              : notice.type === "conflict"
              ? styles.noticeConflict
              : styles.noticeError
          }`}
          role="alert"
        >
          {notice.text}
        </div>
      )}

      <div className={styles.actions}>
        <button
          type="button"
          className={`${styles.button} ${styles.triggerBtn}`}
          onClick={handleTriggerClick}
          disabled={!isOperational || isTriggering}
          title={tooltipText}
          aria-label={`Manual trigger for ${worker.name}`}
        >
          {isTriggering ? (
            <>
              <span className={styles.spinner} />
              <span>Triggering...</span>
            </>
          ) : (
            <>
              <span>⚡</span>
              <span>Trigger</span>
            </>
          )}
        </button>

        <button
          type="button"
          className={`${styles.button} ${styles.detailsBtn}`}
          onClick={handleDetailsClick}
          title={`View details for ${worker.name}`}
          aria-label={`Open details modal for ${worker.name}`}
        >
          <span>Details</span>
        </button>
      </div>
    </div>
  );
}
