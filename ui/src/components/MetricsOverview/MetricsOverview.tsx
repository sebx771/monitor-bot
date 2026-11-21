import type { DashboardStats, StatusFilter } from "../../models/Dashboard";
import { MetricFilterChip } from "./MetricFilterChip";
import styles from "./MetricsOverview.module.css";

export interface MetricsOverviewProps {
  stats: DashboardStats;
  activeFilter?: StatusFilter;
  onToggleFilter?: (filter: StatusFilter) => void;
}

export default function MetricsOverview({
  stats,
  activeFilter = "all",
  onToggleFilter,
}: MetricsOverviewProps) {
  return (
    <section className={styles.overviewSection} aria-label="Cluster metrics and quick filters">
      <MetricFilterChip
        label="Total Workers"
        value={stats.total}
        pillText="Cluster"
        pillClass={styles.pillTotal}
        isActive={activeFilter === "all"}
        onClick={() => { onToggleFilter?.("all"); }}
        ariaLabel={`Filter all workers, total ${String(stats.total)}`}
      />

      <MetricFilterChip
        label="Operational"
        value={stats.operational}
        pillText="Healthy"
        pillClass={styles.pillSuccess}
        isActive={activeFilter === "operational"}
        onClick={() => { onToggleFilter?.("operational"); }}
        ariaLabel={`Filter operational workers, total ${String(stats.operational)}`}
      />

      <MetricFilterChip
        label="Alert / Paused"
        value={stats.alertCount}
        pillText="Attention"
        pillClass={styles.pillWarning}
        isActive={activeFilter === "alert"}
        onClick={() => { onToggleFilter?.("alert"); }}
        ariaLabel={`Filter degraded, paused, or error workers, total ${String(stats.alertCount)}`}
      />
    </section>
  );
}
