import type { TriggerResponse, Worker } from "../models/Worker";
import { API_BASE_URL } from "../core/config";
import { mapApiWorkers } from "./workerMapper";

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

function readErrorMessage(data: Record<string, unknown>, fallback: string): string {
  return typeof data.error === "string" ? data.error : fallback;
}

/**
 * Sends a manual trigger request for a specific worker.
 * Endpoint: POST /api/v1/workers/:id/trigger
 * /this function must be refactorize to complete with the new auth in trigger endpoint
 */
export async function triggerWorker(workerId: string): Promise<TriggerResponse> {
  const response = await fetch(`${API_BASE_URL}/workers/${encodeURIComponent(workerId)}/trigger`, {
    method: "POST",
    credentials: "include",
    headers: {
      Accept: "application/json",
    },
  });

  let data: Record<string, unknown> = {};
  try {
    const responseData: unknown = await response.json();
    if (isRecord(responseData)) {
      data = responseData;
    }
  } catch {
    // Error responses may not include a JSON body.
  }

  if (response.status === 401) {
    throw new Error("Only Admin can do this.");
  }

  if (response.status === 202) {
    return data as unknown as TriggerResponse;
  }

  if (response.status === 409) {
    throw new Error(readErrorMessage(data, "Worker is already busy or a trigger is already pending."));
  }

  if (response.status === 404) {
    throw new Error(readErrorMessage(data, `Worker "${workerId}" not found or inactive.`));
  }

  throw new Error(readErrorMessage(data, `Failed to trigger worker (${String(response.status)})`));
}

/**
 * Fetches the current status of all workers.
 * Endpoint: GET /api/v1/status
 */
export async function fetchWorkerStatus(signal?: AbortSignal): Promise<Worker[]> {
  const response = await fetch(`${API_BASE_URL}/status`, {
    method: "GET",
    credentials: "include",
    headers: {
      Accept: "application/json",
    },
    signal,
  });

  if (!response.ok) {
    throw new Error(`Error ${String(response.status)}: ${response.statusText}`);
  }

  const data: unknown = await response.json();
  return mapApiWorkers(data);
}
