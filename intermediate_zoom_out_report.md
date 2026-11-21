# Stage 8 Intermediate Zoom Out Review Report

*   **Target Document:** [design_draft.md](file:///c:/Users/Usuario/OneDrive/Desktop/monitor-bot/design_draft.md) (`Lines 1-199`)
*   **Reviewer:** `egm-zoom-out`
*   **Date:** 2026-11-18

## Severity Summary

| Severity | Count | Status |
| :--- | :--- | :--- |
| 🔴 **Blocker** | 0 | All previous blockers successfully remediated |
| 🟡 **Major** | 2 | Interaction simplicity and re-render decoupling |
| 🟢 **Nit** | 2 | Subcomponent decomposition and mutation merge clarity |

---

## High-Altitude Architecture & Simplicity Assessment

The remediations incorporated from Stage 6 (Comprehension) and Stage 7 (Critic) have significantly strengthened the design document without adding unnecessary architectural bloat:
1. **Core Problem Alignment**: Strictly centered on operator ergonomics (Zero-CLS, operational density, safe triggering, background telemetry) without drifting into backend re-architecture.
2. **YAGNI & Complexity Control**: Relies on native HTML `<table>`, modular React hooks, and Page Visibility APIs instead of external dependencies (Redux, Zustand, TanStack Table).
3. **Specification Readiness**: Verified 0 Blockers. High specification readiness achieved.

---

## Prioritized Findings

### 1. Close Open Question 1 (Table Navigation to `WorkerDetailModal`) - 🟡 Major
* **Reviewer:** `egm-zoom-out`
* **Affected Scope:** `design_draft.md:40-43`, `design_draft.md:93`
* **Description:** Retaining Open Question 1 regarding whether table badges expand accordion rows or navigate to `WorkerDetailModal` leaves an interaction hazard. Expandable nested accordion rows balloon DOM nodes and CSS focus trapping, risking breaching the `< 120 LOC` rule for `WorkerTable`.
* **Recommendation:** Formally close the question in the specification: clicking a worker row or multi-target badge delegates directly to the existing `WorkerDetailModal`.

### 2. Decouple Toast Queue Re-render Cycle from `useDashboard` - 🟡 Major
* **Reviewer:** `egm-zoom-out`
* **Affected Scope:** `design_draft.md:81-103`, `design_draft.md:129-133`
* **Description:** Composing toast queue state directly in the root `useDashboard` coordinator causes the entire tree (`Header`, `MetricsOverview`, `WorkerGrid`, `WorkerTable`) to re-render whenever a toast mounts, ticks, or auto-dismisses.
* **Recommendation:** Decouple toast state from `useDashboard` via an isolated `ToastContext` / event listener pattern so transient alert lifecycles never cause re-renders in the data views.

### 3. Header Subcomponent Decomposition Safeguard (`RefreshControl.tsx`) - 🟢 Nit
* **Reviewer:** `egm-zoom-out`
* **Affected Scope:** `design_draft.md:89-90`, `design_draft.md:119-121`
* **Description:** Housing the search bar, view switcher, cadence dropdown, sync button, and the 3-failure amber retry badge in `Header.tsx` risks exceeding the `< 120 LOC` limit.
* **Recommendation:** Extract the auto-refresh and network pause controls into a colocated subcomponent (`src/components/Header/RefreshControl.tsx` < 50 LOC).

### 4. Per-Worker Optimistic Mutation Reconciliation Semantics - 🟢 Nit
* **Reviewer:** `egm-zoom-out`
* **Affected Scope:** `design_draft.md:150-152`
* **Description:** Ensure phrasing explicitly specifies a per-worker state merge: background polling updates all unmutated workers normally while selectively retaining the optimistic pending state for `worker.id === activeMutationWorkerId`.

---

## Verdict
**High Specification Readiness (Pass with Minor Revisions)**. The high-level architecture has survived scrutiny with 0 blockers.
