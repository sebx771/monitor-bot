# Changelog

Todas las cambios notables de este proyecto se documentarán en este archivo.

El formato se basa en [Keep a Changelog](https://keepachangelog.com/es/1.1.0/),
y este proyecto adhiere a [Semantic Versioning](https://semver.org/lang/es/).

## [1.3.0] - 2026-08-18

### Added

- **Port `StateStorage`** (`internal/ports/state_storage.go`): interfaz que abstrae el almacenamiento remoto del archivo de estado de sesión con `DownloadState` y `UploadState`.
- **Adapter de GitHub Gist** (`internal/adapters/github_gist.go`): implementa `StateStorage` usando la API de GitHub para descargar y actualizar el estado de sesión en un Gist.
  - Autenticación Bearer, cabeceras `application/vnd.github+json`, timeout de 30 segundos y creación automática de directorios.
  - Pruebas unitarias en `internal/adapters/github_gist_test.go`.
- **Sincronización de sesión con GitHub Gist** (`internal/services/aternos_services.go`): `BotService` descarga el estado de sesión desde el Gist antes de iniciar el navegador y lo sube tras cada uso; reintenta el ciclo ante cookies inválidas.
- **Nuevas variables de entorno** (`internal/config.go`):
  - `HEADLESS`: controla el modo headless de Playwright (`internal/adapters/playwright.go`).
  - `GITHUB_TOKEN` y `GIST_ID`: token y ID del Gist para la sincronización de sesión.
- **Dockerfile multi-etapa**: build con `golang:1.24-bookworm` y runtime `debian:bookworm-slim`, con directorio `/app/storage` para el deploy en producción.

### Changed

- **Modo producción Aiven** (`cmd/main.go`): se comenta la inicialización del bot de Aternos, dejando únicamente activo el worker de Aiven para el deploy.
- **Formateo de código**: se aplica `go fmt` en `internal/aiven`, `internal/config.go`, `internal/minecraft` y `internal/ports`.
- **`.env.example`**: documenta las nuevas variables de entorno para Gist y headless.

## [1.2.0] - 2026-08-09

### Added

- **Módulo Aiven** (`internal/aiven`): cliente HTTP para interactuar con la API de Aiven y checker para revisar el estado de servicios.
  - `client.go`: cliente con métodos `GetServices` y `StartService`, manejo de autenticación Bearer/aivenv1 y timeouts.
  - `checker.go`: lógica de verificación de proyectos que detecta servicios apagados (`POWEROFF`) y los inicia automáticamente.
  - Pruebas unitarias del checker en `internal/aiven/checker_test.go`.
- **Soporte multi-proyecto**: configuración indexada de credenciales (`AIVEN_TOKEN_1`, `AIVEN_PROJECT_1`, etc.) con fallback a variables legacy (`AIVEN_TOKEN`, `AIVEN_PROJECT`).
- **Ejecución paralela de workers** (`cmd/main.go`): orquestación de dos workers concurrentes mediante goroutines y `sync.WaitGroup`.
  - Worker de Minecraft (Aternos): intervalo de 24 horas, cooldown de 70 minutos.
  - Worker de Aiven: intervalo de 60 minutos, cooldown de 30 minutos.
- **Tarea Aiven compuesta**: `buildAivenTask` itera sobre todas las credenciales configuradas; un fallo en un proyecto no aborta el resto, y los errores se combinan con `errors.Join` para activar el cooldown solo si al menos una API falló.

### Changed

- Documentación actualizada en README para reflejar el nuevo módulo de Aiven y la ejecución paralela de workers.
- Mejoras en la configuración centralizada para soportar múltiples credenciales de API.

## [1.1.0] - 2026-08-07

### Added

- **Servicio con el pipeline del bot** (`internal/services`): `BotService` que orquesta la verificación del servidor y su arranque automático.
- **Configuración centralizada** (`internal/config`): struct `Config` con carga de `.env`, validación de variables obligatorias y preparación del archivo de almacenamiento de sesión.
- **Integración del servicio en el arranque del worker**: `cmd/main.go` consume la configuración centralizada y elimina la lectura repetitiva de variables de entorno.

### Changed

- **Tiempos de revisión periódica del servidor**: revisión cada 60 minutos y cooldown de 70 minutos ante fallos de Aternos.
- **Documentación**: mejoras en el README, quitando el enfoque exclusivo en Aternos, y adición de este changelog.

## [1.0.0] - 2026-08-04

### Added

- **Núcleo del bot en Go**: aplicación CLI (`cmd/main.go`) que orquesta la automatización completa del arranque de un servidor de Minecraft.
- **Configuración por variables de entorno**: carga de `.env` con `HOST`, `PORT`, `SERVER_ID` y `STORAGE_PATH`.
- **Verificación de estado del servidor** (`internal/minecraft`): comprobación de si el servidor está online mediante ping con `go-mcping`.
- **Arquitectura limpia**:
  - `internal/ports`: interfaz `BrowserManager` que abstrae la gestión del navegador.
  - `internal/adapters`: implementación de Playwright + Chromium con persistencia de sesión.
  - `internal/automation`: automatización del flujo de arranque del servidor en Aternos.
- **Persistencia de sesión**: guardado y reutilización de cookies/localStorage mediante `StorageState` (`storage/state.json`).
- **Automatización del arranque del servidor**:
  - Apertura de Aternos y navegación al panel de servidores.
  - Selección del servidor por `SERVER_ID`.
  - Clic en el botón de inicio con manejo de banners publicitarios (`Force: true`).
  - Manejo automático de diálogos/avisos (EULA, "Okay", "Accept", "I agree").
  - Detección de sesión activa y error claro si no hay sesión iniciada.
- **Espera de encendido**: sondeo del estado del servidor cada 15 segundos hasta 5 minutos hasta que responda al ping.
- **Worker programable** (`internal/worker`): ejecución periódica de tareas con intervalo y cooldown ante fallos, con protección ante solapamiento de ciclos y cancelación por contexto.
- **Pruebas**: cobertura inicial del worker en `internal/worker/worker_test.go`.

### Support

- **Actualmente automatiza el arranque del servidor en Aternos**; la arquitectura de puertos y adaptadores está diseñada para ampliar el soporte a otros servicios en el futuro.
- Requiere Chromium y sesión iniciada manualmente en la primera ejecución.

### Next steps (no incluidos en esta versión)

- Comandos de control del bot.
- Exposición de una API HTTP.
- Soporte para otras plataformas de hosting.
- Integración de mensajería (p. ej., Discord/Telegram).
