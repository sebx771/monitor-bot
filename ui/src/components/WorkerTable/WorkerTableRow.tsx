import type { MouseEvent, KeyboardEvent } from "react";
import { Blobatar } from "@blobatar/react";
import { idle, mad, sleepy, thinking } from "blobatar/expression";
import { WorkerStatus, type Worker } from "../../models/Worker";
import { useWorkerTrigger } from "../WorkerCard/hooks/useWorkerTrigger";
import styles from "./WorkerTable.module.css";

interface WorkerTableRowProps {
  worker: Worker;
  onSelect?: (worker: Worker) => void;
  activeConfirmationWorkerId?: string | null;
  onRequestConfirmation?: (workerId: string) => void;
  onClearConfirmation?: () => void;
  onTrigger?: (workerId: string) => Promise<void> | void;
}

const expressions = {
  [WorkerStatus.OPERATIONAL]: idle,
  [WorkerStatus.DEGRADED]: thinking,
  [WorkerStatus.ERROR]: mad,
  [WorkerStatus.PAUSED]: sleepy,
} as const;

const ringClassMap = {
  [WorkerStatus.OPERATIONAL]: styles.ringOperational,
  [WorkerStatus.DEGRADED]: styles.ringDegraded,
  [WorkerStatus.ERROR]: styles.ringError,
  [WorkerStatus.PAUSED]: styles.ringPaused,
};

const statusClassMap = {
  [WorkerStatus.OPERATIONAL]: styles.statusOperational,
  [WorkerStatus.DEGRADED]: styles.statusDegraded,
  [WorkerStatus.ERROR]: styles.statusError,
  [WorkerStatus.PAUSED]: styles.statusPaused,
};

export function WorkerTableRow({
  worker,
  onSelect,
  activeConfirmationWorkerId,
  onRequestConfirmation,
  onClearConfirmation,
  onTrigger,
}: WorkerTableRowProps) {
  const {
    isOperational,
    isExecuting,
    isConfirming,
    handleTriggerClick,
    tooltipText,
  } = useWorkerTrigger({
    worker,
    onCustomTrigger: onTrigger,
    activeConfirmationWorkerId,
    onRequestConfirmation,
    onClearConfirmation,
  });

  const handleRowClick = () => { onSelect?.(worker); };

  const handleKeyDown = (e: KeyboardEvent<HTMLTableRowElement>) => {
    if (e.key === "Enter" || e.key === " ") {
      e.preventDefault();
      onSelect?.(worker);
    }
  };

  const onBtnClick = (e: MouseEvent<HTMLButtonElement>) => {
    e.stopPropagation();
    handleTriggerClick();
  };

  const ringStyle = ringClassMap[worker.status] || styles.ringPaused;
  const statusStyle = statusClassMap[worker.status] || styles.statusPaused;
  const formattedLastRun = worker.last_run ? new Date(worker.last_run).toLocaleTimeString() : "Never";

  return (
    <tr
      className={styles.row}
      onClick={handleRowClick}
      onKeyDown={handleKeyDown}
      tabIndex={0}
      role="row"
      aria-label={`Worker ${worker.name}, status ${worker.status}`}
    >
      <td>
        <div className={styles.workerCell}>
          <div className={`${styles.microAvatarWrapper} ${ringStyle}`}>
            <Blobatar name={worker.name} size={28} expression={expressions[worker.status]} />
          </div>
          <div>
            <div className={styles.workerName}>{worker.name}</div>
            <div className={styles.workerId}>{worker.id}</div>
          </div>
        </div>
      </td>

      <td>
        <span className={`${styles.statusPill} ${statusStyle}`}>
          <span className={styles.pillDot} />
          <span>{worker.status}</span>
        </span>
      </td>

      <td>
        <div className={styles.servicesList}>
          {worker.services.slice(0, 3).map((s) => (
            <span key={s.target} className={styles.serviceBadge}>
              {s.target}
            </span>
          ))}
          {worker.services.length > 3 && (
            <span className={styles.serviceBadge}>+{worker.services.length - 3}</span>
          )}
        </div>
      </td>

      <td>
        <span className={styles.timestamp}>{formattedLastRun}</span>
      </td>

      <td className={styles.actionsCell}>
        <button
          type="button"
          className={`${styles.triggerBtn} ${isConfirming ? styles.triggerBtnConfirming : ""}`}
          onClick={onBtnClick}
          disabled={!isOperational || isExecuting}
          title={tooltipText}
          aria-label={`Trigger ${worker.name}`}
        >
          {isExecuting ? "Queuing..." : isConfirming ? "Confirm? (3s)" : "⚡ Run"}
        </button>
      </td>
    </tr>
  );
}

export default WorkerTableRow;
