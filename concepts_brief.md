# UI/UX Overhaul & Observability Enhancements for Monitor Bot

* **Author:** Sebastian
* **Date:** 2026-11-18

## Overview

The current web interface for the Monitor Bot cluster provides baseline operator functionality: worker status visualization via animated avatars (`@blobatar/react`), single-action manual triggers, keyword search, simple aggregate counters, and modal-based token authentication.

However, the dashboard suffers from ergonomic bottlenecks:
1. Low information density caused by large avatar canvases (160px height), forcing constant scrolling when clusters exceed 4–6 workers.
2. Jarring Cumulative Layout Shift (CLS) when alert and sync toasts inject directly into the content stream above the card grid.
3. Lack of automated state synchronization, requiring manual refresh clicks to detect degraded or paused workers.
4. Unguarded execution triggers that expose operators to accidental activation.
5. Passive metric cards that do not link to filtering actions.

This brief proposes a modular overhaul across UI and UX dimensions without altering backend Go APIs, adhering to strict clean architecture constraints (<150 LOC per file, strict TypeScript, modular hook isolation).

---

## Analysis

### 1. Inventive Contradictions & Catalog Transforms

* **Contradiction 1: Expressive Personality vs. Information Density**
  * *Conflict:* Large animated avatars (`@blobatar/react`) provide immediate emotional state cues (`mad`, `sleepy`, `happy`), but consume excessive screen area. Stripping avatars turns the dashboard into a sterile spreadsheet.
  * *Transform (Dual-Density Polymorphism):* Support two distinct views. Retain rich Blobatar cards in default "Comfortable Grid" mode, and provide a "Compact Table" mode where avatars condense into 28px micro-chips or status dots alongside dense tabular metadata (uptime, targets, last run timestamp).

* **Contradiction 2: Real-Time Telemetry vs. Zero Backend Changes**
  * *Conflict:* Operators need up-to-date worker statuses, but Server-Sent Events (SSE) or WebSockets require modifying backend Go server code, which is out of scope. Indiscriminate client polling risks thundering herds against the local API.
  * *Transform (Adaptive Visibility-Aware Polling):* Implement client-side polling with selectable cadences (Off, 15s, 30s, 60s) tied to `document.visibilityState` (Page Visibility API) and `window.navigator.onLine`. When the dashboard tab is in the background or hidden, polling halts immediately and reconciles on refocus.

* **Contradiction 3: Alert Prominence vs. Layout Stability (CLS)**
  * *Conflict:* Critical alerts must be immediately visible, but injecting them above `<main>` shifts the entire layout downwards, breaking interaction target stability.
  * *Transform (Isolated Toast Portal Anchor):* Decouple notifications from document flow using a floating `ToastPortal` mounted to `document.body` via `createPortal`, matching the architecture of the existing `Modal` component. Floating notifications stack with fixed CSS positioning.

---

### 2. Contrarian & Complexity Minimization (YAGNI)

* **No External State Libraries (Redux/Zustand):** All filter, view density, and auto-refresh states belong to local dashboard scope. Composing dedicated custom hooks (`useDashboardFilters`, `useAutoRefresh`, `useToasts`) inside `useDashboard` avoids adding third-party store dependencies.
* **No Heavy Data Grid Engines (TanStack Table / AG Grid):** For clusters of 5–50 nodes, heavy virtualization libraries add bundle weight and maintenance overhead. A semantic, accessible HTML `<table>` styled with CSS Modules satisfies all density requirements while staying under 120 lines of code.
* **No Complex Exponential Backoff State Machines:** Auto-refresh requires only a simple interval with an `AbortController` signal and consecutive failure dampening (pausing sync after 3 consecutive network failures with an amber recovery banner).

---

### 3. Observability & SRE Risks

* **Stale State False Confidence:** If background polling fails silently while an operator is away from their desk, stale data might display a failed worker as operational.
  * *Mitigation:* Render an explicit "Last synced X seconds ago" timestamp and change the refresh control badge to amber if background sync fails.
* **Race Condition Between Trigger and Polling:** If an operator manually triggers a worker during an automated background poll tick, a slow poll response could overwrite the optimistic pending state.
  * *Mitigation:* Track trigger mutations with dedicated sequence locks or request cancellation in `useDashboard` to prevent poll responses from clearing pending mutation indicators.
* **Read Contention on In-Memory Go State:** Excessive concurrent polling across multiple open operator tabs could contend for the mutex on `/api/v1/status`.
  * *Mitigation:* Cap minimum polling frequency at 10 seconds (default to 30s) and enforce tab visibility suspension.

---

### 4. Developer Ergonomics & UX Refinements

* **Metric Cards as Conjoint Filter Toggles:** Metrics cards (`Total`, `Operational`, `Alert / Paused`) act as interactive toggle buttons. Clicking "Alert / Paused" filters the grid to degraded workers; clicking again clears it. When combined with keyword search, filters evaluate as a strict conjunction (`AND`) with an explicit "Clear all filters" badge.
* **Two-Step Inline Trigger Guard:** Rather than interrupting operators with an intrusive modal confirmation for standard triggers, use an inline two-step confirmation pattern: clicking `⚡ Trigger` turns the button into an amber `Confirm? (3s)` countdown state before executing.
* **Operator Preference Persistence:** Store view density (`grid` vs. `compact`) and auto-refresh interval in `localStorage` (`monitor_bot_view_density`, `monitor_bot_refresh_rate`) so returning operators retain their view configuration.

---

### 5. Architectural Decomposition (<150 LOC Rule)

To prevent `useDashboard.ts` (currently 156 LOC) and `App.tsx` (currently 144 LOC) from exceeding architectural limits:

```
ui/src/
├── components/
│   ├── Header/
│   │   ├── Header.tsx                     # View switcher & auto-refresh controls (<120 LOC)
│   │   └── Header.module.css
│   ├── MetricsOverview/
│   │   ├── MetricsOverview.tsx            # Interactive filter container (<90 LOC)
│   │   ├── MetricFilterChip.tsx           # Accessible filter pill (<60 LOC)
│   │   └── MetricsOverview.module.css
│   ├── Toast/
│   │   ├── ToastPortal.tsx                # Portal into document.body (<70 LOC)
│   │   ├── ToastItem.tsx                  # Floating alert item (<80 LOC)
│   │   └── Toast.module.css
│   ├── WorkerTable/                       # Compact list view
│   │   ├── WorkerTable.tsx                # Semantic table view (<120 LOC)
│   │   ├── WorkerTableRow.tsx             # Dense row with micro-avatar (<110 LOC)
│   │   └── WorkerTable.module.css
│   ├── WorkerCard/
│   │   └── hooks/
│   │       └── useWorkerTrigger.ts        # 2-step inline confirmation hook (<70 LOC)
├── hooks/
│   ├── useDashboard.ts                    # Root dashboard coordinator (<100 LOC)
│   ├── useDashboardFilters.ts             # Filter predicates & search logic (<90 LOC)
│   ├── useAutoRefresh.ts                  # Visibility-aware polling loop (<80 LOC)
│   └── useToasts.ts                       # Toast queue management (<75 LOC)
└── models/
    └── Dashboard.ts                       # ViewDensity & FilterState types (<40 LOC)
```

---

### 6. Ideal Final Result (IFR)

An operational cockpit where operators instantly toggle between expressive card monitoring and dense tabular triage, filter cluster states with single-click metric chips, observe background sync status with zero layout jitter, and trigger tasks with high-confidence inline confirmations—delivered with zero backend Go changes and modular, strictly typed React 19 architecture.

---

## Goals vs. Non-Goals

* **In-Scope Goals:**
  * Extract transient alerts to a layout-isolated `ToastPortal` on `document.body` to eliminate CLS.
  * Implement `useAutoRefresh` with configurable intervals (Off, 15s, 30s, 60s) and Page Visibility suspension.
  * Refactor `MetricsOverview` into clickable quick-filter buttons with keyboard accessibility (`role="button"`, `aria-pressed`).
  * Introduce `WorkerTable` compact view with micro-avatars and persist view preference in `localStorage`.
  * Add two-step inline trigger confirmation (`Trigger` -> `Confirm?`) to prevent accidental execution.
  * Decompose existing monolithic files to strictly comply with the `<150 LOC` rule.

* **Out-of-Scope Non-Goals:**
  * Modifying backend Go server endpoints, schemas, or protocols (no SSE or WebSockets).
  * Introducing external state management packages (no Redux, Zustand, or TanStack Table).
  * Changing core `@blobatar/react` avatar designs or visual identity.
  * Implementing multi-tenant role-based access control (RBAC).

---

## Alternatives Weighed

* **Replace Blobatar completely with plain status icons**: Rejected by author. Animated expressions (`mad` for error, `sleepy` for paused) provide immediate emotive status awareness. Retaining avatars in Grid view and using compact micro-avatars in Table view preserves identity while restoring operational density.
* **Server-Sent Events (SSE) / WebSockets**: Rejected by author. Requires rewriting the Go backend routing and event bus. Configurable client-side polling with Page Visibility suspension provides adequate synchronization without expanding backend scope.
* **Intrusive Confirmation Modal for Triggers**: Rejected by author. Spawning a full modal dialogue for routine worker executions adds excessive cognitive and manual overhead. Two-step inline button confirmation provides sufficient accident prevention without modal fatigue.

---

## Key Decisions & Staging for Next Steps

* **Decision 1**: Decouple notifications first (`ToastPortal`) to establish layout stability before introducing view toggles.
* **Decision 2**: Decompose `useDashboard.ts` into isolated single-responsibility hooks (`useAutoRefresh`, `useDashboardFilters`, `useToasts`) before adding table rendering logic.
* **Staging**:
  * **Phase 1**: Stability Foundation (`ToastPortal`, `useToasts`, `useAutoRefresh`).
  * **Phase 2**: Filtering & IA (`useDashboardFilters`, `MetricFilterChip`, empty state recovery).
  * **Phase 3**: Density & Interaction Safety (`WorkerTable`, two-step inline trigger).
