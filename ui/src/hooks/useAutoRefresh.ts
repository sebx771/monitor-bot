import { useState, useEffect, useCallback, useRef } from "react";
import type { RefreshInterval } from "../models/Dashboard";

interface UseAutoRefreshProps {
  interval: RefreshInterval;
  onTick: () => Promise<void>;
}

export function useAutoRefresh({ interval, onTick }: UseAutoRefreshProps) {
  const [consecutiveErrors, setConsecutiveErrors] = useState(0);
  const [isPausedDueToErrors, setIsPausedDueToErrors] = useState(false);
  const onTickRef = useRef(onTick);

  useEffect(() => {
    onTickRef.current = onTick;
  }, [onTick]);

  const handleTick = useCallback(async () => {
    try {
      await onTickRef.current();
      setConsecutiveErrors(0);
      setIsPausedDueToErrors(false);
    } catch {
      setConsecutiveErrors((prev) => {
        const next = prev + 1;
        if (next >= 3) {
          setIsPausedDueToErrors(true);
        }
        return next;
      });
    }
  }, []);

  const resetErrorCountAndResume = useCallback(() => {
    setConsecutiveErrors(0);
    setIsPausedDueToErrors(false);
    void handleTick();
  }, [handleTick]);

  useEffect(() => {
    if (interval === 0 || isPausedDueToErrors) {
      return;
    }

    let timerId: ReturnType<typeof setInterval> | null = null;

    const startTimer = () => {
      if (!timerId && interval > 0) {
        timerId = setInterval(() => {
          void handleTick();
        }, interval * 1000);
      }
    };

    const stopTimer = () => {
      if (timerId) {
        clearInterval(timerId);
        timerId = null;
      }
    };

    const handleVisibilityChange = () => {
      if (document.hidden) {
        stopTimer();
      } else {
        void handleTick();
        startTimer();
      }
    };

    const handleOnline = () => {
      void handleTick();
      startTimer();
    };

    const handleOffline = () => {
      stopTimer();
    };

    if (!document.hidden && navigator.onLine) {
      startTimer();
    }

    document.addEventListener("visibilitychange", handleVisibilityChange);
    window.addEventListener("online", handleOnline);
    window.addEventListener("offline", handleOffline);

    return () => {
      stopTimer();
      document.removeEventListener("visibilitychange", handleVisibilityChange);
      window.removeEventListener("online", handleOnline);
      window.removeEventListener("offline", handleOffline);
    };
  }, [interval, isPausedDueToErrors, handleTick]);

  return {
    isPausedDueToErrors,
    consecutiveErrors,
    resetErrorCountAndResume,
  };
}
