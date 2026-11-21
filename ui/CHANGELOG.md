# UI Changelog

All notable changes to the UI will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/).

## [Unreleased]

### Added

- **Dual-density view mode** (`src/components/WorkerTable`, `src/components/Header`): introduced toggle between comfortable card grid and compact semantic table view with 28px micro-avatars and high-contrast status rings, persisted via `safeStorage`.
- **Zero-CLS Toast Portal** (`src/components/Toast`, `src/hooks/useToasts`): decoupled transient status and error alerts into a fixed floating portal mounted to `document.body`, eliminating layout shifts.
- **Interactive metric quick-filters** (`src/components/MetricsOverview`, `src/hooks/useDashboardFilters`): transformed passive metric cards into keyboard-accessible `aria-pressed` filter chips that evaluate conjointly (`AND`) with keyword search.
- **Adaptive background polling** (`src/hooks/useAutoRefresh`, `src/components/Header/RefreshControl`): added configurable auto-refresh cadence (Off, 15s, 30s, 60s) with automatic suspension when tab is hidden or offline, and visual pause alerts after consecutive network failures.
- **Guarded two-step execution triggers** (`src/components/WorkerCard/hooks/useWorkerTrigger`): added inline 3-second confirmation state (`Trigger` -> `Confirm? (3s)`) with mutual exclusion across cluster workers.
- **Resilient local storage wrapper** (`src/utils/safeStorage`): safe storage access with graceful in-memory Map fallback for restricted or private browsing modes.
- **Admin dashboard controls** (`src`): added sign-in and sign-out, session restoration, and authenticated manual worker triggers.

### Changed

- **Modular hook decomposition** (`src/hooks/useDashboard`): split monolithic dashboard logic into single-responsibility hooks strictly adhering to the `< 150 LOC` rule.

## [1.0.0] - 2026-10-20

### Added

- **Responsive worker dashboard** (`src`): added a dark dashboard layout with worker cards, status overview, and a responsive worker grid.
- **Worker details** (`src/components/WorkerDetailModal`): added a modal showing worker status, monitored services, and last-run information.
- **Manual worker triggers** (`src/components/WorkerCard`, `src/services`): added API integration for triggering workers, with progress and result feedback.
