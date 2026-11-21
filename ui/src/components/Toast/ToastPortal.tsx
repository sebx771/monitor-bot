import { createPortal } from "react-dom";
import type { ToastMessage } from "../../models/Dashboard";
import { ToastItem } from "./ToastItem";
import styles from "./Toast.module.css";

interface ToastPortalProps {
  toasts: ToastMessage[];
  onDismiss: (id: string) => void;
}

export function ToastPortal({ toasts, onDismiss }: ToastPortalProps) {
  if (typeof document === "undefined" || toasts.length === 0) {
    return null;
  }

  return createPortal(
    <section className={styles.portalContainer} aria-label="Notifications" role="region">
      {toasts.map((toast) => (
        <ToastItem key={toast.id} toast={toast} onDismiss={onDismiss} />
      ))}
    </section>,
    document.body
  );
}

export default ToastPortal;
