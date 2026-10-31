import { useState } from "react";
import type { Worker } from "../../models/Worker";
import { triggerWorker } from "../../services/workerApi";
import Modal from "../Modal";
import styles from "./WorkerDetailModal.module.css";

export interface WorkerDetailModalProps {
  worker: Worker;
  isOpen: boolean;
  onClose: () => void;
  isAuthenticated: boolean;
}

export default function WorkerDetailModal({
  worker,
  isOpen,
  onClose,
  isAuthenticated,
}: WorkerDetailModalProps) {
  const [isTriggering, setIsTriggering] = useState(false);
  const [triggerError, setTriggerError] = useState<string | null>(null);
  const [triggerSuccess, setTriggerSuccess] = useState(false);

  const formattedLastRun = worker.last_run
    ? new Date(worker.last_run).toLocaleString()
    : "No records";

  const handleTrigger = async () => {
    if (isTriggering) return;
    setIsTriggering(true);
    setTriggerError(null);
    setTriggerSuccess(false);
    try {
      await triggerWorker(worker.id);
      setTriggerSuccess(true);
    } catch (err: unknown) {
      setTriggerError(err instanceof Error ? err.message : "Failed to trigger worker.");
    } finally {
      setIsTriggering(false);
    }
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose} title={`Worker Details: ${worker.name}`}>
      <div className={styles.container}>
        <div className={styles.statusBadge}>
          <span>Current Status:</span>
          <strong>{worker.status.toUpperCase()}</strong>
        </div>

        {worker.last_run && (
          <div className={styles.infoRow}>
            <span className={styles.infoLabel}>Last Run:</span>
            <span className={styles.infoValue}>{formattedLastRun}</span>
          </div>
        )}

        {worker.message && (
          <div className={styles.messageBox}>
            <span className={styles.infoLabel}>Status Information:</span>
            <p className={styles.messageText}>{worker.message}</p>
          </div>
        )}

        <h4 className={styles.subtitle}>
          Monitored Services ({worker.services.length})
        </h4>

        <ul className={styles.servicesList}>
          {worker.services.length === 0 ? (
            <li className={styles.emptyNotice}>No monitored services registered</li>
          ) : (
            worker.services.map((service, index) => (
              <li key={index} className={styles.serviceItem}>
                <div className={styles.serviceInfo}>
                  <span className={styles.serviceTarget}>{service.target}</span>
                  {service.scope && (
                    <span className={styles.serviceScope}>
                      Scope: {service.scope}
                    </span>
                  )}
                </div>
                <span className={styles.serviceState}>{service.state}</span>
              </li>
            ))
          )}
        </ul>

        <div className={styles.triggerSection}>
          {isAuthenticated ? (
            <>
              <button
                id={`trigger-btn-${worker.id}`}
                type="button"
                className={styles.triggerBtn}
                onClick={() => { void handleTrigger(); }}
                disabled={isTriggering}
                aria-busy={isTriggering}
              >
                {isTriggering ? "Triggering…" : "Trigger Worker"}
              </button>
              {triggerSuccess && (
                <p className={styles.triggerSuccess}>Trigger sent successfully.</p>
              )}
              {triggerError && (
                <p className={styles.triggerError}>{triggerError}</p>
              )}
            </>
          ) : (
            <p className={styles.authHint}>
              Log in to trigger this worker manually.
            </p>
          )}
        </div>
      </div>
    </Modal>
  );
}