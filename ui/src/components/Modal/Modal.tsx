import { useEffect, type ReactNode } from "react";
import { createPortal } from "react-dom";
import styles from "./Modal.module.css";

export interface ModalProps {
    isOpen: boolean;
    onClose: () => void;
    title?: string;
    children: ReactNode; 
}

export default function Modal({ isOpen, onClose, title, children }: ModalProps) {
    useEffect(() => {
        if (!isOpen) return;

        const handleKeyDown = (e: KeyboardEvent) => {
            if (e.key === "Escape") {
                onClose();
            }
        };

        // Prevent body scrolling when modal is active
        const previousOverflow = document.body.style.overflow;
        document.body.style.overflow = "hidden";

        window.addEventListener("keydown", handleKeyDown);
        return () => {
            document.body.style.overflow = previousOverflow;
            window.removeEventListener("keydown", handleKeyDown);
        };
    }, [isOpen, onClose]);

    if (!isOpen || typeof document === "undefined") {
        return null;
    }

    return createPortal(
        <div
            className={styles.backdrop}
            onClick={onClose}
            role="dialog"
            aria-modal="true"
        >
            <div
                className={styles.content}
                onClick={(e) => e.stopPropagation()}
            >
                <header className={styles.header}>
                    {title && <h3 className={styles.title}>{title}</h3>}
                    <button
                        type="button"
                        className={styles.closeButton}
                        onClick={onClose}
                        aria-label="Close modal"
                    >
                        ✕
                    </button>
                </header>
                <div className={styles.body}>{children}</div>
            </div>
        </div>,
        document.body
    );
}