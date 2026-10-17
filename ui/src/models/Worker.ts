
export const WorkerStatus= {
    RUNNING : 'running',
    STOPPED :'stopped',
    ERROR :'error',
    PAUSED :'paused'
} as const

export interface Worker {
    id: string;
    name: string;
    status: typeof WorkerStatus[keyof typeof WorkerStatus];
    services: Service[];
}


export interface Service {
    id: string;
    name: string;
    status: string;
}