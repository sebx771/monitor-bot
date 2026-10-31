import { useState } from "react";
import styles from "./Sidebar.module.css";

interface NavItemData {
  id: string;
  label: string;
  icon: "grid" | "sparkles";
  badge?: string;
  disabled?: boolean;
}

const navItems: NavItemData[] = [
  { id: "workers", label: "Workers", icon: "grid" },
  { id: "coming-soon", label: "Coming Soon", icon: "sparkles", badge: "Soon", disabled: true },
];

interface SidebarProps {
  isOpen?: boolean;
  onClose?: () => void;
  workerCount?: number;
}

export default function Sidebar({ isOpen, workerCount = 4 }: SidebarProps) {
  const [activeItem, setActiveItem] = useState("workers");

  const renderIcon = (type: NavItemData["icon"]) => {
    switch (type) {
      case "grid":
        return (
          <svg className={styles.navIcon} viewBox="0 0 24 24" fill="none">
            <rect x="3" y="3" width="7" height="7" rx="1.5" stroke="currentColor" strokeWidth="2" />
            <rect x="14" y="3" width="7" height="7" rx="1.5" stroke="currentColor" strokeWidth="2" />
            <rect x="14" y="14" width="7" height="7" rx="1.5" stroke="currentColor" strokeWidth="2" />
            <rect x="3" y="14" width="7" height="7" rx="1.5" stroke="currentColor" strokeWidth="2" />
          </svg>
        );
      case "sparkles":
        return (
          <svg className={styles.navIcon} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
            <path d="m12 3-1.9 5.8a2 2 0 0 1-1.3 1.3L3 12l5.8 1.9a2 2 0 0 1 1.3 1.3L12 21l1.9-5.8a2 2 0 0 1 1.3-1.3L21 12l-5.8-1.9a2 2 0 0 1-1.3-1.3z" />
          </svg>
        );
    }
  };

  return (
    <aside className={`${styles.sidebar} ${isOpen ? styles.sidebarOpen : ""}`}>
      <div className={styles.brand}>
        <div className={styles.brandLogo}>
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5">
            <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5" />
          </svg>
        </div>
        <div className={styles.brandInfo}>
          <h1 className={styles.brandTitle}>Monitor Bot</h1>
          <span className={styles.brandSubtitle}>Monitor & Control Hub</span>
        </div>
      </div>

      <nav className={styles.navSection} aria-label="Main navigation">
        {navItems.map((item) => {
          const isActive = activeItem === item.id;
          const displayBadge = item.id === "workers" ? String(workerCount) : item.badge;

          return (
            <button
              key={item.id}
              type="button"
              className={`${styles.navItem} ${isActive ? styles.navItemActive : ""} ${
                item.disabled ? styles.navItemDisabled : ""
              }`}
              onClick={() => {
                if (!item.disabled) {
                  setActiveItem(item.id);
                }
              }}
              disabled={item.disabled}
              aria-current={isActive ? "page" : undefined}
            >
              {renderIcon(item.icon)}
              <span>{item.label}</span>
              {displayBadge && (
                <span
                  className={`${styles.badge} ${
                    item.disabled ? styles.badgeMuted : ""
                  }`}
                >
                  {displayBadge}
                </span>
              )}
            </button>
          );
        })}
      </nav>

      <div className={styles.footer}>
        <div className={styles.systemStatus}>
          <span className={styles.statusIndicator} />
          <span className={styles.systemText}>Cluster Operational</span>
        </div>
        <div className={styles.systemSub}> UI * v1.0.0 </div>
      </div>
    </aside>
  );
}
