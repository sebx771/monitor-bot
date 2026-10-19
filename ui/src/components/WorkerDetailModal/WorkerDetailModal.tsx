import type { Worker } from "../../models/Worker";
import Modal from "../Modal";
import styles from "./WorkerDetailModal.module.css";

interface WorkerDetailModalProps {
    worker: Worker;
    isOpen: boolean;
    onClose: () => void;
}

export default function WorkerDetailModal({ worker, isOpen, onClose }: WorkerDetailModalProps) {
    return (
        <Modal isOpen={isOpen} onClose={onClose} title={worker.name}>
            <div className={styles.container}>
                <div className={styles.statusBadge}>
                    <span>Estado: <strong>{worker.status}</strong></span>
                </div>
                
                <h4 className={styles.subtitle}>Servicios Monitoreados</h4>
                <ul className={styles.servicesList}>
                    {worker.services.map((service, index) => (
                        <li key={index} className={styles.serviceItem}>
                            <span className={styles.serviceTarget}>{service.target}</span>
                            <span className={styles.serviceState}>{service.state}</span>
                        </li>
                    ))}
                </ul>
            </div>
        </Modal>
    );
}