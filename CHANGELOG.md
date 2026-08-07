# Changelog

Todas las cambios notables de este proyecto se documentarán en este archivo.

El formato se basa en [Keep a Changelog](https://keepachangelog.com/es/1.1.0/),
y este proyecto adhiere a [Semantic Versioning](https://semver.org/lang/es/).

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
