import styles from "./Header.module.css";

export interface HeaderProps {
  searchQuery: string;
  onSearchChange: (query: string) => void;
  onRefresh: () => void;
  isRefreshing?: boolean;
  totalWorkers?: number;
  onToggleSidebar?: () => void;
  // Auth props
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
            {totalWorkers === 1
              ? "1 active worker in cluster"
              : `${String(totalWorkers)} active workers in cluster`}
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
            placeholder="Search by name, status, or service..."
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

        {!isAuthLoading && (
          isAuthenticated ? (
            <button
              type="button"
              id="logout-btn"
              className={styles.logoutBtn}
              onClick={onLogout}
              aria-label="Log out"
            >
              Logout
            </button>
          ) : (
            <button
              type="button"
              id="admin-login-btn"
              className={styles.loginBtn}
              onClick={onLoginClick}
              aria-label="Admin login"
            >
              Admin Login
            </button>
          )
        )}
      </div>
    </header>
  );
}
