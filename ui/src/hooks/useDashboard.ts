import { useState, useMemo, useCallback } from "react";
import { WorkerMock } from "../mock/Worker";
import { WorkerStatus, type Worker } from "../models/Worker";

export interface DashboardStats {
  total: number;
  operational: number;
  alertCount: number;
}

export function useDashboard() {
  const [workers, setWorkers] = useState<Worker[]>(WorkerMock);
  const [searchQuery, setSearchQuery] = useState("");
  const [isRefreshing, setIsRefreshing] = useState(false);
  const [selectedWorker, setSelectedWorker] = useState<Worker | null>(null);
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [refreshNotification, setRefreshNotification] = useState<string | null>(null);

  // Filter workers based on query (name, id, status or monitored service targets)
  const filteredWorkers = useMemo(() => {
    const query = searchQuery.trim().toLowerCase();
    if (!query) return workers;

    return workers.filter((worker) => {
      const matchName = worker.name.toLowerCase().includes(query);
      const matchId = worker.id.toLowerCase().includes(query);
      const matchStatus = worker.status.toLowerCase().includes(query);
      const matchServices = worker.services.some(
        (service) =>
          service.target.toLowerCase().includes(query) ||
          service.state.toLowerCase().includes(query) ||
          (service.scope && service.scope.toLowerCase().includes(query))
      );

      return matchName || matchId || matchStatus || matchServices;
    });
  }, [workers, searchQuery]);

  // Reactive cluster statistics
  const stats: DashboardStats = useMemo(() => {
    const total = workers.length;
    const operational = workers.filter((w) => w.status === WorkerStatus.OPERATIONAL).length;
    const alertCount = workers.filter((w) => w.status !== WorkerStatus.OPERATIONAL).length;
    return { total, operational, alertCount };
  }, [workers]);

  // Refresh handler with simulated fetch & notification lifecycle
  const handleRefresh = useCallback(() => {
    setIsRefreshing(true);
    setRefreshNotification(null);

    const refreshTimer = setTimeout(() => {
      setWorkers([...WorkerMock]);
      setIsRefreshing(false);
      setRefreshNotification("Workers successfully synced.");

      const notificationTimer = setTimeout(() => {
        setRefreshNotification(null);
      }, 3000);

      return () => clearTimeout(notificationTimer);
    }, 600);

    return () => clearTimeout(refreshTimer);
  }, []);

  const handleToggleSidebar = useCallback(() => {
    setSidebarOpen((prev) => !prev);
  }, []);

  const handleCloseSidebar = useCallback(() => {
    setSidebarOpen(false);
  }, []);

  const handleSelectWorker = useCallback((worker: Worker) => {
    setSelectedWorker(worker);
  }, []);

  const handleCloseModal = useCallback(() => {
    setSelectedWorker(null);
  }, []);

  return {
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
  };
}
