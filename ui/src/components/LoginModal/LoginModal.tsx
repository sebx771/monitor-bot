import { useEffect } from "react";
import type { LoginCredentials } from "../../models/Auth";
import LoginForm from "../LoginForm/LoginForm";
import styles from "./LoginModal.module.css";

interface LoginModalProps {
  isOpen: boolean;
  onClose: () => void;
  onLogin: (credentials: LoginCredentials) => Promise<void>;
  error: string | null;
}

export default function LoginModal({ isOpen, onClose, onLogin, error }: LoginModalProps) {
  // Close on Escape key
  useEffect(() => {
    if (!isOpen) return;
    const handleKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") onClose();
    };
    document.addEventListener("keydown", handleKey);
    return () => document.removeEventListener("keydown", handleKey);
  }, [isOpen, onClose]);

  if (!isOpen) return null;

  return (
    <div
      className={styles.overlay}
      role="dialog"
      aria-modal="true"
      aria-label="Admin login"
      onClick={onClose}
    >
      {/* Stop propagation so clicking the card doesn't close the modal */}
      <div
        className={styles.content}
        onClick={(e) => e.stopPropagation()}
      >
        <LoginForm onLogin={onLogin} onClose={onClose} error={error} />
      </div>
    </div>
  );
}
