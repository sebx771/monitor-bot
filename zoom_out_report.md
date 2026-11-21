# Stage 5 Initial Zoom Out Review Report

*   **Target Document:** [design_draft.md](file:///c:/Users/Usuario/OneDrive/Desktop/monitor-bot/design_draft.md) (`Lines 1-187`)
*   **Reviewers:** `egm-zoom-out`, `egm-inventive-catalyst`
*   **Date:** 2026-11-18

## Severity Summary

| Severity | Count | Status |
| :--- | :--- | :--- |
| 🔴 **Blocker** | 0 | None detected |
| 🟡 **Major** | 3 | Actionable architectural risks identified |
| 🟢 **Nit** | 2 | Ergonomic and responsive refinements |

---

## Prioritized Findings

### 1. Root Coordinator Hook Re-render Invalidation (`useToasts` Coupling) - 🟡 Major
* **Reviewer:** `egm-zoom-out`
* **Affected Scope:** `design_draft.md:83-110`, `design_draft.md:132-134`
* **Description:** Section 2 claims that `useToasts.ts` *"Auto-removes toasts on timer expiration without re-rendering parent dashboard containers."* However, the architecture diagram routes `ToastHook --> DashboardHook --> AppView`. Composing toast queue state directly inside `useDashboard` causes the root coordinator and its children (`Header`, `MetricsOverview`, `WorkerGrid`, `WorkerTable`) to re-render whenever a toast mounts, ticks, or auto-dismisses.
* **Recommendation:** Decouple toast notification dispatch from data fetching. Expose a lightweight `ToastContext` / event listener pattern, or instantiate `useToasts` directly inside `<ToastPortal />` so that transient alert timers never trigger re-renders in the data views.

### 2. Mutation Lock Coarse-Grained Polling Invalidation - 🟡 Major
* **Reviewer:** `egm-inventive-catalyst`
* **Affected Scope:** `design_draft.md:149-153`
* **Description:** The specification states that during an active worker trigger, *"Polling responses received during an active mutation are filtered so that optimistic worker execution states are not overwritten."* Dropping entire polling responses globally creates an observability blackout across all other cluster workers during long-running tasks.
* **Recommendation:** Apply granular per-worker optimistic reconciliation: allow background poll responses to update the worker list, but preserve the optimistic pending status specifically for the worker matching `activeMutationWorkerId` until its execution resolves or times out.

### 3. Unresolved Ergonomic & Navigation Ambiguities (Open Questions 1 & 2) - 🟡 Major
* **Reviewer:** `egm-zoom-out`
* **Affected Scope:** `design_draft.md:40-44`
* **Description:** Open Questions 1 and 2 define core interaction behaviors for search/filter conjunction and compact table navigation, but remain unresolved in the draft. This introduces interface ambiguity in `useDashboardFilters.ts` and `WorkerTableRow.tsx`.
* **Recommendation:** Lock down both decisions in the design document:
  1. **Status + Search Conjunction:** Maintain strict boolean intersection (`status === activeFilter && matchesQuery`) with a dedicated "Clear all filters" chip.
  2. **Table Navigation:** Clicking a worker row or detail badge navigates directly to `WorkerDetailModal` rather than expanding nested accordion DOM trees, preserving UI simplicity (< 120 LOC).

### 4. Omission of Browser Network Online/Offline Lifecycle Guards - 🟢 Nit
* **Reviewer:** `egm-zoom-out`
* **Affected Scope:** `design_draft.md:120-125`
* **Description:** `useAutoRefresh.ts` tracks `document.visibilityState` but omits explicit integration with `window.navigator.onLine` and `online`/`offline` window events.
* **Recommendation:** Explicitly document `window.addEventListener('offline' | 'online')` inside `useAutoRefresh.ts`. Polling pauses immediately when offline and reconciles immediately upon network restoration.

### 5. Viewport Adaptation Strategy for High-Density Table Mode - 🟢 Nit
* **Reviewer:** `egm-inventive-catalyst`
* **Affected Scope:** `design_draft.md:95`, `design_draft.md:183-187`
* **Description:** The HTML `<table>` mode can clip or force awkward horizontal scrollbars on narrow mobile/tablet screens (< 640px).
* **Recommendation:** Specify an adaptive container rule: wrap `WorkerTable` in an `overflow-x: auto` wrapper, and default to `WorkerGrid` on viewports under 640px.
