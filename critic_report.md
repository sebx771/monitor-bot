# Stage 7 Critic Review Report: `design_draft.md`

*   **Target Document:** [design_draft.md](file:///c:/Users/Usuario/OneDrive/Desktop/monitor-bot/design_draft.md) (`Lines 1-187`)
*   **Reviewers:** `egm-design`, `egm-structural`, `egm-style`
*   **Date:** 2026-11-18

## Severity Summary

| Severity | Count | Status |
| :--- | :--- | :--- |
| 🔴 **Blocker** | 3 | Immediate architectural remediation required |
| 🟡 **Major** | 5 | Structural and ergonomic risks |
| 🟢 **Nit** | 3 | Style, jargon, and visual legibility refinements |
| **Total Findings** | 11 | Deduplicated across all three lenses |

---

## Prioritized Findings

### 1. Contradiction Between Open Question 1 and Hook Specification - 🔴 Blocker
* **Reviewer:** `egm-structural` / `egm-design`
* **Affected Scope:** `design_draft.md:42`, `design_draft.md:126-130`
* **Description:** Line 42 asks whether the active status filter should reset on search or maintain boolean intersection (`AND`), while Line 128 and Line 28 explicitly specify and design around strict boolean intersection.
* **Recommendation:** Remove Question 1 from "Open Questions" and formally commit to strict boolean conjunction (`status === activeFilter && matchesQuery`) with a dedicated "Clear all filters" button for empty states.

### 2. Undefined State Scope and Timer Leaks in `useWorkerTrigger` - 🔴 Blocker
* **Reviewer:** `egm-design` / `egm-structural`
* **Affected Scope:** `design_draft.md:62-64`, `design_draft.md:112-114`, `design_draft.md:136-140`
* **Description:** The technical plan does not define whether `useWorkerTrigger` is instantiated per row or centralized. If instantiated per row, multiple rapid clicks on different workers trigger conflicting unmanaged timers, and unmounting rows leak active timeouts.
* **Recommendation:** Centralize active confirmation state in `useDashboard` (`activeConfirmationId: string | null`), ensuring that confirming Worker B auto-cancels any pending confirmation on Worker A, and enforce `useEffect` cleanup on unmount.

### 3. Silent Auto-Refresh Deadlock on Error Dampening - 🔴 Blocker
* **Reviewer:** `egm-design` / `egm-structural`
* **Affected Scope:** `design_draft.md:124-125`
* **Description:** The plan pauses background polling after 3 consecutive network failures "until manually cleared or retried", but defines no UI affordance to show this paused state or let operators retry. Telemetry silently freezes.
* **Recommendation:** Specify a warning state in the Header refresh control ("Sync Paused: Network Issue [Retry]") that displays after 3 consecutive failures, where clicking retry or refocusing the tab resets the error count and resumes polling.

### 4. Omission of Storage Quota & Private Browsing Fallbacks - 🟡 Major
* **Reviewer:** `egm-design`
* **Affected Scope:** `design_draft.md:27`, `design_draft.md:78`, `design_draft.md:100`
* **Description:** `localStorage.setItem` throws `DOMException: QuotaExceededError` or `SecurityError` in private browsing or strict iframe sandboxes. No fallback is defined.
* **Recommendation:** Wrap storage reads/writes in a safe helper with fallback to in-memory React state.

### 5. Toast Queue Unbounded Stack & Viewport Overflow - 🟡 Major
* **Reviewer:** `egm-design`
* **Affected Scope:** `design_draft.md:141-148`
* **Description:** No maximum queue depth, dismiss button, or container scroll constraints are specified for `ToastPortal`. Cascading cluster alerts will overflow the screen.
* **Recommendation:** Cap visible toasts at 4 (dropping oldest non-critical), add an explicit dismiss button (`×`), and set `max-height: calc(100vh - 3rem); overflow-y: auto;`.

### 6. Missing Mobile Breakpoint & Responsive Table Strategy - 🟡 Major
* **Reviewer:** `egm-design` / `egm-structural`
* **Affected Scope:** `design_draft.md:95`, `design_draft.md:183-187`
* **Description:** HTML table layouts clip or overflow horizontally on screens < 768px without a responsive adaptation plan.
* **Recommendation:** Wrap `WorkerTable` in an `overflow-x: auto` container and automatically fall back to Card Grid on viewports < 640px.

### 7. Omission of Test Strategy and Phased Rollout Plan - 🟡 Major
* **Reviewer:** `egm-structural`
* **Affected Scope:** Document Structure
* **Description:** The draft omitted the 3-phase rollout roadmap from `concepts_brief.md` and contains no test verification criteria.
* **Recommendation:** Restore the 3-phase rollout plan and add a "Verification Strategy" section detailing unit tests for custom hooks.

### 8. Incomplete Screen Reader & ARIA Contracts - 🟡 Major
* **Reviewer:** `egm-design` / `egm-structural`
* **Affected Scope:** `design_draft.md:28`, `design_draft.md:95-96`, `design_draft.md:141-148`
* **Description:** Lacks `aria-live` specifications for `ToastPortal` (polite vs. assertive) and accessible column headers (`<th scope="col">`) / button labels for `WorkerTable`.
* **Recommendation:** Define ARIA roles, live regions, and icon button labeling contracts explicitly.

### 9. Ungrounded Jargon and Promotional Fluff - 🟢 Nit
* **Reviewer:** `egm-style`
* **Affected Scope:** `design_draft.md:13`, `design_draft.md:85`, `design_draft.md:104`, `design_draft.md:128`
* **Description:** Phrases like *"operational cockpit"* and *"predicate conjunction"* add jargon. Replace with *"operational dashboard"* and *"filtering workers matching both search query and status"*.

### 10. Redundant LOC Rule Repetition - 🟢 Nit
* **Reviewer:** `egm-style`
* **Affected Scope:** `design_draft.md:31`, `design_draft.md:118`, `design_draft.md:187`
* **Description:** The `< 150 LOC` constraint is repeated four times. Define once in Goals and remove rest.

### 11. Micro-Avatar (28px) Expression Legibility - 🟢 Nit
* **Reviewer:** `egm-design`
* **Affected Scope:** `design_draft.md:27`, `design_draft.md:95`, `design_draft.md:169`
* **Description:** Scaled-down 28px SVG avatars can blur facial emotional cues (`mad`, `sleepy`, `happy`).
* **Recommendation:** Surround 28px micro-avatars with high-contrast status color rings (green, amber, red) for instant state recognition.
