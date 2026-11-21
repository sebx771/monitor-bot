import { useState, useCallback, useEffect, useRef } from "react";
import type { Worker } from "../models/Worker";
import type { ViewDensity, RefreshInterval } from "../models/Dashboard";
import { fetchWorkerStatus } from "../services/workerApi";
import { safeStorage } from "../utils/safeStorage";
import { useDashboardFilters } from "./useDashboardFilters";
import { useAutoRefresh } from "./useAutoRefresh";
import { useToasts } from "./useToasts";

const DENSITY_STORAGE_KEY = "monitor_bot_view_density";
const CADENCE_STORAGE_KEY = "monitor_bot_refresh_rate";

export function useDashboard() {
  const [workers, setWorkers] = useState<Worker[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isRefreshing, setIsRefreshing] = useState(false);
  const [loadError, setLoadError] = useState<string | null>(null);
  const [selectedWorker, setSelectedWorker] = useState<Worker | null>(null);
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [activeConfirmationWorkerId, setActiveConfirmationWorkerId] = useState<string | null>(null);

  const [viewDensity, setViewDensity] = useState<ViewDensity>(() => {
    const saved = safeStorage.getItem(DENSITY_STORAGE_KEY);
    return saved === "table" || saved === "grid" ? saved : "grid";
  });

  const [refreshInterval, setRefreshInterval] = useState<RefreshInterval>(() => {
    const saved = Number(safeStorage.getItem(CADENCE_STORAGE_KEY));
    return saved === 15 || saved === 30 || saved === 60 ? saved : 30;
  });

  const { toasts, addToast, dismissToast } = useToasts();
  const filterState = useDashboardFilters(workers);
  const abortControllerRef = useRef<AbortController | null>(null);

  const handleDensityChange = useCallback((density: ViewDensity) => {
    setViewDensity(density);
    safeStorage.setItem(DENSITY_STORAGE_KEY, density);
  }, []);

  const handleIntervalChange = useCallback((interval: RefreshInterval) => {
    setRefreshInterval(interval);
    safeStorage.setItem(CADENCE_STORAGE_KEY, String(interval));
  }, []);

  const loadWorkers = useCallback(
    async (isManualRefresh = false) => {
      if (abortControllerRef.current) {
        abortControllerRef.current.abort();
      }
      const controller = new AbortController();
      abortControllerRef.current = controller;

      if (isManualRefresh) setIsRefreshing(true);

      try {
        const nextWorkers = await fetchWorkerStatus(controller.signal);
        if (controller.signal.aborted) return;

        setWorkers(nextWorkers);
        setLoadError(null);
        if (isManualRefresh) {
          addToast("success", "Workers successfully synced.");
        }
      } catch (error: unknown) {
        if (controller.signal.aborted) return;
        const msg = error instanceof Error ? error.message : "Failed to sync workers.";
        setLoadError(msg);
        addToast("error", `Sync error: ${msg}`);
        throw error;
      } finally {
        if (!controller.signal.aborted) {
          setIsLoading(false);
          setIsRefreshing(false);
        }
      }
    },
    [addToast]
  );

  const { isPausedDueToErrors, resetErrorCountAndResume } = useAutoRefresh({
    interval: refreshInterval,
    onTick: async () => { await loadWorkers(false); },
  });

  useEffect(() => {
    let ignore = false;
    const controller = new AbortController();

    async function initialLoad() {
      try {
        const nextWorkers = await fetchWorkerStatus(controller.signal);
        if (!ignore && !controller.signal.aborted) {
          setWorkers(nextWorkers);
          setLoadError(null);
        }
      } catch (error: unknown) {
        if (!ignore && !controller.signal.aborted) {
          setLoadError(error instanceof Error ? error.message : "Failed to load workers.");
        }
      } finally {
        if (!ignore && !controller.signal.aborted) {
          setIsLoading(false);
        }
      }
    }

    void initialLoad();

    return () => {
      ignore = true;
      controller.abort();
    };
  }, []);

  const handleManualRefresh = useCallback(() => {
    if (isLoading || isRefreshing) return;
    void loadWorkers(true);
  }, [isLoading, isRefreshing, loadWorkers]);

  return {
    workers,
    ...filterState,
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
    handleSelectWorker: (w: Worker) => { setSelectedWorker(w); },
    handleCloseModal: () => { setSelectedWorker(null); },
    sidebarOpen,
    handleToggleSidebar: () => { setSidebarOpen((prev) => !prev); },
    handleCloseSidebar: () => { setSidebarOpen(false); },
    activeConfirmationWorkerId,
    onRequestConfirmation: (id: string) => { setActiveConfirmationWorkerId(id); },
    onClearConfirmation: () => { setActiveConfirmationWorkerId(null); },
  };
}
