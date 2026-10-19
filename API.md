# Guía de Consumo de la API — Monitor Bot

> Documentación técnica detallada para la integración y consumo de la API REST de **Monitor Bot**.

---

## 📌 Índice

1. [Información General](#-información-general)
2. [Políticas de CORS y Seguridad](#-políticas-de-cors-y-seguridad)
3. [Códigos de Respuesta HTTP](#-códigos-de-respuesta-http)
4. [Modelos y Esquemas de Datos](#-modelos-y-esquemas-de-datos)
5. [Endpoints de la API](#-endpoints-de-la-api)
   - [Healthcheck](#1-healthcheck)
   - [Consultar Estado General](#2-consultar-estado-general)
   - [Disparar Worker Manualmente](#3-disparar-worker-manualmente)
6. [Ejemplos de Consumo](#-ejemplos-de-consumo)
   - [cURL](#curl)
   - [TypeScript / JavaScript (Fetch API)](#typescript--javascript-fetch-api)
   - [React Hook (Frontend Dashboard)](#react-hook-frontend-dashboard)
   - [Python (requests)](#python-requests)
7. [Buenas Prácticas y Manejo de Errores](#-buenas-prácticas-y-manejo-de-errores)

---

## 🌐 Información General

- **Base URL predeterminada:** `http://localhost:8080`
- **Prefijo de versión:** `/api/v1`
- **Formato de intercambio:** `application/json`
- **Modo del servidor:** Configurable mediante la variable de entorno `GIN_MODE` (`release` por defecto, o `debug`).
- **Autenticación:** Actualmente la API opera en red local/privada sin cabeceras de autorización obligatorias.

---

## 🛡️ Políticas de CORS y Seguridad

El router ([internal/api/router.go](file:///c:/Users/Usuario/OneDrive/Desktop/aternos-bot/internal/api/router.go)) incluye middleware CORS preconfigurado para aplicaciones web:

| Parámetro | Configuración |
|---|---|
| **Orígenes Permitidos (`AllowOrigins`)** | `http://localhost:5173` (Vite dev server)<br/>`http://localhost:3000` |
| **Métodos Permitidos (`AllowMethods`)** | `GET`, `POST`, `OPTIONS` |
| **Cabeceras Permitidas (`AllowHeaders`)** | `Origin`, `Content-Type`, `Accept` |
| **Cabeceras Expuestas (`ExposeHeaders`)** | `Content-Length` |
| **Credenciales (`AllowCredentials`)** | `true` |
| **Tiempo de caché Preflight (`MaxAge`)** | `12 horas` |

> [!NOTE]
> Si despliegas un frontend en un dominio o puerto distinto, debes añadirlo a la lista `AllowOrigins` en `internal/api/router.go`.

---

## 🚦 Códigos de Respuesta HTTP

| Código | Significado | Contexto en Monitor Bot |
|---|---|---|
| `200 OK` | Petición exitosa | Devuelve el healthcheck o la lista de workers activos con sus servicios. |
| `202 Accepted` | Petición aceptada | Se ha recibido y encolado la señal de trigger manual para el worker. |
| `404 Not Found` | Recurso no encontrado | El worker especificado en la URL no existe o no se encuentra activo. |
| `409 Conflict` | Conflicto de concurrencia | El worker ya se encuentra ocupado ejecutando una tarea o ya tiene un trigger pendiente. |
| `500 Internal Server Error` | Error interno del servidor | Fallo inesperado en el servidor HTTP. |

---

## 📦 Modelos y Esquemas de Datos

### 1. `WorkerStatus`

Representa el estado actual de un worker en el almacén en memoria ([internal/status/workers_status.go](file:///c:/Users/Usuario/OneDrive/Desktop/aternos-bot/internal/status/workers_status.go)):

| Campo | Tipo | Requerido | Descripción |
|---|---|---|---|
| `id` | `string` | Sí | Identificador único del worker (`"aiven"`, `"supabase"`). |
| `name` | `string` | Sí | Nombre para mostrar (`"Aiven"`, `"Supabase"`). |
| `status` | `string` | Sí | Estado del worker (`"operational"`, `"degraded"`, `"error"`, `"paused"`). |
| `last_run` | `string` (ISO 8601) | Sí | Timestamp UTC de la última ejecución (ej. `"2026-10-17T16:15:30Z"`). |
| `services` | `MonitoredService[]` | Sí | Colección de servicios o proyectos supervisados por este worker. |
| `message` | `string` | No | Mensaje descriptivo opcional (motivo de degradación o error). |

#### Valores de `Status`:
- `"operational"`: Todos los servicios o proyectos fueron verificados sin errores.
- `"degraded"`: Operación parcial (algunos proyectos o comprobaciones fallaron, pero no todos).
- `"error"`: Fallo crítico en el ciclo de trabajo (todos los proyectos fallaron o error en la API externa).
- `"paused"`: El worker está pausado.

---

### 2. `MonitoredService`

Detalle de un recurso individual dentro de un worker:

| Campo | Tipo | Requerido | Descripción |
|---|---|---|---|
| `target` | `string` | Sí | Nombre o referencia del recurso (ej. `"mysql-13396a35"`, `"my-project-ref"`). |
| `scope` | `string` | No | Ámbito o identificador del proyecto padre (ej. `"terranova1244-6338"`). |
| `state` | `string` | Sí | Estado reportado por el proveedor (ej. `"RUNNING"`, `"POWEROFF"`, `"ACTIVE_HEALTHY"`, `"INACTIVE"`). |
| `message` | `string` | No | Detalle o mensaje de error asociado al servicio si aplica. |

---

### 3. Definiciones TypeScript (para Frontend o Node.js)

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

## 🔌 Endpoints de la API

### 1. Healthcheck

Verifica si la API REST está en línea y lista para recibir peticiones.

- **Método:** `GET`
- **Ruta:** `/api/v1/`
- **Autenticación:** No requerida
- **Headers:**
  - `Accept: application/json`

#### Respuesta Exitosa (`200 OK`):
```json
{
  "status": "ok"
}
```

---

### 2. Consultar Estado General

Obtiene la lista de todos los workers registrados y el estado de sus servicios monitoreados.

- **Método:** `GET`
- **Ruta:** `/api/v1/status`
- **Autenticación:** No requerida
- **Headers:**
  - `Accept: application/json`

#### Respuesta Exitosa (`200 OK`):
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
> Si no hay workers activos (por ejemplo, si todos los módulos están deshabilitados en el `.env`), el endpoint retornará una lista vacía `[]`.

---

### 3. Disparar Worker Manualmente

Fuerza la ejecución inmediata del ciclo de monitoreo/acción del worker indicado, sin esperar a que transcurra el intervalo de chequeo programado.

- **Método:** `POST`
- **Ruta:** `/api/v1/workers/:id/trigger`
- **Parámetros de URL:**
  - `:id` (string, obligatorio): ID del worker a ejecutar. IDs válidos en la configuración actual:
    - `aiven`
    - `supabase`
- **Headers:**
  - `Content-Type: application/json` (opcional)
- **Cuerpo de la petición:** Vacío (no requiere payload).

#### Respuestas Posibles:

##### 🟢 `202 Accepted` — Trigger encolado correctamente
```json
{
  "message": "trigger sent",
  "worker": "aiven"
}
```

##### 🟡 `409 Conflict` — Worker ocupado
Ocurre si el worker ya se encuentra en ejecución o si ya existe una solicitud de trigger pendiente en su canal interno.
```json
{
  "error": "worker is busy or a trigger is already pending"
}
```

##### 🔴 `404 Not Found` — Worker inexistente o inactivo
Ocurre si el identificador `:id` no coincide con ningún worker registrado en el `Store` (o si dicho módulo no está habilitado en `.env`).
```json
{
  "error": "worker not found or not active"
}
```

---

## 💻 Ejemplos de Consumo

### cURL

#### 1. Comprobar salud de la API
```bash
curl -X GET http://localhost:8080/api/v1/
```

#### 2. Consultar el estado de todos los workers
```bash
curl -X GET http://localhost:8080/api/v1/status \
  -H "Accept: application/json"
```

#### 3. Ejecutar manualmente el worker de Aiven
```bash
curl -X POST http://localhost:8080/api/v1/workers/aiven/trigger \
  -H "Accept: application/json"
```

#### 4. Ejecutar manualmente el worker de Supabase
```bash
curl -X POST http://localhost:8080/api/v1/workers/supabase/trigger \
  -H "Accept: application/json"
```

---

### TypeScript / JavaScript (Fetch API)

Cliente reutilizable con manejo tipado de respuestas:

```typescript
const API_BASE_URL = 'http://localhost:8080/api/v1';

/**
 * Obtiene el estado actual de los workers
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
 * Solicita el disparo manual de un worker por ID
 */
export async function triggerWorker(workerId: string): Promise<TriggerResponse> {
  const response = await fetch(`${API_BASE_URL}/workers/${workerId}/trigger`, {
    method: 'POST',
    headers: {
      'Accept': 'application/json',
    },
  });

  const data = await response.json();

  if (response.status === 202) {
    return data as TriggerResponse;
  }

  if (response.status === 409) {
    throw new Error('El worker ya está ocupado o tiene una ejecución pendiente.');
  }

  if (response.status === 404) {
    throw new Error(`El worker "${workerId}" no existe o está inactivo.`);
  }

  throw new Error(data.error || 'Error inesperado al disparar el worker.');
}
```

---

### React Hook (Frontend Dashboard)

Hook personalizado para React (`ui/src/hooks/useWorkerStatus.ts`) con sondeo periódico (polling) y trigger:

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
    } catch (err: any) {
      setError(err.message || 'Error al conectar con la API');
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
    const res = await fetch(`${API_URL}/workers/${workerId}/trigger`, {
      method: 'POST',
    });
    const result = await res.json();

    if (!res.ok) {
      throw new Error(result.error || 'No se pudo activar el worker');
    }

    // Refrescar estado poco después para reflejar cambios
    setTimeout(refreshStatus, 1000);
    return result;
  };

  return { workers, loading, error, refreshStatus, trigger };
}
```

---

### Python (requests)

Script de automatización en Python:

```python
import requests
import time

BASE_URL = "http://localhost:8080/api/v1"

def check_status():
    url = f"{BASE_URL}/status"
    try:
        response = requests.get(url, timeout=5)
        response.raise_for_status()
        workers = response.json()
        print(f"Workers registrados: {len(workers)}")
        for w in workers:
            print(f"- [{w['status'].upper()}] {w['name']} (ID: {w['id']}) - Última ejecución: {w['last_run']}")
            for s in w.get('services', []):
                scope_str = f" ({s['scope']})" if 'scope' in s else ""
                print(f"    * {s['target']}{scope_str}: {s['state']}")
    except requests.RequestException as e:
        print(f"Error al consultar la API: {e}")

def trigger_worker(worker_id: str):
    url = f"{BASE_URL}/workers/{worker_id}/trigger"
    try:
        response = requests.post(url, timeout=5)
        if response.status_code == 202:
            print(f"Worker '{worker_id}' disparado exitosamente: {response.json()}")
        elif response.status_code == 409:
            print(f"Aviso: El worker '{worker_id}' ya está ocupado o pendiente.")
        elif response.status_code == 404:
            print(f"Error: El worker '{worker_id}' no existe o está inhabilitado.")
        else:
            print(f"Error HTTP {response.status_code}: {response.text}")
    except requests.RequestException as e:
        print(f"Error de conexión: {e}")

if __name__ == "__main__":
    check_status()
    print("\nDisparando worker de Aiven...")
    trigger_worker("aiven")
```

---

## 💡 Buenas Prácticas y Manejo de Errores

1. **Gestión de Concurrencia (`409 Conflict`):**
   - El canal de disparo manual (`triggerCh`) tiene una capacidad de búfer de 1.
   - Si un usuario pulsa repetidamente un botón en la UI, capture el error `409` y muestre una notificación amigable como *"El worker ya se encuentra en ejecución"* en lugar de un error crítico.
   - Deshabilite temporalmente el botón de disparo durante 3-5 segundos tras recibir un código `202`.

2. **Intervalo de Sondeo (Polling):**
   - El endpoint `/api/v1/status` lee directamente del almacén en memoria protegido por `sync.RWMutex`, por lo que es extremadamente rápido y ligero.
   - Para dashboards en tiempo real, un intervalo de sondeo entre **3 y 10 segundos** es ideal.

3. **Workers Habilitados dinámicamente:**
   - Si un worker (`AIVEN_ENABLED=false` o `SUPABASE_ENABLED=false`) está deshabilitado en `.env`, no se registrará en el `Store`. Por tanto, el endpoint `/api/v1/workers/:id/trigger` responderá `404 Not Found`.
   - Consulte siempre primero `/api/v1/status` para conocer qué IDs de workers están activos antes de presentar opciones de trigger en su interfaz.
