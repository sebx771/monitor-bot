# API Consumption Guide — Monitor Bot

> Detailed technical documentation for integrating and consuming the **Monitor Bot** REST API.

---

## 📌 Table of Contents

1. [General Information](#-general-information)
2. [CORS Policies & Security](#-cors-policies--security)
3. [HTTP Status Codes](#-http-status-codes)
4. [Data Models & Schemas](#-data-models--schemas)
5. [API Endpoints](#-api-endpoints)
   - [Healthcheck](#1-healthcheck)
   - [Get General Status](#2-get-general-status)
   - [Manually Trigger a Worker](#3-manually-trigger-a-worker)
6. [Consumption Examples](#-consumption-examples)
   - [cURL](#curl)
   - [TypeScript / JavaScript (Fetch API)](#typescript--javascript-fetch-api)
   - [React Hook (Frontend Dashboard)](#react-hook-frontend-dashboard)
   - [Python (requests)](#python-requests)
7. [Best Practices & Error Handling](#-best-practices--error-handling)

---

## 🌐 General Information

- **Default Base URL:** `http://localhost:8080`
- **Version Prefix:** `/api/v1`
- **Exchange Format:** `application/json`
- **Server Mode:** Configurable via the `GIN_MODE` environment variable (`release` by default, or `debug`).
- **Authentication:** The healthcheck and status endpoints are public. Manually triggering a worker requires the `X-API-Key` header, whose value is configured with the `SECRET_KEY` environment variable.

---

## 🛡️ CORS Policies & Security

The router ([internal/api/router.go](file:///c:/Users/Usuario/OneDrive/Desktop/aternos-bot/internal/api/router.go)) includes preconfigured CORS middleware for web applications:

| Parameter | Configuration |
|---|---|
| **Allowed Origins (`AllowOrigins`)** | `http://localhost:5173` (Vite dev server)<br/>`http://localhost:3000` |
| **Allowed Methods (`AllowMethods`)** | `GET`, `POST`, `OPTIONS` |
| **Allowed Headers (`AllowHeaders`)** | `Origin`, `Content-Type`, `Accept`, `X-API-Key` |
| **Exposed Headers (`ExposeHeaders`)** | `Content-Length` |
| **Allow Credentials (`AllowCredentials`)** | `true` |
| **Preflight Cache Duration (`MaxAge`)** | `12 hours` |

> [!NOTE]
> If deploying a frontend on a different domain or port, you must add it to the `AllowOrigins` list in `internal/api/router.go`.

---

## 🚦 HTTP Status Codes

| Code | Meaning | Context in Monitor Bot |
|---|---|---|
| `200 OK` | Successful request | Returns healthcheck status or list of active workers and their monitored services. |
| `202 Accepted` | Accepted request | Manual trigger request received and queued for execution. |
| `401 Unauthorized` | Authentication failed | The trigger request is missing a valid `X-API-Key` header. |
| `404 Not Found` | Resource not found | The specified worker does not exist or is inactive. |
| `409 Conflict` | Concurrency conflict | The worker is already busy executing a cycle or already has a pending trigger. |
| `500 Internal Server Error` | Server error | Unexpected failure within the HTTP server. |

---

## 📦 Data Models & Schemas

### 1. `WorkerStatus`

Represents the current state of a worker stored in the in-memory registry ([internal/status/workers_status.go](file:///c:/Users/Usuario/OneDrive/Desktop/aternos-bot/internal/status/workers_status.go)):

| Field | Type | Required | Description |
|---|---|---|---|
| `id` | `string` | Yes | Unique worker identifier (`"aiven"`, `"supabase"`). |
| `name` | `string` | Yes | Display name (`"Aiven"`, `"Supabase"`). |
| `status` | `string` | Yes | Worker health status (`"operational"`, `"degraded"`, `"error"`, `"paused"`). |
| `last_run` | `string` (ISO 8601) | Yes | UTC timestamp of last cycle execution (e.g. `"2026-10-17T16:15:30Z"`). |
| `services` | `MonitoredService[]` | Yes | Collection of services or projects supervised by this worker. |
| `message` | `string` | No | Optional descriptive message (reason for error or degraded state). |

#### `status` Values:
- `"operational"`: All checked services/projects passed with no errors.
- `"degraded"`: Partial failure (some projects/services failed, but not all).
- `"error"`: Critical failure (all checks failed or external API error).
- `"paused"`: The worker is paused.

---

### 2. `MonitoredService`

Details for an individual resource tracked by a worker:

| Field | Type | Required | Description |
|---|---|---|---|
| `target` | `string` | Yes | Resource name or identifier (e.g. `"mysql-13396a35"`, `"my-project-ref"`). |
| `scope` | `string` | No | Scope or parent project identifier (e.g. `"terranova1244-6338"`). |
| `state` | `string` | Yes | State reported by provider (e.g. `"RUNNING"`, `"POWEROFF"`, `"ACTIVE_HEALTHY"`, `"INACTIVE"`). |
| `message` | `string` | No | Optional error or descriptive message. |

---

### 3. TypeScript Definitions (for Frontend / Node.js)

```typescript
export type WorkerHealthStatus = 'operational' | 'degraded' | 'error' | 'paused';

export interface MonitoredService {
  target: string;
  scope?: string;
  state: string;
  message?: string;
}

export interface WorkerStatus {
  id: string;
  name: string;
  status: WorkerHealthStatus;
  last_run: string;
  services: MonitoredService[];
  message?: string;
}

export interface TriggerResponse {
  message: string;
  worker: string;
}

export interface ErrorResponse {
  error: string;
}
```

---

## 🔌 API Endpoints

### 1. Healthcheck

Check whether the REST API is running and ready to handle requests.

- **Method:** `GET`
- **Route:** `/api/v1/`
- **Authentication:** None
- **Headers:**
  - `Accept: application/json`

#### Success Response (`200 OK`):
```json
{
  "status": "ok"
}
```

---

### 2. Get General Status

Fetches all registered workers along with their monitored service states.

- **Method:** `GET`
- **Route:** `/api/v1/status`
- **Authentication:** None
- **Headers:**
  - `Accept: application/json`

#### Success Response (`200 OK`):
```json
[
  {
    "id": "aiven",
    "name": "Aiven",
    "status": "operational",
    "last_run": "2026-10-17T16:00:05.123456Z",
    "services": [
      {
        "target": "pg-production",
        "scope": "my-aiven-project",
        "state": "RUNNING"
      },
      {
        "target": "redis-cache",
        "scope": "my-aiven-project",
        "state": "POWEROFF"
      }
    ],
    "message": ""
  },
  {
    "id": "supabase",
    "name": "Supabase",
    "status": "operational",
    "last_run": "2026-10-17T16:00:10.987654Z",
    "services": [
      {
        "target": "abcdefghijklmnop",
        "state": "ACTIVE_HEALTHY"
      }
    ],
    "message": ""
  }
]
```

> [!TIP]
> If no workers are active (for example, if all worker modules are disabled in `.env`), the endpoint returns an empty array `[]`.

---

### 3. Manually Trigger a Worker

Forces an immediate cycle for the specified worker without waiting for its scheduled ticker interval.

- **Method:** `POST`
- **Route:** `/api/v1/workers/:id/trigger`
- **Authentication:** Required. Set the `X-API-Key` header to the value of the server's `SECRET_KEY` environment variable.
- **URL Parameters:**
  - `:id` (string, required): ID of the worker to trigger. Valid IDs in the current setup:
    - `aiven`
    - `supabase`
- **Headers:**
  - `X-API-Key: <SECRET_KEY>` (required)
  - `Accept: application/json` (optional)
- **Request Body:** Empty (no payload required).

#### Possible Responses:

##### 🟢 `202 Accepted` — Trigger enqueued
```json
{
  "message": "trigger sent",
  "worker": "aiven"
}
```

##### 🟡 `409 Conflict` — Worker busy
Returned if the worker is currently executing a cycle or already has a pending trigger in its internal channel.
```json
{
  "error": "worker is busy or a trigger is already pending"
}
```

##### 🔒 `401 Unauthorized` — Missing or invalid API key
Returned if the `X-API-Key` header is missing, empty, or does not match the server's `SECRET_KEY`.
```json
{
  "error": "Access denied"
}
```

##### 🔴 `404 Not Found` — Worker not found or inactive
Returned if `:id` does not match any registered worker in the `Store` (or if its module is disabled in `.env`).
```json
{
  "error": "worker not found or not active"
}
```

---

## 💻 Consumption Examples

### cURL

#### 1. Check API Health
```bash
curl -X GET http://localhost:8080/api/v1/
```

#### 2. Get status for all workers
```bash
curl -X GET http://localhost:8080/api/v1/status \
  -H "Accept: application/json"
```

#### 3. Manually trigger the Aiven worker
```bash
curl -X POST http://localhost:8080/api/v1/workers/aiven/trigger \
  -H "Accept: application/json" \
  -H "X-API-Key: $SECRET_KEY"
```

#### 4. Manually trigger the Supabase worker
```bash
curl -X POST http://localhost:8080/api/v1/workers/supabase/trigger \
  -H "Accept: application/json" \
  -H "X-API-Key: $SECRET_KEY"
```

---

### TypeScript / JavaScript (Fetch API)

Reusable client with typed response handling for a trusted server-side runtime. Load `API_KEY` from server-side secret configuration; do not expose `SECRET_KEY` in a browser bundle or a `VITE_*` variable. Browser applications should call their own backend, which can securely add the API key when forwarding the trigger request.

```typescript
const API_BASE_URL = 'http://localhost:8080/api/v1';
const API_KEY = process.env.SECRET_KEY;

if (!API_KEY) {
  throw new Error('SECRET_KEY is not configured');
}

/**
 * Fetches the current status of all workers
 */
export async function fetchWorkerStatus(): Promise<WorkerStatus[]> {
  const response = await fetch(`${API_BASE_URL}/status`, {
    method: 'GET',
    headers: {
      'Accept': 'application/json',
    },
  });

  if (!response.ok) {
    throw new Error(`Error ${response.status}: ${response.statusText}`);
  }

  return response.json();
}

/**
 * Requests an immediate manual trigger for a worker by ID
 */
export async function triggerWorker(workerId: string): Promise<TriggerResponse> {
  const response = await fetch(`${API_BASE_URL}/workers/${workerId}/trigger`, {
    method: 'POST',
    headers: {
      'Accept': 'application/json',
      'X-API-Key': API_KEY,
    },
  });

  const data = await response.json();

  if (response.status === 202) {
    return data as TriggerResponse;
  }

  if (response.status === 409) {
    throw new Error('Worker is already busy or a trigger is already pending.');
  }

  if (response.status === 404) {
    throw new Error(`Worker "${workerId}" not found or inactive.`);
  }

  if (response.status === 401) {
    throw new Error('Unauthorized: missing or invalid API key.');
  }

  throw new Error(data.error || 'Unexpected error while triggering worker.');
}
```

---

### React Hook (Frontend Dashboard)

Custom React hook (`ui/src/hooks/useWorkerStatus.ts`) with polling and trigger dispatch:

The status request can be made directly from the browser. For triggers, send the request to your application's backend proxy and have it add `X-API-Key`; never bundle the server secret in frontend code.

```tsx
import { useState, useEffect, useCallback } from 'react';
import { WorkerStatus } from '../models/Worker';

const API_URL = 'http://localhost:8080/api/v1';

export function useWorkerStatus(pollIntervalMs: number = 5000) {
  const [workers, setWorkers] = useState<WorkerStatus[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  const refreshStatus = useCallback(async () => {
    try {
      const res = await fetch(`${API_URL}/status`);
      if (!res.ok) throw new Error(`HTTP error ${res.status}`);
      const data = await res.json();
      setWorkers(data);
      setError(null);
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Error connecting to API');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    refreshStatus();
    const interval = setInterval(refreshStatus, pollIntervalMs);
    return () => clearInterval(interval);
  }, [refreshStatus, pollIntervalMs]);

  const trigger = async (workerId: string) => {
    const res = await fetch(`/api/trigger/${workerId}`, {
      method: 'POST',
    });
    const result = await res.json();

    if (!res.ok) {
      throw new Error(result.error || 'Failed to trigger worker');
    }

    // Refresh status shortly after to reflect changes
    setTimeout(refreshStatus, 1000);
    return result;
  };

  return { workers, loading, error, refreshStatus, trigger };
}
```

---

### Python (requests)

Automation script in Python:

```python
import requests
import os
import time

BASE_URL = "http://localhost:8080/api/v1"
API_KEY = os.environ["SECRET_KEY"]

def check_status():
    url = f"{BASE_URL}/status"
    try:
        response = requests.get(url, timeout=5)
        response.raise_for_status()
        workers = response.json()
        print(f"Registered workers: {len(workers)}")
        for w in workers:
            print(f"- [{w['status'].upper()}] {w['name']} (ID: {w['id']}) - Last run: {w['last_run']}")
            for s in w.get('services', []):
                scope_str = f" ({s['scope']})" if 'scope' in s else ""
                print(f"    * {s['target']}{scope_str}: {s['state']}")
    except requests.RequestException as e:
        print(f"Error querying API: {e}")

def trigger_worker(worker_id: str):
    url = f"{BASE_URL}/workers/{worker_id}/trigger"
    try:
        response = requests.post(url, headers={"X-API-Key": API_KEY}, timeout=5)
        if response.status_code == 202:
            print(f"Worker '{worker_id}' triggered successfully: {response.json()}")
        elif response.status_code == 409:
            print(f"Notice: Worker '{worker_id}' is already busy or pending.")
        elif response.status_code == 404:
            print(f"Error: Worker '{worker_id}' does not exist or is inactive.")
        elif response.status_code == 401:
            print("Unauthorized: missing or invalid API key.")
        else:
            print(f"HTTP Error {response.status_code}: {response.text}")
    except requests.RequestException as e:
        print(f"Connection error: {e}")

if __name__ == "__main__":
    check_status()
    print("\nTriggering Aiven worker...")
    trigger_worker("aiven")
```

---

## 💡 Best Practices & Error Handling

1. **Concurrency Management (`409 Conflict`):**
   - The trigger channel (`triggerCh`) has a buffer capacity of 1.
   - If a user clicks a button rapidly in the UI, catch the `409` code and display a friendly notification like *"Worker is already running"* instead of a crash or red alert.
   - Temporarily disable the trigger button for 3-5 seconds after receiving a `202` response.

2. **Polling Interval:**
   - The `/api/v1/status` endpoint reads directly from an in-memory store protected by a `sync.RWMutex`, making it extremely fast and lightweight.
   - For real-time dashboards, a polling interval between **3 and 10 seconds** is optimal.

3. **Dynamically Enabled Workers:**
   - If a worker (`AIVEN_ENABLED=false` or `SUPABASE_ENABLED=false`) is disabled in `.env`, it will not register in the `Store`. As a result, the `/api/v1/workers/:id/trigger` endpoint will return `404 Not Found`.
   - Always query `/api/v1/status` first to know which worker IDs are currently active before displaying trigger options in your UI.
