import type { ViewDensity, RefreshInterval } from "../../models/Dashboard";
import { RefreshControl } from "./RefreshControl";
import styles from "./Header.module.css";

export interface HeaderProps {
  searchQuery: string;
  onSearchChange: (query: string) => void;
  onRefresh: () => void;
  isRefreshing?: boolean;
  totalWorkers?: number;
  onToggleSidebar?: () => void;
  viewDensity: ViewDensity;
  onViewDensityChange: (density: ViewDensity) => void;
  refreshInterval: RefreshInterval;
  onIntervalChange: (interval: RefreshInterval) => void;
  isPausedDueToErrors?: boolean;
  onResetErrorsAndResume?: () => void;
  isAuthenticated: boolean;
  isAuthLoading: boolean;
  onLoginClick: () => void;
  onLogout: () => void;
}

export default function Header({
  searchQuery,
  onSearchChange,
  onRefresh,
  isRefreshing = false,
  totalWorkers = 0,
  onToggleSidebar,
  viewDensity,
  onViewDensityChange,
  refreshInterval,
  onIntervalChange,
  isPausedDueToErrors,
  onResetErrorsAndResume,
  isAuthenticated,
  isAuthLoading,
  onLoginClick,
  onLogout,
}: HeaderProps) {
  return (
    <header className={styles.header}>
      <div className={styles.leftSection}>
        {onToggleSidebar && (
          <button
            type="button"
            className={styles.mobileMenuBtn}
            onClick={onToggleSidebar}
            aria-label="Open navigation menu"
          >
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <line x1="3" y1="12" x2="21" y2="12" />
              <line x1="3" y1="6" x2="21" y2="6" />
              <line x1="3" y1="18" x2="21" y2="18" />
            </svg>
          </button>
        )}
        <div className={styles.titleArea}>
          <h2 className={styles.title}>Workers Dashboard</h2>
          <p className={styles.subtitle}>
            {totalWorkers === 1 ? "1 active worker" : `${String(totalWorkers)} active workers`}
          </p>
        </div>
      </div>

      <div className={styles.actionsArea}>
        <div className={styles.searchContainer}>
          <svg className={styles.searchIcon} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
            <circle cx="11" cy="11" r="8" />
            <line x1="21" y1="21" x2="16.65" y2="16.65" />
          </svg>
          <input
            type="text"
            className={styles.searchInput}
            placeholder="Search name, status, target..."
            value={searchQuery}
            onChange={(e) => { onSearchChange(e.target.value); }}
            aria-label="Search workers"
          />
          {searchQuery && (
            <button
              type="button"
              className={styles.clearBtn}
              onClick={() => { onSearchChange(""); }}
              aria-label="Clear search"
            >
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5">
                <line x1="18" y1="6" x2="6" y2="18" />
                <line x1="6" y1="6" x2="18" y2="18" />
              </svg>
            </button>
          )}
        </div>

        <div className={styles.densityToggle} role="group" aria-label="View Density">
          <button
            type="button"
            className={`${styles.densityBtn} ${viewDensity === "grid" ? styles.activeDensity : ""}`}
            onClick={() => { onViewDensityChange("grid"); }}
            aria-label="Grid view (Card)"
            aria-pressed={viewDensity === "grid"}
          >
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <rect x="3" y="3" width="7" height="7" />
              <rect x="14" y="3" width="7" height="7" />
              <rect x="14" y="14" width="7" height="7" />
              <rect x="3" y="14" width="7" height="7" />
            </svg>
          </button>
          <button
            type="button"
            className={`${styles.densityBtn} ${viewDensity === "table" ? styles.activeDensity : ""}`}
            onClick={() => { onViewDensityChange("table"); }}
            aria-label="Table view (Compact)"
            aria-pressed={viewDensity === "table"}
          >
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <line x1="8" y1="6" x2="21" y2="6" />
              <line x1="8" y1="12" x2="21" y2="12" />
              <line x1="8" y1="18" x2="21" y2="18" />
              <line x1="3" y1="6" x2="3.01" y2="6" />
              <line x1="3" y1="12" x2="3.01" y2="12" />
              <line x1="3" y1="18" x2="3.01" y2="18" />
            </svg>
          </button>
        </div>

        <RefreshControl
          onRefresh={onRefresh}
          isRefreshing={isRefreshing}
          refreshInterval={refreshInterval}
          onIntervalChange={onIntervalChange}
          isPausedDueToErrors={isPausedDueToErrors}
          onResetErrorsAndResume={onResetErrorsAndResume}
        />

        {!isAuthLoading && (
          isAuthenticated ? (
            <button type="button" className={styles.logoutBtn} onClick={onLogout} aria-label="Sign out">
              Logout
            </button>
          ) : (
            <button type="button" className={styles.loginBtn} onClick={onLoginClick} aria-label="Sign in">
              Login
            </button>
          )
        )}
      </div>
    </header>
  );
}
