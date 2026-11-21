import type { Worker } from "../../models/Worker";
import { WorkerTableRow } from "./WorkerTableRow";
import styles from "./WorkerTable.module.css";

export interface WorkerTableProps {
  workers: Worker[];
  onSelect?: (worker: Worker) => void;
  activeConfirmationWorkerId?: string | null;
  onRequestConfirmation?: (workerId: string) => void;
  onClearConfirmation?: () => void;
  onTrigger?: (workerId: string) => Promise<void> | void;
}

export function WorkerTable({
  workers,
  onSelect,
  activeConfirmationWorkerId,
  onRequestConfirmation,
  onClearConfirmation,
  onTrigger,
}: WorkerTableProps) {
  return (
    <section className={styles.tableWrapper} aria-label="Workers compact table">
      <table className={styles.table}>
        <thead>
          <tr>
            <th scope="col">Worker</th>
            <th scope="col">Status</th>
            <th scope="col">Targets</th>
            <th scope="col">Last Run</th>
            <th scope="col" className={styles.actionsCell}>Actions</th>
          </tr>
        </thead>
        <tbody>
          {workers.map((worker) => (
            <WorkerTableRow
              key={worker.id}
              worker={worker}
              onSelect={onSelect}
              activeConfirmationWorkerId={activeConfirmationWorkerId}
              onRequestConfirmation={onRequestConfirmation}
              onClearConfirmation={onClearConfirmation}
              onTrigger={onTrigger}
            />
          ))}
        </tbody>
      </table>
    </section>
  );
}

export default WorkerTable;
