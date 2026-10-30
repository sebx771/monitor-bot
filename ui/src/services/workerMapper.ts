import {
  WorkerStatus,
  type MonitoredService,
  type Worker,
  type WorkerHealthStatus,
} from "../models/Worker";

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

function readRequiredString(
  record: Record<string, unknown>,
  key: string,
  path: string,
): string {
  const value = record[key];
  if (typeof value !== "string") {
    throw new Error(`Invalid API response: "${path}.${key}" must be a string.`);
  }
  return value;
}

function readOptionalString(
  record: Record<string, unknown>,
  key: string,
  path: string,
): string | undefined {
  const value = record[key];
  if (value === undefined) {
    return undefined;
  }
  if (typeof value !== "string") {
    throw new Error(`Invalid API response: "${path}.${key}" must be a string.`);
  }
  return value;
}

function isWorkerHealthStatus(value: string): value is WorkerHealthStatus {
  return Object.values(WorkerStatus).some((status) => status === value);
}

function mapApiService(value: unknown, path: string): MonitoredService {
  if (!isRecord(value)) {
    throw new Error(`Invalid API response: "${path}" must be an object.`);
  }

  return {
    target: readRequiredString(value, "target", path),
    state: readRequiredString(value, "state", path),
    scope: readOptionalString(value, "scope", path),
    message: readOptionalString(value, "message", path),
  };
}

function mapApiWorker(value: unknown, index: number): Worker {
  const path = `workers[${String(index)}]`;
  if (!isRecord(value)) {
    throw new Error(`Invalid API response: "${path}" must be an object.`);
  }

  const status = readRequiredString(value, "status", path);
  if (!isWorkerHealthStatus(status)) {
    throw new Error(`Invalid API response: "${path}.status" is unsupported.`);
  }

  const services = value.services;
  if (services !== null && services !== undefined && !Array.isArray(services)) {
    throw new Error(`Invalid API response: "${path}.services" must be an array or null.`);
  }

  return {
    id: readRequiredString(value, "id", path),
    name: readRequiredString(value, "name", path),
    status,
    last_run: readRequiredString(value, "last_run", path),
    services: (services ?? []).map((service: unknown, serviceIndex: number) =>
      mapApiService(service, `${path}.services[${String(serviceIndex)}]`),
    ),
    message: readOptionalString(value, "message", path),
  };
}

export function mapApiWorkers(value: unknown): Worker[] {
  if (!Array.isArray(value)) {
    throw new Error("Invalid API response: expected an array of workers.");
  }

  return value.map(mapApiWorker);
}
