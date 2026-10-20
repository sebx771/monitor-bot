import type { TriggerResponse, Worker } from "../models/Worker";
import { API_BASE_URL } from "../core/config";

/**
 * Sends a manual trigger request for a specific worker.
 * Endpoint: POST /api/v1/workers/:id/trigger
 */
export async function triggerWorker(workerId: string): Promise<TriggerResponse> {
  const response = await fetch(`${API_BASE_URL}/workers/${encodeURIComponent(workerId)}/trigger`, {
    method: "POST",
    headers: {
      Accept: "application/json",
    },
  });

  let data: Record<string, unknown> = {};
  try {
    data = await response.json();
  } catch {
    // If response is not JSON
  }

  if (response.status === 202) {
    return data as unknown as TriggerResponse;
  }

  if (response.status === 409) {
    throw new Error((data.error as string) || "Worker is already busy or a trigger is already pending.");
  }

  if (response.status === 404) {
    throw new Error((data.error as string) || `Worker "${workerId}" not found or inactive.`);
  }

  throw new Error((data.error as string) || `Failed to trigger worker (${response.status})`);
}

/**
 * Fetches the current status of all workers.
 * Endpoint: GET /api/v1/status
 */
export async function fetchWorkerStatus(): Promise<Worker[]> {
  const response = await fetch(`${API_BASE_URL}/status`, {
    method: "GET",
    headers: {
      Accept: "application/json",
    },
  });

  if (!response.ok) {
    throw new Error(`Error ${response.status}: ${response.statusText}`);
  }

  return response.json();
}
