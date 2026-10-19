import { Blobatar } from "@blobatar/react";
import "blobatar/motion.css";
import { idle, mad, sleepy, thinking } from "blobatar/expression";
import { WorkerStatus, type Worker } from "../../models/Worker";
import { useState } from "react";

import styles from "./WorkerCard.module.css";

interface WorkerCardProps {
    worker: Worker;
    size?: number;
}

const expressions = {
    [WorkerStatus.OPERATIONAL]: idle,
    [WorkerStatus.DEGRADED]: thinking,
    [WorkerStatus.ERROR]: mad,
    [WorkerStatus.PAUSED]: sleepy,
} as const;

export default function WorkerCard({ worker, size = 160 }: WorkerCardProps) {
    const [IsHovered , SetIsHover] = useState(false);
    return (
        <div
            className={styles.card}
            onMouseEnter={() => { SetIsHover(true); }}
            onMouseLeave={() => { SetIsHover(false); }}
        >
            <Blobatar
                name={worker.name}
                size={size}
                animate={IsHovered ? "always" : undefined}
                expression={expressions[worker.status]}
            />
            <h2 className={styles.name}>{worker.name}</h2>
            <p className={styles.status}>{worker.status}</p>
            <ul className={styles.servicesList}>
                {worker.services.map((service, index) => (
                    <li key={index} className={styles.serviceItem}>
                        {service.target} ({service.state})
                    </li>
                ))}
            </ul>
        </div>
    );
}
