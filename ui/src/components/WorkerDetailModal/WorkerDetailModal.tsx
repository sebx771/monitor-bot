import type { Worker } from "../../models/Worker";
import Modal from "../Modal";
import styles from "./WorkerDetailModal.module.css";

export interface WorkerDetailModalProps {
    worker: Worker;
    isOpen: boolean;
    onClose: () => void;
}

export default function WorkerDetailModal({ worker, isOpen, onClose }: WorkerDetailModalProps) {
    const formattedLastRun = worker.last_run
        ? new Date(worker.last_run).toLocaleString()
        : "No records";

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
            </div>
        </Modal>
    );
}