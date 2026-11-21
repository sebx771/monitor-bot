# Stage 6 Comprehension Test Report

*   **Target Document:** [design_draft.md](file:///c:/Users/Usuario/OneDrive/Desktop/monitor-bot/design_draft.md) (`Lines 1-187`)
*   **Reviewer:** `egm-comprehension`
*   **Date:** 2026-11-18

--------------------------------------------------------------------------------

## High-Level Paraphrase

### What this proposes
The document proposes a pure frontend refactoring and UX overhaul of the Monitor Bot web dashboard. It replaces the existing layout-shifting alerts with an isolated, body-mounted toast notification portal, introduces an interactive status-filtering mechanism on top of a dual-view density toggle (card grid vs. compact table), adds automated background polling that pauses when the browser tab is hidden, and replaces direct-click manual execution with a guarded two-step confirmation button—all achieved using modular React hooks (< 150 lines of code) without modifying the existing Go backend or adding external state/table dependencies.

### How the existing system works, per this document
According to the document, Monitor Bot supervises microservices and cloud workers through a React web client connected to a Go server backend. The current dashboard presents workers solely as large, fixed-size (160px) cards featuring animated `@blobatar/react` expressive avatars and static metric counts. Data fetching occurs only on initial page mount or through manual page refreshes. Furthermore, system alerts render inline above the worker grid (causing 3-second layout shifts upon auto-dismissal), and manual worker triggers execute immediately on a single click without confirmation.

--------------------------------------------------------------------------------

## Detailed Paraphrase

The design document addresses four operational pain points experienced when managing clusters of workers: low spatial density from 160px avatar cards, visual disorientation caused by Cumulative Layout Shift (CLS) when alert banners appear and disappear, stale data resulting from the lack of automated telemetry polling, and the operational risk of accidental job execution via unguarded trigger buttons.

To solve these issues without touching Go routes, databases, or adding heavyweight external libraries (like Redux, Zustand, TanStack Table, or AG Grid), the document details four architectural improvements:
1. **Zero-CLS Toast Portal**: Transient status and sync notifications are decoupled from the main content stream and rendered via `createPortal` directly to `document.body` in a fixed-position container at the top right, eliminating layout recalculations.
2. **Dual-Density & Conjoint Filtering**: A persistent toggle (`localStorage`) switches between the expressive 160px avatar grid and a dense semantic HTML `<table>` using 28px micro-avatars. Metric counter cards are made interactive (`role="button"`) to filter workers by status (`all`, `operational`, `alert`), functioning conjointly (`AND`) with text search.
3. **Adaptive Polling (`useAutoRefresh`)**: Client polling runs at selectable cadences (Off/0s, 15s, 30s, 60s), hooks into `document.visibilityState` to halt requests while the tab is backgrounded, retries on tab reactivation, and pauses after 3 consecutive network failures. Active polls use `AbortController` to prevent stale response races and respect a mutation lock flag during active manual triggers.
4. **Guarded Trigger FSM (`useWorkerTrigger`)**: Action buttons implement a two-step state machine (`IDLE` -> `CONFIRMING (3s)` -> `EXECUTING` -> `FEEDBACK` -> `IDLE`) where a worker run is dispatched only if confirmed within a 3-second window.

--------------------------------------------------------------------------------

## Comprehension Findings (Blockers, Majors, Nits)

### Blocker Findings
*None identified.* The document's technical narrative and architectural boundaries are fully interpretable, and the core operational mechanisms connect logically.

---

### Major Findings

1. **Severity**: 🟡 Major  
   **Affected Scope**: Section 2: `useDashboardFilters.ts` (Lines 126–129) vs. CUJ 1 (Lines 49–51) vs. Alternative 2 (Line 168)  
   **Description**: Inconsistent Status Taxonomy.  
   In CUJ 1, workers are described in `degraded` state and filtered via an `"Alert / Paused"` metric chip. In `useDashboardFilters`, the filter union is defined as `statusFilter: "all" | "operational" | "alert"`. In Alternative 2, states are described as `mad on failure, sleepy on pause, happy on operational`. A developer cannot determine how backend states (`degraded`, `paused`, `error`, `operational`) map into the 3-state filter union (`"all" | "operational" | "alert"`). The taxonomy must be formally aligned.

2. **Severity**: 🟡 Major  
   **Affected Scope**: Section 4: Data Synchronization & Concurrency Safeguards (Lines 151–152)  
   **Description**: Vague Concurrency / Mutation Lock Specification.  
   The document states that `useDashboard` *"tracks a mutation sequence flag. Polling responses received during an active mutation are filtered so that optimistic worker execution states are not overwritten"*. It is unstated whether this flag is a monotonic integer counter, a timestamp, or a boolean, and what "filtered" means (discarded entirely vs. granularly reconciled).

3. **Severity**: 🟡 Major  
   **Affected Scope**: Section 1: Architecture Diagram & Section 2: `useToasts.ts` (Lines 86, 103, 110, 132–135)  
   **Description**: Architectural Contradiction Regarding Toast State and Re-renders.  
   Line 134 states that `useToasts` *"Auto-removes toasts on timer expiration without re-rendering parent dashboard containers"*. However, the diagram shows `ToastHook --> DashboardHook --> AppView`. Composing `useToasts` directly inside `useDashboard` forces full re-renders whenever a toast mounts or expires.

---

### Nit Findings

1. **Severity**: 🟢 Nit  
   **Affected Scope**: Section: Goals & Section: State Management (Line 31, Line 120, Line 187)  
   **Description**: Undefined Acronym `LOC`. Expand "Lines of Code".

2. **Severity**: 🟢 Nit  
   **Affected Scope**: Section: Non-Goals (Line 38)  
   **Description**: Undefined Acronym `RBAC`. Expand "Role-Based Access Control".

3. **Severity**: 🟢 Nit  
   **Affected Scope**: Section: Overview (Line 13)  
   **Description**: Acronym Order `Zero-CLS`. Introduce "Cumulative Layout Shift (CLS)" before "Zero-CLS".

4. **Severity**: 🟢 Nit  
   **Affected Scope**: Section: Open Questions (Line 43)  
   **Description**: Unintroduced Reference to `WorkerDetailModal`. Define its role in the system.

5. **Severity**: 🟢 Nit  
   **Affected Scope**: Section 2: `useWorkerTrigger.ts` (Lines 138–139)  
   **Description**: Incomplete FSM Transition Definition for `FEEDBACK`. Define the duration/trigger for transitioning from `FEEDBACK` back to `IDLE` (e.g. 3-second timeout).
