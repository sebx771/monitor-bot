import type { DashboardStats } from "../../hooks/useDashboard";
import styles from "./MetricsOverview.module.css";

export interface MetricsOverviewProps {
  stats: DashboardStats;
}

export default function MetricsOverview({ stats }: MetricsOverviewProps) {
  return (
    <section className={styles.overviewSection} aria-label="Cluster metrics overview">
      <div className={styles.statCard}>
        <div className={styles.statInfo}>
          <span className={styles.statLabel}>Total Workers</span>
          <span className={styles.statValue}>{stats.total}</span>
        </div>
        <span className={`${styles.statPill} ${styles.pillTotal}`}>Cluster</span>
      </div>

      <div className={styles.statCard}>
        <div className={styles.statInfo}>
          <span className={styles.statLabel}>Operational</span>
          <span className={styles.statValue}>{stats.operational}</span>
        </div>
        <span className={`${styles.statPill} ${styles.pillSuccess}`}>Healthy</span>
      </div>

      <div className={styles.statCard}>
        <div className={styles.statInfo}>
          <span className={styles.statLabel}>Alert / Paused</span>
          <span className={styles.statValue}>{stats.alertCount}</span>
        </div>
        <span className={`${styles.statPill} ${styles.pillWarning}`}>Attention</span>
      </div>
    </section>
  );
}
