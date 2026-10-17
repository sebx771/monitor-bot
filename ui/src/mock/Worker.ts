import type { Worker}  from "../models/Worker";

export const WorkerMock: Worker[] = [
    {
        id: '1',
        name: 'Aiven Worker',
        status: 'running',
        services: [{id: '1', name: 'Aiven Worker', status: 'running'}]
    },
    {
        id: '2',
        name: 'Supabase Worker',
        status: 'stopped',
        services: [{id: '2', name: 'Supabase Worker', status: 'stopped'}]
    }
]