import { useDashboard } from "./hooks/useDashboard";
import {
  Sidebar,
  Header,
  MetricsOverview,
  WorkerGrid,
  WorkerDetailModal,
} from "./components";
import styles from "./App.module.css";

export default function App() {
  const {
    workers,
    filteredWorkers,
    stats,
    searchQuery,
    setSearchQuery,
    isRefreshing,
    refreshNotification,
    handleRefresh,
    selectedWorker,
    handleSelectWorker,
    handleCloseModal,
    sidebarOpen,
    handleToggleSidebar,
    handleCloseSidebar,
  } = useDashboard();

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
          isRefreshing={isRefreshing}
          totalWorkers={filteredWorkers.length}
          onToggleSidebar={handleToggleSidebar}
        />

        <main className={styles.content}>
          {refreshNotification && (
            <div className={`${styles.toast} ${styles.toastSuccess}`} role="status">
              <svg
                width="18"
                height="18"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="2"
              >
                <polyline points="20 6 9 17 4 12" />
              </svg>
              <span>{refreshNotification}</span>
            </div>
          )}

          <MetricsOverview stats={stats} />

          <WorkerGrid
            workers={filteredWorkers}
            onWorkerClick={handleSelectWorker}
          />
        </main>
      </div>

      {selectedWorker && (
        <WorkerDetailModal
          worker={selectedWorker}
          isOpen={Boolean(selectedWorker)}
          onClose={handleCloseModal}
        />
      )}
    </div>
  );
}