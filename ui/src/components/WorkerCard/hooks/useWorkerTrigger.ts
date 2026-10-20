import { useState, useRef, useEffect, useCallback } from "react";
import { WorkerStatus, type Worker } from "../../../models/Worker";
import { triggerWorker } from "../../../services/workerApi";

export interface TriggerNotice {
  type: "success" | "conflict" | "error";
  text: string;
}

interface UseWorkerTriggerProps {
  worker: Worker;
  onCustomTrigger?: (workerId: string) => Promise<void> | void;
}

export function useWorkerTrigger({ worker, onCustomTrigger }: UseWorkerTriggerProps) {
  const [isTriggering, setIsTriggering] = useState(false);
  const [notice, setNotice] = useState<TriggerNotice | null>(null);
  const timeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  const isOperational = worker.status === WorkerStatus.OPERATIONAL;

  // Clean up any pending timeout when component unmounts
  useEffect(() => {
    return () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
    };
  }, []);

  const clearNoticeAfterDelay = useCallback((delayMs: number = 3500) => {
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
    }
    timeoutRef.current = setTimeout(() => {
      setNotice(null);
    }, delayMs);
  }, []);

  const executeTrigger = useCallback(async () => {
    if (!isOperational || isTriggering) {
      return;
    }

    setIsTriggering(true);
    setNotice(null);

    try {
      if (onCustomTrigger) {
        await onCustomTrigger(worker.id);
      } else {
        await triggerWorker(worker.id);
      }

      setNotice({
        type: "success",
        text: "Cycle queued successfully!",
      });
      clearNoticeAfterDelay(3000);
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : "Failed to trigger worker";
      const isConflict =
        message.toLowerCase().includes("ocupado") ||
        message.toLowerCase().includes("busy");

      setNotice({
        type: isConflict ? "conflict" : "error",
        text: message,
      });
      clearNoticeAfterDelay(4000);
    } finally {
      setIsTriggering(false);
    }
  }, [isOperational, isTriggering, onCustomTrigger, worker.id, clearNoticeAfterDelay]);

  const tooltipText = isTriggering
    ? "Sending request..."
    : !isOperational
    ? `Trigger unavailable: worker in "${worker.status}" state`
    : `Trigger worker ${worker.name} immediately`;

  return {
    isOperational,
    isTriggering,
    notice,
    executeTrigger,
    tooltipText,
  };
}
