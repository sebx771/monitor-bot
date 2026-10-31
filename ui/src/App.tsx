import { useState, useCallback } from "react";
import { useDashboard } from "./hooks/useDashboard";
import { useAuth } from "./hooks/useAuth";
import type { LoginCredentials } from "./models/Auth";
import {
  Sidebar,
  Header,
  MetricsOverview,
  WorkerGrid,
  WorkerDetailModal,
  LoginModal,
} from "./components";
import styles from "./App.module.css";

export default function App() {
  const {
    workers,
    filteredWorkers,
    stats,
    searchQuery,
    setSearchQuery,
    isLoading,
    isRefreshing,
    loadError,
    refreshNotification,
    handleRefresh,
    selectedWorker,
    handleSelectWorker,
    handleCloseModal,
    sidebarOpen,
    handleToggleSidebar,
    handleCloseSidebar,
  } = useDashboard();

  const { isAuthenticated, isLoading: isAuthLoading, error: authError, handleLogin, handleLogout } = useAuth();

  const [isLoginModalOpen, setIsLoginModalOpen] = useState(false);

  const handleLoginSubmit = useCallback(async (credentials: LoginCredentials) => {
    await handleLogin(credentials);
    // handleLogin throws on failure, so reaching this line means success
    setIsLoginModalOpen(false);
  }, [handleLogin]);

  const handleLoginClick = useCallback(() => {
    setIsLoginModalOpen(true);
  }, []);

  const handleCloseLoginModal = useCallback(() => {
    setIsLoginModalOpen(false);
  }, []);

  return (
    <div className={styles.layout}>
      {sidebarOpen && (
        <div
          className={styles.backdrop}
          onClick={handleCloseSidebar}
          aria-hidden="true"
        />
      )}

      <Sidebar
        isOpen={sidebarOpen}
        onClose={handleCloseSidebar}
        workerCount={workers.length}
      />

      <div className={styles.mainArea}>
        <Header
          searchQuery={searchQuery}
          onSearchChange={setSearchQuery}
          onRefresh={handleRefresh}
          isRefreshing={isLoading || isRefreshing}
          totalWorkers={filteredWorkers.length}
          onToggleSidebar={handleToggleSidebar}
          isAuthenticated={isAuthenticated}
          isAuthLoading={isAuthLoading}
          onLoginClick={handleLoginClick}
          onLogout={() => { void handleLogout(); }}
        />

        <main className={styles.content}>
          {refreshNotification && (
            <div className={`${styles.toast} ${styles.toastSuccess}`} role="status">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <polyline points="20 6 9 17 4 12" />
              </svg>
              <span>{refreshNotification}</span>
            </div>
          )}

          {loadError && workers.length > 0 && (
            <div className={`${styles.toast} ${styles.toastError}`} role="alert">
              <span>Could not refresh workers: {loadError}</span>
            </div>
          )}

          {isLoading && workers.length === 0 ? (
            <div className={styles.requestState} role="status">
              Loading workers from the API...
            </div>
          ) : loadError && workers.length === 0 ? (
            <div className={styles.requestState} role="alert">
              <p>Could not load workers: {loadError}</p>
              <button
                type="button"
                className={styles.retryButton}
                onClick={handleRefresh}
                disabled={isRefreshing}
              >
                Retry
              </button>
            </div>
          ) : (
            <>
              <MetricsOverview stats={stats} />
              <WorkerGrid
                workers={filteredWorkers}
                onWorkerClick={handleSelectWorker}
              />
            </>
          )}
        </main>
      </div>

      {selectedWorker && (
        <WorkerDetailModal
          worker={selectedWorker}
          isOpen={Boolean(selectedWorker)}
          onClose={handleCloseModal}
          isAuthenticated={isAuthenticated}
        />
      )}

      <LoginModal
        isOpen={isLoginModalOpen}
        onClose={handleCloseLoginModal}
        onLogin={handleLoginSubmit}
        error={authError}
      />
    </div>
  );
}