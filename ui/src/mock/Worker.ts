import { WorkerStatus, type Worker } from "../models/Worker";

export const WorkerMock: Worker[] = [
    {
        id: 'aiven',
        name: 'Aiven',
        status: WorkerStatus.OPERATIONAL,
        last_run: '2026-10-17T16:00:05.123456Z',
        services: [
            {
                target: 'pg-production',
                scope: 'my-aiven-project',
                state: 'RUNNING',
            },
            {
                target: 'redis-cache',
                scope: 'my-aiven-project',
                state: 'POWEROFF',
            },
        ],
        message: '',
    },
    {
        id: 'supabase',
        name: 'Supabase',
        status: WorkerStatus.OPERATIONAL,
        last_run: '2026-10-17T16:00:10.987654Z',
        services: [
            {
                target: 'abcdefghijklmnop',
                state: 'ACTIVE_HEALTHY',
            },
        ],
        message: '',
    },
    {
        id: 'backup',
        name: 'Backup Worker',
        status: WorkerStatus.DEGRADED,
        last_run: '2026-10-17T15:30:00.000000Z',
        services: [
            {
                target: 's3-bucket-sync',
                state: 'FAILED',
                message: 'Connection timeout after 30s',
            },
        ],
        message: '1 of 1 services failed',
    },
    {
        id: 'crawler',
        name: 'Crawler Worker',
        status: WorkerStatus.PAUSED,
        last_run: '2026-10-17T14:00:00.000000Z',
        services: [
            {
                target: 'aternos-session',
                state: 'INACTIVE',
            },
        ],
        message: 'Worker is paused',
    },
];