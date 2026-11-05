# UI Changelog

All notable changes to the UI will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/).

## [Unreleased]

### Added

- **Admin dashboard controls** (`src`): added sign-in and sign-out, session restoration, and authenticated manual worker triggers.

## [1.0.0] - 2026-10-20

### Added

- **Responsive worker dashboard** (`src`): added a dark dashboard layout with worker cards, status overview, and a responsive worker grid.
- **Worker details** (`src/components/WorkerDetailModal`): added a modal showing worker status, monitored services, and last-run information.
- **Manual worker triggers** (`src/components/WorkerCard`, `src/services`): added API integration for triggering workers, with progress and result feedback.
