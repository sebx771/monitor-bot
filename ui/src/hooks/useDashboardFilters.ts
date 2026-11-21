import { useState, useMemo, useCallback } from "react";
import { WorkerStatus, type Worker } from "../models/Worker";
import type { StatusFilter, DashboardStats } from "../models/Dashboard";

export function useDashboardFilters(workers: Worker[]) {
  const [searchQuery, setSearchQuery] = useState("");
  const [statusFilter, setStatusFilter] = useState<StatusFilter>("all");

  const toggleStatusFilter = useCallback((selected: StatusFilter) => {
    setStatusFilter((prev) => (prev === selected ? "all" : selected));
  }, []);

  const clearAllFilters = useCallback(() => {
    setSearchQuery("");
    setStatusFilter("all");
  }, []);

  const hasActiveFilters = statusFilter !== "all" || searchQuery.trim().length > 0;

  const filteredWorkers = useMemo(() => {
    const query = searchQuery.trim().toLowerCase();

    return workers.filter((worker) => {
      // 1. Status Filter Predicate
      if (statusFilter === "operational" && worker.status !== WorkerStatus.OPERATIONAL) {
        return false;
      }
      if (statusFilter === "alert" && worker.status === WorkerStatus.OPERATIONAL) {
        return false;
      }

      // 2. Search Query Predicate
      if (!query) {
        return true;
      }

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
  }, [workers, searchQuery, statusFilter]);

  const stats: DashboardStats = useMemo(() => {
    const total = workers.length;
    const operational = workers.filter((w) => w.status === WorkerStatus.OPERATIONAL).length;
    const alertCount = workers.filter((w) => w.status !== WorkerStatus.OPERATIONAL).length;
    return { total, operational, alertCount };
  }, [workers]);

  return {
    searchQuery,
    setSearchQuery,
    statusFilter,
    setStatusFilter,
    toggleStatusFilter,
    clearAllFilters,
    hasActiveFilters,
    filteredWorkers,
    stats,
  };
}
