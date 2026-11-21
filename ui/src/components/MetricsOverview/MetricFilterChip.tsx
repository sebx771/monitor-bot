import styles from "./MetricsOverview.module.css";

interface MetricFilterChipProps {
  label: string;
  value: number;
  pillText: string;
  pillClass: string;
  isActive: boolean;
  onClick: () => void;
  ariaLabel: string;
}

export function MetricFilterChip({
  label,
  value,
  pillText,
  pillClass,
  isActive,
  onClick,
  ariaLabel,
}: MetricFilterChipProps) {
  return (
    <button
      type="button"
      className={`${styles.statCard} ${isActive ? styles.activeCard : ""}`}
      onClick={onClick}
      aria-pressed={isActive}
      aria-label={ariaLabel}
    >
      <div className={styles.statInfo}>
        <span className={styles.statLabel}>{label}</span>
        <span className={styles.statValue}>{value}</span>
      </div>
      <span className={`${styles.statPill} ${pillClass}`}>{pillText}</span>
    </button>
  );
}

export default MetricFilterChip;
