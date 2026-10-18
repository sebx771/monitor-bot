export const WorkerStatus = {
    OPERATIONAL: 'operational',
    DEGRADED: 'degraded',
    ERROR: 'error',
    PAUSED: 'paused',
} as const;

export type WorkerHealthStatus = typeof WorkerStatus[keyof typeof WorkerStatus];
export type WorkerStatusValue = WorkerHealthStatus;

export interface MonitoredService {
    target: string;
    scope?: string;
    state: string;
    message?: string;
}

// Alias for compatibility
export type Service = MonitoredService;

export interface Worker {
    id: string;
    name: string;
    status: WorkerHealthStatus;
    last_run: string;
    services: MonitoredService[];
    message?: string;
}

// Allow WorkerStatus to be used both as an enum-like object and as a type
export type WorkerStatus = Worker;

export interface TriggerResponse {
    message: string;
    worker: string;
}

export interface ErrorResponse {
    error: string;
}