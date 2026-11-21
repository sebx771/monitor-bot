import type { RefreshInterval } from "../../models/Dashboard";
import styles from "./Header.module.css";

interface RefreshControlProps {
  onRefresh: () => void;
  isRefreshing: boolean;
  refreshInterval: RefreshInterval;
  onIntervalChange: (interval: RefreshInterval) => void;
  isPausedDueToErrors?: boolean;
  onResetErrorsAndResume?: () => void;
}

export function RefreshControl({
  onRefresh,
  isRefreshing,
  refreshInterval,
  onIntervalChange,
  isPausedDueToErrors = false,
  onResetErrorsAndResume,
}: RefreshControlProps) {
  if (isPausedDueToErrors) {
    return (
      <button
        type="button"
        className={styles.pausedAlertBadge}
        onClick={onResetErrorsAndResume}
        aria-label="Background polling paused due to network issues. Click to retry."
      >
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5">
          <circle cx="12" cy="12" r="10" />
          <line x1="12" y1="8" x2="12" y2="12" />
          <line x1="12" y1="16" x2="12.01" y2="16" />
        </svg>
        <span>Sync Paused (Retry)</span>
      </button>
    );
  }

  return (
    <div className={styles.refreshGroup}>
      <select
        className={styles.cadenceSelect}
        value={refreshInterval}
        onChange={(e) => {
          onIntervalChange(Number(e.target.value) as RefreshInterval);
        }}
        aria-label="Auto-refresh interval"
      >
        <option value={0}>Auto: Off</option>
        <option value={15}>Auto: 15s</option>
        <option value={30}>Auto: 30s</option>
        <option value={60}>Auto: 60s</option>
      </select>

      <button
        type="button"
        className={`${styles.refreshBtn} ${isRefreshing ? styles.refreshing : ""}`}
        onClick={onRefresh}
        disabled={isRefreshing}
        aria-label="Refresh workers status"
      >
        <svg className={styles.refreshIcon} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
          <path d="M21.5 2v6h-6M21.34 15.57a10 10 0 1 1-.57-8.38l5.67-5.19" />
        </svg>
        <span>{isRefreshing ? "Refreshing..." : "Refresh"}</span>
      </button>
    </div>
  );
}

export default RefreshControl;
