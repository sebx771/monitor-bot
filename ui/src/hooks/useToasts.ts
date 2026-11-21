import { useState, useCallback, useRef, useEffect } from "react";
import type { ToastMessage, ToastType } from "../models/Dashboard";

const MAX_TOASTS = 4;
const DEFAULT_DURATION_MS = 3500;

export function useToasts() {
  const [toasts, setToasts] = useState<ToastMessage[]>([]);
  const timeoutsRef = useRef<Map<string, ReturnType<typeof setTimeout>>>(new Map());

  const dismissToast = useCallback((id: string) => {
    const activeTimer = timeoutsRef.current.get(id);
    if (activeTimer) {
      clearTimeout(activeTimer);
      timeoutsRef.current.delete(id);
    }
    setToasts((prev) => prev.filter((toast) => toast.id !== id));
  }, []);

  const addToast = useCallback(
    (type: ToastType, text: string, durationMs = DEFAULT_DURATION_MS) => {
      const id = `${String(Date.now())}-${Math.random().toString(36).substring(2, 9)}`;
      const newToast: ToastMessage = { id, type, text, durationMs };

      setToasts((prev) => {
        const next = [...prev, newToast];
        if (next.length > MAX_TOASTS) {
          const evicted = next.shift();
          if (evicted) {
            const evictedTimer = timeoutsRef.current.get(evicted.id);
            if (evictedTimer) {
              clearTimeout(evictedTimer);
              timeoutsRef.current.delete(evicted.id);
            }
          }
        }
        return next;
      });

      if (durationMs > 0) {
        const timer = setTimeout(() => {
          dismissToast(id);
        }, durationMs);
        timeoutsRef.current.set(id, timer);
      }
    },
    [dismissToast]
  );

  useEffect(() => {
    const currentTimers = timeoutsRef.current;
    return () => {
      currentTimers.forEach((timer) => {
        clearTimeout(timer);
      });
      currentTimers.clear();
    };
  }, []);

  return { toasts, addToast, dismissToast };
}
