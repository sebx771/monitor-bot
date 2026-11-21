import { useState, useCallback } from "react";
import { useDashboard } from "./hooks/useDashboard";
import { useAuth } from "./hooks/useAuth";
import type { LoginCredentials } from "./models/Auth";
import {
  Sidebar,
  Header,
  MetricsOverview,
  WorkerGrid,
  WorkerTable,
  WorkerDetailModal,
  LoginModal,
  ToastPortal,
} from "./components";
import styles from "./App.module.css";

export default function App() {
  const {
    workers,
    filteredWorkers,
    stats,
    searchQuery,
    setSearchQuery,
    statusFilter,
    toggleStatusFilter,
    clearAllFilters,
    hasActiveFilters,
    isLoading,
    isRefreshing,
    loadError,
    toasts,
    dismissToast,
    viewDensity,
    handleDensityChange,
    refreshInterval,
    handleIntervalChange,
    isPausedDueToErrors,
    resetErrorCountAndResume,
    handleManualRefresh,
    selectedWorker,
    handleSelectWorker,
    handleCloseModal,
    sidebarOpen,
    handleToggleSidebar,
    handleCloseSidebar,
    activeConfirmationWorkerId,
    onRequestConfirmation,
    onClearConfirmation,
  } = useDashboard();

  const { isAuthenticated, isLoading: isAuthLoading, error: authError, handleLogin, handleLogout } = useAuth();
  const [isLoginModalOpen, setIsLoginModalOpen] = useState(false);

  const handleLoginSubmit = useCallback(async (credentials: LoginCredentials) => {
    await handleLogin(credentials);
    setIsLoginModalOpen(false);
  }, [handleLogin]);

  const handleLoginClick = useCallback(() => { setIsLoginModalOpen(true); }, []);
  const handleCloseLoginModal = useCallback(() => { setIsLoginModalOpen(false); }, []);

  return (
    <div className={styles.layout}>
      {sidebarOpen && (
        <div className={styles.backdrop} onClick={handleCloseSidebar} aria-hidden="true" />
      )}

      <Sidebar isOpen={sidebarOpen} onClose={handleCloseSidebar} workerCount={workers.length} />

      <div className={styles.mainArea}>
        <Header
          searchQuery={searchQuery}
          onSearchChange={setSearchQuery}
          onRefresh={handleManualRefresh}
          isRefreshing={isLoading || isRefreshing}
          totalWorkers={filteredWorkers.length}
          onToggleSidebar={handleToggleSidebar}
          viewDensity={viewDensity}
          onViewDensityChange={handleDensityChange}
          refreshInterval={refreshInterval}
          onIntervalChange={handleIntervalChange}
          isPausedDueToErrors={isPausedDueToErrors}
          onResetErrorsAndResume={resetErrorCountAndResume}
          isAuthenticated={isAuthenticated}
          isAuthLoading={isAuthLoading}
          onLoginClick={handleLoginClick}
          onLogout={() => { void handleLogout(); }}
        />

        <main className={styles.content}>
          {isLoading && workers.length === 0 ? (
            <div className={styles.requestState} role="status">
              Loading workers from the cluster...
            </div>
          ) : loadError && workers.length === 0 ? (
            <div className={styles.requestState} role="alert">
              <p>Could not load workers: {loadError}</p>
              <button
                type="button"
                className={styles.retryButton}
                onClick={handleManualRefresh}
                disabled={isRefreshing}
              >
                Retry
              </button>
            </div>
          ) : (
            <>
              <MetricsOverview
                stats={stats}
                activeFilter={statusFilter}
                onToggleFilter={toggleStatusFilter}
              />

              {viewDensity === "grid" ? (
                <WorkerGrid
                  workers={filteredWorkers}
                  onWorkerClick={handleSelectWorker}
                  activeConfirmationWorkerId={activeConfirmationWorkerId}
                  onRequestConfirmation={onRequestConfirmation}
                  onClearConfirmation={onClearConfirmation}
                  onClearFilters={clearAllFilters}
                  hasActiveFilters={hasActiveFilters}
                />
              ) : (
                <WorkerTable
                  workers={filteredWorkers}
                  onSelect={handleSelectWorker}
                  activeConfirmationWorkerId={activeConfirmationWorkerId}
                  onRequestConfirmation={onRequestConfirmation}
                  onClearConfirmation={onClearConfirmation}
                />
              )}
            </>
          )}
        </main>
      </div>

      <ToastPortal toasts={toasts} onDismiss={dismissToast} />

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
