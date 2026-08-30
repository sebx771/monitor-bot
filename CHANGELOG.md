# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/).

## [1.4.0] - 2026-08-30

### Added

- **Supabase implementation** (`internal/supabase`): new worker that monitors Supabase projects and automatically restores them when they are inactive.
  - `client.go`: HTTP client for the API (`api.supabase.com`) with Bearer authentication, and methods `ListProjects`, `GetProject` and `RestoreProject`.
  - `checker.go`: detects projects with `INACTIVE` status and restores them; logs results per project.
  - Unit tests for the checker in `internal/supabase/checker_test.go`.
- **Per-module configuration split** (`internal/config`): refactor of the centralized configuration into independent files per module (`aiven.go`, `aternos.go`, `supabase.go`, `helpers.go`) with `Config` getters (`GetAivenConfig`, `GetAternosConfig`, `GetSupaBaseConfig`).
- **Dynamic worker activation** (`cmd/main.go`): workers start conditionally according to the per-module environment variables (`AIVEN_ENABLED`, `ATERNOS_ENABLED`, `SUPABASE_ENABLED`).
- **Enablement helper** (`internal/config/helpers.go`): `GetEnabledVar` reads boolean environment variables and validates that the value is `true`/`false`, with clear error messages.
- **New logging system** (`internal/logger`): centralized logger based on `log/slog` with structured attributes and the emitting module attached to each line.
  - `logger.go`: `slog` wrapper with `NewLogger(module)` and `Debug`, `Info`, `Warn` and `Error` methods.
  - `handler.go` and `base_logger.go`: handler configuration and the shared base logger.
- **Log migration**: all modules (`config`, `worker`, `supabase`, `aiven`, `services`, etc.) migrate their logs to the new centralized system.

### Changed

- **Worker orchestration in `cmd/main.go`**: independent activation per enabled module, avoiding launching unconfigured workers and improving observability with structured logs.
- **Translation to English**: user-facing and developer-facing text (log messages, error strings and code comments) across all modules (`config`, `worker`, `supabase`, `aiven`, `minecraft`, `services`, `adapters`, `automation`, `logger` and tests), plus `README.md`, `.env.example` and this changelog, translated from Spanish to English.

### Support

- The Aternos module is marked as disabled in the production configuration (`ATERNOS_ENABLED=false`), leaving the Aiven and Supabase workers active by default depending on what is enabled.

## [1.3.0] - 2026-08-18

### Added

- **`StateStorage` port** (`internal/ports/state_storage.go`): interface that abstracts the remote storage of the session state file with `DownloadState` and `UploadState`.
- **GitHub Gist adapter** (`internal/adapters/github_gist.go`): implements `StateStorage` using the GitHub API to download and update the session state in a Gist.
  - Bearer authentication, `application/vnd.github+json` headers, 30 second timeout and automatic directory creation.
  - Unit tests in `internal/adapters/github_gist_test.go`.
- **Session sync with GitHub Gist** (`internal/services/aternos_services.go`): `BotService` downloads the session state from the Gist before starting the browser and uploads it after each use; retries the cycle on invalid cookies.
- **New environment variables** (`internal/config.go`):
  - `HEADLESS`: controls Playwright's headless mode (`internal/adapters/playwright.go`).
  - `GITHUB_TOKEN` and `GIST_ID`: token and Gist ID for session sync.
- **Multi-stage Dockerfile**: build with `golang:1.24-bookworm` and `debian:bookworm-slim` runtime, with the `/app/storage` directory for production deploy.

### Changed

- **Aiven production mode** (`cmd/main.go`): the Aternos bot initialization is commented out, leaving only the Aiven worker active for deploy.
- **Code formatting**: `go fmt` applied to `internal/aiven`, `internal/config.go`, `internal/minecraft` and `internal/ports`.
- **`.env.example`**: documents the new environment variables for Gist and headless.

## [1.2.0] - 2026-08-09

### Added

- **Aiven module** (`internal/aiven`): HTTP client to interact with the Aiven API and checker to review the service status.
  - `client.go`: client with `GetServices` and `StartService` methods, Bearer/aivenv1 authentication handling and timeouts.
  - `checker.go`: verification logic per project that detects powered-off services (`POWEROFF`) and starts them automatically.
  - Checker unit tests in `internal/aiven/checker_test.go`.
- **Multi-project support**: indexed credential configuration (`AIVEN_TOKEN_1`, `AIVEN_PROJECT_1`, etc.) with fallback to legacy variables (`AIVEN_TOKEN`, `AIVEN_PROJECT`).
- **Parallel worker execution** (`cmd/main.go`): orchestration of two concurrent workers using goroutines and `sync.WaitGroup`.
  - Minecraft (Aternos) worker: 24 hour interval, 70 minute cooldown.
  - Aiven worker: 60 minute interval, 30 minute cooldown.
- **Composite Aiven task**: `buildAivenTask` iterates over all configured credentials; a failure in one project does not abort the rest, and errors are combined with `errors.Join` to trigger the cooldown only if at least one API failed.

### Changed

- Updated documentation in the README to reflect the new Aiven module and parallel worker execution.
- Improvements to the centralized configuration to support multiple API credentials.

## [1.1.0] - 2026-08-07

### Added

- **Service with the bot pipeline** (`internal/services`): `BotService` that orchestrates the server check and its automatic startup.
- **Centralized configuration** (`internal/config`): `Config` struct with `.env` loading, validation of required variables and preparation of the session storage file.
- **Service integration in the worker startup**: `cmd/main.go` consumes the centralized configuration and removes the repetitive environment variable reading.

### Changed

- **Periodic server check times**: check every 60 minutes and 70 minute cooldown on Aternos failures.
- **Documentation**: README improvements, dropping the exclusive focus on Aternos, and addition of this changelog.

## [1.0.0] - 2026-08-04

### Added

- **Core bot in Go**: CLI application (`cmd/main.go`) that orchestrates the complete automation of starting a Minecraft server.
- **Environment variable configuration**: `.env` loading with `HOST`, `PORT`, `SERVER_ID` and `STORAGE_PATH`.
- **Server status check** (`internal/minecraft`): checking if the server is online using a ping with `go-mcping`.
- **Clean architecture**:
  - `internal/ports`: `BrowserManager` interface that abstracts browser management.
  - `internal/adapters`: Playwright + Chromium implementation with session persistence.
  - `internal/automation`: server startup flow automation on Aternos.
- **Session persistence**: saving and reusing cookies/localStorage through `StorageState` (`storage/state.json`).
- **Server startup automation**:
  - Opening Aternos and navigating to the server panel.
  - Server selection by `SERVER_ID`.
  - Click on the start button with ad banner handling (`Force: true`).
  - Automatic dialog/notice handling (EULA, "Okay", "Accept", "I agree").
  - Active session detection and a clear error if there is no logged-in session.
- **Power-on wait**: server status polling every 15 seconds for up to 5 minutes until it responds to ping.
- **Programmable worker** (`internal/worker`): periodic task execution with interval and cooldown on failures, with protection against overlapping cycles and context cancellation.
- **Tests**: initial worker coverage in `internal/worker/worker_test.go`.

### Support

- **Currently automates server startup on Aternos**; the ports and adapters architecture is designed to extend support to other services in the future.
- Requires Chromium and a manually logged-in session on the first run.

### Next steps (not included in this version)

- Bot control commands.
- HTTP API exposure.
- Support for other hosting platforms.
- Messaging integration (e.g., Discord/Telegram).
