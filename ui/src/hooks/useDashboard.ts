import { useState, useMemo, useCallback, useEffect, useRef } from "react";
import { WorkerStatus, type Worker } from "../models/Worker";
import { fetchWorkerStatus } from "../services/workerApi";

export interface DashboardStats {
  total: number;
  operational: number;
  alertCount: number;
}

export function useDashboard() {
  const [workers, setWorkers] = useState<Worker[]>([]);
  const [searchQuery, setSearchQuery] = useState("");
  const [isLoading, setIsLoading] = useState(true);
  const [isRefreshing, setIsRefreshing] = useState(false);
  const [loadError, setLoadError] = useState<string | null>(null);
  const [selectedWorker, setSelectedWorker] = useState<Worker | null>(null);
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [refreshNotification, setRefreshNotification] = useState<string | null>(null);
  const notificationTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  const loadWorkers = useCallback(async (isRefresh: boolean, signal?: AbortSignal) => {
    if (isRefresh) {
      setIsRefreshing(true);
      setRefreshNotification(null);
    }

    try {
      const nextWorkers = await fetchWorkerStatus(signal);
      if (signal?.aborted) return;

      setWorkers(nextWorkers);
      setLoadError(null);

      if (isRefresh) {
        setRefreshNotification("Workers successfully synced.");
        if (notificationTimeoutRef.current) {
          clearTimeout(notificationTimeoutRef.current);
        }
        notificationTimeoutRef.current = setTimeout(() => {
          setRefreshNotification(null);
          notificationTimeoutRef.current = null;
        }, 3000);
      }
    } catch (error: unknown) {
      if (signal?.aborted) return;
      setRefreshNotification(null);
      setLoadError(error instanceof Error ? error.message : "Failed to load workers.");
    } finally {
      if (!signal?.aborted) {
        setIsLoading(false);
        setIsRefreshing(false);
      }
    }
  }, []);

  useEffect(() => {
    const controller = new AbortController();

    const loadInitialWorkers = async () => {
      try {
        const nextWorkers = await fetchWorkerStatus(controller.signal);
        if (controller.signal.aborted) return;

        setWorkers(nextWorkers);
        setLoadError(null);
      } catch (error: unknown) {
        if (controller.signal.aborted) return;
        setRefreshNotification(null);
        setLoadError(error instanceof Error ? error.message : "Failed to load workers.");
      } finally {
        if (!controller.signal.aborted) {
          setIsLoading(false);
          setIsRefreshing(false);
        }
      }
    };

    void loadInitialWorkers();

    return () => {
      controller.abort();
      if (notificationTimeoutRef.current) {
        clearTimeout(notificationTimeoutRef.current);
      }
    };
  }, [loadWorkers]);

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

  const handleRefresh = useCallback(() => {
    if (isLoading || isRefreshing) return;
    void loadWorkers(true);
  }, [isLoading, isRefreshing, loadWorkers]);

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
  };
}
